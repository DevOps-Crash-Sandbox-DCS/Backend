BEGIN;

-- В dev-режиме проще сбросить тестовые прохождения,
-- потому что sessions/actions могут ссылаться на старые step_id.
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

-- =========================================================
-- 1. permissions-junior
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
) VALUES
(
    'permissions-observe-http',
    'permissions-junior',
    1,
    'Проверить HTTP-ответ',
    'Сначала нужно подтвердить внешний симптом: сайт действительно отвечает ошибкой 403 Forbidden или 500 Internal Error.',
    'Начни с простого HTTP-запроса к локальному сайту.',
    'curl -I http://localhost',
    'HTTP-ответ показывает 403 Forbidden или 500 Internal Server Error.'
),
(
    'permissions-check-nginx-logs',
    'permissions-junior',
    2,
    'Проверить логи Nginx',
    'Нужно понять техническую причину ошибки веб-сервера.',
    'Открой error.log и найди Permission denied.',
    'tail -f /var/log/nginx/error.log',
    'В логах обнаружена ошибка Permission denied при обращении к файлам сайта или SSL-сертификатам.'
),
(
    'permissions-check-owner',
    'permissions-junior',
    3,
    'Проверить владельца файлов',
    'После Permission denied нужно проверить владельца и группу файлов сайта.',
    'Используй ls -la для просмотра владельца и группы.',
    'ls -la /var/www/html',
    'Файлы принадлежат root:root вместо www-data:www-data.'
),
(
    'permissions-fix-owner',
    'permissions-junior',
    4,
    'Исправить владельца web-root',
    'Нужно вернуть корректного владельца директории сайта.',
    'Nginx обычно работает от пользователя www-data.',
    'chown -R www-data:www-data /var/www/html',
    'Владелец файлов исправлен на www-data:www-data.'
),
(
    'permissions-check-nginx-config',
    'permissions-junior',
    5,
    'Проверить конфигурацию Nginx',
    'Перед reload нужно убедиться, что конфигурация nginx валидна.',
    'Безопаснее сначала выполнить nginx -t.',
    'nginx -t',
    'Nginx сообщает, что syntax is ok и test is successful.'
),
(
    'permissions-reload-nginx',
    'permissions-junior',
    6,
    'Перезагрузить Nginx',
    'После исправления прав нужно применить изменения без полного простоя.',
    'Используй reload, если конфигурация валидна.',
    'systemctl reload nginx',
    'Nginx успешно перечитал конфигурацию.'
),
(
    'permissions-final-check',
    'permissions-junior',
    7,
    'Проверить восстановление сайта',
    'Финально убедиться, что сайт снова доступен.',
    'Повтори HTTP-проверку.',
    'curl -I http://localhost',
    'Сайт больше не возвращает 403/500.'
);

-- =========================================================
-- 2. memory-hunter
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
) VALUES
(
    'memory-check-ram',
    'memory-hunter',
    1,
    'Проверить использование памяти',
    'Нужно подтвердить, что проблема связана с оперативной памятью.',
    'Посмотри свободную и занятую память.',
    'free -h',
    'Оперативная память почти полностью занята или недавно резко освободилась после убийства процесса.'
),
(
    'memory-check-processes',
    'memory-hunter',
    2,
    'Проверить процессы',
    'Нужно понять, какие процессы потребляют память и какие сервисы могли исчезнуть.',
    'Используй top, htop или ps.',
    'htop',
    'Видно высокое потребление памяти или отсутствие backend/database процесса.'
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
    'После OOM нужно проверить, жив ли сервис базы данных.',
    'Посмотри статус postgresql.',
    'systemctl status postgresql',
    'PostgreSQL находится в failed/inactive или недавно был перезапущен.'
),
(
    'memory-restart-postgres',
    'memory-hunter',
    5,
    'Перезапустить PostgreSQL',
    'Если база была убита OOM Killer, нужно восстановить ее работу.',
    'Используй systemctl restart.',
    'systemctl restart postgresql',
    'PostgreSQL успешно запущен.'
),
(
    'memory-restart-backend',
    'memory-hunter',
    6,
    'Перезапустить backend',
    'После восстановления базы нужно поднять backend-приложение.',
    'Перезапусти backend-контейнер.',
    'docker compose restart backend',
    'Backend перезапущен и может подключиться к базе.'
),
(
    'memory-final-check',
    'memory-hunter',
    7,
    'Проверить состояние сервисов',
    'Финально убедиться, что сервисы снова работают.',
    'Посмотри состояние docker compose.',
    'docker compose ps',
    'Backend и база данных находятся в рабочем состоянии.'
);

