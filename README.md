# AtonixCorp Platform

## Full-stack deployment

The public site submits the contact form through Next.js to the Go API, which stores messages in PostgreSQL. Start all three services with Docker Compose:

```bash
cp .env.example .env
# Set secure values for POSTGRES_PASSWORD and JWT_SECRET in .env
docker compose up --build
```

Open `http://localhost:3000`. The API is available on `http://localhost:8080` and its readiness endpoint is `http://localhost:8080/health`.

