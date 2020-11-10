FROM golang:1.14.0 AS builder

ARG MAKE_TARGET=all

# Install / Update Dependencies:
#   - CA Certificates
#
RUN apt-get update
RUN apt-get install -y ca-certificates

# Copy the Source Tree
#
WORKDIR /workspace
COPY . .

# Build
#
RUN git status -s > /dev/null
RUN make ${MAKE_TARGET}

# --------------------------------------------------------------

FROM scratch

COPY --from=builder /workspace/bin/* /
COPY --from=builder /tmp /tmp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
