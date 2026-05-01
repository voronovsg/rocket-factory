-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS orders (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    part_uuids UUID[] NOT NULL,
    total_price NUMERIC(10,2) NOT NULL,
    transaction_uuid UUID,
    payment_method TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_orders_user_uuid ON orders (user_uuid);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_orders_user_uuid;
DROP INDEX IF EXISTS idx_orders_status;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd

