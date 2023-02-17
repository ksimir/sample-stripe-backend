FROM golang:1.18.1 as builder

WORKDIR /app

# copy module files first so that they don't need to be downloaded again if no change
COPY go.* ./
RUN go mod download
RUN go mod verify

# copy source files and build the binary
COPY . .
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /app/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .  
EXPOSE 8080
ENTRYPOINT [ "sh", "-c", "./main"]