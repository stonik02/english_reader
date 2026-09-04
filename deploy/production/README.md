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

## Updating the deployed application

The normal update path does not recreate the database or remove uploaded EPUB
files. Their Docker volumes stay in place.

### 1. Check and push changes from the development computer

Run checks relevant to the changed code, then commit and push. Do not add
secrets or local dictionary source archives to Git.

```bash
git status
git add <changed-files>
git commit -m "Describe the change"
git push origin main
```

Never commit any of the following:

- `deploy/production/.env` — production database password and JWT secret;
- `frontend/.env.local` — local browser configuration;
- `backend/en-extract/` or raw Kaikki `.gz` archives;
- `deploy/production/imports/*.jsonl` — the server copy of a dictionary dump.

### 2. Download the new code and restart Reader on the server

Connect through SSH and run these commands from the deployment directory:

```bash
cd /opt/english-reader/app
git status
git pull --ff-only origin main

cd deploy/production
docker compose --env-file .env up --build -d
docker compose --env-file .env run --rm --no-deps --entrypoint /app/reader-migrate api -direction up
curl -fsS http://127.0.0.1:18081/health/ready
docker compose --env-file .env ps
```

`git pull --ff-only` refuses to overwrite unexpected edits on the server. Do
not edit application source files there; keep server-only values in `.env`.
The migration command is safe on every update: it applies only migrations that
have not run yet.

`docker compose up --build -d` briefly replaces Reader containers when their
image changes. PostgreSQL data, EPUB files, reading positions, and vocabulary
remain in their named volumes. It does not restart Remnawave Nginx.

### 3. Update the dictionary only when it actually changed

Most interface or backend changes do **not** require a dictionary import. For
a new converted JSONL dump, copy it to `imports/` and use a new date-like
version (for example `2026-10-01`), never reuse an old version:

```bash
docker compose --env-file .env run --rm --no-deps --entrypoint /app/reader-dictionary-import api \
  -file /imports/en-ru-wiktionary-2026-10-01.jsonl \
  -version 2026-10-01 \
  -source wiktionary
```

This imports dictionary data only. It does not remove books or user data.

### Roll back application code

Before a rollback, make sure the server checkout has no local edits:

```bash
cd /opt/english-reader/app
git status
git log --oneline -10
```

If `git status` is clean, switch to a known-good commit and rebuild:

```bash
git checkout <known-good-commit>
cd deploy/production
docker compose --env-file .env up --build -d
```

Do not roll database migrations back on your own. A code rollback is normally
safe; a database rollback can lose newer data and should be planned separately.

To return to the latest version later:

```bash
cd /opt/english-reader/app
git switch main
git pull --ff-only origin main
```

## Certificate maintenance

Certbot renews the certificate automatically. The server has a deploy hook in
`/etc/letsencrypt/renewal-hooks/deploy/reload-remnawave-nginx` which recreates
only `remnawave-nginx` after a renewal. This lets the container pick up the new
certificate file; Reader, PostgreSQL, and the Remnawave node are not restarted.
