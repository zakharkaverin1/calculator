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
