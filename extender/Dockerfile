# stage 1
FROM golang:1.19.5 AS build

WORKDIR /app

COPY . /app

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o violin-scheduler

# stage 2
FROM alpine AS work

ENV GIN_MODE=release

WORKDIR /app

COPY --from=build /app/violin-scheduler ./violin-scheduler

RUN chmod 777  ./violin-scheduler

ENTRYPOINT ["./violin-scheduler"]
