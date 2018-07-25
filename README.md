# Tanya

## Docker Container

### Building

This will use a multi-stage build to minimize the image size (see [containerize this golang dockerfiles](https://www.cloudreach.com/blog/containerize-this-golang-dockerfiles/))

```
docker build -t tanya .
```

### Testing

```
docker run -it --rm -v $(pwd)/config:/app/config:ro -v $(pwd)/secrets:/app/secrets:ro tanya
```
