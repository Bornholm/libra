# Libra - Partagez vos documents en toute sécurité !

> ⚠ **Ceci est une application factice**
>
> Elle a été créée dans l'objectif d'accompagner le cours sur Docker proposé par OpenClassrooms. Aucun audit de sécurité n'a été réalisé pour valider le fonctionnement de celle-ci et à ce titre elle ne devrait en aucun cas être utilisée en production.

## Démarrer avec les sources

### Prérequis

- [Go >= 1.23](https://go.dev/)
- [Make](https://www.gnu.org/software/make/)

### Commandes

#### Compiler l'application

```bash
make build
```

Le binaire généré sera disponible dans le répertoire `./bin`. Exécuter `./bin/libra -h` pour voir les options de configuration disponibles.

#### Générer les différentes images Docker

```bash
make docker-images
```

Les images générées porteront le nom `libra/<variante>-latest`.

## License

[AGPL-3.0](./LICENCE)
