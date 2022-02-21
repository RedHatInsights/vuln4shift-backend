ARG BUILDIMG=registry.access.redhat.com/ubi8/go-toolset
ARG RUNIMG=registry.access.redhat.com/ubi8-minimal

# ---------------------------------------
# build image
FROM ${BUILDIMG} as buildimg

WORKDIR /vuln4shift
USER 1001

ADD go.mod         /vuln4shift/
ADD database_admin /vuln4shift/database_admin
ADD main.go        /vuln4shift/

RUN go mod download
RUN go build -v main.go

# ---------------------------------------
# runtime image
FROM ${RUNIMG} as runtimeimg

WORKDIR /vuln4shift
USER 1001

COPY --from=buildimg /vuln4shift/main /vuln4shift/
