# Спецификация проекта DCS

## 1. Назначение проекта

DCS — это backend-платформа для интерактивного тренажёра по расследованию и устранению инфраструктурных инцидентов.

Проект моделирует рабочий процесс junior/middle DevOps, SRE или backend-инженера:

1. Пользователь выбирает сценарий инцидента.
2. Backend создаёт пользовательскую сессию.
3. Пользователь подключается к терминалу через WebSocket.
4. Для sandbox-сценариев backend создаёт изолированный Docker-контейнер.
5. Пользователь выполняет команды диагностики и исправления.
6. Backend проверяет действия, начисляет баллы и переводит пользователя по шагам сценария.
7. Пользователь может запросить ML-подсказку.
8. После завершения сессии формируется отчёт.

Проект предназначен для:

- обучения инженеров диагностике инцидентов;
- демонстрации backend-архитектуры;
- демонстрации интеграции Go backend, PostgreSQL, Docker sandbox и ML-сервиса;
- портфолио-проекта уровня middle backend engineer.

---

## 2. Основные возможности

### 2.1 Аутентификация

Пользователь может:

- зарегистрироваться;
- войти в систему;
- получить JWT access token.

JWT используется для защиты API endpoints.

---

### 2.2 Сценарии инцидентов

Сценарий описывает последовательность шагов, которые пользователь должен пройти.

Каждый сценарий содержит:

- идентификатор;
- название;
- описание;
- уровень сложности;
- список шагов;
- ожидаемые команды;
- ожидаемые результаты;
- подсказки;
- баллы за правильные действия.

Пример сценария:
permissions-junior

Суть сценария:
nginx возвращает 403 Forbidden из-за неправильных прав или владельца файлов в /var/www/html.

### 2.3 Сессии пользователя
Сессия — это индивидуальное прохождение сценария конкретным пользователем.

Сессия хранит:

пользователя;
сценарий;
текущий шаг;
статус;
текущий score;
время создания;
время обновления;
время завершения.
Возможные статусы:

in_progress
completed
failed

### 2.4 WebSocket terminal
Для интерактивной работы используется WebSocket endpoint.

Через WebSocket пользователь отправляет команды:

{
  "type": "command",
  "command": "ls -la /var/www/html"
}

Backend возвращает результат:
{
  "type": "output",
  "output": "...",
  "stdout": "...",
  "stderr": "",
  "exitCode": 0
}

Для sandbox-сценариев команда выполняется внутри Docker-контейнера.
Для обычных сценариев команда может обрабатываться симулятором.

### 2.5 Docker sandbox
Для некоторых сценариев backend создаёт Docker-контейнер, в котором пользователь выполняет команды.

Sandbox создаётся при подключении к terminal WebSocket.

Sandbox удаляется:

при завершении сценария;
при отключении WebSocket;
вручную через API;
через cleanup старых sandbox-контейнеров.
Для сценария:
permissions-junior
используется image: dcs-scenario-permissions-junior:dev

Контейнер создаётся с ограничениями:
memory limit;
CPU limit;
pids limit;
no-new-privileges.

### 2.6 Проверка действий пользователя
После выполнения команды backend отправляет действие в actions service.

Action содержит:

session ID;
step ID;
command;
результат проверки;
баллы;
feedback.
Если команда соответствует ожидаемому действию, backend:

начисляет баллы;
переводит сессию на следующий шаг;
завершает сессию, если шаг был последним.

### 2.7 ML hints
Пользователь может запросить подсказку:


POST /api/sessions/:id/hints
Backend собирает контекст:

session;
current step;
историю команд;
score;
статус сессии.
После этого backend отправляет запрос во внешний Python ML service.

Если ML service доступен, backend возвращает ML-подсказку.

Если ML service недоступен, backend возвращает fallback-подсказку из текущего шага сценария.

Все подсказки сохраняются в PostgreSQL.

### 2.8 Reports
После прохождения сессии backend может сформировать отчёт.

