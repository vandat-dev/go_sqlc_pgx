# Go SQLC PGX Project

A simple REST API in Go using SQLC for type-safe queries and PGX as the PostgreSQL driver.

## Features
- CRUD for User
- Type-safe DB queries (SQLC)
- Connection pooling (PGX)
- RESTful API (Gin)
- Database migrations
- HTTP test file included

## Requirements
- Go 1.24+
- PostgreSQL 12+
- SQLC v1.29.0+

## Setup
1. Clone the repo
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Install SQLC (if not already installed):
   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```
4. Configure the database

   Update the connection string in `main.go`:
   ```go
   dbURL := "postgresql://username:password@host:port/database?sslmode=disable"
   ```
5. Run database migrations

   Execute the SQL in `db/migration/001_create_users_table.up.sql`:
   ```sql
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       phone VARCHAR(20),
       email VARCHAR(255) UNIQUE NOT NULL,
       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
   );
   ```
6. Generate SQLC code:
   ```bash
   sqlc generate
   ```
7. Build and run the project:
   ```bash
   go build
   ./go_sqlc_pgx
   ```
   Or run directly:
   ```bash
   go run main.go
   ```
   The server will run at: http://localhost:8080

## API Endpoints

### Users

| Method | Endpoint | Description | Body |
|--------|----------|-------------|------|
| POST | `/users` | Create a new user | `{"name": "string", "phone": "string", "email": "string"}` |
| GET | `/users` | Retrieve all users | - |
| GET | `/users/:id` | Retrieve a user by ID | - |

### Example Request/Response

#### Create a new user
```bash
POST /users
Content-Type: application/json

{
    "name": "anh Dat",
    "phone": "0333322615",
    "email": "dat@gmail.com"
}
```

Response:
```json
{
    "ID": 1,
    "Name": "anh Dat",
    "Phone": {
        "String": "0333322615",
        "Valid": true
    },
    "Email": "dat@gmail.com",
    "CreatedAt": {
        "Time": "2024-01-01T10:00:00Z",
        "Valid": true
    }
}
```

#### Retrieve all users
```bash
GET /users
```

Response:
```json
[
    {
        "ID": 1,
        "Name": "anh Dat",
        "Phone": {
            "String": "0333322615",
            "Valid": true
        },
        "Email": "dat@gmail.com",
        "CreatedAt": {
            "Time": "2024-01-01T10:00:00Z",
            "Valid": true
        }
    }
]
```

## Testing

### Using HTTP file

1. Open `test_rest_api.http` in VS Code/Cursor
2. Click "Send Request" on each test case
3. Or use REST Client extension

### Using curl

```bash
# Create a user
curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{"name":"Test User","phone":"0123456789","email":"test@example.com"}'

# Retrieve all users
curl http://localhost:8080/users

# Retrieve a user by ID
curl http://localhost:8080/users/1
```

## Project Structure

```
go_sqlc_pgx/
├── main.go                           # Entry point
├── go.mod                           # Go modules
├── go.sum                           # Dependencies checksum
├── sqlc.yaml                        # SQLC configuration
├── test_rest_api.http              # HTTP test file
├── README.md                        # Documentation
├── db/
│   ├── migration/
│   │   ├── 001_create_users_table.up.sql    # Up migration
│   │   └── 001_create_users_table.down.sql  # Down migration
│   └── query/
│       └── user.sql                 # SQL queries
└── internal/
    └── db/                          # Generated SQLC code
        ├── db.go                    # Database interface
        ├── models.go                # Data models
        └── user.sql.go              # User queries
```

## Development

### Adding a new migration

1. Create a migration file:
```
db/migration/002_add_new_table.up.sql
db/migration/002_add_new_table.down.sql
```
2. Add SQL queries in `db/query/`
3. Regenerate SQLC code:
   ```bash
   sqlc generate
   ```

### SQLC Configuration

`sqlc.yaml` file:

```yaml
version: "2"
sql:
  - schema: "db/migration"
    queries: "db/query"
    engine: "postgresql"
gen:
  go:
    package: "db"
    out: "internal/db"
    sql_package: "pgx/v5"
```

## Troubleshooting

### "package not found" error
```bash
go mod tidy
sqlc generate
```

### Database connection error
- Ensure PostgreSQL is running
- Check the connection string
- Verify the database exists

### SQLC generate error
- Check SQL syntax in migration files
- Verify sqlc.yaml configuration
- Ensure SQLC is installed

## TODO

- [ ] Add authentication/authorization
- [ ] Add input validation
- [ ] Add pagination for user list
- [ ] Add logging
- [ ] Add unit tests
- [ ] Add Docker support

## Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

---

Made with ❤️ by [Your Name]
