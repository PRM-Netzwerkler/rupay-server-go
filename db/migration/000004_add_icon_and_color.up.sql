

ALTER TABLE "article_type"
ADD COLUMN "icon_codepoint" INT NOT NULL DEFAULT 0,
ADD COLUMN "color" VARCHAR(7) NOT NULL DEFAULT '#FFFFFF';