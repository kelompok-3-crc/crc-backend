CREATE TABLE
    target_produk_bulanan (
        id SERIAL PRIMARY KEY,
        tahun INT NOT NULL,
        bulan INT NOT NULL,
        target_amount BIGINT NOT NULL,
        kantor_cabang_id INT NOT NULL,
        product_id INT NOT NULL,
        target_tahunan_id INT NOT NULL,
        FOREIGN KEY (kantor_cabang_id) REFERENCES kantor_cabang (id) ON UPDATE CASCADE ON DELETE CASCADE,
        FOREIGN KEY (product_id) REFERENCES products (id) ON UPDATE CASCADE ON DELETE CASCADE,
        FOREIGN KEY (target_tahunan_id) REFERENCES target_produk_tahunan (id) ON UPDATE CASCADE ON DELETE CASCADE
    );