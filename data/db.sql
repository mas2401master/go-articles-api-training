-- Table: public.role

-- DROP TABLE IF EXISTS public.role;

CREATE TABLE IF NOT EXISTS public.role
(
    id bigint NOT NULL DEFAULT nextval('role_id_seq'::regclass),
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    status character varying(1) COLLATE pg_catalog."default",
    create_at timestamp without time zone,
    update_at timestamp without time zone,
    CONSTRAINT role_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.role
    OWNER to msilva;

-- Table: public.users

-- DROP TABLE IF EXISTS public.users;

CREATE TABLE IF NOT EXISTS public.users
(
    id bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    username character varying(20) COLLATE pg_catalog."default" NOT NULL,
    firstname character varying(50) COLLATE pg_catalog."default",
    lastname character varying(50) COLLATE pg_catalog."default",
    password text COLLATE pg_catalog."default",
    email character varying(200) COLLATE pg_catalog."default",
    role_id bigint NOT NULL,
    status boolean,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id)
        REFERENCES public.role (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to msilva;


-- Table: public.items

-- DROP TABLE IF EXISTS public.items;

CREATE TABLE IF NOT EXISTS public.items
(
    id bigint NOT NULL DEFAULT nextval('items_id_seq'::regclass),
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    description character varying(300) COLLATE pg_catalog."default",
    price double precision DEFAULT 0,
    available boolean,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT items_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.items
    OWNER to msilva;

-- Table: public.promotion

-- DROP TABLE IF EXISTS public.promotion;

CREATE TABLE IF NOT EXISTS public.promotion
(
    id bigint NOT NULL DEFAULT nextval('promotion_id_seq'::regclass),
    code character varying(20) COLLATE pg_catalog."default" NOT NULL,
    name character varying(50) COLLATE pg_catalog."default",
    used boolean,
    discount double precision DEFAULT 0,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT promotion_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.promotion
    OWNER to msilva;


-- Table: public.orders

-- DROP TABLE IF EXISTS public.orders;

CREATE TABLE IF NOT EXISTS public.orders
(
    id bigint NOT NULL DEFAULT nextval('orders_id_seq'::regclass),
    order_number bigint NOT NULL DEFAULT nextval('orders_order_number_seq'::regclass),
    user_id bigint NOT NULL,
    promotion_id bigint,
    subtotal double precision DEFAULT 0,
    total_discount double precision DEFAULT 0,
    total double precision DEFAULT 0,
    quantity bigint DEFAULT 0,
    status character varying COLLATE pg_catalog."default",
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT orders_pkey PRIMARY KEY (id),
    CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.orders
    OWNER to msilva;


-- Table: public.order_items

-- DROP TABLE IF EXISTS public.order_items;

CREATE TABLE IF NOT EXISTS public.order_items
(
    id bigint NOT NULL DEFAULT nextval('order_items_id_seq'::regclass),
    order_id bigint NOT NULL,
    item_id bigint NOT NULL,
    price double precision DEFAULT 0,
    quantity bigint DEFAULT 0,
    total double precision DEFAULT 0,
    status character varying(20) COLLATE pg_catalog."default",
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT order_items_pkey PRIMARY KEY (id),
    CONSTRAINT order_items_item_id_fkey FOREIGN KEY (item_id)
        REFERENCES public.items (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id)
        REFERENCES public.orders (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.order_items
    OWNER to msilva;


-- Table: public.order_items

INSERT INTO public.role(name, status, create_at, update_at)
	VALUES ( 'ADMIN', 'A', NOW(), NOW());
INSERT INTO public.role(name, status, create_at, update_at)
	VALUES ( 'CLIENTE', 'A', NOW(), NOW());

INSERT INTO public.users(
	username, firstname, lastname, password, email, role_id, status, created_at, updated_at)
	VALUES ('admin', 'Miguel Angel', 'Silva Medina', 'qi94eQI4WIYCuN9p0+cO6gE8C0rAqEvIYu2KFGRVCUY', 'miguels2401@gmail.com', 1, true, now(),now());