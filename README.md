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