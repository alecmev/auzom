ALTER TABLE "user" DROP CONSTRAINT user_email_key;
DROP INDEX user_lower_idx;
CREATE UNIQUE INDEX user_email_idx ON "user" (lower(email));
