FROM golang:1.22-alpine AS golang-builder

RUN addgroup -S ctf && adduser -S ctf -G ctf

WORKDIR /src/app
COPY --chown=ctf:ctf . .

RUN go get
RUN go build

FROM scratch
COPY --from=golang-builder /src/app/ctf /ctf
COPY --from=golang-builder /etc/passwd /etc/passwd
COPY --chown=ctf:ctf static /static
COPY --chown=ctf:ctf templates /templates

USER ctf

EXPOSE 8080

CMD [ "/ctf" ]