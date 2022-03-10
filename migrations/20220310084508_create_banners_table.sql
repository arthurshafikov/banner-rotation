-- +goose Up
-- +goose StatementBegin
CREATE TABLE banners (
    id BIGSERIAL PRIMARY KEY,
    description VARCHAR(300) NULL,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banners;
-- +goose StatementEnd
