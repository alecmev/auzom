ALTER TABLE team ADD COLUMN disbanded_at timestamp with time zone;
ALTER TABLE team ADD COLUMN disbanded_by integer;
ALTER TABLE team
  ADD CONSTRAINT team_disbanded_by_fkey FOREIGN KEY (disbanded_by)
      REFERENCES "user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION;
