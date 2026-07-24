from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
from typing import List, Optional
from datetime import datetime
import requests
import uvicorn
import os

app = FastAPI(
    title="DevOps Crash Sandbox - AI Hint Service",
    description="HTTP API для генерации ML-подсказок студентам"
)

OLLAMA_URL = os.getenv("OLLAMA_URL", "http://ollama:11434/api/chat")
MODEL_NAME = os.getenv("MODEL_NAME", "qwen2.5-coder:7b")


class MLCurrentStep(BaseModel):
    id: str
    title: str
    description: str
    hint: str = ""
    expectedCommand: str = ""
    expectedResult: str = ""


class MLActionEntry(BaseModel):
    stepId: str
    command: str
    isCorrect: bool
    points: int
    feedback: str = ""
    createdAt: datetime


class MLHintRequest(BaseModel):
    userId: str
    sessionId: str
    scenarioId: str
    currentStepId: Optional[str] = None
    sessionStatus: str
    score: int
    hintLevel: str = "basic"
    currentStep: Optional[MLCurrentStep] = None
    history: List[MLActionEntry] = Field(default_factory=list)
    recentTerminalOutput: str = ""


class MLHintResponse(BaseModel):
    hint: str
    confidence: Optional[float] = None
    source: str = "ml"
    reasoning: Optional[str] = None


@app.get("/health")
def health():
    return {
        "status": "ok",
        "model": MODEL_NAME
    }


@app.post("/api/v1/hints", response_model=MLHintResponse)
def get_help_hint(request: MLHintRequest):
    print(f"[AI Сервер] Получен запрос подсказки для sessionId={request.sessionId}, scenarioId={request.scenarioId}")

    if request.history:
        commands_flat = "\n".join(
            [
                f"- Команда: {item.command}\n"
                f"  Корректность: {'да' if item.isCorrect else 'нет'}\n"
                f"  Feedback: {item.feedback or 'нет'}"
                for item in request.history
            ]
        )
    else:
        commands_flat = "Пользователь ещё не вводил команды."

    if request.currentStep:
        current_step_text = (
            f"ID шага: {request.currentStep.id}\n"
            f"Название шага: {request.currentStep.title}\n"
            f"Описание шага: {request.currentStep.description}\n"
            f"Встроенная подсказка шага: {request.currentStep.hint or 'нет'}\n"
            f"Ожидаемая команда: {request.currentStep.expectedCommand or 'не указана'}\n"
            f"Ожидаемый результат: {request.currentStep.expectedResult or 'не указан'}"
        )
    else:
        current_step_text = "Текущий шаг не указан."

    system_instruction = (
        "Вы — опытный DevOps-наставник. Ваша цель — помочь пользователю локализовать и исправить проблему "
        "в Linux/Docker/Kubernetes/Git/CI/CD среде.\n\n"

        "ПРАВИЛА:\n"
        "1. Категорически запрещено писать готовые консольные команды, утилиты или bash-код.\n"
        "2. Запрещено использовать обратные кавычки для команд, имён файлов и кода.\n"
        "3. Не раскрывай прямое решение, если уровень подсказки basic.\n"
        "4. Если уровень detailed — дай более конкретное направление, но без готовой команды.\n"
        "5. Если уровень direct — можно назвать область, файл, сервис или настройку, но всё равно не писать готовую команду.\n"
        "6. Анализируй историю действий: что пользователь уже сделал правильно, а где смотрит не туда.\n"
        "7. Ответ должен быть строго на русском языке.\n"
        "8. Ответ должен быть кратким: 1-3 предложения.\n"
        "9. Обращайся к пользователю напрямую, не называй его студентом.\n"
        "10. Следи за грамотностью. Для символа ; используй термин 'точка с запятой'."
    )

    user_prompt = (
        f"СЦЕНАРИЙ: {request.scenarioId}\n"
        f"СЕССИЯ: {request.sessionId}\n"
        f"СТАТУС СЕССИИ: {request.sessionStatus}\n"
        f"СЧЁТ: {request.score}\n"
        f"УРОВЕНЬ ПОДСКАЗКИ: {request.hintLevel}\n\n"

        f"ТЕКУЩИЙ ШАГ:\n{current_step_text}\n\n"

        f"ИСТОРИЯ ДЕЙСТВИЙ:\n{commands_flat}\n\n"

        f"ПОСЛЕДНИЙ ВЫВОД ТЕРМИНАЛА:\n"
        f"{request.recentTerminalOutput or 'Вывод терминала отсутствует.'}\n\n"

        "Задание: оцени ситуацию и действия пользователя. "
        "Сформулируй одну полезную подсказку, куда двигаться дальше."
    )

    payload = {
        "model": MODEL_NAME,
        "messages": [
            {
                "role": "system",
                "content": system_instruction
            },
            {
                "role": "user",
                "content": user_prompt
            }
        ],
        "stream": False,
        "options": {
            "temperature": 0.3
        }
    }
    try:
        response = requests.post(OLLAMA_URL, json=payload, timeout=90)
        response.raise_for_status()
        response_data = response.json()

        hint_out = (
            response_data
            .get("message", {})
            .get("content", "")
            .strip()
        )

        if not hint_out:
            hint_out = "Посмотри на текущий шаг и историю действий: следующий диагностический фокус должен следовать из последней ошибки."
    
    except Exception as e:
        print(f"[ERROR] Не удалось связаться с Ollama: {str(e)}")
        raise HTTPException(
            status_code=500,
            detail=f"Ошибка на стороне AI-микросервиса: Ollama недоступна ({str(e)})"
        )
    
    return MLHintResponse(
        hint=hint_out,
        confidence=0.75,
        source="ml",
        reasoning="Generated by Ollama model"
    )


if __name__ == "__main__":
    print("==================================================")
    print("AI FastAPI сервер запускается...")
    print(f"Используемая модель Ollama: {MODEL_NAME}")
    print(f"Ollama URL: {OLLAMA_URL}")
    print("==================================================")
    uvicorn.run(app, host="0.0.0.0", port=8000)