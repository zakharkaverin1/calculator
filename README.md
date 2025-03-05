# Параллельный калькулятор на GO

## Краткое описание писание
Данный проект представляет из себя систему, вычисляющую сложные арифметичсекие выражения. Состоит из оркестратора и агентов, выполняющих роль вычислителей.

---

## Установка и запуск

### Шаг 1: Клонировать репозиторий
Вводим в консоль данную команду
```bash
git clone https://github.com/zakharkaverin1/calculator
```

### Шаг 2
```bash
cd calculator
```

### Шаг 3: Установка зависимостей 
```bash
go mod download
```

### Шаг 4: Запускаем оркестратор
```bash
go run .\cmd\orchestrator\main.go
```

### Шаг 5: Открываем вторую консоль
![image](https://github.com/user-attachments/assets/e54daca0-b395-4f3c-ae91-5da4ee645ecf)

### Шаг 6: 
```bash
cd calculator
```

### Шаг 7: Запускаем агента
```bash
go run .\cmd\agent\main.go
```

---

## Архитектура проекта

![dgrm](https://github.com/user-attachments/assets/75c2c4ff-ffaf-4214-b283-2c5ec9a5d5b5)

### Оркестратор
  - принимает выражения
  - разбивает выражения на подзадачи
  - управляет задачами
  - собирает результаты
  - хранит результаты вычислений и их статусы
### Агенты
  - берут задачи с помощью http-запросов
  - вычисляют
  - отправляют результаты на сервер

### Возможности 
  + вычисление сложных арифметических выражений с использованием сложения, вычитания, умножения и деления
  + поддержка скобок
  + параллельное вычисление некоторых подзадач

---

# API Эндпоинты
Чтобы текст красиво отображался на GitHub, используем Markdown-разметку. Вот обновлённая версия раздела `README` с улучшенным форматированием:

---

## API Endpoints

### 1. Создание нового выражения для вычисления

**Запрос:**
- **Метод:** `POST`
- **URL:** `/api/v1/calculate`
- **Тело запроса:** JSON с полем `expression`, содержащим математическое выражение.

**Пример запроса с `curl`:**
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression": "2 + 2 * 2"}'
```

**Ответ:**
- **Статус:** `201 Created`
- **Тело ответа:** JSON с полем `id`, содержащим идентификатор созданного выражения.

```json
{
  "id": 1
}
```

---

### 2. Получение списка всех выражений

**Запрос:**
- **Метод:** `GET`
- **URL:** `/api/v1/expressions`

**Пример запроса с `curl`:**
```bash
curl -X GET http://localhost:8080/api/v1/expressions
```

**Ответ:**
- **Статус:** `200 OK`
- **Тело ответа:** JSON с массивом выражений, каждое из которых содержит `id`, `expression`, `status`, и `result`.

```json
{
  "expressions": [
    {
      "id": 1,
      "expression": "2 + 2 * 2",
      "status": "completed",
      "result": 6
    },
    {
      "id": 2,
      "expression": "3 * (4 - 2)",
      "status": "cooking",
      "result": 0
    }
  ]
}
```

---

### 3. Получение выражения по ID

**Запрос:**
- **Метод:** `GET`
- **URL:** `/api/v1/expressions/{id}`

**Пример запроса с `curl`:**
```bash
curl -X GET http://localhost:8080/api/v1/expressions/1
```

**Ответ:**
- **Статус:** `200 OK`
- **Тело ответа:** JSON с информацией о выражении.

```json
{
  "expression": {
    "id": 1,
    "expression": "2 + 2 * 2",
    "status": "completed",
    "result": 6
  }
}
```

---

### 4. Получение задачи для вычисления (внутренний эндпоинт)

**Запрос:**
- **Метод:** `GET`
- **URL:** `/internal/task`

**Пример запроса с `curl`:**
```bash
curl -X GET http://localhost:8080/internal/task
```

**Ответ:**
- **Статус:** `200 OK`
- **Тело ответа:** JSON с информацией о задаче.

```json
{
  "task": {
    "id": "1",
    "arg1": 2,
    "arg2": 2,
    "operation": "*",
    "operation_time": 1000
  }
}
```

---

### 5. Отправка результата вычисления задачи (внутренний эндпоинт)

**Запрос:**
- **Метод:** `POST`
- **URL:** `/internal/task`
- **Тело запроса:** JSON с полями `id` (идентификатор задачи) и `res` (результат вычисления).

**Пример запроса с `curl`:**
```bash
curl -X POST http://localhost:8080/internal/task \
-H "Content-Type: application/json" \
-d '{"id": 1, "res": 4}'
```

**Ответ:**
- **Статус:** `200 OK`
- **Тело ответа:** JSON с подтверждением принятия результата.

```json
{
  "status": "резльтат принят"
}
```

---

### 6. Пример работы с API

1. **Создание выражения:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/calculate \
   -H "Content-Type: application/json" \
   -d '{"expression": "3 * (4 - 2)"}'
   ```

   **Ответ:**
   ```json
   {
     "id": 2
   }
   ```

2. **Получение списка выражений:**
   ```bash
   curl -X GET http://localhost:8080/api/v1/expressions
   ```

   **Ответ:**
   ```json
   {
     "expressions": [
       {
         "id": 1,
         "expression": "2 + 2 * 2",
         "status": "completed",
         "result": 6
       },
       {
         "id": 2,
         "expression": "3 * (4 - 2)",
         "status": "cooking",
         "result": 0
       }
     ]
   }
   ```

3. **Получение выражения по ID:**
   ```bash
   curl -X GET http://localhost:8080/api/v1/expressions/2
   ```

   **Ответ:**
   ```json
   {
     "expression": {
       "id": 2,
       "expression": "3 * (4 - 2)",
       "status": "cooking",
       "result": 0
     }
   }
   ```

---