Отчёт содержит:

session ID;
user ID;
scenario ID;
итоговый score;
статус;
список действий;
успешные и неуспешные команды;
рекомендации.

### 4. Основные компоненты backend
### 4.1 Auth module
Пакет:
internal/auth
Отвечает за:

регистрацию;
логин;
генерацию JWT;
проверку паролей;
работу с пользователями.

### 4.2 Scenarios module
Пакет:
internal/scenarios
Отвечает за:

список сценариев;
получение сценария по ID;
шаги сценария.

### 4.3 Sessions module
Пакет:
internal/sessions
Отвечает за:

создание сессии;
получение сессии;
список сессий пользователя;
обновление статуса;
переход между шагами;
score.

### 4.4 Actions module
Пакет:
internal/actions
Отвечает за:

приём команд пользователя;
проверку команды на соответствие текущему шагу;
начисление баллов;
сохранение истории действий;
перевод сессии на следующий шаг.

### 4.5 Terminal module
Пакет:
internal/terminal
Отвечает за:

WebSocket terminal;
приём команд;
отправку stdout/stderr/exitCode;
интеграцию с sandbox manager;
fallback-симуляцию команд.

### 4.6 Sandbox module
Пакет:
internal/sandbox
Отвечает за:

создание Docker-контейнера;
выполнение команд внутри контейнера;
удаление контейнера;
хранение состояния sandbox в PostgreSQL;
cleanup старых sandbox.

### 4.7 Hints module
Пакет:
internal/hints
Отвечает за:

сбор контекста сессии;
запрос подсказки у Python ML service;
fallback-подсказки;
сохранение подсказок в PostgreSQL.

### 4.8 Reports module
Пакет:
internal/reports
Отвечает за:

генерацию отчёта по сессии;
агрегацию действий пользователя;
итоговую оценку прохождения.


6. API спецификация
### 6.1 Auth
Регистрация

POST /api/auth/register
Request:
{
  "email": "user@example.com",
  "password": "password123",
  "name": "User"
}


Response:
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "User"
  },
  "accessToken": "jwt-token"
}

Логин
POST /api/auth/login
Request:
{
  "email": "user@example.com",
  "password": "password123"
}
Response:
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "User"
  },
  "accessToken": "jwt-token"
}

### 6.2 Scenarios
Получить список сценариев
GET /api/scenarios
Headers:
Authorization: Bearer <token>

Response:

[
  {
    "id": "permissions-junior",
    "title": "Nginx 403 Forbidden",
    "description": "Диагностика проблемы с правами файлов",
    "difficulty": "junior"
  }
]
Получить сценарий по ID

GET /api/scenarios/:id
Headers:
Authorization: Bearer <token>
Response:
{
  "id": "permissions-junior",
  "title": "Nginx 403 Forbidden",
  "description": "Диагностика проблемы с правами файлов",
  "difficulty": "junior",
  "steps": [
    {
      "id": "permissions-observe-http",
      "title": "Проверить HTTP-ответ",
      "description": "Проверь, какой HTTP статус возвращает nginx",
      "hint": "Используй curl",
      "expectedCommand": "curl -I http://localhost"
    }
  ]
}
6.3 Sessions
Создать сессию

POST /api/sessions
Headers:


Authorization: Bearer <token>
Request:


{
  "scenarioId": "permissions-junior"
}
Response:


{
  "id": "uuid",
  "userId": "uuid",
  "scenarioId": "permissions-junior",
  "currentStepId": "permissions-observe-http",
  "status": "in_progress",
  "score": 0
}
Получить список сессий

GET /api/sessions
Headers:


Authorization: Bearer <token>
Response:


[
  {
    "id": "uuid",
    "scenarioId": "permissions-junior",
    "currentStepId": "permissions-observe-http",
    "status": "in_progress",
    "score": 0
  }
]
Получить сессию по ID

GET /api/sessions/:id
Headers:


