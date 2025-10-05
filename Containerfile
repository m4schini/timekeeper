FROM docker.io/library/golang:1.25 AS build

WORKDIR /go/src/app
COPY . .

RUN go mod tidy
RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian12 AS prod

COPY --from=build /go/bin/app /
CMD ["/app"]