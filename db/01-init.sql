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
    "parentID" int4,
    "currency" VARCHAR(3),
    "balance" float8 NOT NULL DEFAULT 0
);

CREATE SEQUENCE IF NOT EXISTS transfers_id;

CREATE TABLE "transfers" (
    "id" int4 NOT NULL DEFAULT nextval('transfers_id'::regclass),
    "pocket_id_source" int4,
    "pocket_id_target" int4,
    "amount" float8 NOT NULL DEFAULT 0
);