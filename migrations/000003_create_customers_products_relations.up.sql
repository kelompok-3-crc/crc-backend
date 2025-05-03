CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    cif VARCHAR(50) UNIQUE NOT NULL,
    nomor_rekening VARCHAR(50) UNIQUE NOT NULL,
    nama_perusahaan VARCHAR(50) NOT NULL,
    produk_eksisting VARCHAR[] DEFAULT '{}',
    aktivitas_transaksi VARCHAR(100),
    nomor_hp VARCHAR(20) UNIQUE NOT NULL,
    segmen VARCHAR(20),
    address TEXT,
    job VARCHAR(100),
    penghasilan BIGINT DEFAULT 0,
    umur INT NOT NULL,
    gender VARCHAR(10) NOT NULL,
    status_perkawinan BOOLEAN DEFAULT false,
    payroll BOOLEAN DEFAULT false,
    top_produk VARCHAR[] DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(60) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE customer_products (
    customer_id BIGINT NOT NULL,
    product_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (customer_id, product_id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);