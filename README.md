# dockerapp

sample go lang 3 layers app run on top of docker

### Prerequisite
- Docker version 18.03.1-ce or later
- Go 1.10.1 or later
- Put Google Maps API Key in `files/config/main.ini:Key`

### Build & Installation
```sh
$ go get -u
$ make build-image
$ make build-docker
$ make run
```

### Docker stack
- Ubuntu 18.04 + ca-certificates
- Postgres 10.4
- Nginx 1.13