Authorization: Bearer <token>
Response:


{
  "id": "uuid",
  "userId": "uuid",
  "scenarioId": "permissions-junior",
  "currentStepId": "permissions-observe-http",
  "status": "in_progress",
  "score": 0
}
6.4 Terminal WebSocket

GET /api/sessions/:id/terminal
Protocol:


WebSocket
Headers:


Authorization: Bearer <token>
Welcome message:


{
  "type": "welcome",
  "output": "Connected to incident simulator terminal.\nDocker sandbox started.\n",
  "sessionId": "uuid",
  "currentStepId": "permissions-observe-http"
}
Client command:


{
  "type": "command",
  "command": "ls -la /var/www/html"
}
Server output:


{
  "type": "output",
  "output": "total 16\n...",
  "stdout": "total 16\n...",
  "stderr": "",
  "exitCode": 0
}
Action result:


{
  "type": "action_result",
  "sessionId": "uuid",
  "currentStepId": "permissions-fix-owner",
  "actionResult": {
    "isCorrect": true,
    "points": 10,
    "feedback": "Команда выполнена верно"
  }
}
Error:


{
  "type": "error",
  "error": "command is empty"
}
6.5 Sandbox
Запустить sandbox вручную

POST /api/sandbox/start
Headers:


Authorization: Bearer <token>
Request:


{
  "sessionId": "uuid",
  "scenarioId": "permissions-junior",
  "image": "dcs-scenario-permissions-junior:dev"
}
Response:


{
  "id": "uuid",
  "sessionId": "uuid",
  "scenarioId": "permissions-junior",
  "containerName": "dcs-sandbox-uuid",
  "image": "dcs-scenario-permissions-junior:dev",
  "status": "running"
}
Получить sandbox по session ID

GET /api/sandbox/:sessionId
Headers:


Authorization: Bearer <token>
Response:


{
  "id": "uuid",
  "sessionId": "uuid",
  "scenarioId": "permissions-junior",
  "containerName": "dcs-sandbox-uuid",
  "image": "dcs-scenario-permissions-junior:dev",
  "status": "running"
}
Выполнить команду в sandbox

POST /api/sandbox/:sessionId/exec
Headers:


Authorization: Bearer <token>
Request:


{
  "command": "ls -la /var/www/html"
}
Response:


{
  "stdout": "total 16\n...",
  "stderr": "",
  "exitCode": 0
}
Остановить sandbox

POST /api/sandbox/:sessionId/stop
Headers:


Authorization: Bearer <token>
Response:


{
  "status": "stopped"
}
Cleanup старых sandbox

POST /api/sandbox/cleanup
Headers:


Authorization: Bearer <token>
Request:


{
  "olderThanMinutes": 60
}
Response:


{
  "stoppedInDb": 3,
  "removedContainers": 3
}
6.6 Hints
Получить подсказку

POST /api/sessions/:id/hints
Headers:


Authorization: Bearer <token>
Request:


{
  "level": "basic"
}
Допустимые уровни:


basic
detailed
direct
Response:


{
  "hint": "Проверь владельца файлов в /var/www/html.",
  "confidence": 0.82,
  "source": "ml",
  "currentStepId": "permissions-check-owner"
}
Если ML service недоступен:


{
  "hint": "Используй ls -la для просмотра владельца и группы.",
  "confidence": 0.3,
  "source": "fallback",
  "currentStepId": "permissions-check-owner"
}
6.7 Reports
Получить отчёт по сессии

GET /api/reports/sessions/:sessionId
Headers:


Authorization: Bearer <token>
Response:


{
  "sessionId": "uuid",
  "scenarioId": "permissions-junior",
  "status": "completed",
  "score": 100,
  "actions": [
    {
      "stepId": "permissions-observe-http",
      "command": "curl -I http://localhost",
      "isCorrect": true,
      "points": 10,
      "feedback": "Команда выполнена верно"
    }
  ]
}