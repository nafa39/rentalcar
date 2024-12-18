CREATE SCHEMA "rental-car";

CREATE TABLE "rental-car".users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    balance DECIMAL(10, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "rental-car".cars (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    price_per_day DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'available',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "rental-car".reservations (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    car_id INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES "rental-car".users (id),
    FOREIGN KEY (car_id) REFERENCES "rental-car".cars (id)
);

CREATE TABLE "rental-car".payments (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    payment_method VARCHAR(100) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES "rental-car".users (id)
);

-- function for trigger update at
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- create trigger
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON "rental-car".users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON "rental-car".cars
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON "rental-car".reservations
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON "rental-car".payments
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

INSERT INTO "rental-car".users (name, email, password, balance, created_at, updated_at)
VALUES 
('John Doe', 'john.doe@example.com', 'hashedpassword123', 100.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Jane Smith', 'jane.smith@example.com', 'hashedpassword456', 200.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Alice Johnson', 'alice.johnson@example.com', 'hashedpassword789', 300.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO "rental-car".cars (name, category, price_per_day, status, created_at, updated_at)
VALUES 
('Toyota Corolla', 'Sedan', 50.00, 'available', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Honda CR-V', 'SUV', 75.00, 'available', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Ford Mustang', 'Sports', 100.00, 'rented', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Tesla Model 3', 'Electric', 90.00, 'available', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO "rental-car".reservations (user_id, car_id, start_date, end_date, total_price, created_at, updated_at)
VALUES 
(1, 3, '2024-12-20', '2024-12-25', 500.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 2, '2024-12-18', '2024-12-22', 300.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 4, '2024-12-19', '2024-12-21', 180.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO "rental-car".payments (user_id, amount, payment_method, status, created_at, updated_at)
VALUES 
(1, 100.00, 'credit_card', 'completed', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 200.00, 'bank_transfer', 'completed', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 300.00, 'e-wallet', 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);


SELECT * FROM "rental-car".users;
SELECT * FROM "rental-car".cars;
SELECT * FROM "rental-car".reservations;
SELECT * FROM "rental-car".payments;


