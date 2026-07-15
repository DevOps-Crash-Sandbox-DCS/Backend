-- =========================================================
-- Сначала добавляем все шаги в scenario_steps
-- =========================================================

INSERT INTO scenario_steps (
    id,
    scenario_id,
    step_order,
    title,
    description,
    hint,
    expected_command,
    expected_result
)
SELECT * FROM (VALUES
    -- =========================================================
    -- 1. permissions-junior
    -- =========================================================
    (
        'permissions-observe-http',
        'permissions-junior',
        0,
        'Наблюдение за HTTP',
        'Проверьте, что веб-сервер отвечает на запросы',
        'Используйте curl для проверки доступности сайта',
        'curl -I http://localhost',
        'Сайт отвечает с кодом 200 OK'
    ),
    (
        'permissions-check-nginx-logs',
        'permissions-junior',
        1,
        'Проверить логи Nginx',
        'Нужно понять, почему веб-сервер отдает 403 Forbidden или 500 Internal Error.',
        'Открой error.log и найди строки с Permission denied.',
        'tail -f /var/log/nginx/error.log',
        'В логах обнаружена ошибка Permission denied при обращении к файлам сайта или SSL-сертификатам.'
    ),
    (
        'permissions-check-owner',
        'permissions-junior',
        2,
        'Проверить владельца файлов',
        'Проверь, кому принадлежат файлы сайта в /var/www/html.',
        'Используй ls -la для просмотра владельца и группы.',
        'ls -la /var/www/html',
        'Файлы принадлежат root:root вместо www-data:www-data.'
    ),
    (
        'permissions-fix-owner',
        'permissions-junior',
        3,
        'Исправить права доступа',
        'Верни корректного владельца для директории сайта.',
        'Нужно использовать chown.',
        'chown -R www-data:www-data /var/www',
        'Права исправлены, Nginx снова может читать файлы сайта.'
    ),
    (
        'permissions-validate-nginx',
        'permissions-junior',
        4,
        'Проверить конфигурацию Nginx',
        'Убедитесь, что конфигурация Nginx корректна после изменений',
        'Используйте nginx -t для проверки синтаксиса',
        'nginx -t',
        'Конфигурация Nginx валидна'
    ),
    (
        'permissions-reload-nginx',
        'permissions-junior',
        5,
        'Перезагрузить Nginx',
        'Примените изменения конфигурации Nginx',
        'Используйте systemctl или nginx -s reload',
        'systemctl reload nginx',
        'Nginx перезагружен, изменения применены'
    ),
    (
        'permissions-final-check',
        'permissions-junior',
        6,
        'Финальная проверка',
        'Проверьте, что сайт снова работает корректно',
        'Выполните запрос к сайту и проверьте ответ',
        'curl -I http://localhost',
        'Сайт доступен, возвращает 200 OK'
    ),

    -- =========================================================
    -- 2. memory-hunter
    -- =========================================================
    (
        'memory-check-ram',
        'memory-hunter',
        1,
        'Проверить использование памяти',
        'Нужно подтвердить, что проблема связана с оперативной памятью.',
        'Посмотри free или top.',
        'free -m',
        'Оперативная память была полностью занята или недавно резко освободилась после убийства процесса.'
    ),
    (
        'memory-check-processes',
        'memory-hunter',
        2,
        'Проверить запущенные процессы',
        'Посмотрите, какие процессы потребляют больше всего памяти',
        'Используйте top, htop или ps aux --sort=-%mem',
        'top',
        'Обнаружен процесс, потребляющий аномально много памяти'
    ),
    (
        'memory-check-oom',
        'memory-hunter',
        3,
        'Найти события OOM Killer',
        'Проверь системные логи и найди, убивал ли Linux какой-то процесс.',
        'Ищи oom или killed process.',
        'dmesg -T | grep -i oom',
        'В логах есть запись о том, что OOM Killer завершил backend или database process.'
    ),
    (
        'memory-check-postgres-status',
        'memory-hunter',
        4,
        'Проверить статус PostgreSQL',
        'Проверьте, работает ли PostgreSQL после OOM',
        'Используйте systemctl status postgresql',
        'systemctl status postgresql',
        'PostgreSQL остановлен или в состоянии failed'
    ),
    (
        'memory-restart-postgres',
        'memory-hunter',
        5,
        'Перезапустить PostgreSQL',
        'Запустите PostgreSQL заново после OOM Killer',
        'Используйте systemctl start или restart',
        'systemctl restart postgresql',
        'PostgreSQL успешно запущен'
    ),
    (
        'memory-restart-service',
        'memory-hunter',
        6,
        'Перезапустить упавший сервис',
        'После обнаружения причины нужно восстановить работу сервиса.',
        'Используй systemctl или docker compose.',
        'systemctl start backend',
        'Сервис снова запущен, но требуется дальнейшая оптимизация лимитов памяти.'
    ),
    (
        'memory-final-check',
        'memory-hunter',
        7,
        'Финальная проверка',
        'Убедитесь, что все сервисы работают корректно',
        'Проверьте статус всех сервисов',
        'docker compose ps',
        'Все сервисы работают, память в норме'
    ),

    -- =========================================================
    -- 3. dns-disappearance (оригинальное название из 13-й миграции)
    -- =========================================================
    (
        'dns-check-domain',
        'dns-mid',
        1,
        'Проверить доступность внешнего домена',
        'Проверь, может ли контейнер обратиться к внешнему API по доменному имени.',
        'Например, попробуй ping или curl.',
        'ping payment.example.com',
        'Доменное имя не резолвится или запрос завершается timeout.'
    ),
    (
        'dns-check-ip',
        'dns-mid',
        2,
        'Проверить доступность IP-адреса',
        'Нужно понять, проблема в сети или именно в DNS.',
        'Если ping 8.8.8.8 работает, сеть есть.',
        'ping 8.8.8.8',
        'IP-адрес доступен, значит проблема связана с DNS-резолвингом.'
    ),
    (
        'dns-check-ip-connectivity',
        'dns-mid',
        3,
        'Проверить сетевое подключение',
        'Дополнительная проверка доступности внешних IP',
        'Проверьте разные публичные DNS сервера',
        'ping 1.1.1.1',
        'Сеть работает, проблема именно в DNS разрешении'
    ),
    (
        'dns-check-nslookup',
        'dns-mid',
        4,
        'Проверить DNS через nslookup',
        'Используйте nslookup или dig для диагностики DNS',
        'Проверьте, какой DNS сервер используется',
        'nslookup google.com',
        'DNS запрос не проходит или возвращает неправильный IP'
    ),
    (
        'dns-check-resolv-conf',
        'dns-mid',
        5,
        'Проверить /etc/resolv.conf',
        'Посмотрите настройки DNS резолвера',
        'Проверьте, какие nameserver указаны',
        'cat /etc/resolv.conf',
        'В файле указан неверный или недоступный DNS сервер'
    ),
    (
        'dns-fix-resolv-conf',
        'dns-mid',
        6,
        'Исправить DNS-сервер',
        'Открой /etc/resolv.conf и исправь неверный DNS.',
        'Проверь, нет ли опечатки вроде 8.8.3.8.',
        'nano /etc/resolv.conf',
        'DNS исправлен на корректный публичный сервер, например 8.8.8.8 или 1.1.1.1.'
    ),
    (
        'dns-restart-resolved',
        'dns-mid',
        7,
        'Перезапустить systemd-resolved',
        'Если используется systemd-resolved, перезапустите его',
        'Очистите кэш DNS если необходимо',
        'systemctl restart systemd-resolved',
        'systemd-resolved перезапущен, кэш очищен'
    ),
    (
        'dns-final-check',
        'dns-mid',
        8,
        'Финальная проверка DNS',
        'Проверьте, что DNS теперь работает корректно',
        'Выполните запрос к проблемному домену',
        'ping google.com',
        'Домен успешно резолвится, проблема решена'
    ),

    -- =========================================================
    -- 4. rogue-trainee-nginx (оригинальное название из 13-й миграции)
    -- =========================================================
    (
        'rogue-run-nginx-test',
        'rogue-trainee',
        1,
        'Проверить конфигурацию Nginx',
        'Nginx не запускается после reload. Нужно проверить синтаксис конфига.',
        'У Nginx есть специальная команда проверки.',
        'nginx -t',
        'Команда показывает синтаксическую ошибку и номер строки.'
    ),
    (
        'rogue-check-nginx-status',
        'rogue-trainee',
        2,
        'Проверить статус Nginx',
        'Проверьте, в каком состоянии находится Nginx',
        'Используйте systemctl status nginx',
        'systemctl status nginx',
        'Nginx в состоянии failed или inactive'
    ),
    (
        'rogue-check-nginx-config',
        'rogue-trainee',
        3,
        'Проверить конфигурацию Nginx детально',
        'Запустите проверку конфигурации с выводом ошибок',
        'nginx -t покажет точное место ошибки',
        'nginx -t',
        'Обнаружена синтаксическая ошибка в конфигурации'
    ),
    (
        'rogue-show-broken-lines',
        'rogue-trainee',
        4,
        'Показать строки с ошибкой',
        'Откройте конфиг и посмотрите строки вокруг ошибки',
        'Используйте sed или cat для просмотра проблемных строк',
        'sed -n ''40,50p'' /etc/nginx/sites-available/default',
        'Видна пропущенная точка с запятой или закрывающая скобка'
    ),
    (
        'rogue-edit-config',
        'rogue-trainee',
        5,
        'Исправить конфигурационный файл',
        'Открой файл, указанный в ошибке nginx -t, и восстанови пропущенный символ.',
        'Чаще всего пропущена точка с запятой или закрывающая скобка.',
        'nano /etc/nginx/nginx.conf',
        'Конфигурация исправлена.'
    ),
    (
        'rogue-retest-config',
        'rogue-trainee',
        6,
        'Повторно проверить конфигурацию',
        'Убедитесь, что ошибка исправлена',
        'Запустите nginx -t снова',
        'nginx -t',
        'Конфигурация теперь валидна'
    ),
    (
        'rogue-restart-nginx',
        'rogue-trainee',
        7,
        'Запустить Nginx',
        'После исправления конфигурации нужно снова запустить сервис.',
        'Сначала можно еще раз выполнить nginx -t.',
        'systemctl restart nginx',
        'Nginx успешно запущен, сайт снова доступен.'
    ),
    (
        'rogue-final-check',
        'rogue-trainee',
        8,
        'Финальная проверка',
        'Проверьте, что сайт работает',
        'Выполните запрос к сайту',
        'curl -I http://localhost',
        'Сайт доступен, Nginx работает корректно'
    ),

    -- =========================================================
    -- 5. crashloop-auth-service (оригинальное название из 13-й миграции)
    -- =========================================================
    (
        'crashloop-check-containers',
        'crashloop-senior',
        1,
        'Проверить статус контейнеров',
        'Нужно найти контейнер, который циклически перезапускается.',
        'Ищи статус Restarting.',
        'docker ps',
        'Контейнер auth-service находится в статусе Restarting.'
    ),
    (
        'crashloop-check-logs',
        'crashloop-senior',
        2,
        'Посмотреть логи падающего контейнера',
        'Нужно понять, почему контейнер падает при старте.',
        'Используй docker logs.',
        'docker logs auth-service',
        'В логах ошибка подключения к базе данных: Connection refused.'
    ),
    (
        'crashloop-check-env',
        'crashloop-senior',
        3,
        'Проверить переменные окружения',
        'Проверьте настройки подключения к БД',
        'Посмотрите .env или docker-compose.yml',
        'cat .env | grep DB_',
        'Обнаружен неверный порт для подключения к PostgreSQL'
    ),
    (
        'crashloop-search-wrong-port',
        'crashloop-senior',
        4,
        'Найти неверный порт',
        'Найдите в конфигурации неверный порт БД',
        'Ищите 5433 вместо 5432',
        'grep -r "5433" .',
        'Найден неверный порт 5433 в .env или docker-compose.yml'
    ),
    (
        'crashloop-fix-env',
        'crashloop-senior',
        5,
        'Исправить переменные окружения',
        'Проверь порт базы данных в docker-compose или .env.',
        'Сравни DB_PORT с реальным портом Postgres внутри Docker-сети.',
        'nano docker-compose.yml',
        'Неверный порт исправлен на 5432'
    ),
    (
        'crashloop-fix-port',
        'crashloop-senior',
        6,
        'Исправить порт в конфигурации',
        'Замените неверный порт на правильный',
        'Используйте sed или редактор',
        'sed -i ''s/5433/5432/g'' .env',
        'Порт в конфигурации исправлен'
    ),
    (
        'crashloop-recreate-service',
        'crashloop-senior',
        7,
        'Пересоздать контейнер',
        'После исправления конфигурации пересоздайте контейнер',
        'Используйте docker compose up -d --force-recreate',
        'docker compose up -d --force-recreate auth-service',
        'Контейнер пересоздан и успешно запущен'
    ),
    (
        'crashloop-final-check',
        'crashloop-senior',
        8,
        'Финальная проверка',
        'Убедитесь, что контейнер больше не перезапускается',
        'Проверьте статус всех контейнеров',
        'docker ps',
        'Все контейнеры работают стабильно, проблема решена'
    )
) AS v(id, scenario_id, step_order, title, description, hint, expected_command, expected_result)
WHERE NOT EXISTS (SELECT 1 FROM scenario_steps LIMIT 1);

