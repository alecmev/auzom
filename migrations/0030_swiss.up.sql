ALTER TABLE bracket_round ADD COLUMN bye_team_id int4;
ALTER TABLE bracket_round ADD CONSTRAINT bracket_round_bye_team_id_fkey
  FOREIGN KEY (bye_team_id) REFERENCES team(id);

ALTER TABLE bracket DROP CONSTRAINT bracket_type_check;
ALTER TABLE bracket ADD CONSTRAINT bracket_type_check
  CHECK (type IN ('bcl-s8-group-stage', 'bcl-s8-playoffs', 'bcl-sc16-swiss'));

ALTER TABLE match ADD CONSTRAINT match_same_team_check
  CHECK (team_x IS NULL OR team_y IS NULL OR team_x != team_y);
