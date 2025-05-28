FROM golang:1.24-alpine AS golang-builder

RUN addgroup -S ctf && adduser -S ctf -G ctf

WORKDIR /src/app
COPY --chown=ctf:ctf . .

RUN go get
RUN go build

RUN apk add wget unzip
RUN wget https://github.com/twbs/bootstrap/releases/download/v5.3.6/bootstrap-5.3.6-dist.zip
RUN unzip bootstrap-*.zip
RUN mv bootstrap-*-dist static

FROM scratch
COPY --from=golang-builder /src/app/ctf /ctf
COPY --from=golang-builder /etc/passwd /etc/passwd
COPY --from=golang-builder /src/app/static /static
COPY --chown=ctf:ctf templates /templates

USER ctf

EXPOSE 8080

CMD [ "/ctf" ]
