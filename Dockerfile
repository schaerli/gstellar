# run it with
# docker run -v /path/to/some/gstellar.json:/app/gstellar.json -p 5678:5678 --network host -it gstellar
FROM golang:1.17-alpine

WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gstellar

EXPOSE 5678

CMD ["/app/docker-entrypoint.sh"]
