# Usando uma imagem base do Golang
FROM golang:1.22-alpine

# Definindo o diretório de trabalho dentro do container
WORKDIR /app

# Copiando os arquivos do projeto
COPY . .

# Compilando a aplicação
RUN go build -o stress-test

# Definindo o comando de execução
ENTRYPOINT ["./stress-test"]