-- =========================================================
-- Удаляем старые записи для обновляемых сценариев
-- =========================================================

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
-- Добавляем accepted commands для всех шагов
-- =========================================================

INSERT INTO step_accepted_commands (
    step_id,
    command,
    match_type,
    description
) VALUES
-- permissions-junior
('permissions-observe-http', 'curl -I http://localhost', 'exact', 'Проверка HTTP-заголовков локального сайта'),
('permissions-observe-http', 'curl -i http://localhost', 'exact', 'Проверка HTTP-ответа с заголовками и телом'),
('permissions-observe-http', 'curl http://localhost', 'exact', 'Обычный HTTP-запрос к сайту'),

('permissions-check-nginx-logs', 'tail -f /var/log/nginx/error.log', 'exact', 'Просмотр error.log в real-time'),
('permissions-check-nginx-logs', 'tail -n 50 /var/log/nginx/error.log', 'exact', 'Просмотр последних 50 строк error.log'),
('permissions-check-nginx-logs', 'tail /var/log/nginx/error.log', 'exact', 'Просмотр конца error.log'),
('permissions-check-nginx-logs', 'less /var/log/nginx/error.log', 'exact', 'Открыть error.log через less'),
('permissions-check-nginx-logs', '/var/log/nginx/error.log', 'contains', 'Любая разумная команда, которая читает nginx error.log'),

