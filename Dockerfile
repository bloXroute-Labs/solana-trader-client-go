FROM golang:1.18 as builder

# setup environment
RUN mkdir -p /app/serum-client-go
WORKDIR /app/serum-client-go

# copy files and download dependencies
COPY . .
RUN go mod download

# build
RUN rm -rf bin
RUN go build -o bin/ ./benchmark/serumapi

FROM golang:1.18

RUN apt-get update
RUN apt-get install -y net-tools
RUN rm -rf /var/lib/apt/lists/*

# setup user
RUN useradd -ms /bin/bash serum-client-go

# setup environment
RUN mkdir -p /app/serum-client-go
RUN chown -R serum-client-go:serum-client-go /app/serum-client-go

WORKDIR /app/serum-client-go

COPY --from=builder /app/serum-client-go/bin /app/serum-client-go/bin

ENTRYPOINT ["/app/serum-client-go/bin/serumapi"]
