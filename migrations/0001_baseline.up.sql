SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;
CREATE TABLE bracket (
    id integer NOT NULL,
    phase_id integer NOT NULL,
    name text NOT NULL,
    "order" integer NOT NULL,
    structure_type text NOT NULL,
    scoring_type text NOT NULL
);

ALTER TABLE bracket OWNER TO postgres;

CREATE SEQUENCE bracket_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE bracket_id_seq OWNER TO postgres;
ALTER SEQUENCE bracket_id_seq OWNED BY bracket.id;

CREATE TABLE comment (
    id integer NOT NULL,
    comment_group_id integer NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    markdown text NOT NULL,
    last_edited_at timestamp with time zone
);

ALTER TABLE comment OWNER TO postgres;

CREATE TABLE comment_group (
    id integer NOT NULL,
    match_id integer
);

ALTER TABLE comment_group OWNER TO postgres;

CREATE SEQUENCE comment_group_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE comment_group_id_seq OWNER TO postgres;
ALTER SEQUENCE comment_group_id_seq OWNED BY comment_group.id;

CREATE SEQUENCE comment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE comment_id_seq OWNER TO postgres;
ALTER SEQUENCE comment_id_seq OWNED BY comment.id;

CREATE TABLE franchise (
    id integer NOT NULL,
    name text NOT NULL,
    abbr text NOT NULL
);

ALTER TABLE franchise OWNER TO postgres;

CREATE SEQUENCE franchise_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE franchise_id_seq OWNER TO postgres;
ALTER SEQUENCE franchise_id_seq OWNED BY franchise.id;

CREATE TABLE game (
    id integer NOT NULL,
    name text NOT NULL,
    abbr text NOT NULL,
    released_at timestamp with time zone NOT NULL,
    franchise_id integer
);

ALTER TABLE game OWNER TO postgres;

CREATE SEQUENCE game_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE game_id_seq OWNER TO postgres;
ALTER SEQUENCE game_id_seq OWNED BY game.id;

CREATE TABLE game_map (
    id integer NOT NULL,
    game_id integer NOT NULL,
    name text NOT NULL,
    abbr text NOT NULL,
    side_x text NOT NULL,
    side_y text NOT NULL
);

ALTER TABLE game_map OWNER TO postgres;

CREATE SEQUENCE game_map_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE game_map_id_seq OWNER TO postgres;
ALTER SEQUENCE game_map_id_seq OWNED BY game_map.id;

CREATE TABLE match (
    id integer NOT NULL,
    bracket_id integer NOT NULL,
    starts_at timestamp with time zone NOT NULL,
    team_x integer,
    team_y integer,
    parent_match_x integer,
    parent_match_y integer,
    raw_score_x integer,
    raw_score_y integer,
    comp_score_x integer,
    comp_score_y integer
);

ALTER TABLE match OWNER TO postgres;

CREATE SEQUENCE match_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE match_id_seq OWNER TO postgres;
ALTER SEQUENCE match_id_seq OWNED BY match.id;

CREATE TABLE match_map (
    id integer NOT NULL,
    match_id integer NOT NULL,
    game_map_id integer NOT NULL,
    "order" integer NOT NULL
);

ALTER TABLE match_map OWNER TO postgres;

CREATE SEQUENCE match_map_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE match_map_id_seq OWNER TO postgres;
ALTER SEQUENCE match_map_id_seq OWNED BY match.id;

