# Swagger UI

Этот каталог содержит все необходимое для обслуживания Swagger UI для gRPC-сервиса.

## Настройка

Для правильной работы Swagger UI необходимо выполнить несколько шагов по настройке.

### Требования

*   **Node.js и npm**: Для установки зависимостей Swagger UI.
*   **go-bindata**: Для встраивания файлов Swagger UI в приложение.

### Установка

1.  **Установите `go-bindata`:**
    ```sh
    go install github.com/go-bindata/go-bindata/...
    ```

2.  **Установите зависимости Swagger UI:**
    Перейдите в директорию `swagger-ui` и установите npm зависимости.
    ```sh
    cd swagger-ui
    npm install swagger-ui-dist
    cd ..
    ```

3.  **Сгенерируйте `datafile.go`:**
    Эта команда встроит файлы Swagger UI в ваше приложение. Выполните ее из корневой директории проекта.
    ```sh
    go-bindata -o pkg/swagger/datafile.go -pkg swagger pkg/swagger/swagger-ui/...
    ```

## Использование

После выполнения этих шагов и запуска основного приложения, Swagger UI будет доступен по адресу:

[http://localhost:50052/swagger-ui/](http://localhost:50052/swagger-ui/)

Swagger UI автоматически загрузит спецификацию `order.swagger.json` и отобразит документацию по вашему API.