BEGIN;

DELETE FROM step_command_outputs
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

-- =========================================================
-- permissions-junior
-- =========================================================

INSERT INTO step_command_outputs (
    step_id,
    command_pattern,
    match_type,
    stdout,
    stderr,
    exit_code,
    description,
    priority
) VALUES
(
    'permissions-observe-http',
    'curl',
    'contains',
    'HTTP/1.1 403 Forbidden
Server: nginx
Content-Type: text/html

',
    '',
    0,
    'HTTP показывает 403 Forbidden',
    10
),
(
    'permissions-check-nginx-logs',
    '/var/log/nginx/error.log',
    'contains',
    '2026/07/10 22:05:12 [error] 1234#1234: *567 open() "/var/www/html/index.html" failed (13: Permission denied), client: 172.20.0.1
2026/07/10 22:05:15 [error] 1234#1234: *568 directory index of "/var/www/html/" is forbidden, client: 172.20.0.1
2026/07/10 22:05:20 [crit] 1234#1234: *569 SSL: error:02001002: system library:fopen:Permission denied
',
    '',
    0,
    'Nginx error.log с Permission denied',
    10
),
(
    'permissions-check-owner',
    '/var/www/html',
    'contains',
    'total 16
drwxr-xr-x 2 root root 4096 Jul 10 22:03 .
drwxr-xr-x 3 root root 4096 Jul 10 22:03 ..
-rw-r--r-- 1 root root  612 Jul 10 22:03 index.html
',
    '',
    0,
    'Файлы принадлежат root:root',
    10
),
(
    'permissions-fix-owner',
    'chown',
    'contains',
    '',
    '',
    0,
    'chown успешно выполнен',
    10
),
(
    'permissions-check-nginx-config',
    'nginx -t',
    'exact',
    'nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
',
    '',
    0,
    'nginx config валиден',
    10
),
(
    'permissions-reload-nginx',
    'reload nginx',
    'contains',
    '',
    '',
    0,
    'nginx reload выполнен',
    10
),
(
    'permissions-final-check',
    'curl',
    'contains',
    'HTTP/1.1 200 OK
Server: nginx
Content-Type: text/html

',
    '',
    0,
    'Сайт восстановлен',
    10
);

-- =========================================================
-- memory-hunter
-- =========================================================

INSERT INTO step_command_outputs (
    step_id,
    command_pattern,
    match_type,
    stdout,
    stderr,
    exit_code,
    description,
    priority
) VALUES
(
    'memory-check-ram',
    'free',
    'contains',
    '              total        used        free      shared  buff/cache   available
Mem:           1.9Gi       1.8Gi        74Mi        12Mi        60Mi        80Mi
Swap:             0B          0B          0B
',
    '',
    0,
    'Высокое использование RAM',
    10
),
(
    'memory-check-processes',
    'htop',
    'contains',
    'Interactive process viewer opened.
High memory usage detected earlier.
postgres process is not running.
',
    '',
    0,
    'Процессы показывают проблему с памятью',
    10
),
(
    'memory-check-processes',
    'top',
    'contains',
    'top - 22:13:01 up 10 days,  load average: 4.21, 3.80, 2.91
Tasks: 91 total, 1 running, 90 sleeping
MiB Mem : 1996.0 total, 74.0 free, 1840.0 used, 82.0 buff/cache
',
    '',
    0,
    'top показывает высокую память',
    20
),
(
    'memory-check-processes',
    'ps aux',
    'contains',
    'USER       PID %CPU %MEM COMMAND
root         1  0.0  0.1 /sbin/init
backend   3789 92.1 79.4 [killed]
postgres  4567 13.2 65.1 [killed]
',
    '',
    0,
    'ps показывает убитые процессы',
    30
),
(
    'memory-check-oom',
    'oom',
    'contains',
    '[2026-07-10 22:12:45] Out of memory: Killed process 4567 (postgres) total-vm:2845678kB, anon-rss:1245678kB
[2026-07-10 22:12:46] Out of memory: Killed process 3789 (gunicorn)
',
    '',
    0,
    'OOM Killer найден',
    10
),
(
    'memory-check-postgres-status',
    'postgresql',
    'contains',
    '● postgresql.service - PostgreSQL database server
   Loaded: loaded
   Active: failed
',
    '',
    3,
    'PostgreSQL failed',
    10
),
(
    'memory-restart-postgres',
    'restart postgresql',
    'contains',
    '',
    '',
    0,
    'PostgreSQL restarted',
    10
),
(
    'memory-restart-backend',
    'restart backend',
    'contains',
    'Container backend restarted
',
    '',
    0,
    'Backend restarted',
    10
),
(
    'memory-final-check',
    'docker compose ps',
    'exact',
    'NAME        STATUS
backend     running
postgres    running
',
    '',
    0,
    'Сервисы работают',
    10
);

