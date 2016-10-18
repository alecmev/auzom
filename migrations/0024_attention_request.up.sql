CREATE SEQUENCE attention_request_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.attention_request (
  id int4 NOT NULL DEFAULT nextval('attention_request_id_seq'::regclass),
  target text NOT NULL,
  target_id int4 NOT NULL,
  message text NOT NULL,

  created_at timestamptz NOT NULL DEFAULT now(),
  created_by int4 NOT NULL,
  team_by int4,

  claimed_first_at timestamptz,
  claimed_at timestamptz,
  claimed_by int4,
  resolved_at timestamptz,
  is_discarded boolean NOT NULL DEFAULT FALSE,

  updated_at timestamptz,
  updated_by int4,

  CONSTRAINT attention_request_pkey PRIMARY KEY (id),
  CONSTRAINT attention_request_claimed_by_fkey FOREIGN KEY (claimed_by) REFERENCES public."user"(id) ON UPDATE CASCADE,
  CONSTRAINT attention_request_created_by_fkey FOREIGN KEY (created_by) REFERENCES public."user"(id) ON UPDATE CASCADE,
  CONSTRAINT attention_request_team_by_fkey FOREIGN KEY (team_by) REFERENCES public."team"(id),
  CONSTRAINT attention_request_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public."user"(id) ON UPDATE CASCADE,
  CONSTRAINT attention_request_target_check CHECK (target IN ('match'))
)
WITH (
  OIDS=FALSE
);

CREATE INDEX attention_request_target_target_id_idx ON public.attention_request (target, target_id);

CREATE TRIGGER attention_request_updated_at
  BEFORE UPDATE
  ON attention_request
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER attention_request_audit_1
  BEFORE UPDATE
  ON attention_request
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER attention_request_audit_2
  BEFORE UPDATE OF updated_by
  ON attention_request
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER attention_request_audit_3
  BEFORE UPDATE
  ON attention_request
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER attention_request_audit_message
  AFTER UPDATE OF message
  ON attention_request
  FOR EACH ROW
  WHEN (new.message IS DISTINCT FROM old.message)
  EXECUTE PROCEDURE audit('message');

CREATE TRIGGER attention_request_audit_claimed_by
  AFTER UPDATE OF claimed_by
  ON attention_request
  FOR EACH ROW
  WHEN (new.claimed_by IS DISTINCT FROM old.claimed_by)
  EXECUTE PROCEDURE audit('claimed_by');
