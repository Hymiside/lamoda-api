CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    available BOOLEAN NOT NULL DEFAULT TRUE,
    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    part_number TEXT NOT NULL UNIQUE,
    dimensions JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE INDEX part_number_idx ON products (part_number);

CREATE TABLE warehouse_products (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL CONSTRAINT positive_quantity CHECK (quantity >= 0),

    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE reserved_products (
    id SERIAL PRIMARY KEY,
    warehouse_product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL CONSTRAINT positive_quantity CHECK (quantity > 0),

    FOREIGN KEY (warehouse_product_id) REFERENCES warehouse_products(id) ON DELETE CASCADE
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    warehouse_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL CONSTRAINT positive_quantity CHECK (quantity > 0),    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);