BEGIN;

DELETE FROM actions;
DELETE FROM sessions;

DELETE FROM scenario_steps
WHERE scenario_id IN (
    'permissions-junior',
    'memory-hunter',
    'dns-mid',
    'rogue-trainee',
    'crashloop-senior'
);

COMMIT;