PROTO_DIR = ./api
OUT_DIR = ./api

PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

SERVER = main.exe

lint: 
	golangci-lint run -v

# Генерация кода gRPC и grpc-gateway
gen: $(PROTO_FILES)
	protoc -I $(PROTO_DIR) --go_out $(OUT_DIR) --go_opt paths=source_relative \
		--go-grpc_out $(OUT_DIR) --go-grpc_opt paths=source_relative \
		--grpc-gateway_out $(OUT_DIR) --grpc-gateway_opt paths=source_relative $^

# Сборка бинарника
build: gen
	@echo Build Server...
	go build -o $(SERVER) ./cmd/main/main.go

# Запуск сервера
run: build
	@echo Run server...
	.\$(SERVER)

# Очистка
clean:
	@echo Очистка сгенерированных файлов и бинарника...
	del /q $(OUT_DIR)\*.go
	del /q $(SERVER)

.PHONY: gen build run clean