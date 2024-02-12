# Following Steps
## Create or configure a Postgres database
1) If you're using Docker, use command: 
*GOTCHA* Note the server, password, and ports!
docker run --env=POSTGRES_USER=default_user --env=POSTGRES_PASSWORD=1234 --env=POSTGRES_DB=local_rs --env=PG_TRUST_LOCALNET=true -p 5544:5432 -d postgres:latest

2) If your using a local postgres DB, simply confirm the username, password, and name of the db you created.
Hint, if on windows, check "services" by opening up the windows menu and searching "services"

2) You will need to run the following queries in your database to create the necessary tables:
```
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



# Layer8
A suite of  network protocol implementations that sum to create an anonymizing reverse proxy dedicated to dissociating a user's true identity from their online content choices.  

## Why
Thus far, the onus is overwhelmingly placed on end users of the internet to achieve anonymity online which limits scalability (think using Tor, installing a VPN, etc.). There are, however, niche circumstances where user anonymity is desirable and / or necessary. Unfortunately, frictionless solutions using the browser as a platform for end-to-end encryption (think Proton Mail) are easily compromised because of problems associated with trust in the standard client / server model . A market opportunity is available, thus, in the online microservice ecosystem to provide frictionless anonymity services to a company’s end users. In addition to exploring this opportunity, Layer 8 also serves simultaneously as an R&D foundation for future MAIC projects.

## Key Performance Indicator
Realize a complex production system using modern, but already available, technologies to produce a novel microservice. In other words, successfully applied secondary, but not primary, research and development. For example, greenfield implementations of vetted cryptographic primitives in web assembly language.    

## What
Layer 8 is designed to be a scalable internet service platform that enables end-to-end encryption via the browser. This, in turn, enables a user’s true identity to be stripped from their content choices. Very broadly, the proof-of-concept, works as follows: 
1.	A content delivery network serves an in-browser module which exposes the L8 global object. 
2.	The developer invokes client side methods to build an encrypted tunnel to their backend through the Layer 8 reverse proxy ( e.g., L8.registerCitizenship(…) ).
3.	Through an algorithm inspired by OAuth2.0, Layer 8 establishes an encrypted tunnel using JSON Web Tokens. 
4.	By acting as an HTTP reverse proxy, Layer 8 strips identifying header information requests and replaces it with custom metadata suitable for public aggregation.
5.	Metadata can be collated by the Service Provider without fear of deanonymizing their users. 
6.	Because an encrypted tunnel has been established, Layer 8 is ignorant of a user’s content choices whereas the Service Provider is ignorant of the user’s true identity. Only the end user, according to the scheme proposed, is aware of both their true identity and their content choices. 

To succeed, Layer 8 must solve the fundamental problem(s) associated with trusting an unknown server to deliver an uncompromised application to the client with every new page load (see Kobeissi, N. (2021) An Analysis of the ProtonMail Cryptographic Architecture. Retrieved from https://eprint.iacr.org/2018/1121.pdf). To users, it will provide a free, anonymizing, authentication service. 

Layer 8 must be open source so that it can be vetted by the larger internet security community. Crowd scrutiny will be fundamental to establishing its brand identity as a viable, trusted, third party. It should be noted that, technologically, it is already possible to build Layer 8. In fact, the working proof of concept that I wrote using Type Script is available on Github (github.com/satsite13/Layer8.git). Risky primary research and development should not be necessary to realize the project and will be avoided.
