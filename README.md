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
- PostgreSQL 13 (only if running without Docker)

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/airdate.git
   cd airdate
   ```

2. Set up the backend:

   Option A - Using Docker (Recommended):
   ```bash
   cd backend
   
   # Create .env file for Docker
   cat > cmd/server/.env << EOL
   DB_HOST=db
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=password
   DB_NAME=pixel_and_chill
   JWT_SECRET=your_jwt_secret
   EOL

   # Start the services
   docker compose up --build
   ```

   Option B - Running locally (requires PostgreSQL):
   ```bash
   cd backend/cmd/server

   # Create .env.local file
   cat > .env.local << EOL
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_local_password
   DB_NAME=pixel_and_chill
   JWT_SECRET=your_jwt_secret
   EOL

   # Create database in PostgreSQL
   psql -U postgres -c "CREATE DATABASE pixel_and_chill;"

   # Run the server
   go run main.go
   ```

3. Set up the frontend:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

4. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## Development

### Running the Backend

Option A - Using Docker:
```bash
cd backend
docker compose up --build    # Start with logs
# OR
docker compose up --build -d # Start in background
docker logs -f backend-app-1 # View logs when running in background
```

Option B - Running locally:
```bash
cd backend/cmd/server
go run main.go
```

### Running the Frontend
```bash
cd frontend
npm run dev
```

### Environment Files

The backend supports two environment configurations:

1. `cmd/server/.env` - Used with Docker Compose:
```env
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=pixel_and_chill
JWT_SECRET=your_jwt_secret
```

2. `cmd/server/.env.local` - Used for local development:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_local_password
DB_NAME=pixel_and_chill
JWT_SECRET=your_jwt_secret
```

## Features

- User registration and authentication
- Profile management
- Connect Twitch and Discord accounts
- Find gaming partners
- Collaborate on content creation

## API Endpoints

- `GET /api/health` - Health check
- `POST /api/register` - User registration
- `POST /api/login` - User login
- `GET /api/profile` - Get user profile (protected)
- `POST /api/connect/twitch` - Connect Twitch account (protected)
- `POST /api/connect/discord` - Connect Discord account (protected)
- `GET /api/games/search` - Search games (protected)

## Learn More

- [Next.js Documentation](https://nextjs.org/docs)
- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)

