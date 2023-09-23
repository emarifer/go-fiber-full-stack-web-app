# Golang Fiber Full-stack Web App

### Simple web app + API with JWT authentication and Golang Html template rendering.

---

Besides the obvious prerequisite of having Go!, Docker and Nodejs installed, we need to do the following to create the database (2 Docker containers: one for the DB itself, and one for the DB administrator UI [pgAdmin]).

```
$ docker compose up -d  # run at the root of the project.
```

This will create 2 Docker containers (in addition to mounting a volume at the root of the project that will contain the container data) and we will start them with the command:

```
$ docker start go-postgres-db go-pgAdmin  # "docker stop go-postgres-db go-pgAdmin" to stop them.
```

---

Since Tailwindcss is used to style templates, dependencies need to be installed with NPM:

```
$ npm i
```

and run:

```
$ npm run watch-css # in development or "npm run build-css-prod" to generate the minified css for production.
```

Finally, in development, we will only have to execute being located at the root of the project:

```
$ air # Or «go build -ldflags="-s -w" -o ./cmd/main/main ./cmd/main/main.go» in production
```

---

### Test the API using the cURL command.

###### Using the cURL command to test the REST API. Run the following example commands, obviously with your own data:

- Login:

```
curl -v -X POST http://localhost:5500/api/auth/login -d '{ "email": "enrique@enrique.com", "password": "123456" }' -H "content-type: application/json" | json_p
```

- SignUp:

```
curl -v -X POST http://localhost:5500/api/auth/register -d '{ "name": "Enrique Marín",  "email": "enrique@enrique.com", "password": "123456", "passwordConfirm": "123456" }' -H "content-type: application/json" | json_p
```

- Logout:

```
curl --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTUzODg2ODUsImlhdCI6MTY5NTM4NTA4NSwibmJmIjoxNjk1Mzg1MDg1LCJzdWIiOiI3NmE1MGIxNC03MTg2LTQ3YTgtYmQzYi1iYTZhMmY3M2JkMzcifQ.CWOeSLqdmWmP9rgIGgRdS_eNxGCE8fjiIGvL6X6S4yg" -v http://localhost:5500/api/auth/logout -H "content-type: application/json" | json_pp
```

- Get data of the logged in user:

```
curl --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTUzODg2ODUsImlhdCI6MTY5NTM4NTA4NSwibmJmIjoxNjk1Mzg1MDg1LCJzdWIiOiI3NmE1MGIxNC03MTg2LTQ3YTgtYmQzYi1iYTZhMmY3M2JkMzcifQ.CWOeSLqdmWmP9rgIGgRdS_eNxGCE8fjiIGvL6X6S4yg" -v http://localhost:5500/api/users/me -H "content-type: application/json" | json_pp
```

