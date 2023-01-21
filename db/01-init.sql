CREATE SEQUENCE IF NOT EXISTS account_id;

CREATE TABLE "accounts" (
    "id" int4 NOT NULL DEFAULT nextval('account_id'::regclass),
    "balance" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE SEQUENCE IF NOT EXISTS cloud_pocket_id;

CREATE TABLE "cloud_pockets" (
    "id" int4 NOT NULL DEFAULT nextval('cloud_pocket_id'::regclass),
    "name" VARCHAR(255),
    "category" VARCHAR(255),
    "currency" VARCHAR(3),
    "balance" float8 NOT NULL DEFAULT 0
)
