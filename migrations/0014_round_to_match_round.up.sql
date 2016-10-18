ALTER TABLE round RENAME TO match_round;
ALTER TABLE match_round RENAME CONSTRAINT round_pkey TO match_round_pkey;
ALTER TABLE match_round RENAME CONSTRAINT round_match_map_id_fkey TO match_round_match_map_id_fkey;
ALTER SEQUENCE round_id_seq RENAME TO match_round_id_seq;
ALTER TABLE match_round ALTER COLUMN id SET DEFAULT nextval('match_round_id_seq'::regclass);
