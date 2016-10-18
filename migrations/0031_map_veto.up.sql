ALTER TABLE bracket ADD COLUMN map_veto_procedure text NOT NULL DEFAULT '';

CREATE TRIGGER bracket_audit_map_veto_procedure
  AFTER UPDATE OF map_veto_procedure
  ON bracket
  FOR EACH ROW
  WHEN (new.map_veto_procedure IS DISTINCT FROM old.map_veto_procedure)
  EXECUTE PROCEDURE audit('map_veto_procedure');

CREATE SEQUENCE bracket_map_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;

CREATE TABLE bracket_map (
  id int4 NOT NULL DEFAULT nextval('bracket_map_id_seq'::regclass),
  bracket_id int4 NOT NULL,
  game_map_id int4 NOT NULL,
  sub_pool int4 NOT NULL DEFAULT 0,
  is_enabled boolean NOT NULL DEFAULT TRUE,

  created_at timestamptz NOT NULL DEFAULT now(),
  created_by int4 NOT NULL,
  updated_at timestamptz,
  updated_by int4,

  CONSTRAINT bracket_map_pkey PRIMARY KEY (id),
  CONSTRAINT bracket_map_bracket_id_fkey FOREIGN KEY (bracket_id) REFERENCES bracket (id),
  CONSTRAINT bracket_map_created_by_fkey FOREIGN KEY (created_by) REFERENCES "user" (id) ON UPDATE CASCADE,
  CONSTRAINT bracket_map_game_map_id_fkey FOREIGN KEY (game_map_id) REFERENCES game_map (id),
  CONSTRAINT bracket_map_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES "user" (id) ON UPDATE CASCADE,
  CONSTRAINT bracket_map_sub_pool_check CHECK (sub_pool >= 0 AND sub_pool <= 9)
);

CREATE UNIQUE INDEX bracket_map_game_id_token_idx ON bracket_map (bracket_id, game_map_id);

CREATE TRIGGER bracket_map_updated_at
  BEFORE UPDATE
  ON bracket_map
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER bracket_map_audit_1
  BEFORE UPDATE
  ON bracket_map
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER bracket_map_audit_2
  BEFORE UPDATE OF updated_by
  ON bracket_map
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER bracket_map_audit_3
  BEFORE UPDATE
  ON bracket_map
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER bracket_map_audit_sub_pool
  AFTER UPDATE OF sub_pool
  ON bracket_map
  FOR EACH ROW
  WHEN (new.sub_pool IS DISTINCT FROM old.sub_pool)
  EXECUTE PROCEDURE audit('sub_pool');

CREATE TRIGGER bracket_map_audit_is_enabled
  AFTER UPDATE OF is_enabled
  ON bracket_map
  FOR EACH ROW
  WHEN (new.is_enabled IS DISTINCT FROM old.is_enabled)
  EXECUTE PROCEDURE audit('is_enabled');

ALTER TABLE match_map DROP COLUMN "order";
ALTER TABLE match_map ADD COLUMN team_id int4;
ALTER TABLE match_map ADD COLUMN is_ban boolean NOT NULL DEFAULT FALSE;
ALTER TABLE match_map ADD COLUMN discarded_at timestamptz;
ALTER TABLE match_map ADD COLUMN discarded_by int4;

ALTER TABLE match_map ADD CONSTRAINT match_map_discarded_by_fkey
  FOREIGN KEY (discarded_by) REFERENCES "user" (id) ON UPDATE CASCADE;
ALTER TABLE match_map ADD CONSTRAINT match_map_team_id_fkey
  FOREIGN KEY (team_id) REFERENCES team (id) ON UPDATE CASCADE;

ALTER TABLE match ADD COLUMN are_maps_ready boolean NOT NULL DEFAULT FALSE;
