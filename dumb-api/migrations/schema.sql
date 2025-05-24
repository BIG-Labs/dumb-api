--
-- PostgreSQL database dump
--

-- Dumped from database version 14.15 (Homebrew)
-- Dumped by pg_dump version 14.15 (Homebrew)

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
-- Name: chain_states; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chain_states (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    chain_id character varying(255) NOT NULL,
    last_block bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.chain_states OWNER TO postgres;

--
-- Name: edge_states; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.edge_states (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    chain_id character varying(255) NOT NULL,
    token0 character varying(255) NOT NULL,
    token1 character varying(255) NOT NULL,
    pool_id character varying(255) NOT NULL,
    edge_data bytea,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.edge_states OWNER TO admin;

--
-- Name: pool_states; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pool_states (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    token0 character varying(255) NOT NULL,
    token1 character varying(255) NOT NULL,
    pair character varying(255) NOT NULL,
    factory character varying(255) NOT NULL,
    chain_id character varying(255) NOT NULL,
    status character varying(50) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.pool_states OWNER TO postgres;

--
-- Name: price_ticks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.price_ticks (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    token_in character varying(255) NOT NULL,
    token_out character varying(255) NOT NULL,
    chain character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.price_ticks OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: ticks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ticks (
    id uuid NOT NULL,
    pool_address character varying(42) NOT NULL,
    tick_index integer NOT NULL,
    liquidity_gross text NOT NULL,
    liquidity_net text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.ticks OWNER TO postgres;

--
-- Name: tokens; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    address character varying(255) NOT NULL,
    price numeric(24,8) NOT NULL,
    chain_id character varying(255) NOT NULL,
    icon text,
    name character varying(255) NOT NULL,
    symbol character varying(255) NOT NULL,
    decimals integer NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT tokens_decimals_check CHECK ((decimals >= 0))
);


ALTER TABLE public.tokens OWNER TO postgres;

--
-- Name: chain_states chain_states_chain_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chain_states
    ADD CONSTRAINT chain_states_chain_id_key UNIQUE (chain_id);


--
-- Name: chain_states chain_states_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chain_states
    ADD CONSTRAINT chain_states_pkey PRIMARY KEY (id);


--
-- Name: edge_states edge_states_chain_id_token0_token1_pool_id_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.edge_states
    ADD CONSTRAINT edge_states_chain_id_token0_token1_pool_id_key UNIQUE (chain_id, token0, token1, pool_id);


--
-- Name: edge_states edge_states_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.edge_states
    ADD CONSTRAINT edge_states_pkey PRIMARY KEY (id);


--
-- Name: pool_states pool_states_chain_id_pair_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pool_states
    ADD CONSTRAINT pool_states_chain_id_pair_key UNIQUE (chain_id, pair);


--
-- Name: pool_states pool_states_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pool_states
    ADD CONSTRAINT pool_states_pkey PRIMARY KEY (id);


--
-- Name: price_ticks price_ticks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.price_ticks
    ADD CONSTRAINT price_ticks_pkey PRIMARY KEY (id);


--
-- Name: ticks ticks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticks
    ADD CONSTRAINT ticks_pkey PRIMARY KEY (id);


--
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);


--
-- Name: idx_ticks_pool_address; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ticks_pool_address ON public.ticks USING btree (pool_address);


--
-- Name: idx_ticks_pool_tick; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_ticks_pool_tick ON public.ticks USING btree (pool_address, tick_index);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- PostgreSQL database dump complete
--

