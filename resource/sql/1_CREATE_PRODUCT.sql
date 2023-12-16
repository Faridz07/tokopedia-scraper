CREATE TABLE public.products (
	id serial4 NOT NULL,
	url varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
	description text NULL,
	image_link varchar(255) NULL,
	price numeric(10, 2) NULL,
	rating numeric(3, 1) NULL,
	total_rating numeric(10, 2) NULL,
	store_name varchar(255) NULL,
	scraped_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT products_pkey PRIMARY KEY (id),
	CONSTRAINT products_url_key UNIQUE (url)
);