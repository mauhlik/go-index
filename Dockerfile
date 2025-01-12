FROM golang:1.23.4 AS build

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/app ./cmd/go-index

FROM golang:1.23.4

RUN useradd -m app-user

USER app-user

WORKDIR /usr/local/bin/

COPY --from=build /usr/local/bin/app /usr/local/bin/app

ENTRYPOINT ["/usr/local/bin/app"]