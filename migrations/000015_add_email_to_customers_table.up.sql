ALTER TABLE customers
ADD COLUMN IF NOT EXISTS "email" VARCHAR(50) NOT NULL DEFAULT '';