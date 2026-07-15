BEGIN;

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

-- =========================================================
-- permissions-junior
-- =========================================================

INSERT INTO step_accepted_commands (
    step_id,
    command,
    match_type,
    description
) VALUES
('permissions-observe-http', 'curl -I http://localhost', 'exact', 'Проверка HTTP-заголовков'),
('permissions-observe-http', 'curl -i http://localhost', 'exact', 'Проверка HTTP-ответа'),
('permissions-observe-http', 'curl http://localhost', 'exact', 'Обычный HTTP-запрос'),

('permissions-check-nginx-logs', 'tail -f /var/log/nginx/error.log', 'exact', 'Просмотр error.log'),
('permissions-check-nginx-logs', 'tail -n 50 /var/log/nginx/error.log', 'exact', 'Просмотр последних строк error.log'),
('permissions-check-nginx-logs', 'less /var/log/nginx/error.log', 'exact', 'Открыть error.log'),
('permissions-check-nginx-logs', '/var/log/nginx/error.log', 'contains', 'Любая команда чтения nginx error.log'),

('permissions-check-owner', 'ls -la /var/www/html', 'exact', 'Проверить владельца web-root'),
('permissions-check-owner', 'ls -l /var/www/html', 'exact', 'Проверить права web-root'),
('permissions-check-owner', '/var/www/html', 'contains', 'Любая команда просмотра /var/www/html'),

('permissions-fix-owner', 'chown -R www-data:www-data /var/www/html', 'exact', 'Исправить владельца web-root'),
('permissions-fix-owner', 'sudo chown -R www-data:www-data /var/www/html', 'exact', 'Исправить владельца через sudo'),

('permissions-check-nginx-config', 'nginx -t', 'exact', 'Проверить nginx config'),
('permissions-check-nginx-config', 'sudo nginx -t', 'exact', 'Проверить nginx config через sudo'),

('permissions-reload-nginx', 'systemctl reload nginx', 'exact', 'Reload nginx'),
('permissions-reload-nginx', 'sudo systemctl reload nginx', 'exact', 'Reload nginx через sudo'),
('permissions-reload-nginx', 'nginx -s reload', 'exact', 'Reload nginx через signal'),

('permissions-final-check', 'curl -I http://localhost', 'exact', 'Финальная HTTP-проверка'),
('permissions-final-check', 'curl -i http://localhost', 'exact', 'Финальная HTTP-проверка'),
('permissions-final-check', 'curl http://localhost', 'exact', 'Финальный HTTP-запрос');

-- =========================================================
-- memory-hunter
-- =========================================================

INSERT INTO step_accepted_commands (
    step_id,
    command,
    match_type,
    description
) VALUES
('memory-check-ram', 'free -h', 'exact', 'Проверить память'),
('memory-check-ram', 'free -m', 'exact', 'Проверить память в MB'),
('memory-check-ram', 'cat /proc/meminfo', 'exact', 'Проверить память через procfs'),

('memory-check-processes', 'htop', 'exact', 'Проверить процессы'),
('memory-check-processes', 'top', 'exact', 'Проверить процессы через top'),
('memory-check-processes', 'ps aux --sort=-%mem', 'exact', 'Процессы по памяти'),
('memory-check-processes', 'ps aux', 'exact', 'Список процессов'),

('memory-check-oom', 'dmesg -T | grep -i oom', 'exact', 'Найти OOM'),
('memory-check-oom', 'dmesg | grep -i oom', 'exact', 'Найти OOM без human time'),
('memory-check-oom', 'journalctl -k | grep -i oom', 'exact', 'Найти OOM в journal'),
('memory-check-oom', 'grep -i oom', 'contains', 'Команда ищет OOM'),

('memory-check-postgres-status', 'systemctl status postgresql', 'exact', 'Проверить PostgreSQL'),
('memory-check-postgres-status', 'sudo systemctl status postgresql', 'exact', 'Проверить PostgreSQL через sudo'),
('memory-check-postgres-status', 'journalctl -u postgresql -xe', 'exact', 'Логи PostgreSQL'),

