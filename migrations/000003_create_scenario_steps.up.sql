CREATE TABLE IF NOT EXISTS scenario_steps (
    id VARCHAR(120) PRIMARY KEY,
    scenario_id VARCHAR(100) NOT NULL REFERENCES scenarios(id) ON DELETE CASCADE,
    step_order INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    hint TEXT NOT NULL DEFAULT '',
    expected_command TEXT NOT NULL DEFAULT '',
    expected_result TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_scenario_step_order UNIQUE (scenario_id, step_order)
);