# This is Simple Notes app backend

Frontend it's mobile React-Native app now, repo: https://github.com/Saltein/notes-app 
The server supports web and mobile versions of the frontend part.

## How to run local
- Create Postgres DB with name "notes_db"
- Create .env file with the next content:

```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=XXXXXXXX
DB_NAME=notes_db
DB_PORT=5432
SSLMODE=disable

# if ENV == prod -> add HOST==your_host


JWT_SECRET=Aramzamzam
PORT=3001
ENV=dev # or prod

# FOR DOCKER

POSTGRES_DB=notes_db
POSTGRES_USER=postgres
POSTGRES_PASSWORD=XXXXXXXX
```
insert your data ^^^
- Install dependencies with `go mod tidy`
- Run it with the command `go run ./cmd/server/main.go`

- See API swagger doccumentations on `http://localhost:YOUR_ENV_PORT/swagger/index.html#/` in your browser.

### enjoy

## deploy

There is also docker-compose.yaml file you need to deploy server on your VPS(create it outside of project dir):
```
services:
  backend:
    build:
      context:
        ./notes_backend
    container_name: backend
    env_file:
      - ./notes_backend/.env
    ports:
      - "3001:3001"
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    container_name: database
    restart: always
    env_file:
      - ./notes_backend/.env
    volumes:
      - pgdata:/var/lib/postgresql/data
volumes:
  pgdata:
```