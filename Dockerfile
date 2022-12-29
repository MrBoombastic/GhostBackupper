FROM golang:1.19-alpine as builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o ghostbackupper cmd/main.go

FROM alpine:latest as backup

COPY --from=builder /app/ghostbackupper .

CMD [ "sh", "-c", "./ghostbackupper backup --db_host $DB_HOST --db_password $DB_PASSWORD --db_user $DB_USER --db_database $DB_DATABASE --db_port $DB_PORT --content $CONTENT --output $OUTPUT --mega_login $MEGA_LOGIN --mega_password $MEGA_PASSWORD" ]
