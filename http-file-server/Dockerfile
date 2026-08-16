FROM --platform=$BUILDPLATFORM golang:1-alpine AS build

WORKDIR /src
COPY go.mod main.go ./

ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH="$TARGETARCH" \
      go build -trimpath -ldflags='-s -w' -o /http-file-server .

FROM scratch

COPY --from=build /http-file-server /http-file-server

USER 1000:1000
EXPOSE 8000

ENTRYPOINT ["/http-file-server"]