-- =========================================================
-- 3. dns-mid
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
) VALUES
(
    'dns-check-domain',
    'dns-mid',
    1,
    'Проверить доступность домена',
    'Сначала нужно проверить, резолвится ли внешний домен.',
    'Попробуй ping доменного имени.',
    'ping google.com',
    'Доменное имя не резолвится или команда завершается ошибкой DNS.'
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
    'dns-check-nslookup',
    'dns-mid',
    3,
    'Проверить DNS через nslookup',
    'Нужно подтвердить, что DNS-запросы не проходят.',
    'Используй nslookup для проверки домена.',
    'nslookup google.com',
    'nslookup возвращает timeout или ошибку DNS-сервера.'
),
(
    'dns-check-resolv',
    'dns-mid',
    4,
    'Проверить resolv.conf',
    'Нужно посмотреть, какой DNS-сервер прописан в системе.',
    'Проверь nameserver в /etc/resolv.conf.',
    'cat /etc/resolv.conf',
    'В файле указан неправильный DNS-сервер, например 8.8.3.8.'
),
(
    'dns-fix-resolv',
    'dns-mid',
    5,
    'Исправить DNS-сервер',
    'Нужно заменить неправильный DNS на корректный публичный сервер.',
    'Используй 8.8.8.8 или 1.1.1.1.',
    'echo -e "nameserver 8.8.8.8\nnameserver 1.1.1.1" > /etc/resolv.conf',
    'В resolv.conf прописаны корректные DNS-серверы.'
),
(
    'dns-restart-resolved',
    'dns-mid',
    6,
    'Перезапустить resolver',
    'После изменения DNS-настроек нужно применить изменения.',
    'Перезапусти systemd-resolved.',
    'systemctl restart systemd-resolved',
    'DNS resolver перезапущен.'
),
(
    'dns-final-check',
    'dns-mid',
    7,
    'Проверить DNS повторно',
    'Финально убедиться, что домены снова резолвятся.',
    'Повтори проверку домена.',
    'ping google.com',
    'Домен успешно резолвится и отвечает.'
);

-- =========================================================
-- 4. rogue-trainee
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
) VALUES
(
    'rogue-check-nginx-status',
    'rogue-trainee',
    1,
    'Проверить статус Nginx',
    'Нужно понять, работает nginx или упал после reload.',
    'Начни со статуса сервиса.',
    'systemctl status nginx',
    'Nginx находится в failed/inactive или показывает ошибку запуска.'
),
(
    'rogue-run-nginx-test',
    'rogue-trainee',
    2,
    'Проверить конфигурацию Nginx',
    'Nginx не запускается после reload. Нужно проверить синтаксис конфига.',
    'У Nginx есть специальная команда проверки.',
    'nginx -t',
    'Команда показывает синтаксическую ошибку и номер строки.'
),
(
    'rogue-show-broken-lines',
    'rogue-trainee',
    3,
    'Посмотреть проблемный участок конфига',
    'Нужно открыть строки вокруг места ошибки.',
    'Выведи строки около проблемной строки из nginx -t.',
    'sed -n ''40,50p'' /etc/nginx/sites-available/default',
    'Видно лишнюю скобку, пропущенную точку с запятой или некорректную директиву.'
),
(
    'rogue-edit-config',
    'rogue-trainee',
    4,
    'Исправить конфигурационный файл',
    'Открой файл, указанный в ошибке nginx -t, и восстанови пропущенный символ.',
    'Чаще всего пропущена точка с запятой или закрывающая скобка.',
    'nano /etc/nginx/sites-available/default',
    'Конфигурация исправлена.'
),
(
    'rogue-retest-nginx',
    'rogue-trainee',
    5,
    'Повторно проверить конфигурацию',
    'Перед запуском сервиса нужно убедиться, что конфигурация теперь валидна.',
    'Снова выполни nginx -t.',
    'nginx -t',
    'Nginx сообщает, что syntax is ok и test is successful.'
),
(
    'rogue-restart-nginx',
    'rogue-trainee',
    6,
    'Запустить Nginx',
    'После исправления конфигурации нужно снова запустить сервис.',
    'Используй systemctl restart nginx.',
    'systemctl restart nginx',
    'Nginx успешно запущен.'
),
(
    'rogue-final-check',
    'rogue-trainee',
    7,
    'Проверить сайт',
    'Финально убедиться, что сайт снова доступен.',
    'Проверь HTTP-ответ.',
    'curl -I http://localhost',
    'Сайт снова отвечает.'
);

-- =========================================================
-- 5. crashloop-senior
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
) VALUES
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
    'docker logs auth-service --tail 100',
    'В логах ошибка подключения к базе данных: Connection refused.'
),
(
    'crashloop-check-env',
    'crashloop-senior',
    3,
    'Проверить переменные окружения',
    'Нужно проверить параметры подключения к базе данных.',
    'Посмотри DB_* переменные в .env.',
    'cat .env | grep DB_',
    'Видно, что DB_PORT указан неправильно.'
),
(
    'crashloop-search-wrong-port',
    'crashloop-senior',
    4,
    'Найти неверный порт',
    'Нужно найти место, где указан неправильный порт 5433.',
    'Используй grep по проекту.',
    'grep -r "5433" .',
    'Найден неверный порт 5433 в .env или docker-compose.yml.'
),
(
    'crashloop-fix-port',
    'crashloop-senior',
    5,
    'Исправить порт базы данных',
    'PostgreSQL внутри Docker-сети должен быть доступен на 5432, а не 5433.',
    'Замени 5433 на 5432.',
    'sed -i ''s/5433/5432/g'' .env',
    'DB_PORT исправлен на 5432.'
),
(
    'crashloop-recreate-service',
    'crashloop-senior',
    6,
    'Пересоздать auth-service',
    'После изменения .env нужно пересоздать контейнер, чтобы он получил новые переменные.',
    'Используй docker compose up с force-recreate.',
    'docker compose up -d --force-recreate auth-service',
    'auth-service пересоздан с корректными переменными.'
),
(
    'crashloop-final-check',
    'crashloop-senior',
    7,
    'Проверить контейнеры повторно',
    'Финально убедиться, что auth-service больше не падает.',
    'Снова выполни docker ps.',
    'docker ps',
    'auth-service находится в статусе Up, а не Restarting.'
);

COMMIT;