FROM scratch

ARG BINARY_NAME

COPY build/$BINARY_NAME /service

EXPOSE 80 10000 10001
CMD ["/service", "-service", ":80", "-health", ":10000", "-debug", ":10001"]
HEALTHCHECK --interval=2m --timeout=3s CMD curl -f http://localhost:10000/healthz || exit 1
