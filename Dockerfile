FROM golang:1.23 AS build

WORKDIR /build

COPY go.mod go.sum ./     

RUN go mod download      

COPY . .       

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ozontestservice ./cmd/main

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=build /build/ozontestservice .

COPY --from=build /build/internal/database/migrations ./migrations

CMD ["./ozontestservice"]