CREATE TABLE otp (
    token text NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE otp OWNER TO postgres;

COMMENT ON TABLE otp IS 'Used for email verification and password resetting.';

CREATE TABLE phase (
    id integer NOT NULL,
    season_id integer NOT NULL,
    name text NOT NULL,
    starts_at timestamp with time zone NOT NULL,
    abbr text NOT NULL
);

ALTER TABLE phase OWNER TO postgres;

CREATE SEQUENCE phase_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE phase_id_seq OWNER TO postgres;
ALTER SEQUENCE phase_id_seq OWNED BY phase.id;

CREATE TABLE round (
    id integer NOT NULL,
    match_map_id integer NOT NULL,
    raw_score_x integer,
    raw_score_y integer,
    "order" integer NOT NULL,
    is_swapped boolean NOT NULL
);

ALTER TABLE round OWNER TO postgres;

CREATE SEQUENCE round_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE round_id_seq OWNER TO postgres;
ALTER SEQUENCE round_id_seq OWNED BY round.id;

CREATE TABLE season (
    id integer NOT NULL,
    tournament_id integer NOT NULL,
    name text NOT NULL,
    abbr text NOT NULL,
    starts_at timestamp with time zone,
    rules text
);

ALTER TABLE season OWNER TO postgres;

CREATE SEQUENCE season_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE season_id_seq OWNER TO postgres;
ALTER SEQUENCE season_id_seq OWNED BY season.id;

CREATE TABLE session (
    token text NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_used_at timestamp with time zone DEFAULT now() NOT NULL,
    last_used_ip text NOT NULL,
    remember boolean DEFAULT false NOT NULL
);

ALTER TABLE session OWNER TO postgres;

CREATE TABLE team (
    id integer NOT NULL,
    name text NOT NULL,
    abbr text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    created_by integer NOT NULL
);

ALTER TABLE team OWNER TO postgres;

CREATE SEQUENCE team_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE team_id_seq OWNER TO postgres;
ALTER SEQUENCE team_id_seq OWNED BY team.id;

CREATE TABLE tournament (
    id integer NOT NULL,
    game_id integer NOT NULL,
    name text NOT NULL,
    abbr text NOT NULL,
    founded_at timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE tournament OWNER TO postgres;

CREATE SEQUENCE tournament_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE tournament_id_seq OWNER TO postgres;
ALTER SEQUENCE tournament_id_seq OWNED BY tournament.id;

CREATE TABLE "user" (
    id integer NOT NULL,
    email text NOT NULL,
    nickname text NOT NULL,
    fullname text,
    is_admin boolean DEFAULT false NOT NULL,
    password bytea NOT NULL,
    is_email_verified boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_seen_at timestamp with time zone DEFAULT now() NOT NULL,
    gravatar_email text,
    gravatar text NOT NULL
);

ALTER TABLE "user" OWNER TO postgres;

CREATE SEQUENCE user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE user_id_seq OWNER TO postgres;
ALTER SEQUENCE user_id_seq OWNED BY "user".id;

CREATE TABLE user_team (
    id integer NOT NULL,
    user_id integer NOT NULL,
    team_id integer NOT NULL,
    is_leader boolean NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    left_at timestamp with time zone,
    kicked_by integer,
    request_id integer
);

ALTER TABLE user_team OWNER TO postgres;

CREATE SEQUENCE user_team_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE user_team_id_seq OWNER TO postgres;
ALTER SEQUENCE user_team_id_seq OWNED BY user_team.id;

CREATE TABLE user_team_request (
    id integer NOT NULL,
    user_id integer NOT NULL,
    team_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    decision boolean,
    decided_at timestamp with time zone,
    user_decision boolean,
    user_decided_at timestamp with time zone,
    leader_decision boolean,
    leader_decided_at timestamp with time zone,
    leader_decided_by integer,
    admin_decision boolean,
    admin_decided_at timestamp with time zone,
    admin_decided_by integer
);

CREATE SEQUENCE user_team_request_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE user_team_request_id_seq OWNER TO postgres;
ALTER SEQUENCE user_team_request_id_seq OWNED BY user_team.id;

ALTER TABLE user_team_request OWNER TO postgres;

ALTER TABLE ONLY bracket ALTER COLUMN id SET DEFAULT nextval('bracket_id_seq'::regclass);

ALTER TABLE ONLY comment ALTER COLUMN id SET DEFAULT nextval('comment_id_seq'::regclass);

ALTER TABLE ONLY comment_group ALTER COLUMN id SET DEFAULT nextval('comment_group_id_seq'::regclass);

ALTER TABLE ONLY franchise ALTER COLUMN id SET DEFAULT nextval('franchise_id_seq'::regclass);

ALTER TABLE ONLY game ALTER COLUMN id SET DEFAULT nextval('game_id_seq'::regclass);

ALTER TABLE ONLY game_map ALTER COLUMN id SET DEFAULT nextval('game_map_id_seq'::regclass);

ALTER TABLE ONLY match ALTER COLUMN id SET DEFAULT nextval('match_id_seq'::regclass);

ALTER TABLE ONLY match_map ALTER COLUMN id SET DEFAULT nextval('match_map_id_seq'::regclass);

ALTER TABLE ONLY phase ALTER COLUMN id SET DEFAULT nextval('phase_id_seq'::regclass);

ALTER TABLE ONLY round ALTER COLUMN id SET DEFAULT nextval('round_id_seq'::regclass);

ALTER TABLE ONLY season ALTER COLUMN id SET DEFAULT nextval('season_id_seq'::regclass);

ALTER TABLE ONLY team ALTER COLUMN id SET DEFAULT nextval('team_id_seq'::regclass);

ALTER TABLE ONLY tournament ALTER COLUMN id SET DEFAULT nextval('tournament_id_seq'::regclass);

ALTER TABLE ONLY "user" ALTER COLUMN id SET DEFAULT nextval('user_id_seq'::regclass);

ALTER TABLE ONLY user_team ALTER COLUMN id SET DEFAULT nextval('user_team_id_seq'::regclass);

ALTER TABLE ONLY user_team_request ALTER COLUMN id SET DEFAULT nextval('user_team_request_id_seq'::regclass);

ALTER TABLE ONLY bracket
    ADD CONSTRAINT bracket_pkey PRIMARY KEY (id);

ALTER TABLE ONLY comment_group
    ADD CONSTRAINT comment_group_pkey PRIMARY KEY (id);

ALTER TABLE ONLY comment
    ADD CONSTRAINT comment_pkey PRIMARY KEY (id);

ALTER TABLE ONLY franchise
    ADD CONSTRAINT franchise_pkey PRIMARY KEY (id);

ALTER TABLE ONLY game_map
    ADD CONSTRAINT game_map_pkey PRIMARY KEY (id);

ALTER TABLE ONLY game
    ADD CONSTRAINT game_pkey PRIMARY KEY (id);

ALTER TABLE ONLY match_map
    ADD CONSTRAINT match_map_pkey PRIMARY KEY (id);

ALTER TABLE ONLY match
    ADD CONSTRAINT match_pkey PRIMARY KEY (id);

ALTER TABLE ONLY otp
    ADD CONSTRAINT otp_pkey PRIMARY KEY (token);

ALTER TABLE ONLY phase
    ADD CONSTRAINT phase_pkey PRIMARY KEY (id);

ALTER TABLE ONLY round
    ADD CONSTRAINT round_pkey PRIMARY KEY (id);

ALTER TABLE ONLY season
    ADD CONSTRAINT season_pkey PRIMARY KEY (id);

ALTER TABLE ONLY session
    ADD CONSTRAINT session_pkey PRIMARY KEY (token);

ALTER TABLE ONLY team
    ADD CONSTRAINT team_pkey PRIMARY KEY (id);

ALTER TABLE ONLY tournament
    ADD CONSTRAINT tournament_pkey PRIMARY KEY (id);

ALTER TABLE ONLY "user"
    ADD CONSTRAINT user_email_key UNIQUE (email);

ALTER TABLE ONLY "user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);

ALTER TABLE ONLY user_team
    ADD CONSTRAINT user_team_pkey PRIMARY KEY (id);

ALTER TABLE ONLY user_team_request
    ADD CONSTRAINT user_team_request_pkey PRIMARY KEY (id);

CREATE INDEX fki_user_team_request_id_fkey ON user_team USING btree (request_id);

CREATE INDEX user_lower_idx ON "user" USING btree (lower(email));

ALTER TABLE ONLY bracket
    ADD CONSTRAINT bracket_phase_id_fkey FOREIGN KEY (phase_id) REFERENCES phase(id);

ALTER TABLE ONLY comment
    ADD CONSTRAINT comment_comment_group_id_fkey FOREIGN KEY (comment_group_id) REFERENCES comment_group(id);

ALTER TABLE ONLY comment_group
    ADD CONSTRAINT comment_group_match_id_fkey FOREIGN KEY (match_id) REFERENCES match(id);

ALTER TABLE ONLY comment
    ADD CONSTRAINT comment_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY game
    ADD CONSTRAINT game_franchise_id_fkey FOREIGN KEY (franchise_id) REFERENCES franchise(id);

ALTER TABLE ONLY game_map
    ADD CONSTRAINT game_map_game_id_fkey FOREIGN KEY (game_id) REFERENCES game(id);

ALTER TABLE ONLY match
    ADD CONSTRAINT match_bracket_id_fkey FOREIGN KEY (bracket_id) REFERENCES bracket(id);

ALTER TABLE ONLY match_map
    ADD CONSTRAINT match_map_match_id_fkey FOREIGN KEY (match_id) REFERENCES match(id);

ALTER TABLE ONLY match
    ADD CONSTRAINT match_parent_match_x_fkey FOREIGN KEY (parent_match_x) REFERENCES match(id);

ALTER TABLE ONLY match
    ADD CONSTRAINT match_parent_match_y_fkey FOREIGN KEY (parent_match_y) REFERENCES match(id);

ALTER TABLE ONLY match
    ADD CONSTRAINT match_team_x_fkey FOREIGN KEY (team_x) REFERENCES team(id);

ALTER TABLE ONLY match
    ADD CONSTRAINT match_team_y_fkey FOREIGN KEY (team_y) REFERENCES team(id);

ALTER TABLE ONLY otp
    ADD CONSTRAINT otp_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY phase
    ADD CONSTRAINT phase_season_id_fkey FOREIGN KEY (season_id) REFERENCES season(id);

ALTER TABLE ONLY round
    ADD CONSTRAINT round_match_map_id_fkey FOREIGN KEY (match_map_id) REFERENCES match_map(id);

ALTER TABLE ONLY season
    ADD CONSTRAINT season_tournament_id_fkey FOREIGN KEY (tournament_id) REFERENCES tournament(id);

ALTER TABLE ONLY session
    ADD CONSTRAINT session_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY team
    ADD CONSTRAINT team_created_by_fkey FOREIGN KEY (created_by) REFERENCES "user"(id);

ALTER TABLE ONLY tournament
    ADD CONSTRAINT tournament_game_id_fkey FOREIGN KEY (game_id) REFERENCES game(id);

ALTER TABLE ONLY user_team
    ADD CONSTRAINT user_team_kicked_by_fkey FOREIGN KEY (kicked_by) REFERENCES "user"(id);

ALTER TABLE ONLY user_team_request
    ADD CONSTRAINT user_team_request_admin_decided_by_fkey FOREIGN KEY (admin_decided_by) REFERENCES "user"(id);

ALTER TABLE ONLY user_team
    ADD CONSTRAINT user_team_request_id_fkey FOREIGN KEY (request_id) REFERENCES user_team_request(id);

ALTER TABLE ONLY user_team_request
    ADD CONSTRAINT user_team_request_leader_decided_by_fkey FOREIGN KEY (leader_decided_by) REFERENCES "user"(id);

ALTER TABLE ONLY user_team_request
    ADD CONSTRAINT user_team_request_team_id_fkey FOREIGN KEY (team_id) REFERENCES team(id);

ALTER TABLE ONLY user_team_request
    ADD CONSTRAINT user_team_request_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY user_team
    ADD CONSTRAINT user_team_team_id_fkey FOREIGN KEY (team_id) REFERENCES team(id);

ALTER TABLE ONLY user_team
    ADD CONSTRAINT user_team_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user"(id);

