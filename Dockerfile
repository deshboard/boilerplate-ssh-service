FROM deshboard/go:1.7-onbuild

EXPOSE 80

HEALTHCHECK --interval=2m --timeout=3s CMD curl -f http://localhost/_status/healthz || exit 1
