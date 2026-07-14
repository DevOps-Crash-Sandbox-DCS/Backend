CREATE TABLE IF NOT EXISTS step_accepted_commands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    step_id VARCHAR(120) NOT NULL REFERENCES scenario_steps(id) ON DELETE CASCADE,
    command TEXT NOT NULL,
    match_type TEXT NOT NULL DEFAULT 'exact',
    description TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_step_accepted_commands_step_id
ON step_accepted_commands(step_id);

CREATE INDEX IF NOT EXISTS idx_step_accepted_commands_active
ON step_accepted_commands(is_active);