CREATE TABLE diff (
  id serial NOT NULL,
  table_name text NOT NULL,
  row_id int4 NOT NULL,
  column_name text NOT NULL,
  reverse_patch text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  created_by int4 NOT NULL,
  CONSTRAINT diff_pkey PRIMARY KEY (id),
  CONSTRAINT diff_created_by_fkey FOREIGN KEY (created_by) REFERENCES public."user"(id) ON UPDATE CASCADE
);

CREATE INDEX diff_table_name_row_id_column_name_idx ON diff USING btree (
  table_name,
  row_id,
  column_name
);

ALTER TABLE game DROP COLUMN summary_html;
DROP TRIGGER game_audit_summary ON game;

ALTER TABLE tournament DROP COLUMN description_html;
DROP TRIGGER tournament_audit_description ON tournament;

ALTER TABLE season DROP COLUMN rules_html;
DROP TRIGGER season_audit_rules ON season;
ALTER TABLE season ADD COLUMN description text NOT NULL DEFAULT '';

ALTER TABLE comment DROP COLUMN body_html;
DROP TRIGGER comment_audit_body ON comment;
ALTER TABLE comment DROP CONSTRAINT comment_target_check;
ALTER TABLE comment ADD CONSTRAINT comment_target_check
  CHECK (target IN ('news', 'match'));

ALTER TABLE news_item DROP COLUMN preview_html;
ALTER TABLE news_item DROP COLUMN body_html;
ALTER TABLE news_item ADD COLUMN published_at timestamptz;
ALTER TABLE news_item ADD COLUMN is_deleted boolean NOT NULL DEFAULT FALSE;
ALTER TABLE news_item ALTER COLUMN created_at SET DEFAULT now();

ALTER TABLE news_item DROP COLUMN game_id;
ALTER TABLE news_item DROP COLUMN season_id;
ALTER TABLE news_item ADD COLUMN target text NOT NULL;
ALTER TABLE news_item ADD COLUMN target_id int4 NOT NULL;
ALTER TABLE news_item ADD CONSTRAINT news_item_target_check
  CHECK (target IN ('global', 'game', 'tournament', 'season'));
CREATE INDEX news_item_target_target_id_idx ON news_item (target, target_id);

ALTER TABLE news_item DROP COLUMN edited_at;
ALTER TABLE news_item DROP COLUMN edited_by;
ALTER TABLE news_item ADD COLUMN updated_at timestamptz;
ALTER TABLE news_item ADD COLUMN updated_by int4;
ALTER TABLE news_item ADD CONSTRAINT news_item_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES "user"(id) ON UPDATE CASCADE;

CREATE TRIGGER news_item_updated_at
  BEFORE UPDATE
  ON news_item
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER news_item_audit_1
  BEFORE UPDATE
  ON news_item
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER news_item_audit_2
  BEFORE UPDATE OF updated_by
  ON news_item
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER news_item_audit_3
  BEFORE UPDATE
  ON news_item
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER news_item_audit_title
  AFTER UPDATE OF title
  ON news_item
  FOR EACH ROW
  WHEN (new.title IS DISTINCT FROM old.title)
  EXECUTE PROCEDURE audit('title');

CREATE TRIGGER news_item_audit_picture
  AFTER UPDATE OF picture
  ON news_item
  FOR EACH ROW
  WHEN (new.picture IS DISTINCT FROM old.picture)
  EXECUTE PROCEDURE audit('picture');

CREATE TRIGGER news_item_audit_preview
  AFTER UPDATE OF preview
  ON news_item
  FOR EACH ROW
  WHEN (new.preview IS DISTINCT FROM old.preview)
  EXECUTE PROCEDURE audit('preview');

CREATE TRIGGER news_item_audit_published_at
  AFTER UPDATE OF published_at
  ON news_item
  FOR EACH ROW
  WHEN (new.published_at IS DISTINCT FROM old.published_at)
  EXECUTE PROCEDURE audit('published_at');

