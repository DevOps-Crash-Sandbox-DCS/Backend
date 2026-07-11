DELETE FROM step_accepted_commands
WHERE step_id IN (
    SELECT id
    FROM scenario_steps
    WHERE scenario_id IN (
        'permissions-junior',
        'memory-hunter',
        'dns-mid',
        'rogue-trainee',
        'crashloop-senior'
    )
);