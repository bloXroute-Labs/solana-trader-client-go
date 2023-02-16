FROM golang:1.20.1 as builder

# setup environment
RUN mkdir -p /app/solana-trader-client-go
WORKDIR /app/solana-trader-client-go

# copy files and download dependencies
COPY . .
RUN go mod download

# build
RUN rm -rf bin
RUN go build -o bin/ ./benchmark/traderapi

FROM golang:1.20.1

RUN apt-get update
RUN apt-get install -y net-tools
RUN rm -rf /var/lib/apt/lists/*

# setup user
RUN useradd -ms /bin/bash solana-trader-client-go

# setup environment
RUN mkdir -p /app/solana-trader-client-go
RUN chown -R solana-trader-client-go:solana-trader-client-go /app/solana-trader-client-go

WORKDIR /app/solana-trader-client-go

COPY --from=builder /app/solana-trader-client-go/bin /app/solana-trader-client-go/bin

ENTRYPOINT ["/app/solana-trader-client-go/bin/traderapi"]
