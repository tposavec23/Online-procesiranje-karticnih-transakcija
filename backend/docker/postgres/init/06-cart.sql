CREATE TABLE catalogue.cart (
    cart_id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id integer NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE catalogue.cart_item (
    cart_item_id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cart_id integer NOT NULL REFERENCES catalogue.cart (cart_id) ON DELETE CASCADE,
    article_id integer NOT NULL REFERENCES catalogue.article (article_id) ON DELETE CASCADE,
    quantity integer NOT NULL CHECK (quantity > 0),
    added_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (cart_id, article_id)
);