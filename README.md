# keypass

## Запуск сервера 

        docker-compose -f ./docker/docker-compose.yml up -d --build

# Запуск собранного клиента (MacOS/arm)

        Бинарный файл лежит в папке cmd/client/keypass

# Запуск клиента в режиме разработки

        go run cmd/client/main.go

# Сборка клиента

        go build -o cmd/client/keypass cmd/client/main.go

## Сделать файл исполняемым:
        chmod +x cmd/client/keypass

# Запуск сервера в режиме разработки

        go run cmd/server/main.go -k 123 -d 'postgres://postgres:1@localhost:5432/gopherdb?sslmode=disable' -private-key './certs/localhost.key' -public-key './certs/localhost.crt' -r 123
    
