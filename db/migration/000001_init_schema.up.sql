CREATE TABLE "article_type" (
    "uuid" UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "desc" VARCHAR
);

CREATE TABLE "transaction" (
    "uuid" UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "date" TIMESTAMP NOT NULL,
    "price" FLOAT NOT NULL
);

CREATE TABLE "article" (
    "uuid" UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "desc" VARCHAR,
    "purchase_price" FLOAT NOT NULL,
    "resell_price" FLOAT NOT NULL,
    "article_type_uuid" UUID NOT NULL,
    FOREIGN KEY ("article_type_uuid") REFERENCES "article_type"("uuid")
);

CREATE TABLE "article_transaction" (
    "uuid" UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "article_uuid" UUID NOT NULL,
    "transaction_uuid" UUID NOT NULL,
    "amount" INT NOT NULL,
    FOREIGN KEY ("article_uuid") REFERENCES "article"("uuid"),
    FOREIGN KEY ("transaction_uuid") REFERENCES "transaction"("uuid")
);

CREATE TABLE "event" (
    "uuid" UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "desc" VARCHAR,
    "from_date" TIMESTAMP NOT NULL,
    "to_date" TIMESTAMP NOT NULL
);
