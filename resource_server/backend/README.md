# Layer8 Resource Server Backend

## Setup

### Install dependencies

- Clone the repository
- Install dependencies with `go mod tidy`
- Run the server with `go run main.go`

**Note:** Before running, make sure to have `.env` file in the root directory with the following variables:

```bash
JWT_SECRET_KEY=secret
DB_USER=postgres
DB_PASS=
DB_NAME=
DB_HOST=localhost
DB_PORT=5432
SSL_MODE=disable
```

### Database (PostgreSQL)

- Setup a PostgreSQL database and create a table with the following schema:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    salt VARCHAR(255) NOT NULL DEFAULT 'ThisIsARandomSalt123!@#',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE clients (
	id VARCHAR(36) PRIMARY KEY,
	secret VARCHAR NOT NULL,
	name VARCHAR(255) NOT NULL,
	redirect_uri VARCHAR(255) NOT NULL
);

CREATE TABLE user_metadata (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    key VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

```