-- =========================================================
-- dns-mid
-- =========================================================

INSERT INTO step_command_outputs (
    step_id,
    command_pattern,
    match_type,
    stdout,
    stderr,
    exit_code,
    description,
    priority
) VALUES
(
    'dns-check-domain',
    'google.com',
    'contains',
    '',
    'ping: google.com: Temporary failure in name resolution
',
    2,
    'Домен не резолвится',
    10
),
(
    'dns-check-domain',
    'api.payment.com',
    'contains',
    '',
    'curl: (6) Could not resolve host: api.payment.com
',
    6,
    'Payment API не резолвится',
    20
),
(
    'dns-check-ip',
    '8.8.8.8',
    'contains',
    'PING 8.8.8.8 (8.8.8.8): 56 data bytes
64 bytes from 8.8.8.8: icmp_seq=0 ttl=117 time=12.4 ms
',
    '',
    0,
    'IP доступен',
    10
),
(
    'dns-check-ip',
    '1.1.1.1',
    'contains',
    'PING 1.1.1.1 (1.1.1.1): 56 data bytes
64 bytes from 1.1.1.1: icmp_seq=0 ttl=57 time=9.1 ms
',
    '',
    0,
    'IP доступен',
    20
),
(
    'dns-check-nslookup',
    'google.com',
    'contains',
    '',
    ';; connection timed out; no servers could be reached
',
    1,
    'DNS lookup timeout',
    10
),
(
    'dns-check-resolv',
    '/etc/resolv.conf',
    'contains',
    'nameserver 8.8.3.8
',
    '',
    0,
    'Неверный DNS в resolv.conf',
    10
),
(
    'dns-fix-resolv',
    'nameserver 8.8.8.8',
    'contains',
    '',
    '',
    0,
    'DNS исправлен',
    10
),
(
    'dns-restart-resolved',
    'systemd-resolved',
    'contains',
    '',
    '',
    0,
    'resolver restarted',
    10
),
(
    'dns-restart-resolved',
    'resolvectl flush-caches',
    'exact',
    '',
    '',
    0,
    'DNS cache flushed',
    20
),
(
    'dns-final-check',
    'google.com',
    'contains',
    'PING google.com (142.250.185.238): 56 data bytes
64 bytes from 142.250.185.238: icmp_seq=0 ttl=116 time=18.2 ms
',
    '',
    0,
    'DNS восстановлен',
    10
);

-- =========================================================
-- rogue-trainee
-- =========================================================