('permissions-check-owner', 'ls -la /var/www/html', 'exact', 'Проверить владельца и права web-root'),
('permissions-check-owner', 'ls -l /var/www/html', 'exact', 'Короткая проверка прав web-root'),
('permissions-check-owner', '/var/www/html', 'contains', 'Любая команда просмотра /var/www/html'),

('permissions-fix-owner', 'chown -R www-data:www-data /var/www/html', 'exact', 'Исправить владельца web-root'),
('permissions-fix-owner', 'sudo chown -R www-data:www-data /var/www/html', 'exact', 'Исправить владельца через sudo'),

('permissions-validate-nginx', 'nginx -t', 'exact', 'Проверить конфигурацию nginx'),
('permissions-validate-nginx', 'sudo nginx -t', 'exact', 'Проверить конфигурацию nginx через sudo'),

('permissions-reload-nginx', 'systemctl reload nginx', 'exact', 'Reload nginx через systemctl'),
('permissions-reload-nginx', 'sudo systemctl reload nginx', 'exact', 'Reload nginx через sudo'),
('permissions-reload-nginx', 'nginx -s reload', 'exact', 'Reload nginx через nginx signal'),
('permissions-reload-nginx', 'sudo nginx -s reload', 'exact', 'Reload nginx signal через sudo'),

