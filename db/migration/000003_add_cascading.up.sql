-- Up Migration: Allow deletion of articles without affecting article_type

BEGIN;

-- 1. Drop existing foreign key constraints
ALTER TABLE "article" 
DROP CONSTRAINT IF EXISTS "article_article_type_uuid_fkey";

ALTER TABLE "article_transaction" 
DROP CONSTRAINT IF EXISTS "article_transaction_article_uuid_fkey";

-- 2. Re-add foreign key constraints
ALTER TABLE "article" 
ADD CONSTRAINT "article_article_type_uuid_fkey"
FOREIGN KEY ("article_type_uuid") REFERENCES "article_type"("uuid");

-- Adding ON DELETE CASCADE to article_transaction foreign key
ALTER TABLE "article_transaction" 
ADD CONSTRAINT "article_transaction_article_uuid_fkey"
FOREIGN KEY ("article_uuid") REFERENCES "article"("uuid") ON DELETE CASCADE;

COMMIT;
