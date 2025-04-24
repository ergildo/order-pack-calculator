CREATE TABLE IF NOT EXISTS pack_sizes (
	id bigserial NOT NULL,
	product_id bigint NOT NULL,
	"size" bigint NOT NULL,
	active bool DEFAULT true NOT NULL,
	CONSTRAINT pack_sizes_pkey PRIMARY KEY (id)
);