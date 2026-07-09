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

INSERT INTO scenario_steps (
    id,
    scenario_id,
    step_order,
    title,
    description,
    hint,
    expected_command,
    expected_result
) VALUES

-- 1. permissions-junior

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

-- 2. memory-hunter

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
    'memory-check-oom',
    'memory-hunter',
    2,
    'Найти события OOM Killer',
    'Проверь системные логи и найди, убивал ли Linux какой-то процесс.',
    'Ищи oom или killed process.',
    'dmesg -T | grep -i oom',
    'В логах есть запись о том, что OOM Killer завершил backend или database process.'
),
(
    'memory-restart-service',
    'memory-hunter',
    3,
    'Перезапустить упавший сервис',
    'После обнаружения причины нужно восстановить работу сервиса.',
    'Используй systemctl или docker compose.',
    'systemctl start backend',
    'Сервис снова запущен, но требуется дальнейшая оптимизация лимитов памяти.'
),

-- 3. dns-mid

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
    'dns-fix-resolv',
    'dns-mid',
    3,
    'Исправить DNS-сервер',
    'Открой /etc/resolv.conf и исправь неверный DNS.',
    'Проверь, нет ли опечатки вроде 8.8.3.8.',
    'nano /etc/resolv.conf',
    'DNS исправлен на корректный публичный сервер, например 8.8.8.8 или 1.1.1.1.'
),

-- 4. rogue-trainee

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
    'rogue-edit-config',
    'rogue-trainee',
    2,
    'Исправить конфигурационный файл',
    'Открой файл, указанный в ошибке nginx -t, и восстанови пропущенный символ.',
    'Чаще всего пропущена точка с запятой или закрывающая скобка.',
    'nano /etc/nginx/nginx.conf',
    'Конфигурация исправлена.'
),
(
    'rogue-restart-nginx',
    'rogue-trainee',
    3,
    'Запустить Nginx',
    'После исправления конфигурации нужно снова запустить сервис.',
    'Сначала можно еще раз выполнить nginx -t.',
    'systemctl restart nginx',
    'Nginx успешно запущен, сайт снова доступен.'
),

-- 5. crashloop-senior

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
    'crashloop-fix-env',
    'crashloop-senior',
    3,
    'Исправить переменные окружения',
    'Проверь порт базы данных в docker-compose или .env.',
    'Сравни DB_PORT с реальным портом Postgres внутри Docker-сети.',
    'nano docker-compose.yml',
    'Неверный порт исправлен, контейнер больше не падает.'
)

ON CONFLICT (id) DO NOTHING;