FROM golang:1.16-alpine as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o url-short .

WORKDIR /dist

RUN cp /build/url-short .

FROM scratch

COPY --from=builder /dist/url-short /app/

COPY conf.json /app/

WORKDIR /app

CMD ["./url-short"]