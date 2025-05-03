ALTER TABLE users
ADD COLUMN kantor_cabang_id BIGINT NULL DEFAULT NULL;

ALTER TABLE users ADD CONSTRAINT fk_kantor_cabang FOREIGN KEY (kantor_cabang_id) REFERENCES kantor_cabang (id) ON UPDATE CASCADE ON DELETE SET NULL;

CREATE INDEX idx_users_kantor_cabang ON users (kantor_cabang_id);