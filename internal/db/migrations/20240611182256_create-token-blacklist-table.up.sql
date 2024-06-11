CREATE TABLE token_blacklist (
     id SERIAL PRIMARY KEY,
     token TEXT NOT NULL,
     expiry TIMESTAMP NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_token_blacklist_expiry ON token_blacklist (expiry);