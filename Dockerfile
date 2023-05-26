FROM golang:1.20 as build
WORKDIR /src
ADD src /src
RUN go get -d -v ./... \
    && go build -o /app

FROM gcr.io/distroless/base as final
USER 1000
EXPOSE 8080 53835
ENTRYPOINT [ "/app" ]
COPY --from=build /app /
ADD /src/swagger.yaml /swagger.yaml