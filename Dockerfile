FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build cmd/youtube-fetch-search/main.go

EXPOSE 8080

CMD [ "./main" ]