# Gin Gorm PostgreSQL Project Template

## Project Overview

A Go web application template using Gin framework, Gorm ORM, and PostgreSQL database.

## Prerequisites

- Go 1.20+
- PostgreSQL
- Make
- Air (for development hot reloading)

## Creating a New Project from Scratch

### 1. Initialize Go Module

```bash
mkdir your-project-name
cd your-project-name
go mod init your-project-name
```

### 2. Install Core Dependencies

```bash
# Gin Web Framework
go get -u github.com/gin-gonic/gin

# Gorm ORM
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres

# Environment Variables
go get -u github.com/joho/godotenv

# Air for hot reloading (optional)
go install github.com/cosmtrek/air@latest
```

### 3. Create Project Structure

```bash
mkdir -p database models controllers routes middleware
touch main.go .env .air.toml Makefile
```

### 4. Basic Configuration Files

#### .env

```
PORT=8080
SECRET_KEY=a3b56d9f93a3a8e04f37bf8e2b1a3f1d8d83c9f37a11a3f2938a1c9e7f09d2e9
DATABASE_URL=postgresql://{your name}:{your password}@ep-cold-violet-a18cw3am-pooler.ap-southeast-1.aws.neon.tech/{project-name}?sslmode=require
```

#### .air.toml (Hot Reload Configuration)

```toml
root = "."
tmp_dir = "tmp"

[build]
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  stop_on_error = true
```

### 5. Initial Project Setup

- Create database configuration in `database/`
- Define models in `models/`
- Implement routes in `routes/`
- Add controllers in `controllers/`

## Development Workflow

### Run Application

#### Hot Reload (Development)

```bash
make air
```

#### Standard Run

```bash
make run
```

### Build

```bash
make build
```

### Testing

```bash
make test
```

### Code Formatting

```bash
make fmt
```

### Clean Build

#### Linux/macOS

```bash
make clean-linux
```

#### Windows

```bash
make clean-windows
```

## Project Structure

```
.
├── main.go
├── Makefile
├── .env
├── .air.toml
├── go.mod
├── go.sum
├── database/
├── models/
├── controllers/
├── routes/
└── middleware/
```

## Dependencies

- Gin Web Framework
- Gorm ORM
- PostgreSQL Driver
- godotenv

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

### 1. Clone the Repository

```bash
git clone https://github.com/iamalfi/go-template-without-db.git
```
