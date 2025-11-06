FROM golang:latest  

#Установить директорию в которую мы положим все файлы нашего проекта
WORKDIR /develop

#Скопировать все зависимости
COPY go.mod go.sum ./
RUN go mod download

#Скопируем в WORKDIR все файлы нашего проекта
COPY . .
#Порт, на который будем слать запросы
EXPOSE 50052
#Команда для создания .exe файла. Путь к нему в контейнере: ./main. И путь к нему явно: ./cmd/main
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/main
#Запуск .exe файла 
CMD ["/develop/main"]

#Чтобы запустить контейнер с образом, нужно прописать: docker run -d --rm -p 9000:50052 sg/grpc-go