('permissions-final-check', 'curl -I http://localhost', 'exact', 'Финальная проверка HTTP-заголовков'),
('permissions-final-check', 'curl -i http://localhost', 'exact', 'Финальная проверка HTTP-ответа'),
('permissions-final-check', 'curl http://localhost', 'exact', 'Финальный HTTP-запрос'),

-- memory-hunter
('memory-check-ram', 'free -h', 'exact', 'Проверить память в удобном формате'),
('memory-check-ram', 'free -m', 'exact', 'Проверить память в мегабайтах'),
('memory-check-ram', 'cat /proc/meminfo', 'exact', 'Проверить память через procfs'),

('memory-check-processes', 'htop', 'exact', 'Интерактивный просмотр процессов'),
('memory-check-processes', 'top', 'exact', 'Просмотр процессов через top'),
('memory-check-processes', 'ps aux --sort=-%mem', 'exact', 'Отсортировать процессы по памяти'),
('memory-check-processes', 'ps aux', 'exact', 'Посмотреть процессы'),

('memory-check-oom', 'dmesg -T | grep -i oom', 'exact', 'Найти OOM в dmesg'),
('memory-check-oom', 'dmesg | grep -i oom', 'exact', 'Найти OOM в dmesg без human time'),
('memory-check-oom', 'journalctl -k | grep -i oom', 'exact', 'Найти OOM в kernel journal'),
('memory-check-oom', 'grep -i oom', 'contains', 'Команда ищет OOM'),

