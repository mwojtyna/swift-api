FROM golang:1.24.2-bookworm 

WORKDIR /app
EXPOSE ${API_PORT}

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest 

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build-parser && make build-server

CMD [ "./run.sh" ]
