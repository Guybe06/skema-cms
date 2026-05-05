# Skema

Headless CMS open source self-hosted. Connectez votre propre base de données, définissez vos collections et gérez votre contenu via une API REST auto-générée. Vos données restent chez vous.

---

## Fonctionnalités

- **Connexions** — Branchez votre PostgreSQL ou MySQL existant
- **Collections** — Définissez vos structures de données sans écrire de SQL
- **Contenu** — Gérez vos entrées via le dashboard ou l'API
- **API REST auto-générée** — Chaque collection expose automatiquement ses endpoints
- **Clés API** — Contrôle d'accès granulaire par collection et par permission
- **Stockage** — Upload de fichiers en local ou compatible S3 (MinIO, AWS, Cloudflare R2)
- **Serveur MCP** — Intégration native avec les agents IA via le protocole MCP
- **Membres** — Invitez des collaborateurs avec des rôles (owner, admin, member)

---

## Stack

- **Backend** — Go, Gin, PostgreSQL / MySQL, Redis, Swaggo
- **Frontend** — À venir

---

## Démarrage rapide

```bash
git clone https://github.com/your-org/skema-cms.git
cd skema-cms/backend

cp .env.example .env
# Remplissez les variables dans .env

make run
```

L'API sera disponible sur `http://localhost:3000/v1`  
La documentation Swagger sur `http://localhost:3000/docs/v1/index.html`

---

## Configuration

Copiez `.env.example` en `.env` et renseignez les variables suivantes :

| Variable | Description |
|---|---|
| `APP_PORT` | Port du serveur (défaut : 3000) |
| `APP_ENV` | Environnement (`development` / `production`) |
| `CMS_DB_HOST` | Hôte de la base de données Skema |
| `CMS_DB_NAME` | Nom de la base de données Skema |
| `JWT_SECRET` | Clé secrète pour les tokens JWT |
| `ENCRYPTION_KEY` | Clé de chiffrement des credentials clients |
| `REDIS_URL` | URL Redis (optionnel, fallback mémoire si absent) |
| `RESEND_API_KEY` | Clé API Resend pour l'envoi d'emails |
| `STORAGE_DRIVER` | `local` ou `s3` |

---

## Licence

[GNU Affero General Public License v3.0](LICENSE)
