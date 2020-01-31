Banner rotation service

![test](https://github.com/rendau/barot/workflows/test/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/rendau/barot)](https://goreportcard.com/report/github.com/rendau/barot)

migration command example:
    docker run --rm -it -v "$(pwd)/migrations:/migrations:ro" --network network_name migrate/migrate:latest \
        -path /migrations \
        -database postgres://user:psw@host:5432/dbname?sslmode=disable \
        up/down