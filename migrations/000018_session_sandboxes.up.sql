CREATE TABLE IF NOT EXISTS session_sandboxes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    scenario_id VARCHAR(100) NOT NULL REFERENCES scenarios(id) ON DELETE CASCADE,

    container_name VARCHAR(255) NOT NULL UNIQUE,
    image VARCHAR(255) NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'running',

    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    stopped_at TIMESTAMP NULL,
    last_seen_at TIMESTAMP NOT NULL DEFAULT NOW(),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_session_sandbox UNIQUE (session_id)
);

CREATE INDEX IF NOT EXISTS idx_session_sandboxes_session_id
ON session_sandboxes(session_id);

CREATE INDEX IF NOT EXISTS idx_session_sandboxes_status
ON session_sandboxes(status);

CREATE INDEX IF NOT EXISTS idx_session_sandboxes_container_name
ON session_sandboxes(container_name);

CREATE INDEX IF NOT EXISTS idx_session_sandboxes_last_seen_at
ON session_sandboxes(last_seen_at);