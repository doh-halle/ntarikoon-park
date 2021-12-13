--
-- PostgreSQL database dump
--

-- Dumped from database version 12.9 (Ubuntu 12.9-0ubuntu0.20.04.1)
-- Dumped by pg_dump version 12.9 (Ubuntu 12.9-0ubuntu0.20.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: apartment_restrictions; Type: TABLE; Schema: public; Owner: hallecraft
--

CREATE TABLE public.apartment_restrictions (
    id integer NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    apartment_id integer NOT NULL,
    reservation_id integer,
    restriction_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.apartment_restrictions OWNER TO hallecraft;

--
-- Name: apartment_restrictions_id_seq; Type: SEQUENCE; Schema: public; Owner: hallecraft
--

CREATE SEQUENCE public.apartment_restrictions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.apartment_restrictions_id_seq OWNER TO hallecraft;

--
-- Name: apartment_restrictions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: hallecraft
--

ALTER SEQUENCE public.apartment_restrictions_id_seq OWNED BY public.apartment_restrictions.id;


--
-- Name: apartments; Type: TABLE; Schema: public; Owner: hallecraft
--

CREATE TABLE public.apartments (
    id integer NOT NULL,
    apartment_name character varying(255) DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.apartments OWNER TO hallecraft;

--
-- Name: apartments_id_seq; Type: SEQUENCE; Schema: public; Owner: hallecraft
--

CREATE SEQUENCE public.apartments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.apartments_id_seq OWNER TO hallecraft;

--
-- Name: apartments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: hallecraft
--

ALTER SEQUENCE public.apartments_id_seq OWNED BY public.apartments.id;


--
-- Name: reservations; Type: TABLE; Schema: public; Owner: hallecraft
--

CREATE TABLE public.reservations (
    id integer NOT NULL,
    first_name character varying(255) DEFAULT ''::character varying NOT NULL,
    last_name character varying(255) DEFAULT ''::character varying NOT NULL,
    email character varying(255) NOT NULL,
    phone_number character varying(255) DEFAULT ''::character varying NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    apartment_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.reservations OWNER TO hallecraft;

--
-- Name: reservations_id_seq; Type: SEQUENCE; Schema: public; Owner: hallecraft
--

CREATE SEQUENCE public.reservations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.reservations_id_seq OWNER TO hallecraft;

--
-- Name: reservations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: hallecraft
--

ALTER SEQUENCE public.reservations_id_seq OWNED BY public.reservations.id;


--
-- Name: restrictions; Type: TABLE; Schema: public; Owner: hallecraft
--

CREATE TABLE public.restrictions (
    id integer NOT NULL,
    restriction_name character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.restrictions OWNER TO hallecraft;

--
-- Name: restrictions_id_seq; Type: SEQUENCE; Schema: public; Owner: hallecraft
--

CREATE SEQUENCE public.restrictions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.restrictions_id_seq OWNER TO hallecraft;

--
-- Name: restrictions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: hallecraft
--

ALTER SEQUENCE public.restrictions_id_seq OWNED BY public.restrictions.id;


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: hallecraft
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO hallecraft;

--
-- Name: users; Type: TABLE; Schema: public; Owner: hallecraft
--

CREATE TABLE public.users (
    id integer NOT NULL,
    first_name character varying(255) DEFAULT ''::character varying NOT NULL,
    last_name character varying(255) DEFAULT ''::character varying NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(60) NOT NULL,
    access_level integer DEFAULT 1 NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.users OWNER TO hallecraft;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: hallecraft
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO hallecraft;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: hallecraft
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: apartment_restrictions id; Type: DEFAULT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.apartment_restrictions ALTER COLUMN id SET DEFAULT nextval('public.apartment_restrictions_id_seq'::regclass);


--
-- Name: apartments id; Type: DEFAULT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.apartments ALTER COLUMN id SET DEFAULT nextval('public.apartments_id_seq'::regclass);


--
-- Name: reservations id; Type: DEFAULT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.reservations ALTER COLUMN id SET DEFAULT nextval('public.reservations_id_seq'::regclass);


--
-- Name: restrictions id; Type: DEFAULT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.restrictions ALTER COLUMN id SET DEFAULT nextval('public.restrictions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: apartment_restrictions apartment_restrictions_pkey; Type: CONSTRAINT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.apartment_restrictions
    ADD CONSTRAINT apartment_restrictions_pkey PRIMARY KEY (id);


--
-- Name: apartments apartments_pkey; Type: CONSTRAINT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.apartments
    ADD CONSTRAINT apartments_pkey PRIMARY KEY (id);


--
-- Name: reservations reservations_pkey; Type: CONSTRAINT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_pkey PRIMARY KEY (id);


--
-- Name: restrictions restrictions_pkey; Type: CONSTRAINT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.restrictions
    ADD CONSTRAINT restrictions_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: apartment_restrictions_apartment_id_idx; Type: INDEX; Schema: public; Owner: hallecraft
--

CREATE INDEX apartment_restrictions_apartment_id_idx ON public.apartment_restrictions USING btree (apartment_id);


--
-- Name: apartment_restrictions_reservation_id_idx; Type: INDEX; Schema: public; Owner: hallecraft
--

CREATE INDEX apartment_restrictions_reservation_id_idx ON public.apartment_restrictions USING btree (reservation_id);


--
-- Name: apartment_restrictions_start_date_end_date_idx; Type: INDEX; Schema: public; Owner: hallecraft
--

CREATE INDEX apartment_restrictions_start_date_end_date_idx ON public.apartment_restrictions USING btree (start_date, end_date);


--
-- Name: reservations_email_idx; Type: INDEX; Schema: public; Owner: hallecraft
--

CREATE INDEX reservations_email_idx ON public.reservations USING btree (email);


--
-- Name: reservations_last_name_idx; Type: INDEX; Schema: public; Owner: hallecraft
--

CREATE INDEX reservations_last_name_idx ON public.reservations USING btree (last_name);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: hallecraft
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: users_email_idx; Type: INDEX; Schema: public; Owner: hallecraft
--

CREATE UNIQUE INDEX users_email_idx ON public.users USING btree (email);


--
-- Name: apartment_restrictions apartment_restrictions_apartments_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.apartment_restrictions
    ADD CONSTRAINT apartment_restrictions_apartments_id_fk FOREIGN KEY (apartment_id) REFERENCES public.apartments(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: apartment_restrictions apartment_restrictions_reservations_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.apartment_restrictions
    ADD CONSTRAINT apartment_restrictions_reservations_id_fk FOREIGN KEY (reservation_id) REFERENCES public.reservations(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: apartment_restrictions apartment_restrictions_restrictions_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.apartment_restrictions
    ADD CONSTRAINT apartment_restrictions_restrictions_id_fk FOREIGN KEY (restriction_id) REFERENCES public.restrictions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: reservations reservations_apartments_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: hallecraft
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_apartments_id_fk FOREIGN KEY (apartment_id) REFERENCES public.apartments(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

