# betfr-be-revamp
Backend for betting app with friends

## Setup


In the `config` folder, add a file named `default.json` with the following contents: 

```json
{
    "databaseURL": "[PUT_DATABASE_URL_HERE]",
    "secretKey": "[PUT_SECRET_KEY_HERE]",
    "domain": "localhost",
    "port": "8000",
    "debug": true,
    "originFE": "http://localhost:3000",
    "migrationsOnly": false
}
```

Then run `go run .` from the project's root directory to start up the backend. 

API endpoints can be tested with Postman, or used with an instance of the frontend (make sure to set the `originFE` correctly). 