INSERT INTO step_command_outputs (
    step_id,
    command_pattern,
    match_type,
    stdout,
    stderr,
    exit_code,
    description,
    priority
) VALUES
(
    'rogue-check-nginx-status',
    'nginx',
    'contains',
    '● nginx.service - A high performance web server
   Loaded: loaded (/lib/systemd/system/nginx.service; enabled)
   Active: failed (Result: exit-code)
',
    '',
    3,
    'Nginx failed',
    10
),
(
    'rogue-run-nginx-test',
    'nginx -t',
    'exact',
    '',
    'nginx: [emerg] unexpected "}" in /etc/nginx/sites-available/default:45
nginx: configuration file /etc/nginx/nginx.conf test failed
',
    1,
    'nginx config broken',
    10
),
(
    'rogue-show-broken-lines',
    '/etc/nginx/sites-available/default',
    'contains',
    '40: server {
41:     listen 80;
42:     server_name example.com
43:
44:     location / {
45:         root /var/www/html;
46:     }
47: }
',
    '',
    0,
    'Проблемные строки конфига',
    10
),
(
    'rogue-edit-config',
    '/etc/nginx/sites-available/default',
    'contains',
    'File opened and saved.
',
    '',
    0,
    'Config edited',
    10
),
(
    'rogue-retest-nginx',
    'nginx -t',
    'exact',
    'nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
',
    '',
    0,
    'nginx config fixed',
    10
),
(
    'rogue-restart-nginx',
    'restart nginx',
    'contains',
    '',
    '',
    0,
    'nginx restarted',
    10
),
(
    'rogue-final-check',
    'curl',
    'contains',
    'HTTP/1.1 200 OK
Server: nginx
Content-Type: text/html

',
    '',
    0,
    'Сайт доступен',
    10
);

-- =========================================================
-- crashloop-senior
-- =========================================================

INSERT INTO step_command_outputs (
    step_id,
    command_pattern,
    match_type,
    stdout,
    stderr,
    exit_code,
    description,
    priority
) VALUES
(
    'crashloop-check-containers',
    'docker ps',
    'exact',
    'CONTAINER ID   IMAGE          COMMAND                  STATUS
a1b2c3d4e5f6   auth-service   "npm start"              Restarting (1) 42 seconds ago
b2c3d4e5f6a1   postgres       "docker-entrypoint..."   Up 10 minutes
',
    '',
    0,
    'auth-service restarting',
    10
),
(
    'crashloop-check-containers',
    'docker compose ps',
    'exact',
    'NAME           STATUS
auth-service   restarting
postgres       running
',
    '',
    0,
    'compose показывает restarting',
    20
),
(
    'crashloop-check-logs',
    'auth-service',
    'contains',
    'ERROR: connection to database "auth" failed: could not connect to server: Connection refused
Is the server running on host "db" and accepting TCP/IP connections on port 5432?
',
    '',
    1,
    'Ошибка подключения к базе',
    10
),
(
    'crashloop-check-env',
    '.env',
    'contains',
    'DB_HOST=db
DB_PORT=5433
DB_NAME=auth
DB_USER=auth
',
    '',
    0,
    'Неверный DB_PORT',
    10
),
(
    'crashloop-search-wrong-port',
    '5433',
    'contains',
    './.env:DB_PORT=5433
./docker-compose.yml:      DB_PORT: 5433
',
    '',
    0,
    '5433 найден',
    10
),
(
    'crashloop-fix-port',
    '5432',
    'contains',
    '',
    '',
    0,
    'порт исправлен',
    10
),
(
    'crashloop-recreate-service',
    'auth-service',
    'contains',
    'Container auth-service recreated
',
    '',
    0,
    'auth-service recreated',
    10
),
(
    'crashloop-final-check',
    'docker ps',
    'exact',
    'CONTAINER ID   IMAGE          COMMAND                  STATUS
a1b2c3d4e5f6   auth-service   "npm start"              Up 12 seconds
b2c3d4e5f6a1   postgres       "docker-entrypoint..."   Up 13 minutes
',
    '',
    0,
    'auth-service up',
    10
),
(
    'crashloop-final-check',
    'docker compose ps',
    'exact',
    'NAME           STATUS
auth-service   running
postgres       running
',
    '',
    0,
    'compose services running',
    20
);

COMMIT;