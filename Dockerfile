# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-vgs-api-client"
LABEL REPO="https://github.com/zdmytriv/vgs-api-client"

ENV PROJPATH=/go/src/github.com/zdmytriv/vgs-api-client

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/zdmytriv/vgs-api-client
WORKDIR /go/src/github.com/zdmytriv/vgs-api-client

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/zdmytriv/vgs-api-client"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/vgs-api-client/bin

WORKDIR /opt/vgs-api-client/bin

COPY --from=build-stage /go/src/github.com/zdmytriv/vgs-api-client/bin/vgs-api-client /opt/vgs-api-client/bin/
RUN chmod +x /opt/vgs-api-client/bin/vgs-api-client

# Create appuser
RUN adduser -D -g '' vgs-api-client
USER vgs-api-client

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/vgs-api-client/bin/vgs-api-client"]
