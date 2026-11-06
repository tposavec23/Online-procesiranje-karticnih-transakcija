CREATE TABLE account.users (
    user_id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username varchar(50) NOT NULL UNIQUE,
    password_hash varchar(100) NOT NULL,
    role varchar(20) NOT NULL CHECK (role IN ('MERCHANT', 'ADMIN')),
    active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    pdated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);