ARG BUILDIMG=registry.access.redhat.com/ubi8/go-toolset
ARG RUNIMG=registry.access.redhat.com/ubi8-minimal
# ---------------------------------------
# build image
FROM ${BUILDIMG} as buildimg

WORKDIR /vuln4shift

USER root

ADD go.mod                      /vuln4shift/
ADD go.sum                      /vuln4shift/

RUN go mod download

ADD main.go                     /vuln4shift/
ADD base                        /vuln4shift/base
ADD dbadmin                     /vuln4shift/dbadmin
ADD manager                     /vuln4shift/manager
ADD scripts                     /vuln4shift/scripts
ADD test                        /vuln4shift/test

# install swag command to generate swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN mkdir ./manager/docs
RUN bash ./scripts/generate_swagger.sh

RUN go build -v main.go
# ---------------------------------------
# runtime image
FROM ${RUNIMG} as runtimeimg

WORKDIR /vuln4shift
USER 1001

COPY --from=buildimg /vuln4shift/main                       /vuln4shift/
COPY --from=buildimg /vuln4shift/dbadmin/migrations         /vuln4shift/dbadmin/migrations
COPY --from=buildimg /vuln4shift/manager/docs/swagger.json  /vuln4shift/manager/docs/swagger.json
