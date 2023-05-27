FROM golang:1.20 as build
WORKDIR /src
ADD src /src
RUN cd /tmp \
    && apt-get update \
    && apt-get install unzip \
    && curl -L -o protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v23.2/protoc-23.2-linux-x86_64.zip \
    && unzip protoc.zip \
    && cp bin/protoc /bin/protoc \
    && chmod 755 /bin/protoc \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN cd protos && sh ./genproto.sh
RUN go get -d -v ./... \
    && go build -o /app

FROM gcr.io/distroless/base as final
USER 1000
EXPOSE 8080 53835
ENTRYPOINT [ "/app" ]
COPY --from=build /app /
ADD /src/swagger.yaml /swagger.yaml