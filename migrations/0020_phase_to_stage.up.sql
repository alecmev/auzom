ALTER TABLE phase RENAME TO stage;
ALTER TABLE stage RENAME CONSTRAINT phase_pkey TO stage_pkey;
ALTER TABLE stage RENAME CONSTRAINT phase_season_id_fkey TO stage_season_id_fkey;
ALTER SEQUENCE phase_id_seq RENAME TO stage_id_seq;
ALTER TABLE stage ALTER COLUMN id SET DEFAULT nextval('stage_id_seq'::regclass);

ALTER TABLE game DROP CONSTRAINT game_slug_check;
ALTER TABLE game ADD CONSTRAINT game_slug_check CHECK (slug ~ '^[a-z0-9-]+$' AND slug NOT IN ('about', 'admin', 'login', 'password-reset', 'settings', 'signup', 'verify', 'bracket-rounds', 'brackets', 'comments', 'franchises', 'game-maps', 'games', 'match-maps', 'match-rounds', 'matches', 'news-items', 'otps', 'seasons', 'sessions', 'stages', 'teams', 'tournaments', 'user-team-requests', 'user-teams', 'users'));

ALTER TABLE bracket RENAME COLUMN phase_id TO stage_id;
ALTER TABLE bracket RENAME CONSTRAINT bracket_phase_id_fkey TO bracket_stage_id_fkey;
