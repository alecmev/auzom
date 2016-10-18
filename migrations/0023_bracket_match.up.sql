ALTER TABLE game_map ADD COLUMN side_x_abbr text NOT NULL;
ALTER TABLE game_map ADD COLUMN side_y_abbr text NOT NULL;

ALTER TABLE game_map ADD COLUMN created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE game_map ADD COLUMN created_by integer NOT NULL;
ALTER TABLE game_map ADD CONSTRAINT game_map_created_by_fkey FOREIGN KEY (created_by) REFERENCES "user"(id) ON UPDATE CASCADE;

ALTER TABLE game_map ADD COLUMN updated_at timestamp with time zone;
ALTER TABLE game_map ADD COLUMN updated_by integer;
ALTER TABLE game_map ADD CONSTRAINT game_map_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES "user"(id) ON UPDATE CASCADE;

CREATE TRIGGER game_map_updated_at
  BEFORE UPDATE
  ON game_map
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER game_map_audit_1
  BEFORE UPDATE
  ON game_map
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER game_map_audit_2
  BEFORE UPDATE OF updated_by
  ON game_map
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER game_map_audit_3
  BEFORE UPDATE
  ON game_map
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER game_map_audit_name
  AFTER UPDATE OF name
  ON game_map
  FOR EACH ROW
  WHEN (new.name IS DISTINCT FROM old.name)
  EXECUTE PROCEDURE audit('name');

CREATE TRIGGER game_map_audit_abbr
  AFTER UPDATE OF abbr
  ON game_map
  FOR EACH ROW
  WHEN (new.abbr IS DISTINCT FROM old.abbr)
  EXECUTE PROCEDURE audit('abbr');

CREATE TRIGGER game_map_audit_side_x
  AFTER UPDATE OF side_x
  ON game_map
  FOR EACH ROW
  WHEN (new.side_x IS DISTINCT FROM old.side_x)
  EXECUTE PROCEDURE audit('side_x');

CREATE TRIGGER game_map_audit_side_x_abbr
  AFTER UPDATE OF side_x_abbr
  ON game_map
  FOR EACH ROW
  WHEN (new.side_x_abbr IS DISTINCT FROM old.side_x_abbr)
  EXECUTE PROCEDURE audit('side_x_abbr');

CREATE TRIGGER game_map_audit_side_y
  AFTER UPDATE OF side_y
  ON game_map
  FOR EACH ROW
  WHEN (new.side_y IS DISTINCT FROM old.side_y)
  EXECUTE PROCEDURE audit('side_y');

CREATE TRIGGER game_map_audit_side_y_abbr
  AFTER UPDATE OF side_y_abbr
  ON game_map
  FOR EACH ROW
  WHEN (new.side_y_abbr IS DISTINCT FROM old.side_y_abbr)
  EXECUTE PROCEDURE audit('side_y_abbr');

ALTER TABLE bracket ADD COLUMN slug text NOT NULL;
ALTER TABLE bracket ADD CONSTRAINT bracket_slug_check CHECK (slug ~ '^[a-z0-9-]+$' AND slug NOT IN ('settings'));
CREATE UNIQUE INDEX bracket_slug_idx ON bracket (stage_id, slug);

ALTER TABLE bracket DROP COLUMN structure_type;
ALTER TABLE bracket DROP COLUMN scoring_type;
ALTER TABLE bracket ADD COLUMN abbr text NOT NULL;
ALTER TABLE bracket ADD COLUMN type text NOT NULL;
ALTER TABLE bracket ADD CONSTRAINT bracket_type_check CHECK (type IN ('bcl-s8-group-stage', 'bcl-s8-playoffs'));
ALTER TABLE bracket ADD COLUMN size integer NOT NULL;

ALTER TABLE bracket ADD COLUMN created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE bracket ADD COLUMN created_by integer NOT NULL;
ALTER TABLE bracket ADD CONSTRAINT bracket_created_by_fkey FOREIGN KEY (created_by) REFERENCES "user"(id) ON UPDATE CASCADE;

ALTER TABLE bracket ADD COLUMN updated_at timestamp with time zone;
ALTER TABLE bracket ADD COLUMN updated_by integer;
ALTER TABLE bracket ADD CONSTRAINT bracket_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES "user"(id) ON UPDATE CASCADE;

CREATE TRIGGER bracket_updated_at
  BEFORE UPDATE
  ON bracket
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER bracket_audit_1
  BEFORE UPDATE
  ON bracket
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER bracket_audit_2
  BEFORE UPDATE OF updated_by
  ON bracket
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER bracket_audit_3
  BEFORE UPDATE
  ON bracket
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER bracket_audit_slug
  AFTER UPDATE OF slug
  ON bracket
  FOR EACH ROW
  WHEN (new.slug IS DISTINCT FROM old.slug)
  EXECUTE PROCEDURE audit('slug');

CREATE TRIGGER bracket_audit_name
  AFTER UPDATE OF name
  ON bracket
  FOR EACH ROW
  WHEN (new.name IS DISTINCT FROM old.name)
  EXECUTE PROCEDURE audit('name');

