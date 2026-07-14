CREATE TABLE IF NOT EXISTS actions (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    step_id VARCHAR(120) NOT NULL REFERENCES scenario_steps(id) ON DELETE CASCADE,
    command TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    points INT NOT NULL DEFAULT 0,
    feedback TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_actions_session_id ON actions(session_id);
CREATE INDEX IF NOT EXISTS idx_actions_step_id ON actions(step_id);
CREATE INDEX IF NOT EXISTS idx_actions_created_at ON actions(created_at);