('memory-check-postgres-status', 'systemctl status postgresql', 'exact', 'Проверить статус PostgreSQL'),
('memory-check-postgres-status', 'sudo systemctl status postgresql', 'exact', 'Проверить статус PostgreSQL через sudo'),
('memory-check-postgres-status', 'journalctl -u postgresql -xe', 'exact', 'Посмотреть journal PostgreSQL'),

('memory-restart-postgres', 'systemctl restart postgresql', 'exact', 'Перезапустить PostgreSQL'),
('memory-restart-postgres', 'sudo systemctl restart postgresql', 'exact', 'Перезапустить PostgreSQL через sudo'),

('memory-restart-service', 'docker compose restart backend', 'exact', 'Перезапустить backend'),
('memory-restart-service', 'docker-compose restart backend', 'exact', 'Перезапустить backend старым docker-compose'),
('memory-restart-service', 'docker compose restart', 'exact', 'Перезапустить compose stack'),

('memory-final-check', 'docker compose ps', 'exact', 'Проверить compose services'),
('memory-final-check', 'docker-compose ps', 'exact', 'Проверить compose services старой командой'),
('memory-final-check', 'systemctl status postgresql', 'exact', 'Проверить PostgreSQL после восстановления'),

-- dns-disappearance
('dns-check-domain', 'ping google.com', 'exact', 'Проверить домен через ping'),
('dns-check-domain', 'curl https://google.com', 'exact', 'Проверить домен через curl'),
('dns-check-domain', 'curl https://api.payment.com', 'exact', 'Проверить внешний API по домену'),

