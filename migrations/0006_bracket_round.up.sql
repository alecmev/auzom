-- bracket_round

CREATE TABLE bracket_round (
    id integer NOT NULL,
    bracket_id integer NOT NULL,
    sort_number integer NOT NULL,
    name text NOT NULL,
    description text
);

ALTER TABLE bracket_round OWNER TO postgres;

CREATE SEQUENCE bracket_round_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
    
ALTER TABLE ONLY bracket_round
    ADD CONSTRAINT bracket_round_pkey PRIMARY KEY (id);

ALTER TABLE bracket_round_id_seq OWNER TO postgres;
ALTER SEQUENCE bracket_round_id_seq OWNED BY bracket_round.id;
ALTER TABLE ONLY bracket_round ALTER COLUMN id SET DEFAULT nextval('bracket_round_id_seq'::regclass);

ALTER TABLE ONLY bracket_round
    ADD CONSTRAINT bracket_round_bracket_id_fkey FOREIGN KEY (bracket_id) REFERENCES bracket(id);

-- migrate round
alter table match drop column round;

truncate table match cascade; -- Lazy way to avoid issues with null columns and foreign keys in this case..
alter table match add bracket_round_id integer not null;
alter table match add sort_number integer not null;

ALTER TABLE ONLY match
    ADD CONSTRAINT match_bracket_round_id_fkey FOREIGN KEY (bracket_round_id) REFERENCES bracket_round(id);