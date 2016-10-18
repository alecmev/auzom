package mail

import (
	"log"

	"github.com/russross/blackfriday"
	"github.com/sendgrid/sendgrid-go"
)

type Env struct {
	from string
	sg   *sendgrid.SGClient
}

func New(from string, sg *sendgrid.SGClient) *Env {
	var e Env
	e.from = from
	e.sg = sg
	return &e
}

func (e *Env) Send(to, subject, text string) error {
	m := sendgrid.NewMail()
	m.SetFrom(e.from)
	m.AddTo(to)
	m.SetSubject(subject)
	m.SetText(text)
	m.SetHTML(string(blackfriday.MarkdownCommon([]byte(text))))
	if e.sg == nil {
		log.Printf("m.From: %v\n", m.From)
		log.Printf("m.To: %v\n", m.To)
		log.Printf("m.Subject: %v\n", m.Subject)
		log.Printf("m.Text: %v\n", m.Text)
		log.Printf("m.HTML: %v\n", m.HTML)
		return nil
	}

	return e.sg.Send(m)
}
