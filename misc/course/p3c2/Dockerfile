FROM alpine:3.20

# On rend paramétrable la version de l'application Libra à télécharger
# et l'empreinte de vérification d'intégrité de l'archive souhaitée
ARG LIBRA_VERSION=0.0.3
ARG LIBRA_CHECKSUM=0d5a496622a3e6167c43d5a8edfc6d679d5e1603ce78ae48a67daaa0f5aaab5b
ARG LIBRA_REPOSITORY=https://github.com/Bornholm/libra

# On met à jour les dépendances système afin de limiter les risques 
# de présence de failles de sécurité dans l’image finale et
# on installe les dépendance nécessaires à la récupération de l’application Libra
# en spécifiant la version souhaitée
RUN apk upgrade \
    && apk add \
        wget=1.24.5-r0 \
        tar=1.35-r2

# On télécharge l'archive contenant l'exécutable Libra,
# on vérifie l'empreinte de l'archive,
# on créait le répertoire de destination du binaire,
# on extrait celui-ci,
# puis on s'assure qu'il a bien les droits d'exécution.
RUN wget -O libra.tar.gz \
        ${LIBRA_REPOSITORY}/releases/download/v${LIBRA_VERSION}/libra_${LIBRA_VERSION}_linux_amd64.tar.gz \
    && echo "${LIBRA_CHECKSUM} libra.tar.gz" | sha256sum -c \
    && mkdir -p /usr/local/bin \
    && tar -xzf libra.tar.gz -C /usr/local/bin libra \
    && chmod +x /usr/local/bin/libra

# On indique que l’application écoute sur le port 8080
# et que ce port devrait être exposé
EXPOSE 8080

# On indique que le conteneur souhaite utiliser un volume
# pour le répertoire /uploads
VOLUME /uploads

# On créait un utilisateur non privilégié pour exécuter l'application
RUN adduser -H -S -s /dev/null libra

# On bascule sur cet utilisateur
USER libra

# On configure l’exécutable comme processus à lancer au démarrage
# du conteneur
CMD ["/usr/local/bin/libra"]
