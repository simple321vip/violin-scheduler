# stage 1
FROM golang AS build

WORKDIR /app

COPY . /app

##
RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod tidy

RUN go build -o violin-scheduler -ldflags="-w -s"

# stage 2
FROM alpine AS work

ENV GIN_MODE=release

WORKDIR /app

COPY --from=build /app/violin-scheduler ./violin-scheduler

ENTRYPOINT ["./violin-scheduler"]
