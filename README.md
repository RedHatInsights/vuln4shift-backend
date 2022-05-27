[![Tests](https://github.com/RedHatInsights/vuln4shift-backend/actions/workflows/run_tests.yaml/badge.svg)](https://github.com/RedHatInsights/vuln4shift-backend/actions/workflows/run_tests.yaml)
[![codecov](https://codecov.io/gh/RedHatInsights/vuln4shift-backend/branch/master/graph/badge.svg)](https://codecov.io/gh/RedHatInsights/vuln4shift-backend)

# vuln4shift-backend
Vulnerabilities detection for OpenShift images.

## Development environment

    docker-compose build         # Build images
    docker-compose up --build    # Build images and run in foreground
    docker-compose up --build -d # Build images and run in background
    docker-compose down          # Stop and delete containers
    docker-compose down -v       # Stop and delete containers + delete persistent volumes

    # Sync CVE data
    docker-compose run --rm vuln4shift_vmsync

    # Sync Pyxis data
    docker-compose run --rm vuln4shift_pyxis

    # psql console
    docker-compose exec vuln4shift_database psql -U vuln4shift_admin vuln4shift

Manager swagger documentation is running at
```
http://localhost:8000/api/vuln4shift/v1/openapi.json
http://localhost:8000/api/vuln4shift/v1/openapi/index.html
```

## Unit tests
You can run unit tests localy by using.
```
docker-compose run --rm vuln4shift_unit_tests
```

## Deployment to ephemeral environment

Set following apps in `~/.config/bonfire/config.yaml`:

    apps:
    - name: vuln4shift
      components:
      - name: backend
        host: github
        repo: RedHatInsights/vuln4shift-backend
        ref: master
        path: /deploy/clowdapp.yaml

    - name: vuln4shift-local
      components:
      - name: backend-local
        host: local
        repo: ~/work/vuln4shift-backend
        path: /deploy/clowdapp.yaml

Reserve a namespace:

    bonfire namespace reserve -d 4h

Deploy:

    # Using the ClowdApp template from GitHub
    bonfire deploy vuln4shift --namespace <reserved_namespace>

    # Or using to ClowdApp template from local dir
    bonfire deploy vuln4shift-local --namespace <reserved_namespace>

Note that the image tag must exist in the Quay, can be changed to different image by adding e.g.:

    --set-image-tag quay.io/jdobes/vuln4shift-backend=bbbf78b


