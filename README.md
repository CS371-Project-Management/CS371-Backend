## Getting Started
build image
```bash
docker compose build
```
run docker container
```bash
docker-compose up -d
```
## Example of creating a migration
```bash
migrate create -ext sql -dir db/migrations create_users_table
```

## DB check manual
```bash
docker exec -it cs371-db mysql -u root -p cs371db
```

## migration use
```bash
migrate -path db/migrations -database "mysql://root:my-secret-pw@tcp(127.0.0.1:3306)/cs371db" + (up / down)
```