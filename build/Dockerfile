# FROM alpine:latest
FROM registry.access.redhat.com/ubi8-minimal:latest

# set labels for metadata
LABEL maintainer="Franck Grosjean <fgrosjean@redhat.com>" \
  name="lbp-server" \
  description="A Go HTTP server" \
  summary="A Go HTTP server"

# set environment variables
ENV LBP_SERVER=/usr/local/bin/lbp-server \
  USER_UID=1001 \
  USER_NAME=lbp

# install binary
COPY build/_output/linux/bin/lbp-server ${LBP_SERVER}

# copy licenses
RUN mkdir /licenses
COPY LICENSE /licenses

# copy certs
RUN mkdir /etc/certs
COPY build/certs/lbp-*.pem build/certs/ca.crt /etc/certs/
RUN chmod 444 /etc/certs/*

# set entrypoint and port
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/lbp-server", "-tlsCertFile=/etc/certs/lbp-crt.pem", "-tlsKeyFile=/etc/certs/lbp-key.pem"]

# switch to non-root user
USER ${USER_UID}
