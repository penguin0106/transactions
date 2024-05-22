
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL
);


CREATE TABLE wallets (
                         user_id INT PRIMARY KEY,
                         usd NUMERIC(20, 2) NOT NULL DEFAULT 0.00,
                         cryptocurrencies JSONB NOT NULL DEFAULT '{}'::jsonb,
                         FOREIGN KEY (user_id) REFERENCES users(id)
);


CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        seller_id INT NOT NULL,
                        buyer_id INT,
                        cryptocurrency VARCHAR(50) NOT NULL,
                        amount NUMERIC(20, 8) NOT NULL,
                        desired_currency VARCHAR(50) NOT NULL,
                        status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
                        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                        FOREIGN KEY (seller_id) REFERENCES users(id),
                        FOREIGN KEY (buyer_id) REFERENCES users(id)
);


-- Insert users
INSERT INTO users (username) VALUES ('seller1'), ('buyer1');

-- Insert wallets
INSERT INTO wallets (user_id, usd, cryptocurrencies) VALUES
                                                         (1, 1000.00, '{"BTC": 2.0, "ETH": 10.0}'),
                                                         (2, 500.00, '{"BTC": 1.0, "ETH": 5.0}');

-- Insert orders
INSERT INTO orders (seller_id, cryptocurrency, amount, desired_currency, status) VALUES
                                                                                     (1, 'BTC', 0.5, 'USD', 'PENDING'),
                                                                                     (2, 'ETH', 2.0, 'BTC', 'PENDING');