FROM golang:latest

WORKDIR /go/src/backend

# COPY go.mod .
# COPY go.sum .

COPY . .
RUN go mod download

# // これは本番環境のやつ
RUN go build -o from-docker main.go
# CMD [ "go run main.go" ]

# // これは本番環境のやつ
CMD [ "./from-docker" ]
