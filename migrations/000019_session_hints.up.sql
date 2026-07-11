CREATE TABLE IF NOT EXISTS session_hints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    scenario_id VARCHAR(100) NOT NULL REFERENCES scenarios(id) ON DELETE CASCADE,
    step_id VARCHAR(100) NULL REFERENCES scenario_steps(id) ON DELETE SET NULL,

    hint_level VARCHAR(30) NOT NULL DEFAULT 'basic',

    request_payload JSONB NOT NULL,
    response_payload JSONB NOT NULL,

    hint TEXT NOT NULL,
    confidence NUMERIC(5, 4) NULL,
    source VARCHAR(50) NOT NULL DEFAULT 'ml',

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_session_hints_session_id
ON session_hints(session_id);

CREATE INDEX IF NOT EXISTS idx_session_hints_user_id
ON session_hints(user_id);

CREATE INDEX IF NOT EXISTS idx_session_hints_step_id
ON session_hints(step_id);

CREATE INDEX IF NOT EXISTS idx_session_hints_created_at
ON session_hints(created_at);