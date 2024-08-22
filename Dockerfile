# Etapa de build
FROM golang:1.22 AS build

# Definir diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiar os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# Baixar as dependências
RUN go mod download

# Copiar o código-fonte para o contêiner
COPY . .

# Compilar o binário
RUN go build -o main .

# Etapa final: cria uma imagem menor apenas com o binário
FROM gcr.io/distroless/base-debian10

# Copiar o binário compilado para a nova imagem
COPY --from=build /app/main /main

# Informar qual porta a aplicação vai usar
EXPOSE 8080

# Comando para rodar a aplicação
ENTRYPOINT ["/main"]