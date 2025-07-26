# 🔗 Brevity URL Shortener - Backend Server

[![Go Version](https://img.shields.io/badge/Go-1.24.5+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/imraushankr/brevityp)

A high-performance REST API [![API Version](https://img.shields.io/badge/API-v1-orange.svg)](https://your-domain.com/api/v1)
server for URL shortening built with Go. This backend service provides comprehensive analytics, user management, subscription handling, and enterprise-grade monitoring capabilities for URL shortening applications.

## 🌟 Overview

Brevity Backend Server is a robust REST API service that powers URL shortening applications. Built with Go for maximum performance, it provides comprehensive backend functionality including user authentication, detailed analytics, credit systems, subscription management, and production-ready monitoring. This server is designed to be consumed by web frontends, mobile applications, or any client that needs reliable URL shortening capabilities.

## 📚 Table of Contents

- [🚀 Features](#-features)
- [🏗️ Architecture](#-architecture)
- [📝 API Documentation](#-api-documentation)
- [⚙️ Installation](#-installation)
- [🔧 Configuration](#-configuration)
- [💻 Development](#-development)
- [🚢 Deployment](#-deployment)
- [🧪 Testing](#-testing)
- [📊 Monitoring](#-monitoring)
- [🔒 Security](#-security)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)
- [👥 Contributors](#-contributors)

## 🚀 Features

### Core API Features
- ✂️ **URL Shortening API**: RESTful endpoints for creating and managing short URLs
- 🔀 **Redirect Service**: High-performance URL redirection with caching
- 📊 **Analytics API**: Comprehensive click tracking and reporting endpoints
- ⚡ **High Performance**: Built with Go's concurrency model for maximum throughput

### User Management API
- 👤 **Authentication Endpoints**: JWT-based auth with login/register/refresh
- 🔐 **Password Security**: Bcrypt hashing with secure password management
- ✉️ **Email Verification**: Account verification workflow via API
- 🔄 **Password Reset**: Secure password recovery endpoints
- 🖼️ **Avatar Management**: Cloudinary integration for profile picture uploads

### Business Logic API
- 💰 **Credit System**: RESTful endpoints for credit management and tracking
- 🎟️ **Promo Codes**: API for promotional credit distribution
- 🔄 **Subscription Management**: Complete subscription lifecycle API
- 💳 **Payment Integration**: Transaction tracking and payment history
- 📈 **Usage Analytics**: Detailed usage reports and insights via API

### Server Infrastructure
- 🩺 **Health Check Endpoints**: Comprehensive health and status monitoring
- 📈 **Prometheus Metrics**: Built-in metrics collection and monitoring endpoints
- 🏗️ **Database Migrations**: Version-controlled schema management system
- 🔍 **Structured Logging**: Zap-based logging with configurable levels
- 🛡️ **API Validation**: Robust request validation and sanitization middleware

### 📁 Key Architecture Components

- **`src/cmd/`**: Application entry points and CLI commands
- **`src/configs/`**: Configuration management and validation
- **`src/internal/app/`**: Application initialization and dependency injection
- **`src/internal/handlers/v1/`**: HTTP request handlers for API v1
- **`src/internal/middleware/`**: HTTP middleware (auth, logging, CORS, rate limiting)
- **`src/internal/models/`**: Data models and database schemas
- **`src/internal/pkg/`**: Internal packages and utilities
  - **`auth/`**: JWT handling and authentication logic
  - **`database/`**: Database connection and configuration
  - **`email/`**: Email service integration (SMTP)
  - **`interfaces/`**: Interface definitions for dependency injection
  - **`logger/`**: Structured logging setup
  - **`storage/`**: File storage integration (Cloudinary)
- **`src/internal/repository/`**: Data access layer with database operations
- **`src/internal/routes/v1/`**: API route definitions and grouping
- **`src/internal/services/`**: Business logic and service layer
- **`src/internal/utils/`**: Shared utility functions
- **`src/migrations/`**: Database schema migrations
- **`scripts/`**: Build, deployment, and utility scripts


## 🏗️ Server Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Load Balancer │    │  REST API Server│    │    Database     │
│     (Nginx)     │───▶│      (Gin)      │───▶│    (SQLite)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   Monitoring    │
                       │  (Prometheus)   │
                       └─────────────────┘
```

### Backend Components
- **REST API Server**: Gin-based HTTP server with middleware for auth, logging, and metrics
- **Database Layer**: GORM with SQLite for development, easily adaptable to PostgreSQL/MySQL
- **Authentication Service**: JWT tokens with refresh token support
- **File Storage Service**: Cloudinary integration for avatar uploads
- **Monitoring Service**: Prometheus metrics with custom collectors
- **Redirect Service**: High-performance URL redirection with caching

## 📝 API Documentation

### 🌐 Base URL
```
Development: http://localhost:8080/api/v1
Production:  https://your-domain.com/api/v1
```

### 🔐 Authentication
All protected endpoints require an `Authorization` header:
```http
Authorization: Bearer <your-jwt-token>
```

### 📡 API Endpoints

#### 🖥️ System Routes

| Method | Endpoint           | Description                     | Auth Required | Response Format |
|--------|--------------------|---------------------------------|---------------|-----------------|
| GET    | `/system/health`   | Health check endpoint           | No            | JSON            |
| GET    | `/system/status`   | System status information       | No            | JSON            |
| GET    | `/system/metrics`  | Prometheus metrics endpoint     | No            | Text/Plain      |
| GET    | `/system/stats`    | Application statistics          | No            | JSON            |
| GET    | `/system/config`   | Configuration details           | Yes           | JSON            |

#### 🔑 Authentication Routes

| Method | Endpoint                     | Description                          | Auth Required | Body Required |
|--------|------------------------------|--------------------------------------|---------------|---------------|
| POST   | `/auth/signup`               | Register new user                    | No            | Yes           |
| POST   | `/auth/signin`               | User login                           | No            | Yes           |
| POST   | `/auth/signout`              | User logout                          | No            | No            |
| GET    | `/auth/verify-email`         | Verify email address                 | No            | Query param   |
| POST   | `/auth/forgot-password`      | Initiate password reset              | No            | Yes           |
| PATCH  | `/auth/reset-password/:token`| Complete password reset              | No            | Yes           |
| PATCH  | `/auth/change-password`      | Change password (authenticated)      | Yes           | Yes           |
| POST   | `/auth/refresh`              | Refresh access token                 | Refresh token | Yes           |

#### 👤 User Routes

| Method | Endpoint           | Description                     | Auth Required | Body Required |
|--------|--------------------|---------------------------------|---------------|---------------|
| GET    | `/users/me`        | Get user profile                | Yes           | No            |
| PUT    | `/users/me`        | Update user profile             | Yes           | Yes           |
| POST   | `/users/avatar`    | Upload user avatar              | Yes           | Multipart     |
| DELETE | `/users/me`        | Delete user account             | Yes           | No            |

#### ✂️ URL Routes

| Method | Endpoint               | Description                     | Auth Required | Body Required |
|--------|------------------------|---------------------------------|---------------|---------------|
| POST   | `/urls`                | Create new short URL            | Optional*     | Yes           |
| GET    | `/r/:code`             | Redirect to original URL        | No            | No            |
| GET    | `/urls`                | Get user's URLs                 | Yes           | No            |
| GET    | `/urls/:id`            | Get URL details                 | Yes           | No            |
| PUT    | `/urls/:id`            | Update URL                      | Yes           | Yes           |
| DELETE | `/urls/:id`            | Delete URL                      | Yes           | No            |
| GET    | `/urls/:id/analytics`  | Get URL analytics               | Yes           | No            |

*Anonymous users have limited URL creation capabilities*

#### 💰 Credit Routes

| Method | Endpoint               | Description                     | Auth Required | Body Required |
|--------|------------------------|---------------------------------|---------------|---------------|
| GET    | `/credits/balance`     | Get user credit balance         | Yes           | No            |
| POST   | `/credits/apply-promo` | Apply promo code                | Yes           | Yes           |
| GET    | `/credits/usage`       | Get credit usage history        | Yes           | No            |

#### 🔄 Subscription Routes

| Method | Endpoint                  | Description                     | Auth Required | Body Required |
|--------|---------------------------|---------------------------------|---------------|---------------|
| POST   | `/subscriptions`          | Create new subscription         | Yes           | Yes           |
| GET    | `/subscriptions`          | Get user subscription           | Yes           | No            |
| PUT    | `/subscriptions`          | Update subscription             | Yes           | Yes           |
| DELETE | `/subscriptions`          | Cancel subscription             | Yes           | No            |
| GET    | `/subscriptions/plans`    | Get available subscription plans| Yes           | No            |
| GET    | `/subscriptions/payments` | Get payment history             | Yes           | No            |

## 📦 Prerequisites & Dependencies

### ⚙️ System Requirements

- **Go**: 1.24.5 or higher
- **SQLite**: 3.x (development) / PostgreSQL 12+ (production recommended)
- **Git**: Latest version
- **Docker**: Optional, for containerized deployment

### 📦 Core Dependencies

| Package | Version | Purpose | Documentation |
|---------|---------|---------|---------------|
| `github.com/gin-gonic/gin` | v1.10.1 | HTTP web framework | [Docs](https://gin-gonic.com/) |
| `gorm.io/gorm` | v1.30.0 | ORM library | [Docs](https://gorm.io/) |
| `gorm.io/driver/sqlite` | v1.6.0 | SQLite driver for GORM | [Docs](https://gorm.io/docs/connecting_to_the_database.html#SQLite) |
| `github.com/golang-jwt/jwt/v5` | v5.2.3 | JWT authentication | [Docs](https://github.com/golang-jwt/jwt) |
| `github.com/teris-io/shortid` | latest | Short ID generation | [Docs](https://github.com/teris-io/shortid) |
| `github.com/cloudinary/cloudinary-go/v2` | v2.10.1 | Cloudinary integration | [Docs](https://cloudinary.com/documentation/go_integration) |
| `github.com/prometheus/client_golang` | v1.22.0 | Metrics collection | [Docs](https://prometheus.io/docs/guides/go-application/) |
| `go.uber.org/zap` | v1.27.0 | Structured logging | [Docs](https://pkg.go.dev/go.uber.org/zap) |

### 🔧 Development Dependencies

| Package | Version | Purpose | Documentation |
|---------|---------|---------|---------------|
| `github.com/golang-migrate/migrate/v4` | v4.18.3 | Database migrations | [Docs](https://github.com/golang-migrate/migrate) |
| `github.com/fsnotify/fsnotify` | v1.9.0 | Filesystem watching | [Docs](https://github.com/fsnotify/fsnotify) |
| `github.com/spf13/viper` | v1.20.1 | Configuration management | [Docs](https://github.com/spf13/viper) |
| `github.com/joho/godotenv` | v1.5.1 | Environment variables | [Docs](https://github.com/joho/godotenv) |

### 🔐 Security Packages

| Package | Version | Purpose | Documentation |
|---------|---------|---------|---------------|
| `golang.org/x/crypto` | v0.40.0 | Cryptographic functions | [Docs](https://pkg.go.dev/golang.org/x/crypto) |
| `github.com/go-playground/validator/v10` | v10.27.0 | Input validation | [Docs](https://github.com/go-playground/validator) |

## ⚙️ Installation

### 🚀 Quick Start

1. **Clone the repository**:
   ```bash
   # HTTPS
   git clone https://github.com/imraushankr/brevityp.git
   
   # SSH
   git clone git@github.com:imraushankr/brevityp.git
   ```

2. **Navigate to server directory**:
   ```bash
   cd brevity/server
   ```

3. **Set up environment configuration**:
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` file with your configuration (detailed breakdown below):

### 📋 Environment Configuration

The `.env.example` file provides a comprehensive configuration template. Here's what each section controls:

#### 🔧 Application Settings
```env
# Runtime environment and debugging
APP_ENV=development                # development/staging/production  
APP_DEBUG=true                    # Enable debug mode for development
ANON_URL_LIMIT=5                  # Max URLs for anonymous users
AUTH_URL_LIMIT=15                 # Max URLs for authenticated users
```

#### 🌐 Server Settings
```env
# Server binding and timeouts
SERVER_HOST=0.0.0.0               # Host to bind to
SERVER_PORT=8080                  # API server port
SERVER_READ_TIMEOUT=10s           # Request read timeout
SERVER_WRITE_TIMEOUT=10s          # Response write timeout
SERVER_SHUTDOWN_TIMEOUT=15s       # Graceful shutdown timeout
```

#### 🗃️ Database Settings
```env
# SQLite configuration with performance tuning
DB_SQLITE_PATH=./data/brevity.db          # Database file location
DB_SQLITE_BUSY_TIMEOUT=5000               # Connection timeout (ms)
DB_SQLITE_FOREIGN_KEYS=true               # Enable FK constraints
DB_SQLITE_JOURNAL_MODE=WAL                # Write-Ahead Logging
DB_SQLITE_CACHE_SIZE=-2000                # Cache size in KB
```

#### 🔐 Security Settings
```env
# JWT Configuration - Use strong secrets (32+ characters)
JWT_ACCESS_SECRET=your_strong_access_secret_here     # Access token secret
JWT_ACCESS_EXPIRY=15m                                # 15 minutes
JWT_REFRESH_SECRET=your_strong_refresh_secret_here   # Refresh token secret  
JWT_REFRESH_EXPIRY=168h                              # 7 days
JWT_RESET_SECRET=your_strong_reset_secret_here       # Password reset secret
JWT_ISSUER=brevity-service                           # Token issuer
JWT_SECURE_COOKIE=true                               # HTTPS-only cookies
```

#### 📧 Email Settings
```env
# SMTP configuration for notifications
EMAIL_PROVIDER=smtp                       # Email provider type
SMTP_HOST=smtp.example.com               # SMTP server host
SMTP_PORT=587                            # 587 for TLS, 465 for SSL
SMTP_USERNAME=your_email@example.com     # SMTP username
SMTP_PASSWORD=your_email_password        # SMTP password
SMTP_FROM_EMAIL=noreply@example.com      # From address
SMTP_FROM_NAME=Brevity Service           # From name
SMTP_USE_TLS=true                        # Enable TLS encryption
```

#### ☁️ Cloudinary Settings
```env
# Required for avatar uploads
CLOUDINARY_CLOUD_NAME=your_cloud_name    # Cloudinary cloud name
CLOUDINARY_API_KEY=your_api_key          # API key
CLOUDINARY_API_SECRET=your_api_secret    # API secret
```

#### 💾 Storage Settings
```env
# File upload configuration
STORAGE_MAX_AVATAR_SIZE=5242880          # 5MB max avatar size
STORAGE_UPLOAD_DIR=./uploads             # Local upload directory
```

#### 📝 Logger Settings
```env
# Logging configuration
LOG_LEVEL=debug                          # debug/info/warn/error
LOG_FORMAT=console                       # console/json
LOG_FILE_PATH=./logs/brevity.log         # Log file location
```

#### 🌍 CORS Settings
```env
# Cross-origin resource sharing
CORS_ENABLED=true                        # Enable CORS
CORS_ALLOW_ORIGINS=*                     # Allowed origins (comma-separated)
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS  # Allowed methods
CORS_MAX_AGE=12h                         # Preflight cache duration
```

#### ⚡ Rate Limiting
```env
# API rate limiting
RATE_LIMIT_ENABLED=true                  # Enable rate limiting
RATE_LIMIT_REQUESTS=100                  # Requests per window
RATE_LIMIT_WINDOW=1m                     # Time window (1 minute)
```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

5. **Set up development environment**:
   ```bash
   task setup
   ```

6. **Start the development server**:
   ```bash
   task dev
   ```

Your API server will be accessible at: `http://localhost:8080`
**API Base URL**: `http://localhost:8080/api/v1`

### 🐳 Docker Installation

1. **Build the Docker image**:
   ```bash
   docker build -t brevity-api-server .
   ```

2. **Run with Docker Compose**:
   ```bash
   docker-compose up -d
   ```

## 🔧 Configuration

### Environment Variables Reference

| Category | Variable | Description | Default | Required |
|----------|----------|-------------|---------|----------|
| **App** | `APP_ENV` | Runtime environment | `development` | No |
| **App** | `APP_DEBUG` | Debug mode | `true` | No |
| **App** | `ANON_URL_LIMIT` | Anonymous URL limit | `5` | No |
| **App** | `AUTH_URL_LIMIT` | Authenticated URL limit | `15` | No |
| **Server** | `SERVER_HOST` | Server host | `0.0.0.0` | No |
| **Server** | `SERVER_PORT` | Server port | `8080` | No |
| **Server** | `SERVER_READ_TIMEOUT` | Read timeout | `10s` | No |
| **Server** | `SERVER_WRITE_TIMEOUT` | Write timeout | `10s` | No |
| **Server** | `SERVER_SHUTDOWN_TIMEOUT` | Shutdown timeout | `15s` | No |
| **Database** | `DB_SQLITE_PATH` | SQLite database path | `./data/brevity.db` | No |
| **Database** | `DB_SQLITE_BUSY_TIMEOUT` | Connection timeout (ms) | `5000` | No |
| **Database** | `DB_SQLITE_FOREIGN_KEYS` | Enable FK constraints | `true` | No |
| **Database** | `DB_SQLITE_JOURNAL_MODE` | Journal mode | `WAL` | No |
| **Database** | `DB_SQLITE_CACHE_SIZE` | Cache size (KB) | `-2000` | No |
| **JWT** | `JWT_ACCESS_SECRET` | Access token secret | - | **Yes** |
| **JWT** | `JWT_ACCESS_EXPIRY` | Access token expiry | `15m` | No |
| **JWT** | `JWT_REFRESH_SECRET` | Refresh token secret | - | **Yes** |
| **JWT** | `JWT_REFRESH_EXPIRY` | Refresh token expiry | `168h` | No |
| **JWT** | `JWT_RESET_SECRET` | Reset token secret | - | **Yes** |
| **JWT** | `JWT_ISSUER` | Token issuer | `brevity-service` | No |
| **JWT** | `JWT_SECURE_COOKIE` | Secure cookies | `true` | No |
| **Email** | `EMAIL_PROVIDER` | Email provider | `smtp` | For email features |
| **Email** | `SMTP_HOST` | SMTP server host | - | For email features |
| **Email** | `SMTP_PORT` | SMTP server port | `587` | For email features |
| **Email** | `SMTP_USERNAME` | SMTP username | - | For email features |
| **Email** | `SMTP_PASSWORD` | SMTP password | - | For email features |
| **Email** | `SMTP_FROM_EMAIL` | From email address | - | For email features |
| **Email** | `SMTP_FROM_NAME` | From display name | `Brevity Service` | No |
| **Email** | `SMTP_USE_TLS` | Enable TLS | `true` | No |
| **Storage** | `CLOUDINARY_CLOUD_NAME` | Cloudinary cloud name | - | For avatar uploads |
| **Storage** | `CLOUDINARY_API_KEY` | Cloudinary API key | - | For avatar uploads |
| **Storage** | `CLOUDINARY_API_SECRET` | Cloudinary API secret | - | For avatar uploads |
| **Storage** | `STORAGE_MAX_AVATAR_SIZE` | Max avatar size (bytes) | `5242880` | No |
| **Storage** | `STORAGE_UPLOAD_DIR` | Local upload directory | `./uploads` | No |
| **Logging** | `LOG_LEVEL` | Logging level | `debug` | No |
| **Logging** | `LOG_FORMAT` | Log format | `console` | No |
| **Logging** | `LOG_FILE_PATH` | Log file path | `./logs/brevity.log` | No |
| **CORS** | `CORS_ENABLED` | Enable CORS | `true` | No |
| **CORS** | `CORS_ALLOW_ORIGINS` | Allowed origins | `*` | No |
| **CORS** | `CORS_ALLOW_METHODS` | Allowed methods | `GET,POST,PUT,DELETE,OPTIONS` | No |
| **CORS** | `CORS_MAX_AGE` | Preflight cache duration | `12h` | No |
| **Rate Limit** | `RATE_LIMIT_ENABLED` | Enable rate limiting | `true` | No |
| **Rate Limit** | `RATE_LIMIT_REQUESTS` | Requests per window | `100` | No |
| **Rate Limit** | `RATE_LIMIT_WINDOW` | Time window | `1m` | No |

### ⚙️ Production Environment Configuration

For production deployment, ensure these critical settings:

```env
# Production settings
APP_ENV=production
APP_DEBUG=false
JWT_SECURE_COOKIE=true

# Use strong, unique secrets (32+ characters)
JWT_ACCESS_SECRET=your_production_access_secret_minimum_32_chars
JWT_REFRESH_SECRET=your_production_refresh_secret_minimum_32_chars
JWT_RESET_SECRET=your_production_reset_secret_minimum_32_chars

# Restrict CORS origins
CORS_ALLOW_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Database optimization for production
DB_SQLITE_CACHE_SIZE=-8000        # Increase cache for better performance
DB_SQLITE_BUSY_TIMEOUT=10000      # Increase timeout for high concurrency

# Enhanced rate limiting
RATE_LIMIT_REQUESTS=60            # More restrictive for production
RATE_LIMIT_WINDOW=1m

# Structured logging for production
LOG_FORMAT=json
LOG_LEVEL=info
```

## 🛠️ Task Commands

### 🚀 Server Management

| Command | Description | Usage |
|---------|-------------|-------|
| `task server` | Run development server with hot reload (Air) | Development |
| `task server:prod` | Run production server | Production |
| `task dev` | Complete dev workflow (migrate + server) | Development |

### 🗃️ Database Management

| Command | Description | Usage |
|---------|-------------|-------|
| `task db:reset` | Reset database (cross-platform) | Development |
| `task db:migrate` | Apply all pending migrations | All environments |
| `task db:rollback` | Rollback last migration | Development |
| `task db:create-migration` | Create new migration files | Development |
| `task db:seed` | Seed database with sample data | Development |

**Example**: Create a new migration
```bash
task db:create-migration -- create_users_table
```

### 🏗️ Development Setup

| Command | Description | Usage |
|---------|-------------|-------|
| `task setup` | Setup complete development environment | Initial setup |
| `task deps` | Download and verify dependencies | Development |
| `task clean` | Clean build artifacts and cache | Development |
| `task mod:tidy` | Clean up go.mod and go.sum | Development |

### 🧪 Testing & Quality

| Command | Description | Usage |
|---------|-------------|-------|
| `task test` | Run all tests with coverage | Development |
| `task test:unit` | Run unit tests only | Development |
| `task test:integration` | Run integration tests | Development |
| `task test:watch` | Run tests in watch mode | Development |
| `task coverage` | Generate detailed coverage report | Development |
| `task lint` | Run golangci-lint | Development |
| `task fmt` | Format code (gofmt + goimports) | Development |

### 🩺 Health & Monitoring

| Command | Description | Usage |
|---------|-------------|-------|
| `task health` | Check server health status | All environments |
| `task metrics` | Display current metrics | Monitoring |
| `task logs` | Show server logs | Debugging |

### 📦 Build & Deploy

| Command | Description | Usage |
|---------|-------------|-------|
| `task build` | Build production binary | Production |
| `task build:docker` | Build Docker image | Production |
| `task deploy:staging` | Deploy to staging environment | Staging |
| `task deploy:prod` | Deploy to production | Production |

### 🛠️ Development Tools

| Command | Description | Usage |
|---------|-------------|-------|
| `task air:init` | Initialize Air configuration | Setup |
| `task gen:docs` | Generate API documentation | Development |
| `task gen:mocks` | Generate test mocks | Development |

### ℹ️ Help

| Command | Description |
|---------|-------------|
| `task --list` | Show all available commands |
| `task --help` | Show detailed help |

## 💻 Development

### 🔄 Development Workflow

1. **Setup development environment**:
   ```bash
   task setup
   ```

2. **Start development server with hot reload**:
   ```bash
   task dev
   ```

3. **Make your changes** - Air will automatically reload the server

4. **Run tests and checks**:
   ```bash
   task test
   task lint
   task fmt
   ```

5. **Create database migrations if needed**:
   ```bash
   task db:create-migration -- add_new_table
   ```

### 🏗️ Project Structure

```
brevity/
└── server/                     # Backend API Server
    ├── src/                    # Source code
    │   ├── cmd/                # Application entrypoints
    │   ├── configs/            # Configuration management
    │   ├── internal/           # Private application code
    │   │   ├── app/            # Application initialization
    │   │   ├── handlers/       # HTTP API handlers
    │   │   │   └── v1/         # API version 1 handlers
    │   │   ├── middleware/     # HTTP middleware
    │   │   ├── models/         # Data models and schemas
    │   │   ├── pkg/            # Internal packages
    │   │   │   ├── auth/       # Authentication utilities
    │   │   │   ├── database/   # Database connection & config
    │   │   │   ├── email/      # Email service integration
    │   │   │   ├── interfaces/ # Interface definitions
    │   │   │   ├── logger/     # Logging utilities
    │   │   │   └── storage/    # File storage (Cloudinary)
    │   │   ├── repository/     # Data access layer
    │   │   ├── routes/         # Route definitions
    │   │   │   └── v1/         # API version 1 routes
    │   │   ├── services/       # Business logic services
    │   │   └── utils/          # Utility functions
    │   └── migrations/         # Database migrations
    ├── scripts/                # Build and deployment scripts
    ├── .air.toml              # Air hot reload configuration
    ├── .dockerignore          # Docker ignore file
    ├── .env.example           # Environment template
    ├── .gitignore             # Git ignore file
    ├── docker-compose.yml     # Docker composition
    ├── Dockerfile             # Container definition
    ├── go.mod                 # Go module definition
    ├── go.sum                 # Go module checksums
    ├── README.md              # This file
    └── Taskfile.yaml          # Task runner configuration
```

### 🧪 Testing

The server includes comprehensive testing:

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test API endpoints and database interactions
- **Load Tests**: Performance testing for high-traffic scenarios

Run specific test types:
```bash
# Run all tests
task test

# Run with coverage
task coverage

# Run specific package tests
go test ./internal/handlers/...

# Run tests with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

## 🚢 Deployment

### 🌐 Production API Server Deployment

#### Option 1: Direct Server Deployment

1. **Build the API server**:
   ```bash
   task build
   ```

2. **Set up production environment**:
   ```bash
   export GIN_MODE=release
   export PORT=8080
   # Set other production environment variables
   ```

3. **Run the API server**:
   ```bash
   ./brevity-server
   ```

#### Option 2: Docker Deployment

1. **Build Docker image**:
   ```bash
   docker build -t brevity-api-server:latest .
   ```

2. **Run with Docker**:
   ```bash
   docker run -d \
     --name brevity-api-server \
     -p 8080:8080 \
     --env-file .env.production \
     brevity-api-server:latest
   ```

#### Option 3: Docker Compose

```yaml
version: '3.8'

services:
  brevity-api-server:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
      - DB_PATH=/data/brevity.db
    volumes:
      - ./data:/data
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - brevity-api-server
    restart: unless-stopped
```

### 🔧 Production Configuration

For production API server deployment, ensure:

1. **Database**: Use PostgreSQL or MySQL instead of SQLite for production loads
2. **SSL/TLS**: Configure HTTPS with proper certificates for API security
3. **Reverse Proxy**: Use Nginx or similar for load balancing and SSL termination
4. **Monitoring**: Set up Prometheus and Grafana for API metrics
5. **Logging**: Configure structured logging with log rotation for API requests
6. **Backup**: Implement regular database backups

### 📊 Monitoring

#### Prometheus Metrics

Access server metrics at: `http://your-domain:8080/api/v1/system/metrics`

Key server metrics monitored:
- HTTP API request duration and count
- Database query performance
- URL creation and click rates via API
- User registration and authentication API calls
- Credit usage and subscription API interactions
- Server health and uptime

#### Health Checks

- **Server Liveness**: `GET /api/v1/system/health`
- **Server Readiness**: `GET /api/v1/system/status`

## 🔒 Security

### 🛡️ API Security Features

- **JWT Authentication**: Secure token-based API authentication
- **Password Hashing**: Bcrypt with configurable rounds for user passwords
- **Input Validation**: Comprehensive API request validation
- **Rate Limiting**: Configurable rate limiting per API endpoint
- **CORS**: Cross-origin resource sharing configuration for web clients
- **Security Headers**: Comprehensive security headers for API responses
- **SQL Injection Protection**: Parameterized queries via GORM

### 🔐 Security Best Practices

1. **Environment Variables**: Never commit API secrets to version control
2. **JWT Secrets**: Use strong, randomly generated secrets for token signing
3. **Database**: Use prepared statements and ORM for query safety
4. **Input Validation**: Validate and sanitize all API request data
5. **HTTPS**: Always use HTTPS for API endpoints in production
6. **Regular Updates**: Keep server dependencies updated

### 🚨 Security Checklist

- [ ] JWT secret is strong and randomly generated
- [ ] All sensitive data is stored in environment variables
- [ ] Database credentials are secure
- [ ] API rate limiting is enabled
- [ ] Input validation is comprehensive for all endpoints
- [ ] HTTPS is configured for all API endpoints in production
- [ ] Security headers are set for API responses
- [ ] Dependencies are regularly updated

## 🤝 Contributing

We welcome contributions to the Brevity API server! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### 📋 Development Process

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** your changes: `git commit -m 'Add amazing feature'`
4. **Push** to the branch: `git push origin feature/amazing-feature`
5. **Open** a Pull Request

### 📝 Code Style

- Follow Go conventions and best practices
- Use `gofmt` for code formatting
- Add tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting

### 🐛 Bug Reports

When reporting server bugs, please include:
- Go version
- Operating system
- API endpoint affected
- Steps to reproduce
- Expected vs actual API response
- Relevant server logs or error messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Contributors

<table>
  <tr>
    <td align="center">
      <a href="https://github.com/imraushankr">
        <img src="https://github.com/imraushankr.png" width="100px;" alt=""/>
        <br />
        <sub><b>Raushan Kumar</b></sub>
      </a>
      <br />
      <a href="https://github.com/imraushankr/brevityp/commits?author=imraushankr" title="Code">💻</a>
      <a href="#design-imraushankr" title="Design">🎨</a>
      <a href="#ideas-imraushankr" title="Ideas">🤔</a>
    </td>
  </tr>
</table>

### 🌟 Acknowledgments

- Thanks to all contributors who have helped make this API server better
- Inspired by modern URL shortening services and RESTful API design
- Built with love using Go and the amazing open-source community

---

<div align="center">

**[⬆ Back to Top](#-brevity-url-shortener---backend-server)**

Made with ❤️ by the Brevity Backend Team

</div>