# TextSave
Web server that stores 5000 characters under a small code for sharing things. Codes expire after 30 minutes.

### Environmental Variables
- `REDIS_ADDRESS`
    - URI with port.
    - e.g. `localhost:6379`
- `REDIS_PASS`
    - Empty for no password
- `REDIS_DB`
    - DB number
- `PORT`
    - Default `8080`

### Run via Docker
```bash
docker build -t text-save .
docker run -d -p 8080:8080 --env-file .env text-save
```
