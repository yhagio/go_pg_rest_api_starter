# Go + Postgres REST API

Reddit like REST API

### TODOs
- [x] User signup
- [x] User login
- [x] Test protected route (/secret) that requires user-logged-in
- [x] Routes path `/api/`
- [x] Signup notification email
- [x] Forgot password + reset password via email
- [x] Post (CRUD)

- [ ] Stripe payment integration
- [ ] Google OAuth
- [ ] Deployment (i.e. Heroku, Digital Ocean, AWS, GCP)
- [ ] Admin can delete any posts


### curl commands

```bash
curl -X "POST" "http://localhost:3000/api/signup" -H 'Content-Type: application/json; charset=utf-8' -d $'{"username":"alice", "email":"alice@example.com", "password":"password123"}'

curl -X "POST" "http://localhost:3000/api/login" -H 'Content-Type: application/json; charset=utf-8' -d $'{"email":"alice@example.com", "password":"password123"}'

curl -H "Authorization: Bearer <JWT_TOKEN>" -H 'Content-Type: application/json; charset=utf-8' http://localhost:3000/api/me -w "\n"

curl -X "POST" "http://localhost:3000/api/forgot_password" -H 'Content-Type: application/json; charset=utf-8' -d $'{"email":"alice@example.com"}'

curl -X "POST" "http://localhost:3000/api/update_password?token=<PROVIDED_TOKEN>" -H 'Content-Type: application/json; charset=utf-8' -d $'{"email":"alice@example.com", "password":"updatedPassword"}'


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
export APP_PORT=3000
```