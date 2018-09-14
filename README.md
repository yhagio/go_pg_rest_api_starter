# Go + Postgres REST API


### Demo API (Heroku)

https://powerful-oasis-83075.herokuapp.com/api/
(Mailgun is disabled for demo though)


### CRUD REST API example / starter

- [Go - language](https://golang.org/)
- [Postgres - Databse](https://www.postgresql.org/)
- [gorm - Go ORM, Migration tasks](https://github.com/jinzhu/gorm)
- [Gorilla Mux - Router](https://github.com/gorilla/mux)

### TODOs
- [x] User Sign Up
- [x] User Login
- [x] Protected route (/api/me) User Profile
- [x] Signup notification email (Mailgun)
- [x] Forgot password + reset password via email (Mailgun)
- [x] Post (CRUD)

Future consideration (idea)
- [ ] Stripe payment integration
- [ ] Google, Facebook, etc OAuth
- [ ] Deployment (i.e. Heroku, Digital Ocean, AWS, GCP)
- [ ] Admin tasks
- [ ] Cron jobs


### Available API and curl commands

```bash
# Signup /api/signup
curl -X "POST" "http://localhost:3000/api/signup" -H 'Content-Type: application/json; charset=utf-8' -d $'{"username":"alice", "email":"alice@example.com", "password":"password123"}'

# Login /api/login
curl -X "POST" "http://localhost:3000/api/login" -H 'Content-Type: application/json; charset=utf-8' -d $'{"email":"alice@example.com", "password":"password123"}'

# User profile /api/me
curl -H "Authorization: Bearer <JWT_TOKEN>" -H 'Content-Type: application/json; charset=utf-8' http://localhost:3000/api/me -w "\n"

# Forgot password /api/forgot_password
curl -X "POST" "http://localhost:3000/api/forgot_password" -H 'Content-Type: application/json; charset=utf-8' -d $'{"email":"alice@example.com"}'

# Update / Reset apssword /api/update_password
curl -X "POST" "http://localhost:3000/api/update_password?token=<PROVIDED_TOKEN>" -H 'Content-Type: application/json; charset=utf-8' -d $'{"email":"alice@example.com", "password":"updatedPassword"}'

# Create a post /api/posts
curl -X "POST" "http://localhost:3000/api/posts?token=<PROVIDED_TOKEN>" -H 'Content-Type: application/json; charset=utf-8' -d $'{"title":"Hello World", "description":"Hello everyone, this is my first post"}'

# Fetch a post /api/posts/:id
curl -H 'Content-Type: application/json; charset=utf-8' http://localhost:3000/api/posts/1

# Update a post /api/posts/:id/update
curl -X "PUT" "http://localhost:3000/api/posts/1/update?token=<PROVIDED_TOKEN>" -H 'Content-Type: application/json; charset=utf-8' -d $'{"title":"Hello Again", "description":"Hello everyone, this is modified"}'

# Delete a post /api/posts/:id/delete
curl -X "DELETE" "http://localhost:3000/api/posts/1/delete?token=<PROVIDED_TOKEN>" -H 'Content-Type: application/json; charset=utf-8'
```

### Environmental Variables (example)
```bash
export MAILGUN_APIKEY=key-b9cc7f4l39ee5ca2e7cc6cb89328afd
export MAILGUN_PUBLIC_KEY=pubkey-123439481f695175063232234hjsdf9
export MAILGUN_DOMAIN=example.mailgun.org

export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=dev_user
export DB_PASSWORD=password123
export DB_NAME=my_dev_db

export JWT_SIGN_KEY=123JWT-key
export HAMC_KEY=super-secret-key
export PEPPER=super-pepper-key

export APP_ENV=development
export PORT=3000
```