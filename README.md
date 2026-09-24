# AtonixCorp Platform

## Full-stack deployment

The public site submits the contact form through Next.js to the Go API, which stores messages in PostgreSQL. Start all three services with Docker Compose:

```bash
cp .env.example .env
# Set secure values for POSTGRES_PASSWORD and JWT_SECRET in .env
docker compose up --build
```

## Apache2 domain deployment

Compose binds the website and API only to the server loopback interface: `127.0.0.1:3000` and `127.0.0.1:8080`. Your existing Apache2 instance is the public HTTPS edge proxy and routes the domains to those local services.

Before starting on the production server, configure these DNS records to the server's public IPv4 address:

```text
atonixcorp.com      A      <server-public-ip>
www.atonixcorp.com  A      <server-public-ip>
api.atonixcorp.com  A      <server-public-ip>
```

Set `POSTGRES_PASSWORD` and `JWT_SECRET` in `.env`, then run:

```bash
docker compose up -d --build
```

## Contact email delivery

The contact form stores each inquiry in PostgreSQL and sends a notification to `contact@atonixcorp.com` through Google SMTP. Set these values in `.env` before starting the API:

```text
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-google-workspace-or-gmail-address
SMTP_PASSWORD=your-google-app-password
SMTP_RECIPIENT=contact@atonixcorp.com
```

Use a Google App Password, not the account sign-in password. The SMTP account must be allowed to send mail to `contact@atonixcorp.com`.

Install the supplied Apache virtual-host configuration from `deploy/apache/atonixcorp.conf`, enable the Apache proxy and SSL modules, and point its certificate paths to your existing certificates. The API readiness endpoint is available at `https://api.atonixcorp.com/health`.

```bash
sudo a2enmod proxy proxy_http headers ssl rewrite
sudo cp deploy/apache/atonixcorp.conf /etc/apache2/sites-available/atonixcorp.conf
sudo a2ensite atonixcorp.conf
sudo apache2ctl configtest
sudo systemctl reload apache2
```

Use Certbot or your existing certificate management workflow to provision certificates for `atonixcorp.com`, `www.atonixcorp.com`, and `api.atonixcorp.com` before enabling the TLS virtual hosts.

