# Notes API

A production-oriented REST API for note management, backed by PostgreSQL.

## Concepts (learning)

- PostgreSQL with `pgx/v5` and connection pooling (`pgxpool`)
- `database/sql` patterns (query, scan, `sql.ErrNoRows`)
- Repository pattern (separating database logic from handlers)
- SQL migrations (plain `.sql` files)
- Docker Compose for local database setup
- Structured logging with `log/slog`
- Graceful shutdown with `os/signal` and `context`

## Endpoints

| Method   | Path            | Description        |
|----------|-----------------|--------------------|
| `GET`    | `/health`       | Health check       |
| `POST`   | `/notes/`       | Create a note      |
| `GET`    | `/notes/`       | List all notes     |
| `GET`    | `/notes/{id}`   | Get note by ID     |
| `DELETE` | `/notes/{id}`   | Delete a note      |

## Project layout

```
notes-api/
├── cmd/
│   └── main.go          # Entry point, server setup
├── migrations/
│   └── 001_create_notes.sql
├── go.mod
├── docker-compose.yml
├── .gitignore
└── README.md
```

*(Layout will grow as the project evolves.)*

## Usage

```bash
git clone https://github.com/theusualdeveloper/notes-api.git
cd notes-api
docker compose up -d
go run ./cmd
```

Server listens on `http://localhost:8080`.

## Requirements

- Go 1.22+
- Docker

## License

MIT
