ALTER TABLE products
ADD COLUMN IF NOT EXISTS ikon VARCHAR(50) DEFAULT NULL;

ALTER TABLE products
ADD COLUMN IF NOT EXISTS prediksi VARCHAR(50) DEFAULT NULL;

INSERT INTO
    products (nama, ikon, prediksi)
VALUES
    ('Mitraguna', 'mitraguna', 'mitraguna'),
    ('Griya', 'griya', 'griya'),
    ('Pensiun', 'pensiun', 'pensiun'),
    ('Prapensiun', 'prapensiun', 'prapensiun'),
    ('OTO', 'oto', 'oto'),
    ('Hasanah Card', 'hasanah-card', 'hasanahcard');