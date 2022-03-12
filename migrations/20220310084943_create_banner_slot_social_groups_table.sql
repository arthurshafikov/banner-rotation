-- +goose Up
-- +goose StatementBegin
CREATE TABLE banner_slots_social_groups (
    banner_slot_id BIGSERIAL,
    CONSTRAINT banner_slots_social_groups_banner_slot_id_foreign_key
        FOREIGN KEY(banner_slot_id) 
        REFERENCES banner_slots(id),
    social_group_id BIGSERIAL,
    CONSTRAINT banner_slots_social_groups_social_group_id_foreign_key
        FOREIGN KEY(social_group_id) 
        REFERENCES social_groups(id),
    views BIGINT DEFAULT 0 NOT NULL,
    clicks BIGINT DEFAULT 0 NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banner_slots_social_groups;
-- +goose StatementEnd
