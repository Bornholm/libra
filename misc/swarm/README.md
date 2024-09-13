# Docker Swarm dans Docker

Le fichier `docker-compose.yml` contenu dans ce répertoire permet de démarrer un cluster [Docker Swarm](https://docs.docker.com/engine/swarm/) dans un environnement [Docker Compose](https://docs.docker.com/compose/), autrement dit sans avoir à démarrer plusieurs machines avec Docker.

## Démarrer le cluster Swarm

1. Se positionner dans le répertoire:

   ```bash
   cd misc/swarm
   ```

2. Exécuter la commande `docker compose up`

Après quelques secondes, le cluster Swarm devrait être initialisé et prêt à l'emploi.

## Manipuler le cluster

Pour manipuler le cluster et exécuter des commandes Docker sur celui-ci, deux options sont possibles.

### Entrer dans le conteneur

Vous pouvez entrer les commandes Docker directement en entrant dans le conteneur `manager`. Pour ce faire, faites:

```bash
docker compose exec manager /bin/sh
```

Puis tapez les commandes Docker souhaitées, par exemple:

```bash
docker node ls
```

### Avec le client de sa machine hôte

Vous pouvez également interroger votre cluster Swarm en utilisant le client Docker de votre machine hôte. Pour ce faire:

```bash
# Récupérez les droits sur le répertoire `certs/client`
sudo chown -R $(whoami): certs/client

# Récupérer la liste des noeuds constituant votre cluster Swarm
docker --tlscacert ./certs/client/ca.pem --tlscert ./certs/client/cert.pem  --tlskey ./certs/client/key.pem --tls -H localhost:22376 node ls
```

Il vous faudra préciser à chaque fois les flags `--tlscacert`, `--tlscert`, `--tlskey`, `--tls` et `-H` afin de communiquer avec le daemon Docker de votre noeud `manager`.

## Nettoyer l'environnement

Pour nettoyer l'environnement, vous pouvez faire:

```bash
docker compose down --remove-orphans -v
```
