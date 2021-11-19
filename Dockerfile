FROM golang:1.17.0-buster as build

ENV CGO_ENABLED=0

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN go build -o ./checker ./cmd/checker


FROM scratch
COPY --from=build /app/checker /
ENTRYPOINT ["/checker"]