('memory-restart-postgres', 'systemctl restart postgresql', 'exact', 'Перезапустить PostgreSQL'),
('memory-restart-postgres', 'sudo systemctl restart postgresql', 'exact', 'Перезапустить PostgreSQL через sudo'),

('memory-restart-backend', 'docker compose restart backend', 'exact', 'Перезапустить backend'),
('memory-restart-backend', 'docker-compose restart backend', 'exact', 'Перезапустить backend старым docker-compose'),

('memory-final-check', 'docker compose ps', 'exact', 'Проверить сервисы'),
('memory-final-check', 'docker-compose ps', 'exact', 'Проверить сервисы старым docker-compose'),
('memory-final-check', 'systemctl status postgresql', 'exact', 'Проверить PostgreSQL');

-- =========================================================
-- dns-mid
-- =========================================================

INSERT INTO step_accepted_commands (
    step_id,
    command,
    match_type,
    description
) VALUES
('dns-check-domain', 'ping google.com', 'exact', 'Проверить домен'),
('dns-check-domain', 'curl https://google.com', 'exact', 'Проверить домен через curl'),
('dns-check-domain', 'curl https://api.payment.com', 'exact', 'Проверить payment API'),

('dns-check-ip', 'ping 8.8.8.8', 'exact', 'Проверить IP'),
('dns-check-ip', 'ping 1.1.1.1', 'exact', 'Проверить IP Cloudflare'),

('dns-check-nslookup', 'nslookup google.com', 'exact', 'Проверить DNS через nslookup'),
('dns-check-nslookup', 'dig google.com', 'exact', 'Проверить DNS через dig'),
('dns-check-nslookup', 'host google.com', 'exact', 'Проверить DNS через host'),

('dns-check-resolv', 'cat /etc/resolv.conf', 'exact', 'Посмотреть resolv.conf'),
('dns-check-resolv', 'cat /etc/resolv.conf | grep nameserver', 'exact', 'Посмотреть nameserver'),
('dns-check-resolv', 'grep nameserver /etc/resolv.conf', 'exact', 'Найти nameserver'),
('dns-check-resolv', '/etc/resolv.conf', 'contains', 'Любая команда чтения resolv.conf'),

('dns-fix-resolv', 'echo -e "nameserver 8.8.8.8\nnameserver 1.1.1.1" > /etc/resolv.conf', 'exact', 'Записать корректные DNS'),
('dns-fix-resolv', 'echo "nameserver 8.8.8.8" > /etc/resolv.conf', 'exact', 'Записать Google DNS'),
('dns-fix-resolv', 'nameserver 8.8.8.8', 'contains', 'Команда содержит корректный nameserver'),

('dns-restart-resolved', 'systemctl restart systemd-resolved', 'exact', 'Перезапустить resolver'),
('dns-restart-resolved', 'sudo systemctl restart systemd-resolved', 'exact', 'Перезапустить resolver через sudo'),
('dns-restart-resolved', 'resolvectl flush-caches', 'exact', 'Очистить DNS cache'),

('dns-final-check', 'ping google.com', 'exact', 'Проверить домен после исправления'),
('dns-final-check', 'nslookup google.com', 'exact', 'Проверить DNS после исправления'),
('dns-final-check', 'dig google.com', 'exact', 'Проверить DNS через dig');

-- =========================================================
-- rogue-trainee
-- =========================================================

INSERT INTO step_accepted_commands (
    step_id,
    command,
    match_type,
    description
) VALUES
('rogue-check-nginx-status', 'systemctl status nginx', 'exact', 'Проверить nginx status'),
('rogue-check-nginx-status', 'sudo systemctl status nginx', 'exact', 'Проверить nginx status через sudo'),
('rogue-check-nginx-status', 'docker logs nginx --tail 50', 'exact', 'Посмотреть docker logs nginx'),

('rogue-run-nginx-test', 'nginx -t', 'exact', 'Проверить nginx config'),
('rogue-run-nginx-test', 'sudo nginx -t', 'exact', 'Проверить nginx config через sudo'),

