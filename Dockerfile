FROM golang:1.19-alpine
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

ENV HOST localhost
ENV PASSWORD password
ENV DATABASE ghost
ENV PORT 3306
ENV CONTENT /var/www/ghost/content
ENV OUTPUT backup.tar.gz

RUN go build -o /ghostbackupper

CMD [ "/ghostbackupper --host $HOST --password $PASSWORD --database $DATABASE --port $PORT --content $CONTENT --output $OUTPUT" ]