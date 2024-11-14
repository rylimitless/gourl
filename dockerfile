FROM ubuntu

WORKDIR /app

COPY . .

RUN apt update && \
    apt install -y golang-go ca-certificates sqlite3 && \
    update-ca-certificates

# Add CGO_ENABLED here, before building/running Go
# ENV CGO_ENABLED=1

EXPOSE 9032

VOLUME [ "/app/database" ]

RUN useradd app

RUN mkdir -p /home/app && chown -R app:app /home/app

RUN  chown -R app:app /app/database && \
    chmod 777 /app/database


USER app

CMD ["go" ,"run", "main.go"]