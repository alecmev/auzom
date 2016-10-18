
CREATE TABLE news_item (
    id integer NOT NULL,
    title text NOT NULL,
    preview text NOT NULL,
    preview_raw text NOT NULL,
    body text NOT NULL,
    body_raw text NOT NULL,
    picture text NOT NULL,
    game_id integer,
    season_id integer,
    created_by integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    edited_by integer,
    edited_at timestamp with time zone
);

ALTER TABLE news_item OWNER TO postgres;

CREATE SEQUENCE news_item_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
    
ALTER TABLE ONLY news_item
    ADD CONSTRAINT news_item_pkey PRIMARY KEY (id);

ALTER TABLE news_item_id_seq OWNER TO postgres;
ALTER SEQUENCE news_item_id_seq OWNED BY news_item.id;
ALTER TABLE ONLY news_item ALTER COLUMN id SET DEFAULT nextval('news_item_id_seq'::regclass);

ALTER TABLE ONLY news_item
    ADD CONSTRAINT news_item_game_id_fkey FOREIGN KEY (game_id) REFERENCES game(id);
ALTER TABLE ONLY news_item
    ADD CONSTRAINT news_item_season_id_fkey FOREIGN KEY (season_id) REFERENCES season(id);
ALTER TABLE ONLY news_item
    ADD CONSTRAINT news_item_created_by_fkey FOREIGN KEY (created_by) REFERENCES "user"(id);
ALTER TABLE ONLY news_item
    ADD CONSTRAINT news_item_edited_by_fkey FOREIGN KEY (edited_by) REFERENCES "user"(id);
