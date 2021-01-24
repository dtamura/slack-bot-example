FROM golang:1.15 as build

WORKDIR /build
COPY . .
RUN make

FROM alpine:3.13
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /build/bin/test-app app

CMD [ "./app" ]