FROM golang:1.23.0-alpine3.20 AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
RUN echo $(pwd)
RUN echo $(ls -a)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o device-management main.go
RUN echo $(ls -a)
RUN echo $(ls -a device-management)
RUN echo $(ls -a config/*json)

FROM scratch
COPY --from=build /app/device-management /device-management
COPY --from=build /app/config/*.json /config/
EXPOSE 8080
ENTRYPOINT ["/device-management"]