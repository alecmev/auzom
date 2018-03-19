### DevOps

Day-to-day is 100% automated. But initial deployment - not so much.

#### Initial deployment

`auzom-legacy-jeremejevs` is the current name of the GCP project. When switching
projects, you gotta grep through this repository and replace all occurrences by
hand, unfortunately. You gotta create the project manually too, in GCP.

Terraform is used for provisioning. For GCP access, it'll look at your default
credentials, and for Cloudflare - for a file called `terraform.tfvars`, with the
following contents:

```hcl
cloudflare_email = "EMAIL"
cloudflare_token = "TOKEN"
```

Look at Terraform documentation for further info.

Just in case, `us-east1` is the only fitting "Always Free" region.

Watch out for Terraform provider version pinning, it might result in odd things
down the line. Make sure to update regularly.

Terraform's state is stored in a bucket managed by Terraform itself. This is a
bit of a chicken'n'egg, but I think it should actually work on the first
`apply`. But not sure. Can't be bothered to find out.

Once provisioned, everything is ready (including DNS), except for the data disk.
It needs to be formatted to ext4. That is a bit of an ordeal on COOS (which you
need to SSH into), but can be done. It's already mounted at `/mnt/disks/data`.

After that, you also need to supply the secret stuff, by the means of creating a
local file called `.secrets` of the following shape:

```bash
SENDGRID=TOKEN
SLACK=https://hooks.slack.com/services/TOKEN
SENTRY=https://TOKEN@app.getsentry.com/TOKEN
LOADER=loaderio-TOKEN
BF4_ADDRESS=IP:PORT
BF4_PASSWORD=PASSWORD
```

And then uploading it by calling `scripts/secrets u`. This isn't the perfect
solution (something like Vault would be amazing), but a pretty good bang-for-
buck simplicity- and security-wise. FYI, there's a symmetric command,
`scripts/secrets d`, which pulls the `.secrets` from the bucket and overwrites
your local copy with it.

Something that Terraform can't do yet is setup a build trigger (GCP's CI). You
need to do that through the GUI or `gcloud`. Just create one that points at this
repository, and gets triggered by every commit on `master`, and uses
`cloudbuild.yaml` for configuration. That's it.

Another quick gotcha (that Terraform may or may not be able to take care of,
it's a bit complicated) is that `PROJECT_ID@cloudbuild.gserviceaccount.com`
service account needs to be given two extra roles: `Compute Instance Admin
(beta)` and `Service Account User`. This is needed so that `scripts/deploy` can
SSH into an instance in the context of a build job. Done through IAM.

And trigger a build manually. This will go through the whole process of building
all images and deploying.

There's a slight quirk with Let's Encrypt. Due to some sort of a race condition
in the way the containers are started/initialized, that I have zero desire to
get to the bottom of, LE's validation requests might not go through on the first
try. To fix that, you need to do a `. scripts/docker-compose restart` through a
`gcloud compute ssh user@main` after it's done generating the Diffie-Hellman
(watch the logs). This needs to be only on the initial deployment.

The instance self-recovers when rebooted.

At this point, the app is ready, but the database is empty. To import production
data run `scripts/db-replace path/to/your/dump.sql`.

And voila!

#### Maintenance

The only thing that I haven't automated is recurring data disk snapshotting. One
reason is that it's a whole lot of bother for no gain (nobody is using this
website right now), and another reason is that the data does evolve little by
little (due to the worker updating BF4 usernames and such), so the snapshots
would keep accumulating, potentially into several dollars per month, despite
them being incremental.

However, it's super easy to do manually. Just run `gcloud compute disks snapshot
data`, when needed (e.g. before deploying a new migration).

Something to look out for is log accumulation. Currently, I'm not sure how/if
Compose handles rotation, and the logs don't go anywhere (everything sensible is
paid), so the instance can run out of space in some distant future.

#### Random gotchas

Previously, there was an HTTP-to-HTTPS rule in Cloudflare. But there's no need
in it with the new setup, since the LE helper takes care of it (it redirects to
HTTPS immediately, with an HSTS header, both the client and the API).

I wanted to take it a step further and outright refuse HTTP on the API, but it
isn't necessary, and, more importantly, causes an unresolvable issue with LE.
[More info](https://github.com/JrCs/docker-letsencrypt-nginx-proxy-companion/issues/290).
