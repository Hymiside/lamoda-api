CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    available BOOLEAN NOT NULL DEFAULT FAlSE,
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

CREATE TABLE shipped_products (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    warehouse_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL CONSTRAINT positive_quantity CHECK (quantity > 0),    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);

INSERT INTO warehouses (title, available, lat, lng) VALUES 
    ('Warehouse G', true, 52.5200, 13.4050),
    ('Warehouse H', true, 55.7558, 37.6176),
    ('Warehouse I', false, 35.6895, 139.6917);

-- Пример заполнения таблицы products
INSERT INTO products (title, part_number, dimensions) VALUES 
    ('Product 5', 'P97531', '{"width": 10, "height": 10, "depth": 10}'),
    ('Product 6', 'P13579', '{"width": 15, "height": 20, "depth": 12}');
