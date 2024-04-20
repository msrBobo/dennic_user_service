# workspace (GOPATH) configured at /go
FROM golang:1.18 as builder

WORKDIR /app

ADD . /app

# Copy the local package files to the container's workspace.
COPY . ./

# installing depends and build
RUN make build-linux

FROM alpine:latest

COPY --from=builder app/bin/dennic_user_service ./dennic_user_service

RUN chmod +x ./dennic_user_service

CMD ["./dennic_user_service"]
