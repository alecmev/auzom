alter table bracket drop column structure_type;
alter table bracket drop column scoring_type;

alter table bracket add structure_type integer;
alter table bracket add scoring_type integer;

update bracket set structure_type=0 where structure_type is null;
update bracket set scoring_type=0 where scoring_type is null;

alter table bracket alter column structure_type set not null;
alter table bracket alter column scoring_type set not null;

alter table match add seed_x integer;
alter table match add seed_y integer;

alter table match add round integer;
update match set round = 0;
alter table match alter column round set not null;
alter table match add parent_match_x_is_loser boolean;
update match set parent_match_x_is_loser = false;
alter table match alter column parent_match_x_is_loser set not null;
alter table match add parent_match_y_is_loser boolean;
update match set parent_match_y_is_loser = false;
alter table match alter column parent_match_y_is_loser set not null;