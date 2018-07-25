# Tanya

## Docker Container

### Building

```
docker run --rm -v "$PWD":/go/src/github.com/matthew-parlette/tanya -w /go/src/github.com/matthew-parlette/tanya iron/go:dev go build -o tanya
docker build -t tanya .
```

### Testing

```
docker run -it --rm -v $(pwd)/config:/app/config:ro -v $(pwd)/secrets:/app/secrets:ro tanya
```
