
alter table team add picture text;

update team set picture='';
alter table team alter column picture set not null;