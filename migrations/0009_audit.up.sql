CREATE TABLE audit
(
  id serial NOT NULL,
  table_name text NOT NULL,
  row_id integer NOT NULL,
  column_name text NOT NULL,
  changed_at timestamp with time zone NOT NULL DEFAULT now(),
  changed_from text,
  changed_by integer,
  CONSTRAINT audit_pkey PRIMARY KEY (id),
  CONSTRAINT audit_changed_by_fkey FOREIGN KEY (changed_by)
      REFERENCES "user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
)
WITH (
  OIDS=FALSE
);
ALTER TABLE audit OWNER TO postgres;

CREATE INDEX audit_table_name_row_id_column_name_idx
  ON audit
  USING btree
  (table_name COLLATE pg_catalog."default", row_id, column_name COLLATE pg_catalog."default");

CREATE FUNCTION audit_1()
  RETURNS trigger AS
$BODY$
BEGIN
  IF OLD.updated_by IS NOT DISTINCT FROM NEW.updated_by THEN
    NEW.updated_by := -1;
  END IF;
  RETURN NEW;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE SECURITY DEFINER
  COST 100;
ALTER FUNCTION audit_1() OWNER TO postgres;

CREATE FUNCTION audit_2()
  RETURNS trigger AS
$BODY$
BEGIN
  IF NEW.updated_by = -1 THEN
    NEW.updated_by := OLD.updated_by;
  END IF;
  RETURN NEW;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE SECURITY DEFINER
  COST 100;
ALTER FUNCTION audit_2() OWNER TO postgres;

CREATE FUNCTION audit_3()
  RETURNS trigger AS
$BODY$
BEGIN
  IF NEW.updated_by = -1 THEN
    NEW.updated_by := NULL;
  END IF;
  RETURN NEW;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE SECURITY DEFINER
  COST 100;
ALTER FUNCTION audit_3() OWNER TO postgres;

CREATE FUNCTION audit()
  RETURNS trigger AS
$BODY$
DECLARE
  column_name text;
  changed_from text;
BEGIN
  column_name = TG_ARGV[0]::text;
  EXECUTE 'SELECT ($1).' || column_name || '::text' INTO changed_from USING OLD;
  INSERT INTO public.audit (table_name, row_id, column_name, changed_from, changed_by)
    VALUES (TG_TABLE_NAME, OLD.id, column_name, changed_from, NEW.updated_by);
  RETURN NULL;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE SECURITY DEFINER
  COST 100;
ALTER FUNCTION audit() SET search_path=pg_catalog, public;
ALTER FUNCTION audit() OWNER TO postgres;

CREATE FUNCTION updated_at()
  RETURNS trigger AS
$BODY$
BEGIN
  NEW.updated_at := now();
  RETURN NEW;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE SECURITY DEFINER
  COST 100;
ALTER FUNCTION updated_at()
  OWNER TO postgres;

ALTER TABLE team RENAME COLUMN picture TO logo;
ALTER TABLE team ALTER COLUMN logo DROP NOT NULL;
UPDATE team SET logo=NULL WHERE logo='';

ALTER TABLE team ADD COLUMN updated_at timestamp with time zone;
ALTER TABLE team ADD COLUMN updated_by integer;
ALTER TABLE team
  ADD CONSTRAINT team_updated_by_fkey FOREIGN KEY (updated_by)
      REFERENCES "user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION;

CREATE TRIGGER team_audit_1
  BEFORE UPDATE
  ON team
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER team_audit_2
  BEFORE UPDATE OF updated_by
  ON team
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER team_audit_3
  BEFORE UPDATE
  ON team
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER team_audit_abbr
  AFTER UPDATE OF abbr
  ON team
  FOR EACH ROW
  WHEN (new.abbr IS DISTINCT FROM old.abbr)
  EXECUTE PROCEDURE audit('abbr');

CREATE TRIGGER team_audit_logo
  AFTER UPDATE OF logo
  ON team
  FOR EACH ROW
  WHEN (new.logo IS DISTINCT FROM old.logo)
  EXECUTE PROCEDURE audit('logo');

CREATE TRIGGER team_audit_name
  AFTER UPDATE OF name
  ON team
  FOR EACH ROW
  WHEN (new.name IS DISTINCT FROM old.name)
  EXECUTE PROCEDURE audit('name');

CREATE TRIGGER team_updated_at
  BEFORE UPDATE
  ON team
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

ALTER TABLE "user" ADD COLUMN updated_at timestamp with time zone;
ALTER TABLE "user" ADD COLUMN updated_by integer;
ALTER TABLE "user"
  ADD CONSTRAINT user_updated_by_fkey FOREIGN KEY (updated_by)
      REFERENCES "user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION;

CREATE TRIGGER user_audit_1
  BEFORE UPDATE
  ON "user"
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER user_audit_2
  BEFORE UPDATE OF updated_by
  ON "user"
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER user_audit_3
  BEFORE UPDATE
  ON "user"
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER user_audit_email
  AFTER UPDATE OF email
  ON "user"
  FOR EACH ROW
  WHEN (new.email IS DISTINCT FROM old.email)
  EXECUTE PROCEDURE audit('email');

CREATE TRIGGER user_audit_fullname
  AFTER UPDATE OF fullname
  ON "user"
  FOR EACH ROW
  WHEN (new.fullname IS DISTINCT FROM old.fullname)
  EXECUTE PROCEDURE audit('fullname');

CREATE TRIGGER user_audit_gravatar_email
  AFTER UPDATE OF gravatar_email
  ON "user"
  FOR EACH ROW
  WHEN (new.gravatar_email IS DISTINCT FROM old.gravatar_email)
  EXECUTE PROCEDURE audit('gravatar_email');

CREATE TRIGGER user_audit_is_admin
  AFTER UPDATE OF is_admin
  ON "user"
  FOR EACH ROW
  WHEN (new.is_admin IS DISTINCT FROM old.is_admin)
  EXECUTE PROCEDURE audit('is_admin');

CREATE TRIGGER user_audit_is_email_verified
  AFTER UPDATE OF is_email_verified
  ON "user"
  FOR EACH ROW
  WHEN (new.is_email_verified IS DISTINCT FROM old.is_email_verified)
  EXECUTE PROCEDURE audit('is_email_verified');

CREATE TRIGGER user_audit_nickname
  AFTER UPDATE OF nickname
  ON "user"
  FOR EACH ROW
  WHEN (new.nickname IS DISTINCT FROM old.nickname)
  EXECUTE PROCEDURE audit('nickname');

CREATE TRIGGER user_updated_at
  BEFORE UPDATE
  ON "user"
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();
