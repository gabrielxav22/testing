# Build stage
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Copia os arquivos de módulo.
COPY go.mod go.sum ./

# Baixa as dependências. Esta camada será armazenada em cache
# e só será reconstruída se go.mod ou go.sum mudarem.
RUN go mod download

# (Opcional, mas boa prática) Verifica a integridade das dependências.
RUN go mod verify

# Copia o restante do código-fonte.
# Esta camada será invalidada com mais frequência.
COPY . .

# Build the binary.
# O -mod=vendor é removido pois agora usamos os módulos baixados.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Final stage
# Usamos a tag :nonroot para uma imagem que já roda sem privilégios.
FROM gcr.io/distroless/static-debian12:nonroot

# Definimos um diretório de trabalho limpo e sem privilégios.
WORKDIR /app

# Copy the pre-built binary from the build stage
# O binário é copiado para o /app definido acima.
COPY --from=builder /app/main .

# Command to run the executable
# O comando é executado a partir do WORKDIR /app.
CMD ["./main"] 