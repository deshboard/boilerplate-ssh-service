FROM scratch

ARG BUILD_DIR
ARG BINARY_NAME

COPY $BUILD_DIR/$BINARY_NAME /service

EXPOSE 10000
CMD ["/service", "-debug.addr", ":10000"]
HEALTHCHECK CMD curl -f http://localhost:10000/healthz || exit 1
