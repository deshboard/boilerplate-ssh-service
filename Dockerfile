FROM deshboard/go:1.7-onbuild

EXPOSE 80

HEALTHCHECK --interval=2m --timeout=3s CMD curl -f http://localhost:81/healthz || exit 1
