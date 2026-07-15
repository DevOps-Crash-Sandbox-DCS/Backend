from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List
import requests
import uvicorn

app = FastAPI(
    title="DevOps Crash Sandbox - AI Hint Service",
    description="HTTP API для генерации ML-подсказок студентам"
)

OLLAMA_URL = "http://ollama:11434/api/generate"
MODEL_NAME = "qwen2.5-coder:7b"


class HintRequest(BaseModel):
    crash_name: str
    crash_description: str
    command_history: List[str]
    system_logs: str


class HintResponse(BaseModel):
    hint_text: str


@app.post("/api/v1/hint", response_model=HintResponse)
def get_help_hint(request: HintRequest):
    print(f"[AI Сервер] Получен HTTP-запрос для сценария: {request.crash_name}")
    if request.command_history:
        commands_flat = "\n".join([f"- {cmd}" for cmd in request.command_history])
    else:
        commands_flat = "Студент еще не вводил команды в терминале."

    system_instruction = (
        "Вы — опытный DevOps-наставник. Ваша цель — помочь студенту локализовать и исправить аварию в Linux/Docker среде.\n"
        "ПРАВИЛА ИГРЫ:\n"
        "1. КАТЕГОРИЧЕСКИ, АБСОЛЮТНО ЗАПРЕЩЕНО писать готовые консольные команды, утилиты или bash-код для решения или диагностики (например, вместо 'выполни docker stats...', скажи 'посмотри на динамику потребления ресурсов контейнерами').\n"
        "2. ЗАПРЕЩЕНО использовать форматирование кода (обратные кавычки ` `) для написания названий команд.\n"
        "3. Всегда анализируй утилиты, которые использует пользователь в истории команд (например, если он использует Git, оцени, правильный ли файл он пытается откатить или проверить через diff)\n"
        "4. Дай краткую, наталкивающую на правильную мысль подсказку (строго 1-3 предложения), анализируя, что студент УЖЕ сделал правильно.\n"
        "5. Проанализируй историю команд студента: если он делает что-то не то или смотрит не те файлы, аккуратно укажи на это.\n"
        "6. Отвечай строго на русском языке, кратко и профессионально.\n"
        "7. СТРОГО СЛЕДИ ЗА ГРАММАТИКОЙ: правильно склоняй русские слова. Никогда не путай термины: в конфигурациях Linux/Nginx используй 'точка с запятой' для символа ;, а не просто 'запятая'\n"
        "8. Обращайся к пользователю напрямую (на 'ты' или 'вы'), не говори о нём в третьем лице ('студент')"
    )

    user_prompt = (
        f"КОНТЕКСТ ТЕКУЩЕГО ИНЦИДЕНТА: {request.crash_description}\n\n"
        f"ИСТОРИЯ КОМАНД СТУДЕНТА В ТЕРМИНАЛЕ:\n{commands_flat}\n\n"
        f"ЛОГИ ИЗ АВАРИЙНОГО КОНТЕЙНЕРА:\n{request.system_logs or 'Логи пусты или отсутствуют.'}\n\n"
        f"Задание: Оцени ситуацию и действия студента. Сформулируй одну точечную подсказку, куда ему двигаться дальше."
    )

    payload = {
        "model": MODEL_NAME,
        "prompt": user_prompt,
        "system": system_instruction,
        "stream": False
    }

    try:
        response = requests.post(OLLAMA_URL, json=payload, timeout=90)
        response.raise_for_status()
        response_data = response.json()
        hint_out = response_data.get("response", "Не удалось получить текст подсказки от нейросети.")
    except Exception as e:
        print(f"[ERROR] Не удалось связаться с Ollama: {str(e)}")
        raise HTTPException(
            status_code=500,
            detail=f"Ошибка на стороне AI-микросервиса: Оллама недоступна ({str(e)})"
        )

    return HintResponse(hint_text=hint_out)


if __name__ == '__main__':
    print("==================================================")
    print("AI FastAPI Сервер успешно запускается...")
    print(f"Используемая модель Ollama: {MODEL_NAME}")
    print("==================================================")
    uvicorn.run(app, host="0.0.0.0", port=8000)
