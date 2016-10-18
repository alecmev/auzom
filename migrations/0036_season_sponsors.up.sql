ALTER TABLE season ADD COLUMN sponsors text NOT NULL DEFAULT '';
ALTER TABLE season ALTER COLUMN sponsors DROP DEFAULT;

CREATE TRIGGER season_audit_sponsors
  AFTER UPDATE OF sponsors
  ON season
  FOR EACH ROW
  WHEN (new.sponsors IS DISTINCT FROM old.sponsors)
  EXECUTE PROCEDURE audit('sponsors');
