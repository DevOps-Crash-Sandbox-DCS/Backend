CREATE TABLE IF NOT EXISTS step_command_outputs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    step_id VARCHAR(120) NOT NULL REFERENCES scenario_steps(id) ON DELETE CASCADE,

    command_pattern TEXT NOT NULL,
    match_type VARCHAR(30) NOT NULL DEFAULT 'exact',

    stdout TEXT NOT NULL DEFAULT '',
    stderr TEXT NOT NULL DEFAULT '',
    exit_code INT NOT NULL DEFAULT 0,

    description TEXT NOT NULL DEFAULT '',
    priority INT NOT NULL DEFAULT 100,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_step_command_outputs_step_id
ON step_command_outputs(step_id);

CREATE INDEX IF NOT EXISTS idx_step_command_outputs_active
ON step_command_outputs(is_active);

CREATE INDEX IF NOT EXISTS idx_step_command_outputs_priority
ON step_command_outputs(priority);