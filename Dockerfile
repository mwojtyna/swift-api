FROM golang:1.24.2-bookworm AS builder
WORKDIR /app
EXPOSE ${API_PORT}
COPY . .
RUN go mod download
RUN make build-parser && make build-server

FROM builder AS runner

CMD make parse && make serve
