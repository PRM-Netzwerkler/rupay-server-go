-- Down Migration: Revert changes to restore previous constraints

BEGIN;

-- 1. Drop existing foreign key constraints
ALTER TABLE "article" 
DROP CONSTRAINT IF EXISTS "article_article_type_uuid_fkey";

ALTER TABLE "article_transaction" 
DROP CONSTRAINT IF EXISTS "article_transaction_article_uuid_fkey";

-- 2. Re-add foreign key constraints to previous state
ALTER TABLE "article" 
ADD CONSTRAINT "article_article_type_uuid_fkey"
FOREIGN KEY ("article_type_uuid") REFERENCES "article_type"("uuid");

-- Restore the foreign key without cascading deletes
ALTER TABLE "article_transaction" 
ADD CONSTRAINT "article_transaction_article_uuid_fkey"
FOREIGN KEY ("article_uuid") REFERENCES "article"("uuid");

COMMIT;
