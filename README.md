# Skema

Le CMS headless qui se connecte à votre base de données existante. Vos données restent chez vous.

Skema transforme votre base de données en un espace de gestion simple et intuitif. Connectez votre espace de stockage existant et obtenez instantanément un tableau de bord moderne pour gérer tout votre contenu. Simple, rapide, et sans jamais perdre le contrôle.

---

## Pourquoi Skema ?

- **Vos données restent chez vous** - Skema se connecte à votre base de données existante. Vos informations ne bougent pas, ne sont pas copiées. Si vous arrêtez Skema demain, tout est toujours là, exactement comme avant.
- **Gain de temps** - Plus besoin de recréer un espace de gestion à chaque projet. Skema vous fait gagner 2-3 jours de développement par projet.
- **Sécurité maximale** - Vos accès sont protégés par un chiffrement bancaire. Personne ne peut voir vos informations, même pas nous.
- **Intégration IA native** - Grâce au protocole MCP, connectez Claude, ChatGPT ou tout autre assistant IA directement à vos données.
- **Multi-projets** - Gérez tous vos projets depuis un seul tableau de bord. Plus besoin de jongler entre plusieurs outils.

---

## Fonctionnalités

### Core

- **Connexions** - Branchez votre PostgreSQL ou MySQL existant
- **Collections** - Définissez vos structures de données sans écrire de SQL
- **Contenu** - Gérez vos entrées via le dashboard ou l'API
- **API REST auto-générée** - Chaque collection expose automatiquement ses endpoints
- **Clés API** - Contrôle d'accès granulaire par collection et par permission
- **Stockage** - Upload de fichiers en local ou compatible S3 (MinIO, AWS, Cloudflare R2)
- **Membres** - Invitez des collaborateurs avec des rôles (owner, admin, member)

### Intégration IA (MCP)

Skema implémente le protocole MCP (Model Context Protocol) pour permettre aux assistants IA d'interagir avec vos données de manière sécurisée.

Outils MCP disponibles :

- `get_collections` - Lister les collections
- `get_collection_items` - Récupérer les données
- `create_collection_item` - Créer du contenu
- `update_collection_item` - Modifier du contenu
- `search_collection_items` - Rechercher
- `batch_create_items` - Import en masse

Cas d'usage populaires :

- "Génère 50 descriptions produits pour mon catalogue"
- "Analyse les commandes du mois et donne-moi les tendances"
- "Traduis tous les articles de blog en anglais"
- "Crée un rapport des clients inactifs depuis 3 mois"

---

## Stack technique

- **Backend** - Go, Gin, PostgreSQL / MySQL, Redis, Swaggo
- **Frontend** - À venir

---

## Déploiement

Skema peut être déployé de deux manières :

### Auto-hébergement

Hébergez Skema sur vos propres serveurs. Vous gardez le contrôle complet de votre infrastructure et de vos données.

```bash
git clone https://github.com/Guybe06/skema-cms.git
cd skema-cms/backend

cp .env.example .env
# Remplissez les variables dans .env

make run
```

L'API sera disponible sur `http://localhost:3000/v1`  
La documentation Swagger sur `http://localhost:3000/docs/v1/index.html`

### Service hébergé

Utilisez Skema directement sur notre plateforme. Aucune installation nécessaire, nous gérons l'infrastructure pour vous.

