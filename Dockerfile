# ===== build stage ====
FROM golang:1.19.13-bullseye as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -lgflags="-w -s" -o app

# ===== deploy stage ====
FROM golang:1.19.13-alpine as deploy

RUN apk update

COPY --from=build /app/app ./app

CMD ["./app"]

# ===== dev ====
FROM golang:1.19.13-bullseye as dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest
CMD ["air"]
