# 🧩 Chorvo

**Chorvo** is a collaborative project and team management tool designed for modern teams. It enables organizations to plan, track, manage, budget, and report across projects and categories like software, finance, HR, operations, and more — all in one place.

---

## 🚀 Features

- 🔐 Multi-tenant organization and user management
- 🛠️ Project planning, timelines, and budgeting
- 📊 Reports & progress tracking
- 👥 Role-based access control
- ⚡ REST API-first architecture with versioning
- 📄 API documentation using swagger
- 🐳 Docker-based development setup

---

## 🏗️ Tech Stack

| Layer      | Stack              |
| ---------- | ------------------ |
| Backend    | Go (Gin framework) |
| API Docs   | Swagger (v1)       |
| Database   | PostgreSQL         |
| Auth       | JWT / OTP Email    |
| Deployment | Docker             |

---

## 🧱 Project Structure

```bash
chorvo/
├── backend/
│ ├── cmd/ # App entrypoint (main.go)
│ ├── config/ # Configuration loading
│ ├── database/ # DB connection setup
│ ├── docs/ # API docs (Swagger)
│ ├── internal/
│ │ ├── api/v1/
│ │ │ ├── handlers/ # Handlers for versioned APIs
│ │ │ ├── routes/ # Routes of the APIs
│ │ │ └── services/ # Business logic
│ │ ├── models/
│ │ │ ├── entities/ # DB models
│ │ │ ├── requests/ # Request body structs
│ │ │ └── responses/ # Response body structs
│ │ └── repositories/ # DB interaction logic
│ ├── router/ # Router setup
│ ├── .env
│ ├── docker-compose.yml # docker-compose
│ ├── Dockerfile # chorvo server image
│ ├── go.mod
│ └── go.sum
└── README.md
```
