FROM golang:latest

WORKDIR /go/src/backend

# COPY go.mod .
# COPY go.sum .

COPY . .
RUN go mod download

# RUN go build -o from-docker main.go // これは本番環境のやつ
CMD [ "go run main.go" ]

# CMD [ "./from-docker" ] // これは本番環境のやつ
