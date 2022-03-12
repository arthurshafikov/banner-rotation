-- +goose Up
-- +goose StatementBegin
CREATE TABLE slots (
    id BIGSERIAL PRIMARY KEY,
    description VARCHAR(300) NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE slots;
-- +goose StatementEnd
