-- +goose Up
CREATE TABLE IF NOT EXISTS raw_footages (
    id             BIGSERIAL PRIMARY KEY,
    created_at     TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ,
    title          TEXT NOT NULL,
    assigned_to_id BIGINT NOT NULL REFERENCES employees (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_raw_footages_deleted_at ON raw_footages (deleted_at);
CREATE INDEX IF NOT EXISTS idx_raw_footages_assigned_to_id ON raw_footages (assigned_to_id);

CREATE TABLE IF NOT EXISTS raw_footage_files (
    raw_footage_id BIGINT NOT NULL REFERENCES raw_footages (id) ON UPDATE CASCADE ON DELETE CASCADE,
    v_node_id      BIGINT NOT NULL REFERENCES filesystem_nodes (id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (raw_footage_id, v_node_id)
);

CREATE TABLE IF NOT EXISTS edited_videos (
    id               BIGSERIAL PRIMARY KEY,
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,
    raw_footage_id   BIGINT NOT NULL REFERENCES raw_footages (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    edited_v_node_id BIGINT NOT NULL REFERENCES filesystem_nodes (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_edited_videos_deleted_at ON edited_videos (deleted_at);

CREATE TABLE IF NOT EXISTS published_videos (
    id                BIGSERIAL PRIMARY KEY,
    created_at        TIMESTAMPTZ,
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,
    edited_video_id   BIGINT NOT NULL REFERENCES edited_videos (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    you_tube_video_id VARCHAR(32) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_published_videos_deleted_at ON published_videos (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS published_videos;
DROP TABLE IF EXISTS edited_videos;
DROP TABLE IF EXISTS raw_footage_files;
DROP TABLE IF EXISTS raw_footages;
