CREATE TABLE catalogue.category (
    category_id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name varchar(100) NOT NULL,
    description text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE catalogue.article (
    article_id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    category_id integer NOT NULL REFERENCES catalogue.category (category_id),
    code varchar(50) NOT NULL UNIQUE,
    name varchar(150) NOT NULL,
    description text,
    price numeric(10,2) NOT NULL CHECK (price >= 0),
    available_quantity integer NOT NULL DEFAULT 0 CHECK (available_quantity >= 0),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);