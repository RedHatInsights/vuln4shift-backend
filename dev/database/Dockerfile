FROM quay.io/cloudservices/postgresql-rds:16-4649c84

# install pg_repack
USER root
RUN curl -o /etc/yum.repos.d/postgresql.repo \
        https://copr.fedorainfracloud.org/coprs/g/insights/postgresql-16/repo/epel-8/group_insights-postgresql-16-epel-8.repo
RUN dnf install -y pg_repack

USER postgres
