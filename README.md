# Simple JWT Auth Server

Simple auth server to verify JWT token, intended to be used with Traefik using ForwardAuth middleware.

https://doc.traefik.io/traefik/v2.3/middlewares/forwardauth/

### Basic Usage

Docker compose example
```yaml
services:
  traefik:
    # [...]
    
  jwt-server:
    build: .
    ports:
      - "8080:8080"
    environment:
      JWT_SIGNING_KEY: SuperSecretKey
```

Traefik config example (YAML)
```yaml
http:
  routers:
    router1:
      service: myService
      middlewares:
        - jwtAuth
      rule: "Path(`/admin`)"
      
  middlewares:
    jwtAuth:
      forwardAuth:
        address: "http://jwt-server:8080/auth"

  services:
    myService:
      loadBalancer:
        servers:
          - url: "http://127.0.0.1:80"
```

```shell
curl \
  --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhdXRoIiwibmFtZSI6IkpvaG4gRG9lIiwicmFuayI6ImFkbWluIiwiaWF0IjoxNTE2MjM5MDIyfQ.A0UcP9CVAysMVXuh4oC0wZohmLXAjH-aybLUxp6UVJQ" \
  http://localhost/admin
````

Default settings
```yaml
JWT_HEADER: "Authorization"
JWT_AUTH_SCHEME: "Bearer"
JWT_ALGORITHM: "HS256"
SERVER_ADDR: ":8080"
SERVER_PATH: "/auth"
```

### Filters

You can verify claims with basic filters query param following this format:  
`claim1:value1|claim2:value2`

Example:  
`http://jwt-server/auth?filters=sub:auth|rank:admin`