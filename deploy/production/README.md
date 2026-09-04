# Production deployment (Ubuntu + existing Remnawave)

This directory deploys English Reader without changing Remnawave containers.
Only loopback ports are opened by this Compose project: frontend `18080`, HTTP
API `18081`, and gRPC-Web Envoy `18082`. The existing Remnawave Nginx receives
public HTTPS traffic and proxies the `krylov.reader.mooo.com` host to them.

## Before deployment

1. Point `krylov.reader.mooo.com` to the server public IPv4 address and wait
   until `getent ahostsv4 krylov.reader.mooo.com` returns that address.
2. Ensure inbound TCP port 80 is allowed temporarily for Let's Encrypt HTTP
   validation. Port 443 must remain owned by `rw-core`.
3. Install Docker Engine with Compose plugin, Git, and Certbot on Ubuntu.

## Server installation

```bash
sudo install -d -m 0755 /opt/english-reader
sudo chown "$USER":"$USER" /opt/english-reader
git clone https://github.com/stonik02/english_reader.git /opt/english-reader/app
cd /opt/english-reader/app/deploy/production
cp .env.example .env
```

Generate secrets on the server, edit `.env`, and replace both placeholder
values. `POSTGRES_PASSWORD` must contain only letters and digits because it is
part of a PostgreSQL connection URL.

```bash
openssl rand -hex 32
```

Create the certificate before adding the Nginx host:

```bash
sudo certbot certonly --standalone -d krylov.reader.mooo.com
```

Copy the two resulting certificate files into the existing Nginx container by
adding the following mounts to its Compose definition, then recreate only that
container:

```text
/etc/letsencrypt/live/krylov.reader.mooo.com/fullchain.pem -> /etc/nginx/ssl/krylov.reader.mooo.com/fullchain.pem:ro
/etc/letsencrypt/live/krylov.reader.mooo.com/privkey.pem -> /etc/nginx/ssl/krylov.reader.mooo.com/privkey.pem:ro
```

Back up `/opt/remnawave/nginx.conf`, append
`remnawave-reader.nginx.conf`, validate with `docker exec remnawave-nginx nginx
-t`, and reload with `docker exec remnawave-nginx nginx -s reload`.

## Start Reader and import only the dictionary

Copy the converted dictionary JSONL to `imports/` on the server, then run:

```bash
mkdir -p imports
docker compose --env-file .env up --build -d
docker compose --env-file .env run --rm --no-deps --entrypoint /app/reader-migrate api -direction up
docker compose --env-file .env run --rm --no-deps --entrypoint /app/reader-dictionary-import api \
  -file /imports/en-ru-wiktionary-2026-09-04.jsonl -version 2026-09-04 -source wiktionary
```

The import contains dictionary entries only: no local users, books, EPUB files,
reading positions, or vocabulary are copied.

## Verification

```bash
curl -fsS http://127.0.0.1:18081/health/ready
docker compose --env-file .env ps
```

Then open `https://krylov.reader.mooo.com`, register the first account, and
look up a known word such as `been`.
