FROM r.deso.tech/dockerhub/library/golang:1.17-alpine as builder

WORKDIR /whoami

COPY . .

ENV CGO_ENABLED=0

RUN apk add --no-cache --update ca-certificates make git && make build && sleep 5

# Create a minimal container to run a Golang static binary
FROM r.deso.tech/dockerhub/library/alpine:3

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /whoami/template/   /whoami/template/
COPY --from=builder /whoami/static/     /whoami/static/
COPY --from=builder /whoami/bin/whoami* /whoami/whoami

WORKDIR /whoami

RUN touch readiness

ENTRYPOINT ["./whoami"]

EXPOSE 8080
