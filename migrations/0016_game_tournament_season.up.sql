ALTER TABLE game ADD COLUMN slug text NOT NULL;
ALTER TABLE game ADD CONSTRAINT game_slug_check CHECK (slug ~ '^[a-z0-9-]+$' AND slug NOT IN ('about', 'admin', 'login', 'password-reset', 'settings', 'signup', 'verify', 'bracket-rounds', 'brackets', 'comments', 'franchises', 'game-maps', 'games', 'match-maps', 'match-rounds', 'matches', 'news-items', 'otps', 'phases', 'seasons', 'sessions', 'teams', 'tournaments', 'user-team-requests', 'user-teams', 'users'));
CREATE UNIQUE INDEX game_slug_idx ON game (slug);

ALTER TABLE game ADD COLUMN created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE game ADD COLUMN created_by integer NOT NULL;
ALTER TABLE game ADD CONSTRAINT game_created_by_fkey FOREIGN KEY (created_by) REFERENCES "user"(id);

ALTER TABLE game ADD COLUMN updated_at timestamp with time zone;
ALTER TABLE game ADD COLUMN updated_by integer;
ALTER TABLE game
  ADD CONSTRAINT game_updated_by_fkey FOREIGN KEY (updated_by)
      REFERENCES "user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE game ALTER COLUMN released_at DROP NOT NULL;
ALTER TABLE game ADD COLUMN cover text NOT NULL;
ALTER TABLE game ADD COLUMN summary text NOT NULL;
ALTER TABLE game ADD COLUMN summary_raw text NOT NULL;

CREATE TRIGGER game_updated_at
  BEFORE UPDATE
  ON game
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER game_audit_1
  BEFORE UPDATE
  ON game
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER game_audit_2
  BEFORE UPDATE OF updated_by
  ON game
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER game_audit_3
  BEFORE UPDATE
  ON game
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER game_audit_name
  AFTER UPDATE OF name
  ON game
  FOR EACH ROW
  WHEN (new.name IS DISTINCT FROM old.name)
  EXECUTE PROCEDURE audit('name');

CREATE TRIGGER game_audit_abbr
  AFTER UPDATE OF abbr
  ON game
  FOR EACH ROW
  WHEN (new.abbr IS DISTINCT FROM old.abbr)
  EXECUTE PROCEDURE audit('abbr');

CREATE TRIGGER game_audit_released_at
  AFTER UPDATE OF released_at
  ON game
  FOR EACH ROW
  WHEN (new.released_at IS DISTINCT FROM old.released_at)
  EXECUTE PROCEDURE audit('released_at');

CREATE TRIGGER game_audit_franchise_id
  AFTER UPDATE OF franchise_id
  ON game
  FOR EACH ROW
  WHEN (new.franchise_id IS DISTINCT FROM old.franchise_id)
  EXECUTE PROCEDURE audit('franchise_id');

CREATE TRIGGER game_audit_slug
  AFTER UPDATE OF slug
  ON game
  FOR EACH ROW
  WHEN (new.slug IS DISTINCT FROM old.slug)
  EXECUTE PROCEDURE audit('slug');

CREATE TRIGGER game_audit_cover
  AFTER UPDATE OF cover
  ON game
  FOR EACH ROW
  WHEN (new.cover IS DISTINCT FROM old.cover)
  EXECUTE PROCEDURE audit('cover');

CREATE TRIGGER game_audit_summary
  AFTER UPDATE OF summary
  ON game
  FOR EACH ROW
  WHEN (new.summary IS DISTINCT FROM old.summary)
  EXECUTE PROCEDURE audit('summary');

CREATE TRIGGER game_audit_summary_raw
  AFTER UPDATE OF summary_raw
  ON game
  FOR EACH ROW
  WHEN (new.summary_raw IS DISTINCT FROM old.summary_raw)
  EXECUTE PROCEDURE audit('summary_raw');

ALTER TABLE tournament ADD COLUMN slug text NOT NULL;
ALTER TABLE tournament ADD CONSTRAINT tournament_slug_check CHECK (slug ~ '^[a-z0-9-]+$' AND slug NOT IN ('settings'));
CREATE UNIQUE INDEX tournament_slug_idx ON tournament (game_id, slug);

ALTER TABLE tournament ADD COLUMN created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE tournament ADD COLUMN created_by integer NOT NULL;
ALTER TABLE tournament ADD CONSTRAINT tournament_created_by_fkey FOREIGN KEY (created_by) REFERENCES "user"(id);

ALTER TABLE tournament ADD COLUMN updated_at timestamp with time zone;
ALTER TABLE tournament ADD COLUMN updated_by integer;
ALTER TABLE tournament
  ADD CONSTRAINT tournament_updated_by_fkey FOREIGN KEY (updated_by)
      REFERENCES "user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION;

