FROM golang:1.19-alpine as compiler

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /service


FROM ghcr.io/surnet/alpine-wkhtmltopdf:3.17.0-0.12.6-full
WORKDIR /
COPY --from=compiler /service /

EXPOSE 80

ENTRYPOINT ["/service" ]