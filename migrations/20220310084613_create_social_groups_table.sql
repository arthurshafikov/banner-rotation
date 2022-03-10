-- +goose Up
-- +goose StatementBegin
CREATE TABLE social_groups (
    id BIGSERIAL PRIMARY KEY,
    description VARCHAR(300) NULL,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE social_groups;
-- +goose StatementEnd
