FROM --platform=$BUILDPLATFORM golang:1.27 AS builder

WORKDIR /go/src/github.com/micromdm/micromdm/

ARG TARGETARCH
ARG TARGETOS

ENV CGO_ENABLED=0 \
	GOARCH=$TARGETARCH \
	GOOS=$TARGETOS \
	GOTOOLCHAIN=local

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make deps
RUN make CURRENT_PLATFORM=linux


FROM alpine:3.22

RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/micromdm/micromdm/build/linux/micromdm /usr/bin/
COPY --from=builder /go/src/github.com/micromdm/micromdm/build/linux/mdmctl /usr/bin/

EXPOSE 80 443
VOLUME ["/var/db/micromdm"]
CMD ["micromdm", "serve"]