CREATE TRIGGER news_item_audit_is_deleted
  AFTER UPDATE OF is_deleted
  ON news_item
  FOR EACH ROW
  WHEN (new.is_deleted IS DISTINCT FROM old.is_deleted)
  EXECUTE PROCEDURE audit('is_deleted');

ALTER TABLE season ADD COLUMN team_size int4 NOT NULL DEFAULT 8;
ALTER TABLE season ADD COLUMN team_size_max int4 NOT NULL DEFAULT 0;
ALTER TABLE season ADD COLUMN capacity int4 NOT NULL DEFAULT 0;
ALTER TABLE season ADD COLUMN duration int4 NOT NULL DEFAULT 60;
ALTER TABLE season ADD COLUMN youtube_playlist text NOT NULL DEFAULT '';
ALTER TABLE season ALTER COLUMN team_size DROP DEFAULT;
ALTER TABLE season ALTER COLUMN team_size_max DROP DEFAULT;
ALTER TABLE season ALTER COLUMN capacity DROP DEFAULT;
ALTER TABLE season ALTER COLUMN duration DROP DEFAULT;
ALTER TABLE season ALTER COLUMN youtube_playlist DROP DEFAULT;

CREATE TRIGGER season_audit_team_size
  AFTER UPDATE OF team_size
  ON season
  FOR EACH ROW
  WHEN (new.team_size IS DISTINCT FROM old.team_size)
  EXECUTE PROCEDURE audit('team_size');

CREATE TRIGGER season_audit_team_size_max
  AFTER UPDATE OF team_size_max
  ON season
  FOR EACH ROW
  WHEN (new.team_size_max IS DISTINCT FROM old.team_size_max)
  EXECUTE PROCEDURE audit('team_size_max');

CREATE TRIGGER season_audit_capacity
  AFTER UPDATE OF capacity
  ON season
  FOR EACH ROW
  WHEN (new.capacity IS DISTINCT FROM old.capacity)
  EXECUTE PROCEDURE audit('capacity');

CREATE TRIGGER season_audit_duration
  AFTER UPDATE OF duration
  ON season
  FOR EACH ROW
  WHEN (new.duration IS DISTINCT FROM old.duration)
  EXECUTE PROCEDURE audit('duration');

CREATE TRIGGER season_audit_youtube_playlist
  AFTER UPDATE OF youtube_playlist
  ON season
  FOR EACH ROW
  WHEN (new.youtube_playlist IS DISTINCT FROM old.youtube_playlist)
  EXECUTE PROCEDURE audit('youtube_playlist');

ALTER TABLE tournament ADD COLUMN email text NOT NULL DEFAULT '';
ALTER TABLE tournament ADD COLUMN twitch text NOT NULL DEFAULT '';
ALTER TABLE tournament ADD COLUMN youtube text NOT NULL DEFAULT '';
ALTER TABLE tournament ADD COLUMN twitter text NOT NULL DEFAULT '';
ALTER TABLE tournament ADD COLUMN facebook text NOT NULL DEFAULT '';
ALTER TABLE tournament ADD COLUMN discord text NOT NULL DEFAULT '';
ALTER TABLE tournament ADD COLUMN web text NOT NULL DEFAULT '';
ALTER TABLE tournament ADD COLUMN twitch_live text NOT NULL DEFAULT '';
ALTER TABLE tournament ADD COLUMN blur text NOT NULL DEFAULT '';
ALTER TABLE tournament ADD COLUMN logo text NOT NULL DEFAULT '';
ALTER TABLE tournament ADD COLUMN logo_has_text boolean NOT NULL DEFAULT FALSE;
ALTER TABLE tournament ALTER COLUMN email DROP DEFAULT;
ALTER TABLE tournament ALTER COLUMN twitch DROP DEFAULT;
ALTER TABLE tournament ALTER COLUMN youtube DROP DEFAULT;
ALTER TABLE tournament ALTER COLUMN twitter DROP DEFAULT;
ALTER TABLE tournament ALTER COLUMN facebook DROP DEFAULT;
ALTER TABLE tournament ALTER COLUMN discord DROP DEFAULT;
ALTER TABLE tournament ALTER COLUMN web DROP DEFAULT;
ALTER TABLE tournament ALTER COLUMN twitch_live DROP DEFAULT;
ALTER TABLE tournament ALTER COLUMN blur DROP DEFAULT;
ALTER TABLE tournament ALTER COLUMN logo DROP DEFAULT;
ALTER TABLE tournament ALTER COLUMN logo_has_text DROP DEFAULT;

