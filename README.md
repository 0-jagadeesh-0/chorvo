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
├── server/                    # Backend application root
│   ├── cmd/
│   │   └── server/            # Application entry points
│   │       └── main.go        # Main application entry point
│   │
│   ├── config/                # Configuration management
│   │   └── database.go        # Database configuration
│   │
│   ├── internal/              # Internal packages
│   │   ├── api/
│   │   │   └── v1/           # API version 1
│   │   │       ├── handlers/  # HTTP request handlers
│   │   │       ├── middleware/# HTTP middleware
│   │   │       ├── routes/    # Route definitions
│   │   │       ├── services/  # Business logic
│   │   │       └── utils/     # Utility functions
│   │   │
│   │   ├── domain/           # Domain layer
│   │   │   └── models/       # Domain models
│   │   │
│   │   └── repositories/     # Data access layer
│   │
│   ├── api/                  # API documentation
│   │   └── openapi.yaml      # OpenAPI/Swagger specification
│   │
│   ├── .env                  # Environment variables
│   ├── .gitignore           # Git ignore rules
│   ├── docker-compose.yml    # Docker compose configuration
│   ├── Dockerfile           # Docker build instructions
│   ├── go.mod               # Go module definition
│   └── go.sum               # Go module checksums
│
└── README.md                # Project documentation
```
