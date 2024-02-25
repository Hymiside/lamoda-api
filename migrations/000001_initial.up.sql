CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

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
    width INTEGER,
    height INTEGER,
    depth INTEGER
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
    reservation_id UUID,
    warehouse_product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL CONSTRAINT positive_quantity CHECK (quantity > 0),
    status INT NOT NULL DEFAULT 0, -- 0 - reserved, 1 - cancelled, 2 - shipped
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (warehouse_product_id) REFERENCES warehouse_products(id) ON DELETE CASCADE
);

INSERT INTO warehouses (title, available, lat, lng) VALUES 
    ('Warehouse G', true, 52.5200, 13.4050),
    ('Warehouse H', true, 55.7558, 37.6176),
    ('Warehouse I', false, 35.6895, 139.6917),
    ('Warehouse J', false, 48.8566, 2.3522),
    ('Warehouse K', true, 51.5074, 0.1278);

INSERT INTO products (title, part_number, width, height, depth) VALUES 
    ('Product 5', 'P97531', 10, 23, 15),
    ('Product 6', 'P13579', 10, 23, 15),
    ('Product 7', 'P97431', 10, 23, 15),
    ('Product 8', 'P13279', 10, 23, 15);

INSERT INTO warehouse_products (warehouse_id, product_id, quantity) VALUES 
    (1, 1, 1),
    (1, 2, 23),
    (1, 3, 1),
    (1, 4, 4),
    (2, 1, 2),
    (2, 2, 12),
    (2, 3, 2),
    (2, 4, 2),
    (3, 1, 3),
    (3, 2, 1),
    (3, 3, 3),
    (3, 4, 34),
    (4, 1, 2),
    (4, 2, 234),
    (4, 3, 23),
    (4, 4, 3),
    (5, 1, 234),
    (5, 2, 1),
    (5, 3, 4),
    (5, 4, 7);