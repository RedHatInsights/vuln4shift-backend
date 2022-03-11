# vuln4shift-backend
Vulnerabilities detection for OpenShift images.

## Development environment

    docker-compose build         # Build images
    docker-compose up --build    # Build images and run in foreground
    docker-compose up --build -d # Build images and run in background
    docker-compose down          # Stop and delete containers
    docker-compose down -v       # Stop and delete containers + delete persistent volumes

Manager swagger documentation is running at
```
http://localhost:8000/api/vuln4shift/v1/openapi.json
http://localhost:8000/openapi/index.html
```

## Unit tests
You can run unit tests localy by using.
```
docker-compose run --rm vuln4shift_unit_tests
```
