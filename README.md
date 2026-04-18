# API Go — Documentation technique

## Stack technique

- **Langage** : Go
- **Framework HTTP** : Gin
- **Base de données** : PostgreSQL
- **Driver DB** : `jackc/pgx/v5`
- **Authentification** : JWT (`golang-jwt/jwt/v5`)
- **Chiffrement** : bcrypt (`golang.org/x/crypto/bcrypt`)
- **Variables d'environnement** : `joho/godotenv`
- **CORS** : `gin-contrib/cors`

---

## Structure du projet

```
├── main.go
├── .env
├── database/
│   └── database.go          # Connexion PostgreSQL + migrations
├── handlers/
│   ├── user_handler.go      # Handlers HTTP utilisateurs
│   └── auth_handler.go      # Handlers HTTP authentification
├── services/
│   ├── user_service.go      # Logique métier utilisateurs
│   └── auth_service.go      # Logique métier auth (tokens, bcrypt)
├── repository/
│   ├── user_repository.go   # Requêtes SQL utilisateurs
│   └── auth_repository.go   # Requêtes SQL refresh tokens
├── middleware/
│   ├── auth_middleware.go   # Vérification JWT sur les routes protégées
|   └── logger_middleware.go # Log des requêtes entrantes
├── ressources/
│   ├── user.go              # Modèles et DTOs utilisateur
│   └── auth.go              # DTOs authentification
├── router/
│   └── router.go            # Définition des routes
└── utils/
    └── logger.go            # Définition du logger (zap)
```

### Rôle de chaque couche

| Couche | Rôle |
|---|---|
| `handlers/` | Reçoit la requête HTTP, valide le body, appelle le service, renvoie la réponse |
| `services/` | Logique métier (ex : vérifier qu'un email n'existe pas, hasher un mot de passe) |
| `repository/` | Requêtes SQL brutes vers PostgreSQL |
| `middleware/` | Intercepte les requêtes avant les handlers (vérification JWT) |
| `ressources/` | Définition des structs Go et DTOs (Data Transfer Objects) |

---

## Variables d'environnement

Fichier `.env` à la racine du projet :

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=madb
ALLOWED_ORIGINS=http://localhost:5173
JWT_SECRET=un-secret-long-et-random
```


---

## Base de données

### lancer PostgreSQL avec Docker

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=madb \
  -p 5432:5432 \
  postgres:latest
```

### Migrations

Les tables sont créées automatiquement au démarrage dans `database/database.go` via la fonction `migrate()`. Deux tables sont créées :

**`users`**
```sql
CREATE TABLE IF NOT EXISTS users (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(150) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    role       VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**`refresh_tokens`**
```sql
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## Routes

### Routes publiques

| Méthode | Route | Description |
|---|---|---|
| `POST` | `/auth/register` | Créer un compte |
| `POST` | `/auth/login` | Se connecter |
| `POST` | `/auth/refresh` | Renouveler l'access token |
| `POST` | `/auth/logout` | Se déconnecter |

### Routes protégées (JWT requis)

| Méthode | Route | Description |
|---|---|---|
| `GET` | `/me` | Profil de l'utilisateur connecté |
| `GET` | `/users` | Liste tous les utilisateurs |
| `GET` | `/users/:id` | Récupère un utilisateur |
| `POST` | `/users` | Créer un utilisateur |
| `PUT` | `/users/:id` | Modifier un utilisateur |
| `DELETE` | `/users/:id` | Supprimer un utilisateur |

Les routes protégées nécessitent le header :
```
Authorization: Bearer <access_token>
```

---

## Authentification JWT

### Fonctionnement des deux tokens

| | Access Token | Refresh Token |
|---|---|---|
| Durée | 15 minutes | 7 jours |
| Stockage côté client | Mémoire (variable JS) | Cookie `httpOnly` |
| Rôle | Authentifier chaque requête | Renouveler l'access token |
| Stocké en base | Non | Oui |

Le refresh token est envoyé via un cookie `httpOnly` — il est **inaccessible au JavaScript**

### Contenu du JWT (payload)

```json
{
  "user_id": 42,
  "role": "user",
  "exp": 1714000000
}
```

### Flow d'authentification

```
POST /auth/login ou /auth/register
        │
        ▼
Vérifie email + password (bcrypt)
        │
        ▼
Génère access_token (15min) + refresh_token (7j)
        │
        ├── access_token → renvoyé dans le body JSON
        └── refresh_token → posé en cookie httpOnly
                │
                ▼
Requêtes suivantes : Authorization: Bearer <access_token>
                │
                ▼ (token expiré)
POST /auth/refresh (cookie envoyé automatiquement par le navigateur)
                │
                ▼
Nouvel access_token renvoyé dans le body
```

### Refresh Token Rotation

À chaque appel à `/auth/refresh` :
1. L'ancien refresh token est **supprimé en base**
2. Un nouveau refresh token est **généré et stocké**
3. Le cookie est **remplacé**

Si un refresh token déjà utilisé est présenté, toutes les sessions de l'utilisateur sont invalidées.

### Contexte Gin

Le middleware injecte les données du token dans le contexte Gin, accessible dans tous les handlers des routes protégées :

```go
// Dans le middleware — injection
c.Set("user_id", int(claims["user_id"].(float64)))
c.Set("role", claims["role"].(string))

// Dans un handler — lecture
userID := c.MustGet("user_id").(int)
role := c.MustGet("role").(string)
```

---

## DTOs (Data Transfer Objects)

Les DTOs permettent de contrôler ce qui entre et sort de l'API, sans exposer le modèle interne (ex : le mot de passe hashé).

| DTO | Usage |
|---|---|
| `UserCreateDTO` | Corps de la requête à la création |
| `UserUpdateDTO` | Corps de la requête à la modification |
| `UserPublicDTO` | Réponse renvoyée au client (sans password) |
| `LoginDTO` | Corps de la requête de login |
| `RegisterDTO` | Corps de la requête de register |
| `TokenResponseDTO` | Réponse contenant access_token + infos user |

---

## CORS

La liste des origines autorisées est configurée via la variable d'environnement `ALLOWED_ORIGINS` (séparées par des virgules).

| Environnement | `ALLOWED_ORIGINS` |
|---|---|
| Dev local | `http://localhost:5173` |
| Docker Compose | `http://localhost:5173,http://frontend:5173` |
| Production | `https://monsite.com` |

---

## Lancer le projet

```bash
# Installer les dépendances
go mod tidy

# Lancer PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=madb \
  -p 5432:5432 postgres:latest

# Lancer l'API
go run main.go
```

L'API est disponible sur `http://localhost:8690`.