### [auzom](https://legacy.auzom.gg/)

Esports tournament platform by Auzom Entertainment.

#### Prerequisites

1. [Install
   Docker](https://docs.docker.com/engine/installation/linux/ubuntulinux/).
1. [Install Docker Compose](https://docs.docker.com/compose/install/).
1. Make sure ports `3000` and `5432` aren't taken.

#### Getting started

1. Run `docker-compose up -d` to launch the app.
1. Run `docker-compose logs -f` and wait until it stops printing stuff rapidly.
   Exit with `ctrl + c`.
1. **Optional.** Run `scripts/db-replace-dev somedump.sql` to import a database
   dump.
1. Run `scripts/db-migrate`, even if you've imported a dump.
1. The app should be available at [`localhost:3000`](http://localhost:3000/).

#### Wrapping up

1. Run `docker-compose kill` to shut down the app.
1. **Optional.** Run `docker-compose rm` to clean up (data isn't preserved).

#### Special thanks

... to Michiel De Backker AKA [Wieweet](https://github.com/wieweet), who has
helped me out with the workload and implemented a few important API endpoints /
some key model logic in the early days of Auzom.

#### Legal

The original screenshots, which background images on the landing page are
derived from, aren't made by me, but rather
[ShadowSix](https://twitter.com/Shadow6ix) and
[Berdu](https://twitter.com/Berduu), who were kind enough to allow us to use
them for Battlefield Conquest League, Aces High and Auzom itself.

All SVG files in `frontend/src/assets`, besides `logo.svg`, belong to the
respective brands/creators.

Auzom's logo is owned by Olegs Jeremejevs and has no license. You can get
permission to use it via email in my profile.

The rest of the repository is covered by [MIT License
(Expat)](https://tldrlegal.com/license/mit-license), as seen in `LICENSE.md`.
