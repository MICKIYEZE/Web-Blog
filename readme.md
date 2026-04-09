# Personal Go Web Blog (Dockerized)

A lightweight personal blog built in Go, featuring templates, authentication, JWT middleware, and persistent article storage. The entire project runs inside Docker using Docker Compose.

## 🚀 Technologies Used
- Go 1.22
- Docker
- Docker Compose
- HTML Templates
- JWT Authentication

## 📦 Project Structure
articles/            # Blog posts (persisted via Docker volume)
handlers/            # HTTP handlers and middleware
templates/           # HTML templates
main.go              # Application entry point
Dockerfile
docker-compose.yml
.env.example

## 🔧 Requirements
- Docker 20+
- Docker Compose v2+

## ▶️ Running the Project

1. Clone the repository
git clone https://github.com/<your-username>/WEB-BLOG-SUT.git
cd WEB-BLOG-SUT

2. Create your .env file
cp .env.example .env
Edit the values as needed.

3. Build and start the containers
docker-compose build
docker-compose up -d

4. Open the application
http://localhost:8080

## 💾 Persistence
Articles are stored in a Docker named volume:
volumes:
  - blog_articles:/articles
This ensures your JSON article files survive container restarts.

## ❤️ Health Check
Docker automatically checks the app every 30 seconds:
healthcheck:
  test: ["CMD", "wget", "-qO-", "http://localhost:8080/health"]
Your Go app must expose /health and return HTTP 200.

## 🔐 Environment Variables
PORT           Internal port (default: 8080)
AUTH_USER      Login username
AUTH_PASS      Login password
ARTICLES_DIR   Directory for article storage
JWT_SECRET     Secret key for JWT middleware

## 🛠 Useful Docker Commands
docker-compose up -d        # Start containers
docker-compose down         # Stop and remove containers
docker-compose logs -f      # View logs
docker-compose ps           # Check container & health status
docker-compose build        # Rebuild images
