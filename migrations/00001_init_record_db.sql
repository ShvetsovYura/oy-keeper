-- +goose Up
-- +goose StatementBegin
-- public.item определение



CREATE TABLE item (
	"uuid" uuid NOT NULL,
	"name" varchar(512) NOT NULL,
	username varchar NULL,
	"password" varchar NULL,
	url varchar NULL,
	expired_at date NULL,
	created_at date NOT NULL,
	updated_at date NOT NULL,
	user_uuid uuid NOT NULL,
	cardnum varchar NULL,
	description text NULL,
	"version" int4 NOT NULL,
	CONSTRAINT item_pk PRIMARY KEY (uuid)
);

-- Permissions
-- ALTER TABLE item OWNER TO lbman;
-- GRANT ALL ON TABLE item TO lbman;

-- file определение
CREATE TABLE file (
	"uuid" uuid NOT NULL,
	item_uuid uuid NULL,
	"path" varchar NULL,
	hash varchar NOT NULL,
	"size" int4 NOT NULL,
	"name" varchar NOT NULL,
	created_at date NOT NULL,
	updated_at date NOT NULL,
	meta text NULL,
	CONSTRAINT file_pk PRIMARY KEY (uuid)
);
-- public.file внешние ключи
ALTER TABLE "file" ADD CONSTRAINT file_item_fk FOREIGN KEY (item_uuid) REFERENCES item("uuid") ON DELETE CASCADE;

-- "attribute" определение
CREATE TABLE "attribute" (
	"uuid" uuid NOT NULL,
	item_uuid uuid NULL,
	"name" varchar NULL,
	value text NULL,
	CONSTRAINT attribute_pk PRIMARY KEY (uuid)
);
-- "attribute" внешние ключи
ALTER TABLE "attribute" ADD CONSTRAINT attribute_item_fk FOREIGN KEY (item_uuid) REFERENCES item("uuid") ON DELETE CASCADE;


CREATE TABLE "user" (
	"uuid" uuid NOT NULL,
	login varchar NOT NULL,
	otp_secret varchar NOT NULL,
	otp_auth varchar NOT NULL,
	otp_verified bool DEFAULT false NOT NULL,
	is_active varchar DEFAULT true NOT NULL,
	created_at date NOT NULL,
	updated_at date NOT NULL,
	CONSTRAINT user_pk PRIMARY KEY (uuid),
	CONSTRAINT user_unique UNIQUE (login)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "file";
DROP TABLE "attribute";
DROP TABLE "item";
DROP TABLE "user";
-- +goose StatementEnd
