UPDATE "user" SET fullname='';
ALTER TABLE "user" ALTER COLUMN fullname SET NOT NULL;

UPDATE "user" SET gravatar_email=email;
ALTER TABLE "user" ALTER COLUMN gravatar_email SET NOT NULL;

CREATE FUNCTION user_gravatar_email()
  RETURNS trigger AS
$BODY$
BEGIN
  NEW.gravatar_email := NEW.email;
  RETURN NEW;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE SECURITY DEFINER
  COST 100;
ALTER FUNCTION user_gravatar_email()
  OWNER TO postgres;

CREATE TRIGGER user_gravatar_email
  BEFORE INSERT
  ON "user"
  FOR EACH ROW
  EXECUTE PROCEDURE user_gravatar_email();
