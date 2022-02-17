FROM registry.access.redhat.com/ubi8/ubi-minimal

WORKDIR /vuln4shift
USER 1001

ADD entrypoint.sh /vuln4shift/
