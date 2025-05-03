ALTER TABLE kantor_cabang
DROP COLUMN created_at,
DROP COLUMN updated_at,
DROP COLUMN deleted_at;

ALTER TABLE target_produk_tahunan
DROP COLUMN created_at,
DROP COLUMN updated_at,
DROP COLUMN deleted_at;

ALTER TABLE target_produk_bulanan
DROP COLUMN created_at,
DROP COLUMN updated_at,
DROP COLUMN deleted_at;
