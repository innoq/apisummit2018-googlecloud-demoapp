FROM golang AS build
COPY go.mod go.sum todo.go /todo/
WORKDIR /todo
RUN go build

FROM debian:stretch-slim
COPY --from=build /todo/todo todo
EXPOSE 8080
CMD ["./todo", "-addr=0.0.0.0:8080"]
