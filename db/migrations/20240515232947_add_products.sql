-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id              BIGSERIAL NOT NULL PRIMARY KEY
    ,name           TEXT NOT NULL
    ,sku            TEXT NOT NULL
    ,qty            INTEGER NOT NULL
    ,created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
    ,updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    ,UNIQUE(name)
);

CREATE INDEX products_name_idx ON products (name);
CREATE INDEX products_sku_idx ON products (sku);

-- +goose Down
DROP TABLE IF EXISTS products;
