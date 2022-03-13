-- +goose Up
-- +goose StatementBegin
CREATE TABLE banner_slot_social_groups (
    id BIGSERIAL PRIMARY KEY,
    banner_slot_id BIGSERIAL,
    CONSTRAINT banner_slot_social_groups_banner_slot_id_foreign_key
        FOREIGN KEY(banner_slot_id) 
        REFERENCES banner_slots(id)
        ON DELETE CASCADE,
    social_group_id BIGSERIAL,
    CONSTRAINT banner_slot_social_groups_social_group_id_foreign_key
        FOREIGN KEY(social_group_id) 
        REFERENCES social_groups(id)
        ON DELETE CASCADE,
    views BIGINT DEFAULT 0 NOT NULL,
    clicks BIGINT DEFAULT 0 NOT NULL,
    UNIQUE (banner_slot_id, social_group_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banner_slot_social_groups;
-- +goose StatementEnd
