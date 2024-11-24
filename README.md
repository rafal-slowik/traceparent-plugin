# traceparent-plugin

Traefik traceparent-plugin plugin

test with docker-compose setup:

```
version: "3.8"

services:
  traefik:
    image: traefik:v3.2
    container_name: traefik
    command:
      - --api.insecure=true
      - "--log.level=DEBUG"
      - --providers.docker=true
      - --providers.docker.exposedByDefault=false
      - --entrypoints.web.address=:80
      - --experimental.plugins.traceparent-plugin.moduleName=github.com/rafal-slowik/traceparent-plugin
      - --experimental.plugins.traceparent-plugin.version=v0.0.2
    ports:
      - "80:80"
      - "8080:8080" # Traefik dashboard
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock

  test-service:
    image: traefik/whoami
    container_name: whoami
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.test-service.rule=Host(`localhost`)"
      - "traefik.http.routers.test-service.entrypoints=web"
      - "traefik.http.middlewares.test-middleware.plugin.traceparent-plugin.HeaderName=X-Appgw-Trace-Id"
      - "traefik.http.routers.test-service.middlewares=test-middleware"
      - "traefik.http.services.test-service.loadbalancer.server.port=80"
```