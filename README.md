# NATS Streaming L0

Учебный проект по созданию сервиса для работы с данными заказов.

### Видео-демонстрация

[Ссылка на видео](https://drive.google.com/file/d/1WQOKvknVZ7LhSpSM44n6vzS4-3Fjhq_q/view?usp=drive_link)

### Стэк

- Go
- PostgreSQL
- NATS Streaming (STAN)
- Docker, Docker Compose
- net/http, chi
- go-playground/validator

### Схема Базы Данных

SQL-запрос для создания необходимых таблиц находится в файле `sql.sql` в корневой папке проекта.

### Как запустить проект

Для запуска потребуется: Go, Docker, Docker Compose, NATS Streaming Server.

1)  Склонировать репозиторий:
    ```bash
    git clone https://github.com/DenisDavidov-creator/wb_l0_task.git
    cd wb_l0_task
    ```

2)  Запустить базу данных в Docker:
    ```bash
    docker compose up -d
    ```

3)  Запустить NATS Streaming Server:
    В новом терминале перейти в папку с сервером и выполнить:
    ```bash
    ./nats-streaming-server
    ```

4)  Запустить сервис:
    В новом терминале запустить главный файл сервиса:
    ```bash
    go run main.go
    ```

5)  Опубликовать тестовые данные:
    В еще одном терминале запустить паблишер. Он отправит сообщения из файла `orders.json`.
    ```bash
    go run ./cmd/publisher/main.go
    ```

6)  Проверить результат:
    Открыть файл `index.html` в браузере (например, через Live Server). Вставить `order_uid` одного из отправленных заказов и нажать на кнопку, чтобы получить данные.
