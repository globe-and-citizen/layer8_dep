# Contributing
We welcome the analysis! Start but running the project locally. 

# Run the Project
## Create or configure a Postgres database
1) If you're using Docker, use command: 
*GOTCHA* Note the server, password, and ports!
docker run --env=POSTGRES_USER=default_user --env=POSTGRES_PASSWORD=1234 --env=POSTGRES_DB=local_rs --env=PG_TRUST_LOCALNET=true -p 5544:5432 -d postgres:latest

2) If your using a local postgres DB, simply confirm the username, password, and name of the db you created.
Hint, if on windows, check "services" by opening up the windows menu and searching "services"

2) You will need to run the following queries in your database to create the necessary tables:

``` sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    -- phone_number VARCHAR(50) NOT NULL,
    -- address VARCHAR(255) NOT NULL,
    -- email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    -- phone_number_verified BOOLEAN NOT NULL DEFAULT FALSE,
    -- location_verified BOOLEAN NOT NULL DEFAULT FALSE,
    -- national_id_verified BOOLEAN NOT NULL DEFAULT FALSE,
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

## Other Utilities and Programs You'll Need
1) Download and instal Golang
2) Install the GNU Make utility program by running `$winget install GnuWin32.Make`. You will need to add `C:\Program Files (x86)\GnuWin32\bin` to your environment Path variable.

## Run the sp_mock frontend & backend
1) In any terminal, run make `$make npm_install_all`
2) In any terminal, run make `$make go_mod_tidy`
3) Navigate to `./server` and run `go get "github.com/globe-and-citizen/layer8-utils"`
4) Clone `.env.dev` to `.env` in `sp_mock/backend`
5) Navigate to the root directory. Run `$make run_backend`
4) Clone `.env.dev` to `.env` in `sp_mock/frontend`
8) Navigate to the root directory. Run `$make run_frontend`
9) Configure the `.env` in `/server` to connect to your local PG implementation
    Example: 
    ```
    JWT_SECRET_KEY=secret
    SSL_MODE=disable
    DB_USER=postgres # YOUR DB MAY DIFFER!
    DB_PASS=1234 # YOUR DB MAY DIFFER!
    DB_NAME=local_rs # YOUR DB MAY DIFFER!
    SERVER_PORT=5001
    DB_HOST=localhost
    DB_PORT=5432 # YOUR DB MAY DIFFER!
    UP_999_SECRET_KEY=
    MP_123_SECRET_KEY=
    SSL_ROOT_CERT=
    ```
10) With Golang installed, run `$make run_server` from the root directory
11) Navigate to `http://localhost:5001`. Register a new Layer8 user.
12) Navigate to `http://localhost:5173` to register and login a user of the sp_mock.
