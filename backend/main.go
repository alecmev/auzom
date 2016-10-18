package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	slacklib "github.com/bluele/slack"
	"github.com/codegangsta/cli"
	"github.com/elithrar/simple-scrypt"
	"github.com/getsentry/raven-go"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web/middleware"

	"app/api"
	"app/mail"
	"app/models"
	"app/slack"
	"app/utils"
	"app/worker"
)

func main() {
	// without this, the output contains full date / milliseconds
	log.SetFlags(log.Ltime)
	log.Println("launching...")

	staticHost := os.Getenv("STATIC_HOST")
	if staticHost == "" {
		log.Fatalln("ERROR: no static host supplied, aborting")
	}

	// the error doesn't matter, the application shouldn't fail to launch if the
	// database happens to be down, which is perfectly plausible
	db, _ := sqlx.Open(
		"postgres", `
    sslmode=`+os.Getenv("POSTGRES_SSLMODE")+`
    host=`+os.Getenv("POSTGRES_HOST")+`
    port=`+os.Getenv("POSTGRES_PORT")+`
    dbname=`+os.Getenv("POSTGRES_DBNAME")+`
    user=`+os.Getenv("POSTGRES_USER")+`
    password=`+os.Getenv("POSTGRES_PASSWORD"),
	)

	var sg *sendgrid.SGClient
	if sgKey := os.Getenv("SENDGRID"); sgKey != "" {
		sg = sendgrid.NewSendGridClientWithApiKey(sgKey)
	} else {
		log.Println("WARNING: no SendGrid key found, printing to stdout")
	}

	var hook *slacklib.WebHook
	if hookURL := os.Getenv("SLACK"); hookURL != "" {
		hook = slacklib.NewWebHook(hookURL)
	} else {
		log.Println("WARNING: no Slack hook URL found, printing to stdout")
	}

	var err error
	var sentry *raven.Client
	if sentryKey := os.Getenv("SENTRY"); sentryKey != "" {
		sentry, err = raven.New(sentryKey)
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("WARNING: no Sentry key found, printing to stdout")
	}

	env := api.New(
		staticHost,
		scrypt.Params{
			N:       16384,
			R:       8,
			P:       1,
			SaltLen: 128,
			DKLen:   256,
		},
		models.New(db),
		mail.New("auzom <support@auzom.gg>", sg),
		slack.New(staticHost, hook),
		sentry,
	)

	// without this, the "from" IP is of nginx-proxy, not the real IP
	goji.Insert(middleware.RealIP, middleware.Logger)
	goji.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"https://" + staticHost},
		AllowedMethods:   []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler)

	goji.Get("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "auzom")
	})

	// This is a potential security hole, but if somebody is able to modify our
	// environment variables, then we have much bigger problems...
	loaderToken := os.Getenv("LOADER")
	if loaderToken != "" {
		goji.Get("/"+loaderToken+"/", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, loaderToken)
		})
	}

	flag.Set("bind", ":80") // it's 8000 by default

	app := cli.NewApp()
	app.Name = "app"
	app.Usage = "some auzom stuff"
	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "launch API server",
			Action: func(c *cli.Context) {
				setupRoutes(env)
				goji.Serve()
			},
		},
		{
			Name:  "work",
			Usage: "launch background worker",
			Action: func(c *cli.Context) {
				address := os.Getenv("BF4_ADDRESS")
				password := os.Getenv("BF4_PASSWORD")
				if address == "" || password == "" {
					log.Fatalln("ERROR: battlefield 4 verification server info missing")
				}

				go goji.Serve()
				worker.Run(address, password, env)
			},
		},
		{
			Name:  "user",
			Usage: "user-related actions",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					Usage:     "create a new user",
					ArgsUsage: "<email> <password>",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "admin, a",
							Usage: "make this user an admin",
						},
					},
					Action: func(c *cli.Context) {
						args := c.Args()
						if len(args) < 2 {
							log.Fatalf("expected 2 arguments, received %d", len(args))
						}

						email, plainPassword := args[0], args[1]
						_, err = env.M.GetUserByEmail(email)
						if err == nil {
							log.Fatalf("%s is taken", email)
						}
						if err != utils.ErrNotFound {
							log.Fatal(err.Error())
						}

						password, err := scrypt.GenerateFromPassword(
							[]byte(plainPassword),
							env.Scrypt,
						)
						if err != nil {
							log.Fatal(err.Error())
						}

						err = env.M.CreateUser(&models.User{
							Email:           email,
							Password:        password,
							IsEmailVerified: true,
							UserPublic: models.UserPublic{
								IsAdmin: c.Bool("admin"),
							},
						})
						if err != nil {
							log.Fatal(err.Error())
						}
					},
				},
			},
		},
	}

	app.Run(os.Args)
	log.Println("done")
}
