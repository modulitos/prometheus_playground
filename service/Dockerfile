FROM golang:1.19-alpine
RUN apk --no-cache add git make bash

WORKDIR /go/src/app
COPY . /go/src/app

# Because we are running distroless we dont want to depend on C libraries(CGO_ENABLED=0)
# Compile version into binary for logging, errors - "-X ${GO_PROJECT}/pkg/version.version=${APP_VERSION}"
# RUN CGO_ENABLED=0 go build --ldflags="-X ${GO_PROJECT}/pkg/version.version=${APP_VERSION}" -o /go/bin/ cmd/app.go

RUN go build -o /go/bin/ cmd/app.go
USER 1001

ARG BUILD_WITH_PORT="1234"
ARG BUILD_WITH_DEBUG="true"

ENV PORT=$BUILD_WITH_PORT
ENV DEBUG=$BUILD_WITH_DEBUG

CMD ["/go/bin/app"]


# For distroless:
# Now copy it into our base image.
# FROM gcr.io/distroless/static:nonroot
# COPY --from=build /go/bin/app /
# EXPOSE 8080
# USER 1001

# CMD ["/app"]