# FROM golang:1.21

# WORKDIR /app

# COPY go.mod ./
# RUN go mod download

# COPY *.go ./

# # RUN go build -o /server
# RUN go build -v -o /

# EXPOSE 8080

# CMD [ "/go-docker-demo" ]
FROM golang:1.21-alpine

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./server
EXPOSE 8080

CMD ["app"]