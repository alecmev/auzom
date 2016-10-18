ALTER TABLE season ADD COLUMN ended_at timestamptz;

CREATE TRIGGER season_audit_ended_at
  AFTER UPDATE OF ended_at
  ON season
  FOR EACH ROW
  WHEN (new.ended_at IS DISTINCT FROM old.ended_at)
  EXECUTE PROCEDURE audit('ended_at');

ALTER TABLE season ADD CONSTRAINT season_timeline_check CHECK (
  (
    signups_opened_at IS NULL AND
    signups_closed_at IS NULL AND
    ended_at IS NULL
  ) OR (
    published_at IS NOT NULL AND
    signups_opened_at IS NOT NULL AND
    signups_closed_at IS NULL AND
    ended_at IS NULL AND

    published_at <= signups_opened_at
  ) OR (
    published_at IS NOT NULL AND
    signups_opened_at IS NOT NULL AND
    signups_closed_at IS NOT NULL AND
    ended_at IS NULL AND

    published_at <= signups_opened_at AND
    signups_opened_at < signups_closed_at
  ) OR (
    published_at IS NOT NULL AND
    signups_opened_at IS NOT NULL AND
    signups_closed_at IS NOT NULL AND
    ended_at IS NOT NULL AND

    published_at <= signups_opened_at AND
    signups_opened_at < signups_closed_at AND
    signups_closed_at < ended_at
  )
);

ALTER TABLE season DROP CONSTRAINT season_slug_check;
ALTER TABLE season ADD CONSTRAINT season_slug_check CHECK (
  slug ~ '^[a-z0-9-]+$' AND slug NOT IN (
    'about', 'overview', 'news-and-videos', 'rules', 'past-seasons'
  )
);

ALTER TABLE tournament DROP CONSTRAINT tournament_slug_check;
ALTER TABLE tournament ADD CONSTRAINT tournament_slug_check CHECK (
  slug ~ '^[a-z0-9-]+$'
);
