ALTER TABLE user_team ADD COLUMN updated_at timestamp with time zone;
ALTER TABLE user_team ADD COLUMN updated_by integer;
ALTER TABLE user_team
  ADD CONSTRAINT user_team_updated_by_fkey FOREIGN KEY (updated_by)
      REFERENCES "user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION;

CREATE TRIGGER user_team_updated_at
  BEFORE UPDATE
  ON user_team
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER user_team_audit_1
  BEFORE UPDATE
  ON user_team
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER user_team_audit_2
  BEFORE UPDATE OF updated_by
  ON user_team
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER user_team_audit_3
  BEFORE UPDATE
  ON user_team
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER user_team_audit_is_leader
  AFTER UPDATE OF is_leader
  ON user_team
  FOR EACH ROW
  WHEN (new.is_leader IS DISTINCT FROM old.is_leader)
  EXECUTE PROCEDURE audit('is_leader');
