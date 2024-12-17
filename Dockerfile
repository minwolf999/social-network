# Utiliser l'image officielle de Go
FROM golang:1.23.1

# Définir le répertoire de travail
WORKDIR /social-network

# Installer SQLite3
RUN apt-get update && \
    apt-get install -y sqlite3 libsqlite3-dev && \
    apt-get install -y zsh && \
    sh -c "$(wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh -O -)" && \
    go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    ln -s /go/bin/linux_amd64/migrate /usr/local/bin/migrate

# Copier les fichiers Go et autres fichiers nécessaires
COPY . .

RUN go mod tidy && \
    make ServBuild

# Exposer le port
EXPOSE 8080

# Laisser le container ouvert une fois lancer
CMD ["make", "RunBuild"]