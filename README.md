# vuln4shift-backend
Vulnerabilities detection for OpenShift images.

## Development environment

    docker-compose build         # Build images
    docker-compose up --build    # Build images and run in foreground
    docker-compose up --build -d # Build images and run in background
    docker-compose down          # Stop and delete containers
    docker-compose down -v       # Stop and delete containers + delete persistent volumes

