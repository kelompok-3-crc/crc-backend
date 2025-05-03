DROP INDEX IF EXISTS idx_users_kantor_cabang;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_kantor_cabang;
ALTER TABLE users DROP COLUMN IF EXISTS kantor_cabang_id;