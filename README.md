### Проект "Магазин для сотрудников" (тестовое задание для компании "Avito")

## Как запустить проект:

1. **Склонируйте репозиторий и перейдите в папку проекта:**
   ```sh
   git clone https://github.com/yourusername/yourproject.git
   cd yourproject
   ```

2. **Запустите Docker Compose:**
   ```sh
   docker-compose up -d
   ```

3. **Запустите миграции:**
   ```sh
   go run migrations/auto.go
   ```

4. **Запустите приложение:**
   ```sh
   go run cmd/main.go
   ```
   
5. **API Документация и конфиг для Postman:**
   ```
   в файле api/schema.json
   ```

 Проект должен теперь работать и быть доступным по адресу `http://localhost:8080`.