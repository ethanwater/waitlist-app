FROM golang:1.16-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/waitlist

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/waitlist .

CMD ["./waitlist"]
