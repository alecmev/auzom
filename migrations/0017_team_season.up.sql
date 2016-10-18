CREATE SEQUENCE public.team_season_request_id_seq
  INCREMENT 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;
ALTER TABLE public.team_season_request_id_seq
  OWNER TO postgres;

CREATE TABLE public.team_season_request
(
  id integer NOT NULL DEFAULT nextval('team_season_request_id_seq'::regclass),
  team_id integer NOT NULL,
  season_id integer NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  created_by integer NOT NULL,
  decision boolean,
  decided_at timestamp with time zone,
  decided_by integer,
  cancelled_by integer,
  CONSTRAINT team_season_request_pkey PRIMARY KEY (id),
  CONSTRAINT team_season_request_decided_by_fkey FOREIGN KEY (decided_by)
      REFERENCES public."user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT team_season_request_cancelled_by_fkey FOREIGN KEY (cancelled_by)
      REFERENCES public."user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT team_season_request_season_id_fkey FOREIGN KEY (season_id)
      REFERENCES public.season (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT team_season_request_team_id_fkey FOREIGN KEY (team_id)
      REFERENCES public.team (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT team_season_request_created_by_fkey FOREIGN KEY (created_by)
      REFERENCES public."user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.team_season_request
  OWNER TO postgres;

CREATE SEQUENCE public.team_season_id_seq
  INCREMENT 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;
ALTER TABLE public.team_season_id_seq
  OWNER TO postgres;

CREATE TABLE public.team_season
(
  id integer NOT NULL DEFAULT nextval('team_season_id_seq'::regclass),
  team_id integer NOT NULL,
  season_id integer NOT NULL,
  request_id integer NOT NULL,
  left_at timestamp with time zone,
  is_done boolean NOT NULL DEFAULT FALSE,
  kicked_by integer,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone,
  updated_by integer,
  CONSTRAINT team_season_pkey PRIMARY KEY (id),
  CONSTRAINT team_season_kicked_by_fkey FOREIGN KEY (kicked_by)
      REFERENCES public."user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT team_season_request_id_fkey FOREIGN KEY (request_id)
      REFERENCES public.team_season_request (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT team_season_season_id_fkey FOREIGN KEY (season_id)
      REFERENCES public.season (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT team_season_updated_by_fkey FOREIGN KEY (updated_by)
      REFERENCES public."user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT team_season_team_id_fkey FOREIGN KEY (team_id)
      REFERENCES public.team (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.team_season
  OWNER TO postgres;

CREATE TRIGGER team_season_audit_1
  BEFORE UPDATE
  ON public.team_season
  FOR EACH ROW
  EXECUTE PROCEDURE public.audit_1();

CREATE TRIGGER team_season_audit_2
  BEFORE UPDATE OF updated_by
  ON public.team_season
  FOR EACH ROW
  EXECUTE PROCEDURE public.audit_2();

CREATE TRIGGER team_season_audit_3
  BEFORE UPDATE
  ON public.team_season
  FOR EACH ROW
  EXECUTE PROCEDURE public.audit_3();

CREATE TRIGGER team_season_audit_left_at
  AFTER UPDATE OF left_at
  ON public.team_season
  FOR EACH ROW
  WHEN ((new.left_at IS DISTINCT FROM old.left_at))
  EXECUTE PROCEDURE public.audit('left_at');

CREATE TRIGGER team_season_audit_is_done
  AFTER UPDATE OF is_done
  ON public.team_season
  FOR EACH ROW
  WHEN ((new.is_done IS DISTINCT FROM old.is_done))
  EXECUTE PROCEDURE public.audit('is_done');

CREATE TRIGGER team_season_updated_at
  BEFORE UPDATE
  ON public.team_season
  FOR EACH ROW
  EXECUTE PROCEDURE public.updated_at();
