CREATE TABLE IF NOT EXISTS scenarios (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    difficulty VARCHAR(50) NOT NULL,
    category VARCHAR(100) NOT NULL,
    estimated_minutes INT NOT NULL DEFAULT 15,
    user_notification TEXT NOT NULL,
    desktop_symptoms TEXT NOT NULL,
    terminal_solution TEXT NOT NULL,
    quick_fix TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO scenarios (
    id,
    title,
    description,
    difficulty,
    category,
    estimated_minutes,
    user_notification,
    desktop_symptoms,
    terminal_solution,
    quick_fix,
    is_active
) VALUES
(
    'permissions-junior',
    'Наперекосяк с правами',
    'После обновления конфигурации или деплоя права на SSL-сертификаты или на корневую папку сайта /var/www/html случайно изменились на root:root вместо www-data:www-data.',
    'junior',
    'linux/nginx',
    15,
    'Критический сбой веб-сервера! Ошибка 403 Forbidden / 500 Internal Error.',
    'Трафик резко падает до нуля, графики процессора на нуле, сервер простаивает.',
    'Студент должен зайти в контейнер, проверить логи Nginx через tail -f /var/log/nginx/error.log, увидеть Permission denied, вспомнить команду chown -R www-data:www-data /var/www и исправить права.',
    'Сбросить права доступа веб-сервера на дефолтные',
    TRUE
),
(
    'memory-hunter',
    'Охотник за памятью',
    'В backend-приложении есть утечка памяти. В момент пиковой нагрузки оперативная память заканчивается, и Linux OOM Killer завершает процесс базы данных или backend.',
    'mid',
    'linux/backend',
    20,
    'База данных недоступна. Ошибка 502 Bad Gateway.',
    'График RAM плавно рос до 100%, затем резко упал. Процесс backend отсутствует в списке запущенных.',
    'Студент должен проверить системные логи через dmesg -T | grep -i oom или grep -i kill /var/log/syslog, обнаружить, что процесс был убит системой, запустить его заново и оптимизировать лимиты памяти.',
    'Экстренная перезагрузка СУБД и очистка кэша RAM',
    TRUE
),
(
    'dns-mid',
    'Таинственное исчезновение DNS',
    'Контейнер с backend потерял связь с внешним миром, потому что в /etc/resolv.conf прописан неверный или недоступный DNS-сервер, например 8.8.3.8 вместо 8.8.8.8.',
    'mid',
    'networking/dns',
    20,
    'Таймаут операций оплаты. Пользователи не могут совершить покупки.',
    'График сетевых ошибок и timeout резко растет. Появляется ошибка 504 Gateway Timeout.',
    'Студент пытается сделать ping внешнего ресурса — не работает. Делает ping 8.8.8.8 — работает. Понимает, что проблема в DNS, открывает /etc/resolv.conf, находит опечатку и исправляет ее.',
    'Применить публичные DNS Google/Cloudflare',
    TRUE
),
(
    'rogue-trainee',
    'Каскадный обвал из-за Rogue Trainee',
    'Стажер случайно удалил точку с запятой или закрывающую скобку в конфигурации Nginx, после чего выполнил nginx -s reload. Сервер упал и не поднимается.',
    'senior',
    'nginx/config',
    25,
    'Полный отказ инфраструктуры. Сайт лежит.',
    'Полный ноль по всем графикам. Контейнер Nginx имеет статус Exited.',
    'Студент должен запустить nginx -t, прочитать отчет синтаксического анализатора, найти строку с ошибкой, открыть конфиг через nano или vim, вернуть точку с запятой или скобку и запустить сервис.',
    'Откатить конфигурацию Nginx к последнему стабильному Git-коммиту',
    TRUE
),
(
    'crashloop-senior',
    'Петля смерти в контейнерах: CrashLoopBackOff',
    'Контейнер с микросервисом циклически падает при старте. Причина — в конфигурации или переменных окружения указан неверный порт для подключения к базе данных.',
    'senior',
    'docker/kubernetes',
    30,
    'Микросервис авторизации нестабилен. Циклическая перезагрузка.',
    'Скачкообразные графики CPU: пик при старте, затем падение в ноль. Контейнер постоянно находится в Restarting.',
    'Студент должен посмотреть docker logs container_id, увидеть Connection refused к базе данных, проверить переменные окружения, найти неправильный порт и исправить docker-compose или Kubernetes manifest.',
    'Форсировать перезапуск всего стека с дефолтными .env-параметрами',
    TRUE
)
ON CONFLICT (id) DO NOTHING;