('dns-check-ip', 'ping 8.8.8.8', 'exact', 'Проверить сеть по IP'),
('dns-check-ip', 'ping 1.1.1.1', 'exact', 'Проверить сеть по IP Cloudflare DNS'),

('dns-check-ip-connectivity', 'ping 8.8.8.8', 'exact', 'Проверить сеть по IP'),
('dns-check-ip-connectivity', 'ping 1.1.1.1', 'exact', 'Проверить сеть по IP Cloudflare DNS'),
('dns-check-ip-connectivity', 'curl http://1.1.1.1', 'exact', 'Проверить доступность IP через curl'),

('dns-check-nslookup', 'nslookup google.com', 'exact', 'Проверить DNS через nslookup'),
('dns-check-nslookup', 'dig google.com', 'exact', 'Проверить DNS через dig'),
('dns-check-nslookup', 'host google.com', 'exact', 'Проверить DNS через host'),

('dns-check-resolv-conf', 'cat /etc/resolv.conf', 'exact', 'Посмотреть resolv.conf'),
('dns-check-resolv-conf', 'cat /etc/resolv.conf | grep nameserver', 'exact', 'Посмотреть nameserver'),
('dns-check-resolv-conf', 'grep nameserver /etc/resolv.conf', 'exact', 'Найти nameserver'),
('dns-check-resolv-conf', '/etc/resolv.conf', 'contains', 'Любая команда чтения resolv.conf'),

('dns-fix-resolv-conf', 'echo -e "nameserver 8.8.8.8\nnameserver 1.1.1.1" > /etc/resolv.conf', 'exact', 'Записать корректные DNS'),
('dns-fix-resolv-conf', 'echo "nameserver 8.8.8.8" > /etc/resolv.conf', 'exact', 'Записать Google DNS'),
('dns-fix-resolv-conf', 'nameserver 8.8.8.8', 'contains', 'Команда содержит корректный nameserver'),

('dns-restart-resolved', 'systemctl restart systemd-resolved', 'exact', 'Перезапустить systemd-resolved'),
('dns-restart-resolved', 'sudo systemctl restart systemd-resolved', 'exact', 'Перезапустить systemd-resolved через sudo'),
('dns-restart-resolved', 'resolvectl flush-caches', 'exact', 'Очистить DNS cache'),

('dns-final-check', 'ping google.com', 'exact', 'Проверить домен после исправления'),
('dns-final-check', 'nslookup google.com', 'exact', 'Проверить DNS lookup после исправления'),
('dns-final-check', 'dig google.com', 'exact', 'Проверить DNS через dig'),

-- rogue-trainee-nginx
('rogue-run-nginx-test', 'nginx -t', 'exact', 'Проверить конфигурацию nginx'),
('rogue-run-nginx-test', 'sudo nginx -t', 'exact', 'Проверить конфигурацию nginx через sudo'),

('rogue-check-nginx-status', 'systemctl status nginx', 'exact', 'Проверить nginx status'),
('rogue-check-nginx-status', 'sudo systemctl status nginx', 'exact', 'Проверить nginx status через sudo'),
('rogue-check-nginx-status', 'docker logs nginx --tail 50', 'exact', 'Посмотреть логи nginx container'),

('rogue-check-nginx-config', 'nginx -t', 'exact', 'Проверить конфигурацию nginx'),
('rogue-check-nginx-config', 'sudo nginx -t', 'exact', 'Проверить конфигурацию nginx через sudo'),

('rogue-show-broken-lines', 'sed -n ''40,50p'' /etc/nginx/sites-available/default', 'exact', 'Показать строки вокруг ошибки'),
('rogue-show-broken-lines', 'cat /etc/nginx/sites-available/default | sed -n ''40,50p''', 'exact', 'Показать строки через cat и sed'),
('rogue-show-broken-lines', '/etc/nginx/sites-available/default', 'contains', 'Команда читает nginx site config'),

