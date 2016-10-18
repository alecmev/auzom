ALTER TABLE season DROP CONSTRAINT season_slug_check;
ALTER TABLE season ADD CONSTRAINT season_slug_check CHECK (
  slug ~ '^[a-z0-9-]+$' AND slug NOT IN (
    'about', 'overview', 'news-and-videos', 'rules', 'seasons'
  )
);
