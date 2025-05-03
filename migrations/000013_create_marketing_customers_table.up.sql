CREATE TABLE
    marketing_customers (
        id SERIAL PRIMARY KEY,
        customer_id BIGINT NOT NULL,
        marketing_id INT NOT NULL,
        status VARCHAR(20) NOT NULL DEFAULT 'new',
        product_id INT,
        amount BIGINT,
        notes TEXT,
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP
        WITH
            TIME ZONE,
            CONSTRAINT fk_marketing_customer_customer FOREIGN KEY (customer_id) REFERENCES customers (id)  ON UPDATE CASCADE ON DELETE CASCADE,
            CONSTRAINT fk_marketing_customer_marketing FOREIGN KEY (marketing_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
            CONSTRAINT fk_marketing_customer_product FOREIGN KEY (product_id) REFERENCES products (id) ON UPDATE CASCADE ON DELETE CASCADE,
            CONSTRAINT chk_status CHECK (status IN ('new','contacted', 'closed', 'rejected')),
            CONSTRAINT chk_closed_status CHECK (
                (
                    status = 'closed'
                    AND product_id IS NOT NULL
                    AND amount IS NOT NULL
                )
                OR (status != 'closed')
            )
    );

-- Create partial unique index for active customers
CREATE UNIQUE INDEX idx_unique_active_customer ON marketing_customers (customer_id)
WHERE
    deleted_at IS NULL;

CREATE INDEX idx_marketing_customers_status ON marketing_customers (status);

CREATE INDEX idx_marketing_customers_marketing_id ON marketing_customers (marketing_id);