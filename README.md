![image](https://github.com/user-attachments/assets/1efbf6ca-9f43-4468-8e9f-d38c20242588)# Параллельный калькулятор на GO

# Описание
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
![image](https://github.com/user-attachments/assets/fd53ef69-7b50-461b-96ea-48cba43a7bb7)


### Шаг 6: 
```bash
  cd calculator
```

### Шаг 4: Запускаем агента
```bash
  go run .\cmd\agent\main.go
```
