Migrations and the `migrate` CLI

This project uses the golang-migrate CLI to run SQL migrations from
`internal/database/migrations`.

Required dependency
- The `migrate` CLI (github.com/golang-migrate/migrate/v4/cmd/migrate) must be
  installed and available in your PATH to use the `make` migration targets.

Quick install
- Using the Makefile (recommended for this repo):

```bash
make deps
```

- Manually via `go install` (ensure you include the `postgres` build tag):

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

PATH expectations
- `go install` places the binary in `$GOBIN` (if set) or `$GOPATH/bin`.
  Make sure that directory is in your PATH so `migrate` is discoverable.

Add the Go bin to your PATH (example for bash):

```bash
# Add GOPATH/bin to PATH for this session
export PATH="$(go env GOPATH)/bin:$PATH"

# If you use GOBIN, prefer this form (falls back to GOPATH/bin):
export PATH="$(go env GOBIN 2>/dev/null || echo $(go env GOPATH)/bin):$PATH"
```

Verify installation

```bash
# prints the location of the migrate binary (if in PATH)
command -v migrate

# or print the migrate version
migrate -version
```

Using Makefile migration targets
- Once `migrate` is installed and in PATH, run:

```bash
# run up migrations
make migrate-up

# run down migrations
make migrate-down
```

Environment
- The Makefile uses DB_* variables or values from a local `.env` file to build the
  Postgres connection string. See `Makefile` top variables (`DB_HOST`, `DB_USER`, etc.).

Docker alternative
- If you prefer not to install the CLI locally, you can use the official docker image.
  Example (adjust the database connection string as needed):

```bash
# mount the migrations directory and run the 'up' command
docker run --rm \
  -v "$(pwd)/internal/database/migrations":/migrations \
  --network host \
  migrate/migrate \
  -path=/migrations -database "postgres://postgres:password@localhost:5432/campus_core?sslmode=disable" up
```

Notes
- The Makefile now attempts to auto-install `migrate` if it's not found, but you still
  need to ensure your Go bin (GOBIN/GOPATH/bin) is in PATH so the binary becomes
  discoverable in subsequent shells or CI steps.
- If using CI, prefer installing `migrate` in the CI image or using the Docker approach above.

