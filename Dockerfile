FROM golang:1.22.1-alpine3.19 AS builder

WORKDIR /app/auth

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY . .

RUN go build -o ./bin/sso ./cmd/sso/main.go

FROM alpine:3.19 AS runner

WORKDIR /etc/app/auth

COPY --from=builder /app/auth/bin/sso ./
COPY --from=builder /app/auth/storage ./storage
COPY --from=builder /app/auth/configs ./configs

RUN adduser -DH ssousr && chown -R ssousr: /etc/app/auth && chmod -R 700 /etc/app/auth

USER ssousr
 
CMD [ "./sso" ]
