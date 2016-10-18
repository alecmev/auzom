alter table comment rename markdown to body;
alter table comment add body_raw text not null;

-- Denormalize comment_group into comment table
drop table comment_group cascade;
alter table comment add match_id integer;
alter table comment add news_id integer;
alter table comment drop column comment_group_id;

ALTER TABLE ONLY comment
    ADD CONSTRAINT comment_match_id_fkey FOREIGN KEY (match_id) REFERENCES match(id);
ALTER TABLE ONLY comment
    ADD CONSTRAINT comment_news_id_fkey FOREIGN KEY (news_id) REFERENCES news_item(id);

-- Could add index on match_id and news_id since they will be used to look up the comments.