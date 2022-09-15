# Gogo Space

## docs/
Includes a Postman collection.

## Environment variables
To override config variables (from app.env file), set the environment values to be different than:
```bash
export PORT=8080
export API_KEY=DEMO_KEY
export CONCURRENT_REQUESTS=5
```

## Makefile
- Run server
    ```bash
    make server
    ```

- Run test
    ```bash
    make test
    ```

- Run server from a docker container
    ```bash
    make serverfromdocker
    ```
