FROM scratch

ARG BINARY_NAME

COPY build/$BINARY_NAME /service

EXPOSE 80 90
CMD ["/service"]
HEALTHCHECK --interval=2m --timeout=3s CMD curl -f http://localhost:90/healthz || exit 1
