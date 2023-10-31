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
    phone_number VARCHAR(50) NOT NULL,
    address VARCHAR(255) NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    phone_number_verified BOOLEAN NOT NULL DEFAULT FALSE,
    location_verified BOOLEAN NOT NULL DEFAULT FALSE,
    national_id_verified BOOLEAN NOT NULL DEFAULT FALSE,
    salt VARCHAR(255) NOT NULL DEFAULT 'ThisIsARandomSalt123!@#',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```