('rogue-show-broken-lines', 'sed -n ''40,50p'' /etc/nginx/sites-available/default', 'exact', 'Показать строки вокруг ошибки'),
('rogue-show-broken-lines', 'cat /etc/nginx/sites-available/default | sed -n ''40,50p''', 'exact', 'Показать строки через cat и sed'),
('rogue-show-broken-lines', '/etc/nginx/sites-available/default', 'contains', 'Команда читает nginx config'),

('rogue-edit-config', 'nano /etc/nginx/sites-available/default', 'exact', 'Открыть config в nano'),
('rogue-edit-config', 'vim /etc/nginx/sites-available/default', 'exact', 'Открыть config в vim'),
('rogue-edit-config', 'vim +45 /etc/nginx/sites-available/default', 'exact', 'Открыть config на строке 45'),

('rogue-retest-nginx', 'nginx -t', 'exact', 'Повторно проверить nginx config'),
('rogue-retest-nginx', 'sudo nginx -t', 'exact', 'Повторно проверить nginx config через sudo'),

('rogue-restart-nginx', 'systemctl restart nginx', 'exact', 'Перезапустить nginx'),
('rogue-restart-nginx', 'sudo systemctl restart nginx', 'exact', 'Перезапустить nginx через sudo'),
('rogue-restart-nginx', 'nginx -s reload', 'exact', 'Reload nginx'),

('rogue-final-check', 'curl -I http://localhost', 'exact', 'Проверить сайт'),
('rogue-final-check', 'curl http://localhost', 'exact', 'Проверить сайт');

-- =========================================================
-- crashloop-senior
-- =========================================================

INSERT INTO step_accepted_commands (
    step_id,
    command,
    match_type,
    description
) VALUES
('crashloop-check-containers', 'docker ps', 'exact', 'Проверить контейнеры'),
('crashloop-check-containers', 'docker compose ps', 'exact', 'Проверить compose services'),
('crashloop-check-containers', 'docker-compose ps', 'exact', 'Проверить compose services старым docker-compose'),

('crashloop-check-logs', 'docker logs auth-service --tail 100', 'exact', 'Посмотреть последние логи auth-service'),
('crashloop-check-logs', 'docker logs --since 5m auth-service', 'exact', 'Посмотреть свежие логи auth-service'),
('crashloop-check-logs', 'docker logs auth-service', 'exact', 'Посмотреть все логи auth-service'),
('crashloop-check-logs', 'docker-compose logs --tail=50 auth', 'exact', 'Посмотреть логи auth'),

('crashloop-check-env', 'cat .env | grep DB_', 'exact', 'Посмотреть DB переменные'),
('crashloop-check-env', 'grep DB_ .env', 'exact', 'Найти DB переменные'),
('crashloop-check-env', 'cat .env', 'exact', 'Посмотреть .env'),

('crashloop-search-wrong-port', 'grep -r "5433" .', 'exact', 'Найти 5433'),
('crashloop-search-wrong-port', 'grep -R "5433" .', 'exact', 'Найти 5433 рекурсивно'),
('crashloop-search-wrong-port', '5433', 'contains', 'Команда ищет 5433'),

('crashloop-fix-port', 'sed -i ''s/5433/5432/g'' .env', 'exact', 'Исправить порт в .env'),
('crashloop-fix-port', 'DB_PORT=5432', 'contains', 'Команда устанавливает корректный DB_PORT'),

('crashloop-recreate-service', 'docker compose up -d --force-recreate auth-service', 'exact', 'Пересоздать auth-service'),
('crashloop-recreate-service', 'docker-compose up -d --force-recreate auth-service', 'exact', 'Пересоздать auth-service старым docker-compose'),
('crashloop-recreate-service', 'docker compose restart auth-service', 'exact', 'Перезапустить auth-service'),

('crashloop-final-check', 'docker ps', 'exact', 'Финально проверить контейнеры'),
('crashloop-final-check', 'docker compose ps', 'exact', 'Финально проверить compose services'),
('crashloop-final-check', 'docker-compose ps', 'exact', 'Финально проверить compose services старым docker-compose');

COMMIT;