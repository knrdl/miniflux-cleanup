FROM golang:1.23.2-alpine3.19 as builder

WORKDIR /go/src/app

COPY src .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /miniflux-cleanup





FROM scratch

EXPOSE 8080/tcp

VOLUME /data

COPY www /www

ADD https://unpkg.com/spectre.css/dist/spectre-icons.min.css /www
ADD https://cdn.jsdelivr.net/npm/@json-editor/json-editor@latest/dist/jsoneditor.min.js /www
ADD https://maxcdn.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css /www

COPY --from=builder /miniflux-cleanup /miniflux-cleanup

CMD ["/miniflux-cleanup"]
