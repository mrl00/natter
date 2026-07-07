# natter

Natter -- the social network for coffee mornings, book groups, and other small gatherings. In truth, this project is a study about API security.

## Project Structure

```
cmd/
├── server/                 # API entrypoint
└── healthcheck/            # healthcheck binary
internal/
├── api/                    # HTTP layer: handlers, router, JSON responses
├── model/                  # domain types (Space, Message, request payloads)
├── repository/             # Repository interface (persistence contract)
├── service/                # business logic and typed errors
└── infra/
    ├── memory/             # in-memory implementation
    └── postgres/           # PostgreSQL implementation
build/docker/
├── Dockerfile              # prod (builder + distroless)
└── Dockerfile.dev          # dev (alpine + air)
migrations/                 # SQL migrations
tests/e2e/                  # hurl end-to-end tests
```

The architecture follows Clean Architecture principles:
- **handlers** depend on `service.Service`
- **service** depends on `repository.Repository`
- **infra/** implements `repository.Repository`

This keeps layers decoupled — swapping persistence only requires a new implementation under `internal/infra`.

## Endpoints

| Method | Path                                     | Description                       |
| ------ | ---------------------------------------- | --------------------------------- |
| GET    | `/health`                                | Health check                      |
| POST   | `/spaces`                                | Create a space                    |
| POST   | `/spaces/{spaceId}/messages`             | Add a message to a space          |
| GET    | `/spaces/{spaceId}/messages`             | List messages (`?since=<RFC3339>`) |
| GET    | `/spaces/{spaceId}/messages/{messageId}` | Get a single message              |

## Running

```sh
go run ./cmd/server
```

Set `DATABASE_URL` to use PostgreSQL (falls back to in-memory):

```sh
DATABASE_URL="postgres://natter:natter@localhost:5432/natter_db?sslmode=disable" go run ./cmd/server
```

### Docker Compose

```sh
# dev with hot reload
docker compose up --build

# prod
docker build -f build/docker/Dockerfile -t natter .
docker run -p 8080:8080 natter
```

## Tests

```sh
# run all e2e tests
hurl --test tests/e2e/*.hurl
```
