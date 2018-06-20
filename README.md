# Tanya

## Docker Container

### Building

```
docker build -t tanya .
```

### Testing

```
docker run -it --rm -v $(pwd)/config:/app/config:ro -v $(pwd)/secrets:/app/secrets:ro tanya
```
