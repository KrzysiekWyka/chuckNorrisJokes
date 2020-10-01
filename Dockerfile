FROM golang:1.15.2-alpine AS build

WORKDIR /app
COPY . /app/
RUN CGO_ENABLED=0 go build -o bin

FROM scratch

WORKDIR /app

COPY --from=build /app/bin /app/bin
COPY --from=build /app/views /app/views

ARG PORT=3000
ENV PORT $PORT

ARG GIN_MODE=release
ENV GIN_MODE $GIN_MODE

ENTRYPOINT ["/app/bin"]