# This is Simple Notes app backend
## How to run
- Create Postgres DB with name "notes_db"
- Create .env file with the next content:

```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=XXXXXXXX
DB_NAME=notes_db
DB_PORT=5432
SSLMODE=disable

JWT_SECRET=Aramzamzam
PORT=3001
ENV=dev # or prod
```
insert your data ^^^
- Install dependencies with `go mod tidy`
- Run it with the command `go run ./cmd/server/main.go`

- See API swagger doccumentations on `http://localhost:YOUR_ENV_PORT/swagger/index.html#/`

### enjoy
