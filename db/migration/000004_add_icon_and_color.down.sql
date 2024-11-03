-- Down Migration: Remove icon_codepoint and color from article_type

ALTER TABLE "article_type"
DROP COLUMN "icon_codepoint",
DROP COLUMN "color";