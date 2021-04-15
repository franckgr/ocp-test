# FROM alpine:latest
FROM ubi8-minimal:latest

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

# set entrypoint
ENTRYPOINT ["/usr/local/bin/lbp-server"]

# switch to non-root user
USER ${USER_UID}
