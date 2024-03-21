## Udacity Go CRM Backend

Course Link: https://www.udacity.com/course/golang--cd11970

Given the times that I have worked with Go in a professional capacity have been sporadic, I started the Udacity course to make sure I didn't miss out on any of the fundamentals.  I had never built a Go project from scratch and was always stepping into projects that were already mature.

So, with that, I wanted to stick to Go standard library as much as possible for this project.  I know I could have reached for a lot of different packages to make my code more simplistic, but I wanted to be sure I understood the standard library concepts before getting into the abstractions too heavily.

Following https://github.com/golang-standards/project-layout for project layout recommendations.

### External Packages Used

- `air`: Needed for live reloading.
- `pgx`: Needed for interfacing with postgres.
- `godotenv`: Needed for .env file import.

I would definitely reach for a HTTP library to make building out of these REST API pieces easier and more concise.  I would look at Gin or Chi mainly.

I am not interested in libraries that are bringing the "Javscript Problem" along like Fiber or Echo.  

### `air` Setup Note

Needed `alias air='$(go env GOPATH)/bin/air'` to get `air init` to run properly.

### Running Project

1. Import the `seed_data.sql` in your local Postgres client.
2. Create your .env file with `cp .env.example .env` .
3. Set your `POSTGRES_DB_URL` in the `.env`.
4. Run `go get .` to download all external packages.
5. Run `./cmd/udacity-go-crm-backend` from the root of the project.
6. To build the binary, run `go build -o ./bin ./cmd/udacity-go-crm-backend`.