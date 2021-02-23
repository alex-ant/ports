# ports

### Run tests

The following script will run unit tests in a docker container spinning up the database:

`./run-tests.sh`

### Run locally

Use docker-compose to run the client, the server and the database instances on your machine:

`docker-compose build && docker-compose up`

Once lauched, the PostgreSQL database running on `localhost:5432` will be populated with  
the contents of `ports.json` file. Database name is `ports`.  
Then the REST API server running on port `8080` will become available for fetching port info.
The endpoint returns all the available port data.

Sample HTTP request:

```
curl -X GET 'http://localhost:8080/fetch-port-info'
```

Sample response:

```json
{
    "msg": "ok",
    "port-info": [
        {
            "id": "CAQUE",
            "name": "Quebec",
            "city": "Quebec",
            "country": "Canada",
            "coordinates": [
                -71.2428,
                46.803284
            ],
            "province": "Quebec",
            "timezone": "America/Toronto",
            "unlocs": [
                "CAQUE"
            ],
            "code": "80109"
        },
        {
            "id": "CUHAV",
            "name": "La Habana",
            "city": "La Habana",
            "country": "Cuba",
            "coordinates": [
                -82.35,
                23.12
            ],
            "province": "Ciudad de La Habana",
            "timezone": "America/Havana",
            "unlocs": [
                "CUHAV"
            ],
            "code": "23944"
        }
    ],
    "status": 200
}
```
