ALTER TABLE bracket DROP CONSTRAINT bracket_type_check;
ALTER TABLE bracket ADD CONSTRAINT bracket_type_check CHECK (type IN (
  'bcl-s8-group-stage', 'bcl-s8-playoffs', 'bcl-sc16-swiss', 'ace-pre-swiss'
));
