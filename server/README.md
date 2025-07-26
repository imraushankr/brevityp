# 🔗 Brevity URL Shortener Server

A high-performance URL shortening service with analytics, user management, and subscription features.

## 📚 Table of Contents
- 🚀 [Features](#features)
- 📝 [API Documentation](#api-documentation)
- ⚙️ [Installation](#installation)
- 🔧 [Configuration](#configuration)
- 💻 [Development](#development)
- 🚢 [Deployment](#deployment)

## 🚀 Features
- ✂️ URL shortening with custom codes
- 👤 User authentication and management
- 📊 URL analytics and click tracking
- 💰 Credit system for premium features
- 🔄 Subscription management
- 🩺 Health monitoring and metrics

## 📝 API Documentation

### 🌐 Base URL
`https://your-domain.com/api/v1`

## 📦 Prerequisites & Packages

### ⚙️ System Requirements
- Go 1.24.5+
- SQLite 3.x
- Git

### 📦 Core Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/gin-gonic/gin` | v1.10.1 | HTTP web framework |
| `gorm.io/gorm` | v1.30.0 | ORM library |
| `gorm.io/driver/sqlite` | v1.6.0 | SQLite driver for GORM |
| `github.com/golang-jwt/jwt/v5` | v5.2.3 | JWT authentication |
| `github.com/teris-io/shortid` | latest | Short ID generation |
| `github.com/cloudinary/cloudinary-go/v2` | v2.10.1 | Cloudinary integration |
| `github.com/prometheus/client_golang` | v1.22.0 | Metrics collection |
| `go.uber.org/zap` | v1.27.0 | Structured logging |

### 🔧 Development Dependencies
| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/golang-migrate/migrate/v4` | v4.18.3 | Database migrations |
| `github.com/fsnotify/fsnotify` | v1.9.0 | Filesystem watching |
| `github.com/spf13/viper` | v1.20.1 | Configuration management |
| `github.com/joho/godotenv` | v1.5.1 | Environment variables |

### 🔐 Security Packages
| Package | Version | Purpose |
|---------|---------|---------|
| `golang.org/x/crypto` | v0.40.0 | Cryptographic functions |
| `github.com/go-playground/validator/v10` | v10.27.0 | Input validation |

### 🔐 Authentication
All endpoints except public ones require an `Authorization` header with a valid JWT token.

### 📡 API Endpoints

#### 🖥️ System Routes

| Method | Endpoint           | Description                     | Auth Required | Icon |
|--------|--------------------|---------------------------------|---------------|------|
| GET    | `/system/health`   | Health check endpoint           | No            | 🩺   |
| GET    | `/system/status`   | System status information       | No            | 📊   |
| GET    | `/system/metrics`  | Prometheus metrics endpoint     | No            | 📈   |
| GET    | `/system/stats`    | Application statistics          | No            | 📊   |
| GET    | `/system/config`   | Configuration details           | Yes           | ⚙️   |

#### 🔑 Authentication Routes

| Method | Endpoint                     | Description                          | Auth Required | Icon |
|--------|------------------------------|--------------------------------------|---------------|------|
| POST   | `/auth/signup`               | Register new user                    | No            | ✏️   |
| POST   | `/auth/signin`               | User login                           | No            | 🔑   |
| POST   | `/auth/signout`              | User logout                          | No            | 🚪   |
| GET    | `/auth/verify-email`         | Verify email address                 | No            | ✉️   |
| POST   | `/auth/forgot-password`      | Initiate password reset              | No            | 🔓   |
| PATCH  | `/auth/reset-password/:token`| Complete password reset              | No            | 🔄   |
| PATCH  | `/auth/change-password`      | Change password (authenticated)      | Yes           | 🔐   |
| POST   | `/auth/refresh`              | Refresh access token                 | Refresh token | 🔄   |

#### 👤 User Routes

| Method | Endpoint           | Description                     | Auth Required | Icon |
|--------|--------------------|---------------------------------|---------------|------|
| GET    | `/users/me`        | Get user profile                | Yes           | 👤   |
| PUT    | `/users/me`        | Update user profile             | Yes           | ✏️   |
| POST   | `/users/avatar`    | Upload user avatar              | Yes           | 🖼️  |
| DELETE | `/users/me`        | Delete user account             | Yes           | 🗑️  |

#### ✂️ URL Routes

| Method | Endpoint               | Description                     | Auth Required | Icon |
|--------|------------------------|---------------------------------|---------------|------|
| POST   | `/urls`                | Create new short URL            | Optional*     | ➕   |
| GET    | `/r/:code`             | Redirect to original URL        | No            | 🔀   |
| GET    | `/urls`                | Get user's URLs                 | Yes           | 📋   |
| GET    | `/urls/:id`            | Get URL details                 | Yes           | 🔍   |
| PUT    | `/urls/:id`            | Update URL                      | Yes           | ✏️   |
| DELETE | `/urls/:id`            | Delete URL                      | Yes           | 🗑️  |
| GET    | `/urls/:id/analytics`  | Get URL analytics               | Yes           | 📊   |

\* Anonymous users have limited URL creation

#### 💰 Credit Routes

| Method | Endpoint               | Description                     | Auth Required | Icon |
|--------|------------------------|---------------------------------|---------------|------|
| GET    | `/credits/balance`     | Get user credit balance         | Yes           | 💰   |
| POST   | `/credits/apply-promo` | Apply promo code                | Yes           | 🎟️  |
| GET    | `/credits/usage`       | Get credit usage history        | Yes           | 📜   |

#### 🔄 Subscription Routes

| Method | Endpoint               | Description                     | Auth Required | Icon |
|--------|------------------------|---------------------------------|---------------|------|
| POST   | `/subscriptions`       | Create new subscription         | Yes           | ➕   |
| GET    | `/subscriptions`       | Get user subscription           | Yes           | 🔍   |
| PUT    | `/subscriptions`       | Update subscription             | Yes           | ✏️   |
| DELETE | `/subscriptions`       | Cancel subscription             | Yes           | 🗑️  |
| GET    | `/subscriptions/plans` | Get available subscription plans| Yes           | 📋   |
| GET    | `/subscriptions/payments` | Get payment history          | Yes           | 💳   |

## 🛠️ Task Commands

### 🚀 Server Management
| Command | Description | Example |
|---------|-------------|---------|
| `task server` | Run development server with hot reload | `task server` |
| `task server:prod` | Run production server | `task server:prod` |

### 🗃️ Database Management
| Command | Description | Example |
|---------|-------------|---------|
| `task db.reset` | Reset database (Windows compatible) | `task db.reset` |
| `task db.reset.go` | Reset database (Cross-platform Go version) | `task db.reset.go` |
| `task migrate.up` | Apply all pending migrations | `task migrate.up` |
| `task migrate.down` | Rollback last migration | `task migrate.down` |
| `task migrate.create` | Create new migration files | `task migrate.create create_users_table` |

### 🏗️ Development Setup
| Command | Description | Example |
|---------|-------------|---------|
| `task setup` | Setup development environment | `task setup` |
| `task dev` | Complete dev workflow (migrate + server) | `task dev` |

### 🩺 Health & Monitoring
| Command | Description | Example |
|---------|-------------|---------|
| `task health` | Check server health status | `task health` |

### ℹ️ Help
| Command | Description | Example |
|---------|-------------|---------|
| `task list` | Show all available commands | `task list` |
