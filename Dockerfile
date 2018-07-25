FROM iron/go:dev as builder
WORKDIR /app
ENV SRC_DIR=/go/src/github.com/matthew-parlette/tanya/
ADD . $SRC_DIR
RUN cd $SRC_DIR;go build -o tanya;cp tanya /app/tanya
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o tanya .

FROM iron/go
WORKDIR /app
COPY --from=builder /app/tanya /app/
RUN chmod +x /app/tanya
CMD ["./tanya"]