CREATE TRIGGER bracket_audit_abbr
  AFTER UPDATE OF abbr
  ON bracket
  FOR EACH ROW
  WHEN (new.abbr IS DISTINCT FROM old.abbr)
  EXECUTE PROCEDURE audit('abbr');

CREATE TRIGGER bracket_audit_order
  AFTER UPDATE OF "order"
  ON bracket
  FOR EACH ROW
  WHEN (new.order IS DISTINCT FROM old.order)
  EXECUTE PROCEDURE audit('order');

ALTER TABLE bracket_round RENAME COLUMN sort_number TO number;
CREATE UNIQUE INDEX bracket_round_idx ON bracket_round (bracket_id, number);

ALTER TABLE bracket_round ADD COLUMN created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE bracket_round ADD COLUMN created_by integer NOT NULL;
ALTER TABLE bracket_round ADD CONSTRAINT bracket_round_created_by_fkey FOREIGN KEY (created_by) REFERENCES "user"(id) ON UPDATE CASCADE;

ALTER TABLE bracket_round ADD COLUMN updated_at timestamp with time zone;
ALTER TABLE bracket_round ADD COLUMN updated_by integer;
ALTER TABLE bracket_round ADD CONSTRAINT bracket_round_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES "user"(id) ON UPDATE CASCADE;

CREATE TRIGGER bracket_round_updated_at
  BEFORE UPDATE
  ON bracket_round
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER bracket_round_audit_1
  BEFORE UPDATE
  ON bracket_round
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER bracket_round_audit_2
  BEFORE UPDATE OF updated_by
  ON bracket_round
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER bracket_round_audit_3
  BEFORE UPDATE
  ON bracket_round
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER bracket_round_audit_name
  AFTER UPDATE OF name
  ON bracket_round
  FOR EACH ROW
  WHEN (new.name IS DISTINCT FROM old.name)
  EXECUTE PROCEDURE audit('name');

CREATE TRIGGER bracket_round_audit_description
  AFTER UPDATE OF description
  ON bracket_round
  FOR EACH ROW
  WHEN (new.description IS DISTINCT FROM old.description)
  EXECUTE PROCEDURE audit('description');

ALTER TABLE match DROP COLUMN bracket_round_id;
ALTER TABLE match ADD COLUMN bracket_round integer NOT NULL;
ALTER TABLE match RENAME COLUMN starts_at TO started_at;
ALTER TABLE match RENAME COLUMN parent_match_x TO parent_x;
ALTER TABLE match RENAME COLUMN parent_match_x_is_loser TO parent_x_is_loser;
ALTER TABLE match RENAME COLUMN parent_match_y TO parent_y;
ALTER TABLE match RENAME COLUMN parent_match_y_is_loser TO parent_y_is_loser;
ALTER TABLE match RENAME COLUMN comp_score_x TO score_x;
ALTER TABLE match RENAME COLUMN comp_score_y TO score_y;

ALTER TABLE match ALTER COLUMN score_x TYPE float8;
ALTER TABLE match ALTER COLUMN score_y TYPE float8;
ALTER TABLE match ALTER COLUMN raw_score_x TYPE float8;
ALTER TABLE match ALTER COLUMN raw_score_y TYPE float8;

ALTER TABLE match ADD COLUMN is_overridden boolean NOT NULL DEFAULT FALSE;
ALTER TABLE match ADD COLUMN is_penalized boolean NOT NULL DEFAULT FALSE;
ALTER TABLE match ADD COLUMN reporting_closed_at timestamptz;

ALTER TABLE match ADD COLUMN created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE match ADD COLUMN created_by integer NOT NULL;
ALTER TABLE match ADD CONSTRAINT match_created_by_fkey FOREIGN KEY (created_by) REFERENCES "user"(id) ON UPDATE CASCADE;

ALTER TABLE match ADD COLUMN updated_at timestamp with time zone;
ALTER TABLE match ADD COLUMN updated_by integer;
ALTER TABLE match ADD CONSTRAINT match_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES "user"(id) ON UPDATE CASCADE;

CREATE TRIGGER match_updated_at
  BEFORE UPDATE
  ON match
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER match_audit_1
  BEFORE UPDATE
  ON match
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER match_audit_2
  BEFORE UPDATE OF updated_by
  ON match
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER match_audit_3
  BEFORE UPDATE
  ON match
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER match_audit_started_at
  AFTER UPDATE OF started_at
  ON match
  FOR EACH ROW
  WHEN (new.started_at IS DISTINCT FROM old.started_at)
  EXECUTE PROCEDURE audit('started_at');

CREATE TRIGGER match_audit_reporting_closed_at
  AFTER UPDATE OF reporting_closed_at
  ON match
  FOR EACH ROW
  WHEN (new.reporting_closed_at IS DISTINCT FROM old.reporting_closed_at)
  EXECUTE PROCEDURE audit('reporting_closed_at');

ALTER TABLE match_map ADD CONSTRAINT match_map_game_map_id_fkey FOREIGN KEY (game_map_id) REFERENCES game_map(id);

