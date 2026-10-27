ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_tasks_active_status_created_at
    ON tasks (status, created_at DESC, id DESC)
    WHERE deleted_at IS NULL;
