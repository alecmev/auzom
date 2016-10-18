ALTER TABLE bracket_round ADD COLUMN map_veto_procedure text NOT NULL DEFAULT '';

CREATE TRIGGER bracket_round_audit_map_veto_procedure
  AFTER UPDATE OF map_veto_procedure
  ON bracket_round
  FOR EACH ROW
  WHEN (new.map_veto_procedure IS DISTINCT FROM old.map_veto_procedure)
  EXECUTE PROCEDURE audit('map_veto_procedure');