ALTER TABLE match_map ADD COLUMN created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE match_map ADD COLUMN created_by integer NOT NULL;
ALTER TABLE match_map ADD CONSTRAINT match_map_created_by_fkey FOREIGN KEY (created_by) REFERENCES "user"(id) ON UPDATE CASCADE;

CREATE SEQUENCE match_report_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.match_report (
  id int4 NOT NULL DEFAULT nextval('match_report_id_seq'::regclass),
  match_id int4 NOT NULL,

  score_x float8 NOT NULL,
  score_y float8 NOT NULL,
  raw_score_x float8 NOT NULL,
  raw_score_y float8 NOT NULL,

  override_reason text NOT NULL,
  is_penal_override boolean NOT NULL,
  score_x_override float8,
  score_y_override float8,
  raw_score_x_override float8,
  raw_score_y_override float8,

  maps_played integer NOT NULL,
  maps_x integer NOT NULL,
  maps_y integer NOT NULL,
  rounds_played integer NOT NULL,
  rounds_x integer NOT NULL,
  rounds_y integer NOT NULL,

  team_by integer,
  agreed_upon_at timestamptz,
  agreed_upon_by int4,
  created_at timestamptz NOT NULL DEFAULT now(),
  created_by int4 NOT NULL,
  CONSTRAINT match_report_pkey PRIMARY KEY (id),
  CONSTRAINT match_report_agreed_upon_by_fkey FOREIGN KEY (agreed_upon_by) REFERENCES public."user"(id),
  CONSTRAINT match_report_created_by_fkey FOREIGN KEY (created_by) REFERENCES public."user"(id),
  CONSTRAINT match_report_match_id_fkey FOREIGN KEY (match_id) REFERENCES public."match"(id),
  CONSTRAINT match_report_team_by_fkey FOREIGN KEY (team_by) REFERENCES public."team"(id)
)
WITH (
  OIDS=FALSE
);

ALTER TABLE match ADD COLUMN match_report_id integer;
ALTER TABLE match ADD CONSTRAINT match_match_report_id_fkey FOREIGN KEY (match_report_id) REFERENCES match_report(id);

ALTER TABLE match_round DROP COLUMN match_map_id;
ALTER TABLE match_round ADD COLUMN match_report_id integer NOT NULL;
ALTER TABLE match_round ADD CONSTRAINT match_round_match_report_id_fkey FOREIGN KEY (match_report_id) REFERENCES match_report(id);
ALTER TABLE match_round ADD COLUMN game_map_id integer NOT NULL;
ALTER TABLE match_round ADD CONSTRAINT match_round_game_map_id_fkey FOREIGN KEY (game_map_id) REFERENCES game_map(id);
ALTER TABLE match_round DROP COLUMN "order";

ALTER TABLE match_round RENAME COLUMN is_swapped TO is_team_x_on_side_y;
ALTER TABLE match_round ADD COLUMN is_not_played boolean NOT NULL;
ALTER TABLE match_round ALTER COLUMN raw_score_x SET NOT NULL;
ALTER TABLE match_round ALTER COLUMN raw_score_y SET NOT NULL;
ALTER TABLE match_round ALTER COLUMN raw_score_x TYPE float8;
ALTER TABLE match_round ALTER COLUMN raw_score_y TYPE float8;

ALTER TABLE match_round ADD COLUMN override_reason text NOT NULL;
ALTER TABLE match_round ADD COLUMN is_penal_override boolean NOT NULL;
ALTER TABLE match_round ADD COLUMN raw_score_x_override float8;
ALTER TABLE match_round ADD COLUMN raw_score_y_override float8;

ALTER TABLE match_round ADD COLUMN created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE match_round ADD COLUMN created_by integer NOT NULL;
ALTER TABLE match_round ADD CONSTRAINT match_round_created_by_fkey FOREIGN KEY (created_by) REFERENCES "user"(id) ON UPDATE CASCADE;

CREATE SEQUENCE match_penalty_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.match_penalty (
  id int4 NOT NULL DEFAULT nextval('match_penalty_id_seq'::regclass),
  match_report_id int4 NOT NULL,
  match_round_id int4,

  reason text NOT NULL,
  score_x float8 NOT NULL,
  score_y float8 NOT NULL,
  raw_score_x float8 NOT NULL,
  raw_score_y float8 NOT NULL,

  created_at timestamptz NOT NULL DEFAULT now(),
  created_by int4 NOT NULL,
  CONSTRAINT match_penalty_pkey PRIMARY KEY (id),
  CONSTRAINT match_penalty_created_by_fkey FOREIGN KEY (created_by) REFERENCES public."user"(id) ON UPDATE CASCADE,
  CONSTRAINT match_penalty_match_report_id_fkey FOREIGN KEY (match_report_id) REFERENCES public.match_report(id),
  CONSTRAINT match_penalty_match_round_id_fkey FOREIGN KEY (match_round_id) REFERENCES public.match_round(id)
)
WITH (
  OIDS=FALSE
);