CREATE TRIGGER tournament_audit_email
  AFTER UPDATE OF email
  ON tournament
  FOR EACH ROW
  WHEN (new.email IS DISTINCT FROM old.email)
  EXECUTE PROCEDURE audit('email');

CREATE TRIGGER tournament_audit_twitch
  AFTER UPDATE OF twitch
  ON tournament
  FOR EACH ROW
  WHEN (new.twitch IS DISTINCT FROM old.twitch)
  EXECUTE PROCEDURE audit('twitch');

CREATE TRIGGER tournament_audit_youtube
  AFTER UPDATE OF youtube
  ON tournament
  FOR EACH ROW
  WHEN (new.youtube IS DISTINCT FROM old.youtube)
  EXECUTE PROCEDURE audit('youtube');

CREATE TRIGGER tournament_audit_twitter
  AFTER UPDATE OF twitter
  ON tournament
  FOR EACH ROW
  WHEN (new.twitter IS DISTINCT FROM old.twitter)
  EXECUTE PROCEDURE audit('twitter');

CREATE TRIGGER tournament_audit_facebook
  AFTER UPDATE OF facebook
  ON tournament
  FOR EACH ROW
  WHEN (new.facebook IS DISTINCT FROM old.facebook)
  EXECUTE PROCEDURE audit('facebook');

CREATE TRIGGER tournament_audit_discord
  AFTER UPDATE OF discord
  ON tournament
  FOR EACH ROW
  WHEN (new.discord IS DISTINCT FROM old.discord)
  EXECUTE PROCEDURE audit('discord');

CREATE TRIGGER tournament_audit_web
  AFTER UPDATE OF web
  ON tournament
  FOR EACH ROW
  WHEN (new.web IS DISTINCT FROM old.web)
  EXECUTE PROCEDURE audit('web');

CREATE TRIGGER tournament_audit_twitch_live
  AFTER UPDATE OF twitch_live
  ON tournament
  FOR EACH ROW
  WHEN (new.twitch_live IS DISTINCT FROM old.twitch_live)
  EXECUTE PROCEDURE audit('twitch_live');

CREATE TRIGGER tournament_audit_blur
  AFTER UPDATE OF blur
  ON tournament
  FOR EACH ROW
  WHEN (new.blur IS DISTINCT FROM old.blur)
  EXECUTE PROCEDURE audit('blur');

CREATE TRIGGER tournament_audit_logo
  AFTER UPDATE OF logo
  ON tournament
  FOR EACH ROW
  WHEN (new.logo IS DISTINCT FROM old.logo)
  EXECUTE PROCEDURE audit('logo');

CREATE TRIGGER tournament_audit_logo_has_text
  AFTER UPDATE OF logo_has_text
  ON tournament
  FOR EACH ROW
  WHEN (new.logo_has_text IS DISTINCT FROM old.logo_has_text)
  EXECUTE PROCEDURE audit('logo_has_text');

DROP TRIGGER team_audit_logo ON team;
ALTER TABLE team DROP COLUMN logo;
ALTER TABLE team ADD COLUMN logo text NOT NULL DEFAULT '';
ALTER TABLE team ALTER COLUMN logo DROP DEFAULT;

CREATE TRIGGER team_audit_logo
  AFTER UPDATE OF logo
  ON team
  FOR EACH ROW
  WHEN (new.logo IS DISTINCT FROM old.logo)
  EXECUTE PROCEDURE audit('logo');
