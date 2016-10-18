ALTER TABLE game ADD COLUMN verification_handle text NULL;
CREATE UNIQUE INDEX game_verification_handle_idx ON public.game (verification_handle);

CREATE TRIGGER game_audit_verification_handle
  AFTER UPDATE OF verification_handle
  ON game
  FOR EACH ROW
  WHEN (new.verification_handle IS DISTINCT FROM old.verification_handle)
  EXECUTE PROCEDURE audit('verification_handle');

CREATE SEQUENCE user_game_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.user_game (
  id int4 NOT NULL DEFAULT nextval('user_game_id_seq'::regclass),
  user_id int4 NOT NULL,
  game_id int4 NOT NULL,
  token text NULL,
  data jsonb NULL,
  name text NULL,
  link text NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  verified_at timestamptz NULL,
  updated_at timestamptz NULL,
  updated_by int4,
  data_updated_at timestamptz NULL,
  data_update_requested_at timestamptz NULL,
  nullified_at timestamptz NULL,
  nullified_by int4 NULL,
  CONSTRAINT user_game_pkey PRIMARY KEY (id),
  CONSTRAINT user_game_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
  CONSTRAINT user_game_nullified_by_fkey FOREIGN KEY (nullified_by) REFERENCES public."user"(id) ON UPDATE CASCADE,
  CONSTRAINT user_game_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public."user"(id) ON UPDATE CASCADE,
  CONSTRAINT user_game_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON UPDATE CASCADE
)
WITH (
  OIDS=FALSE
);

CREATE UNIQUE INDEX user_game_user_id_game_id_idx ON public.user_game (user_id, game_id) WHERE nullified_at IS NULL;
CREATE UNIQUE INDEX user_game_game_id_token_idx ON public.user_game (game_id, token);

CREATE TRIGGER user_game_updated_at
  BEFORE UPDATE
  ON user_game
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER user_game_audit_1
  BEFORE UPDATE
  ON user_game
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER user_game_audit_2
  BEFORE UPDATE OF updated_by
  ON user_game
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER user_game_audit_3
  BEFORE UPDATE
  ON user_game
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER user_game_audit_data
  AFTER UPDATE OF data
  ON user_game
  FOR EACH ROW
  WHEN (new.data IS DISTINCT FROM old.data)
  EXECUTE PROCEDURE audit('data');

CREATE TRIGGER user_game_audit_name
  AFTER UPDATE OF name
  ON user_game
  FOR EACH ROW
  WHEN (new.name IS DISTINCT FROM old.name)
  EXECUTE PROCEDURE audit('name');

CREATE TRIGGER user_game_audit_link
  AFTER UPDATE OF link
  ON user_game
  FOR EACH ROW
  WHEN (new.link IS DISTINCT FROM old.link)
  EXECUTE PROCEDURE audit('link');