CREATE TRIGGER tournament_updated_at
  BEFORE UPDATE
  ON tournament
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER tournament_audit_1
  BEFORE UPDATE
  ON tournament
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER tournament_audit_2
  BEFORE UPDATE OF updated_by
  ON tournament
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER tournament_audit_3
  BEFORE UPDATE
  ON tournament
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER tournament_audit_name
  AFTER UPDATE OF name
  ON tournament
  FOR EACH ROW
  WHEN (new.name IS DISTINCT FROM old.name)
  EXECUTE PROCEDURE audit('name');

CREATE TRIGGER tournament_audit_abbr
  AFTER UPDATE OF abbr
  ON tournament
  FOR EACH ROW
  WHEN (new.abbr IS DISTINCT FROM old.abbr)
  EXECUTE PROCEDURE audit('abbr');

CREATE TRIGGER tournament_audit_founded_at
  AFTER UPDATE OF founded_at
  ON tournament
  FOR EACH ROW
  WHEN (new.founded_at IS DISTINCT FROM old.founded_at)
  EXECUTE PROCEDURE audit('founded_at');

CREATE TRIGGER tournament_audit_description_raw
  AFTER UPDATE OF description_raw
  ON tournament
  FOR EACH ROW
  WHEN (new.description_raw IS DISTINCT FROM old.description_raw)
  EXECUTE PROCEDURE audit('description_raw');

CREATE TRIGGER tournament_audit_slug
  AFTER UPDATE OF slug
  ON tournament
  FOR EACH ROW
  WHEN (new.slug IS DISTINCT FROM old.slug)
  EXECUTE PROCEDURE audit('slug');

ALTER TABLE season ADD COLUMN slug text NOT NULL;
ALTER TABLE season ADD CONSTRAINT season_slug_check CHECK (slug ~ '^[a-z0-9-]+$' AND slug NOT IN ('about', 'matches', 'news-and-videos', 'overview', 'rules', 'settings'));
CREATE UNIQUE INDEX season_slug_idx ON season (tournament_id, slug);

ALTER TABLE season ADD COLUMN created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE season ADD COLUMN created_by integer NOT NULL;
ALTER TABLE season ADD CONSTRAINT season_created_by_fkey FOREIGN KEY (created_by) REFERENCES "user"(id);

ALTER TABLE season ADD COLUMN updated_at timestamp with time zone;
ALTER TABLE season ADD COLUMN updated_by integer;
ALTER TABLE season
  ADD CONSTRAINT season_updated_by_fkey FOREIGN KEY (updated_by)
      REFERENCES "user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE season RENAME COLUMN starts_at TO published_at;
ALTER TABLE season ADD COLUMN signups_opened_at timestamp with time zone;
ALTER TABLE season ADD COLUMN signups_closed_at timestamp with time zone;

CREATE TRIGGER season_updated_at
  BEFORE UPDATE
  ON season
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER season_audit_1
  BEFORE UPDATE
  ON season
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER season_audit_2
  BEFORE UPDATE OF updated_by
  ON season
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER season_audit_3
  BEFORE UPDATE
  ON season
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER season_audit_name
  AFTER UPDATE OF name
  ON season
  FOR EACH ROW
  WHEN (new.name IS DISTINCT FROM old.name)
  EXECUTE PROCEDURE audit('name');

CREATE TRIGGER season_audit_abbr
  AFTER UPDATE OF abbr
  ON season
  FOR EACH ROW
  WHEN (new.abbr IS DISTINCT FROM old.abbr)
  EXECUTE PROCEDURE audit('abbr');

CREATE TRIGGER season_audit_published_at
  AFTER UPDATE OF published_at
  ON season
  FOR EACH ROW
  WHEN (new.published_at IS DISTINCT FROM old.published_at)
  EXECUTE PROCEDURE audit('published_at');

CREATE TRIGGER season_audit_rules_raw
  AFTER UPDATE OF rules_raw
  ON season
  FOR EACH ROW
  WHEN (new.rules_raw IS DISTINCT FROM old.rules_raw)
  EXECUTE PROCEDURE audit('rules_raw');

CREATE TRIGGER season_audit_slug
  AFTER UPDATE OF slug
  ON season
  FOR EACH ROW
  WHEN (new.slug IS DISTINCT FROM old.slug)
  EXECUTE PROCEDURE audit('slug');

CREATE TRIGGER season_audit_signups_opened_at
  AFTER UPDATE OF signups_opened_at
  ON season
  FOR EACH ROW
  WHEN (new.signups_opened_at IS DISTINCT FROM old.signups_opened_at)
  EXECUTE PROCEDURE audit('signups_opened_at');

CREATE TRIGGER season_audit_signups_closed_at
  AFTER UPDATE OF signups_closed_at
  ON season
  FOR EACH ROW
  WHEN (new.signups_closed_at IS DISTINCT FROM old.signups_closed_at)
  EXECUTE PROCEDURE audit('signups_closed_at');
