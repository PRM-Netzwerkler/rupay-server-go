-- Add the "price" column to "article_transaction"
ALTER TABLE "article_transaction"
ADD COLUMN "price" FLOAT NOT NULL DEFAULT 0;

-- Optionally populate the "price" with current article resell prices
UPDATE "article_transaction" at
SET "price" = a.resell_price
FROM "article" a
WHERE at.article_uuid = a.uuid;

-- Optionally remove the default if you want the column to always require a value
ALTER TABLE "article_transaction"
ALTER COLUMN "price" DROP DEFAULT;