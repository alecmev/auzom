ALTER TABLE stage ADD COLUMN slug text NOT NULL;
ALTER TABLE stage ADD CONSTRAINT stage_slug_check CHECK (slug ~ '^[a-z0-9-]+$');
CREATE UNIQUE INDEX stage_slug_idx ON stage (season_id, slug);

ALTER TABLE stage ADD COLUMN created_at timestamp with time zone DEFAULT now() NOT NULL;
ALTER TABLE stage ADD COLUMN created_by integer NOT NULL;
ALTER TABLE stage
  ADD CONSTRAINT stage_created_by_fkey FOREIGN KEY (created_by)
      REFERENCES "user" (id) MATCH SIMPLE
      ON UPDATE CASCADE ON DELETE NO ACTION;

ALTER TABLE stage ADD COLUMN updated_at timestamp with time zone;
ALTER TABLE stage ADD COLUMN updated_by integer;
ALTER TABLE stage
  ADD CONSTRAINT stage_updated_by_fkey FOREIGN KEY (updated_by)
      REFERENCES "user" (id) MATCH SIMPLE
      ON UPDATE CASCADE ON DELETE NO ACTION;

ALTER TABLE stage RENAME COLUMN starts_at TO started_at;

CREATE TRIGGER stage_updated_at
  BEFORE UPDATE
  ON stage
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER stage_audit_1
  BEFORE UPDATE
  ON stage
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER stage_audit_2
  BEFORE UPDATE OF updated_by
  ON stage
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER stage_audit_3
  BEFORE UPDATE
  ON stage
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER stage_audit_name
  AFTER UPDATE OF name
  ON stage
  FOR EACH ROW
  WHEN (new.name IS DISTINCT FROM old.name)
  EXECUTE PROCEDURE audit('name');

CREATE TRIGGER stage_audit_abbr
  AFTER UPDATE OF abbr
  ON stage
  FOR EACH ROW
  WHEN (new.abbr IS DISTINCT FROM old.abbr)
  EXECUTE PROCEDURE audit('abbr');

CREATE TRIGGER stage_audit_started_at
  AFTER UPDATE OF started_at
  ON stage
  FOR EACH ROW
  WHEN (new.started_at IS DISTINCT FROM old.started_at)
  EXECUTE PROCEDURE audit('started_at');

CREATE TRIGGER stage_audit_slug
  AFTER UPDATE OF slug
  ON stage
  FOR EACH ROW
  WHEN (new.slug IS DISTINCT FROM old.slug)
  EXECUTE PROCEDURE audit('slug');