Disponible sur [skemacms.com](https://skemacms.com)

---

## Configuration

Copiez `.env.example` en `.env` et renseignez les variables suivantes :

| Variable         | Description                                       |
| ---------------- | ------------------------------------------------- |
| `APP_PORT`       | Port du serveur (défaut : 3000)                   |
| `APP_ENV`        | Environnement (development / production)          |
| `CMS_DB_HOST`    | Hôte de la base de données Skema                  |
| `CMS_DB_NAME`    | Nom de la base de données Skema                   |
| `JWT_SECRET`     | Clé secrète pour les tokens JWT                   |
| `ENCRYPTION_KEY` | Clé de chiffrement des credentials clients        |
| `REDIS_URL`      | URL Redis (optionnel, fallback mémoire si absent) |
| `RESEND_API_KEY` | Clé API Resend pour l'envoi d'emails              |
| `STORAGE_DRIVER` | local ou s3                                       |
| `S3_ENDPOINT`    | Endpoint S3 (si STORAGE_DRIVER=s3)                |
| `S3_ACCESS_KEY`  | Clé d'accès S3 (si STORAGE_DRIVER=s3)             |
| `S3_SECRET_KEY`  | Clé secrète S3 (si STORAGE_DRIVER=s3)             |
| `S3_BUCKET`      | Nom du bucket S3 (si STORAGE_DRIVER=s3)           |
| `S3_REGION`      | Région S3 (si STORAGE_DRIVER=s3)                  |

---

## Roadmap

Skema est un projet en constante évolution. Voici ce que nous prévoyons pour 2026 :

### Disponible aujourd'hui

- Connexion PostgreSQL et MySQL
- Collections avec tous les types de champs
- Relations entre collections
- API REST avec clés API granulaires
- Système d'organisations et de rôles
- Dashboard avec métriques
- MCP Tools

### Q1 2026

- Types avancés (workflows, calculs automatiques, agrégations, audit logs)
- Recherche full-text (intégration Elasticsearch)

### Q2 2026

- API GraphQL (en complément de l'API REST existante)
- Webhooks (notifications temps réel pour vos intégrations)

### Q3 2026

- App mobile (iOS et Android)

### Q4 2026

- Support MongoDB & Firebase (bases NoSQL)

---

## Communauté

Rejoignez la communauté Skema et contribuez au projet :

- **GitHub** - Issues ouvertes, feature requests, contributions
- **Discord** - Parlez directement avec l'équipe et la communauté
- **Community** - Rejoignez la communauté sur [community.skemacms.com](https://community.skemacms.com)
- **Documentation** - Guides complets et référence API sur [docs.skemacms.com](https://docs.skemacms.com)
- **Email** - Contactez-nous à contributors@skemacms.com

---

## Sécurité

Vos données ne passent jamais par nos serveurs. Nous ne les voyons pas, nous ne les stockons pas.

### Ce que nous gardons

- Vos informations de connexion (protégées)
- La structure de vos contenus
- Les comptes de votre équipe

### Ce que nous ne touchons PAS

- Le contenu de vos données
- Vos textes, images, informations clients
- Aucune copie, jamais

---

## Licence

[GNU Affero General Public License v3.0](LICENSE)

---

## Contribuer

Les contributions sont les bienvenues ! Voici comment participer au projet :

### Via Pull Request

1. Fork le projet
2. Créez une branche pour votre fonctionnalité (`git checkout -b feature/ma-fonctionnalite`)
3. Commit vos changements (`git commit -m 'Ajoute ma fonctionnalité'`)
4. Push vers la branche (`git push origin feature/ma-fonctionnalite`)
5. Ouvrez une Pull Request

### Via l'application Skema

Utilisez le module "Boîte à suggestions" directement dans l'application Skema pour proposer des améliorations, signaler des bugs ou suggérer de nouvelles fonctionnalités.

### Guidelines de contribution

- Suivez le style de code existant
- Ajoutez des tests pour les nouvelles fonctionnalités
- Mettez à jour la documentation si nécessaire
- Assurez-vous que les tests passent avant de soumettre
- Une PR par fonctionnalité ou bugfix

### Communication

Rejoignez notre Discord pour discuter avec l'équipe et la communauté avant de commencer des contributions importantes.

N'hésitez pas à faire évoluer la plateforme. Vive la liberté de nos données !

---

## À propos

Skema est développé par **Winsa Ltd**.

**Développeur principal** : Guylain Béni (guybe)
