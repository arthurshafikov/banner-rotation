-- +goose Up
-- +goose StatementBegin
CREATE TABLE banner_slots (
    id BIGSERIAL PRIMARY KEY,
    banner_id BIGSERIAL,
    slot_id BIGSERIAL,
    CONSTRAINT banner_slots_banner_id_foreign_key
        FOREIGN KEY(banner_id) 
        REFERENCES banners(id),
    CONSTRAINT banner_slots_slot_id_foreign_key
        FOREIGN KEY(slot_id) 
        REFERENCES slots(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banner_slots;
-- +goose StatementEnd