('rogue-edit-config', 'vim +45 /etc/nginx/sites-available/default', 'exact', 'Открыть конфиг в vim на строке 45'),
('rogue-edit-config', 'nano /etc/nginx/sites-available/default', 'exact', 'Открыть конфиг в nano'),
('rogue-edit-config', 'vim /etc/nginx/sites-available/default', 'exact', 'Открыть конфиг в vim'),

('rogue-retest-config', 'nginx -t', 'exact', 'Повторно проверить nginx config'),
('rogue-retest-config', 'sudo nginx -t', 'exact', 'Повторно проверить nginx config через sudo'),

('rogue-restart-nginx', 'systemctl restart nginx', 'exact', 'Перезапустить nginx'),
('rogue-restart-nginx', 'sudo systemctl restart nginx', 'exact', 'Перезапустить nginx через sudo'),
('rogue-restart-nginx', 'nginx -s reload', 'exact', 'Reload nginx после исправления'),

('rogue-final-check', 'curl -I http://localhost', 'exact', 'Проверить сайт'),
('rogue-final-check', 'curl http://localhost', 'exact', 'Проверить сайт обычным curl'),

-- crashloop-auth-service
('crashloop-check-containers', 'docker ps', 'exact', 'Проверить контейнеры'),
('crashloop-check-containers', 'docker compose ps', 'exact', 'Проверить compose services'),
('crashloop-check-containers', 'docker-compose ps', 'exact', 'Проверить compose services старой командой'),

('crashloop-check-logs', 'docker logs auth-service --tail 100', 'exact', 'Посмотреть последние логи auth-service'),
('crashloop-check-logs', 'docker logs --since 5m auth-service', 'exact', 'Посмотреть свежие логи auth-service'),
('crashloop-check-logs', 'docker-compose logs --tail=50 auth', 'exact', 'Посмотреть логи auth через docker-compose'),
('crashloop-check-logs', 'docker logs auth-service', 'exact', 'Посмотреть все логи auth-service'),

('crashloop-check-env', 'cat .env | grep DB_', 'exact', 'Посмотреть DB переменные'),
('crashloop-check-env', 'grep DB_ .env', 'exact', 'Найти DB переменные'),
('crashloop-check-env', 'cat .env', 'exact', 'Посмотреть .env'),

('crashloop-search-wrong-port', 'grep -r "5433" .', 'exact', 'Найти неверный порт 5433'),
('crashloop-search-wrong-port', 'grep -R "5433" .', 'exact', 'Найти неверный порт 5433 рекурсивно'),
('crashloop-search-wrong-port', '5433', 'contains', 'Команда ищет 5433'),

('crashloop-fix-env', 'nano docker-compose.yml', 'exact', 'Открыть docker-compose.yml в nano'),
('crashloop-fix-env', 'vim docker-compose.yml', 'exact', 'Открыть docker-compose.yml в vim'),

('crashloop-fix-port', 'sed -i ''s/5433/5432/g'' .env', 'exact', 'Исправить порт в .env'),
('crashloop-fix-port', 'DB_PORT=5432', 'contains', 'Команда устанавливает корректный DB_PORT'),

('crashloop-recreate-service', 'docker compose up -d --force-recreate auth-service', 'exact', 'Пересоздать auth-service'),
('crashloop-recreate-service', 'docker-compose up -d --force-recreate auth-service', 'exact', 'Пересоздать auth-service старым compose'),
('crashloop-recreate-service', 'docker compose restart auth-service', 'exact', 'Перезапустить auth-service'),

('crashloop-final-check', 'docker ps', 'exact', 'Финально проверить контейнеры'),
('crashloop-final-check', 'docker compose ps', 'exact', 'Финально проверить compose services'),
('crashloop-final-check', 'docker-compose ps', 'exact', 'Финально проверить compose services старой командой');