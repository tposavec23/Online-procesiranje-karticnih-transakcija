CREATE TABLE catalogue.cart (
    cart_id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id integer NOT NULL REFERENCES account.users (user_id),
    status varchar(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'COMPLETED')),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- A user can have many carts over time, but only one active at once.
CREATE UNIQUE INDEX one_active_cart_per_user
    ON catalogue.cart (user_id) WHERE status = 'ACTIVE';

CREATE TABLE catalogue.cart_item (
    cart_item_id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cart_id integer NOT NULL REFERENCES catalogue.cart (cart_id) ON DELETE CASCADE,
    article_id integer NOT NULL REFERENCES catalogue.article (article_id),
    quantity integer NOT NULL CHECK (quantity > 0),
    unit_price numeric(10,2) NOT NULL CHECK (unit_price >= 0),
    added_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (cart_id, article_id)
);