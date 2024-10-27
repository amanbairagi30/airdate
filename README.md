# airdate

airdate is a platform for gamers and content creators to connect, collaborate, and find gaming partners. This project uses a modern full-stack architecture with Next.js for the frontend and Go for the backend.

## Tech Stack

### Frontend
- Next.js 13+ (React framework)
- TypeScript
- Tailwind CSS

### Backend
- Go 1.22.6
- PostgreSQL 13
- Docker

## Getting Started

### Prerequisites
- Node.js 14+ and npm
- Go 1.22.6+
- Docker and Docker Compose
- PostgreSQL 13 (if running without Docker)

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/pixel-and-chill.git
   cd pixel-and-chill
   ```

2. Set up the backend:
   ```bash
   cd backend
   cp cmd/server/.env.example cmd/server/.env
   # Edit the .env file with your database credentials and JWT secret
   docker-compose up --build
   ```

3. Set up the frontend:
   ```bash
   cd frontend
   npm install
   ```

4. Run the development servers:
   
   For the backend (if not using Docker):
   ```bash
   cd backend/cmd/server
   go run main.go
   ```

   For the frontend:
   ```bash
   npm run dev
   ```

5. Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Features

- User registration and authentication
- Profile management
- Connect Twitch and Discord accounts
- Find gaming partners
- Collaborate on content creation

## Learn More

To learn more about the technologies used in this project, check out the following resources:

- [Next.js Documentation](https://nextjs.org/docs)
- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)

