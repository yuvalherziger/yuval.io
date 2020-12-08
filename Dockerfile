ARG GO_VERSION=1.15.3
ARG ALPINE_VERSION=3.12

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as go-base
RUN apk add --no-cache git build-base

FROM node:12-alpine as node-base

COPY ui/ /ui/
WORKDIR /ui
RUN yarn install --non-interactive

FROM node-base as ui-build
RUN yarn run build:prod


# Build bio CLI (1)
# -----------------
FROM go-base AS cli-source
ARG CLI_PROJECT_PATH=/usr/src/bio

WORKDIR ${CLI_PROJECT_PATH}/

COPY bio/go.mod bio/go.sum ${CLI_PROJECT_PATH}/
RUN go mod download

COPY bio/cmd ${CLI_PROJECT_PATH}/cmd

# Build bio CLI (2)
# -----------------
FROM cli-source AS cli-build
RUN go build -o /usr/bin/bio ./cmd/bio

# Build REST API (1)
# ------------------
FROM go-base AS api-source
ARG API_PROJECT_PATH=/usr/src/api

WORKDIR ${API_PROJECT_PATH}/

COPY api/cmd ${API_PROJECT_PATH}/cmd

# Build REST API (2)
# ------------------
FROM api-source AS api-build
RUN go build -o /usr/bin/api ./cmd/api

# Final image:
FROM alpine:${ALPINE_VERSION} AS release

RUN adduser -D -u 10000 bio

COPY --from=cli-build /usr/bin/bio /usr/bin/
COPY --from=api-build /usr/bin/api /usr/bin/
COPY --from=ui-build /ui/dist/bio /opt/bio/app/static/

ARG CONFIG_PATH=config/
ARG CONFIG_FILE=config.yaml
COPY ${CONFIG_PATH}/${CONFIG_FILE} /etc/config/${CONFIG_FILE}


USER 10000
WORKDIR /var/lib/api
EXPOSE 8090

CMD ["/usr/bin/api"]
