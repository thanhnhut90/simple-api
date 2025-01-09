# simple-api

To add a string:    `curl -sXPOST http://localhost:8080/api/create -d '{"value": "string"}'`

To update a string: `curl -sXPUT http://localhost:8080/api/update -d '{"id": 2, "value": "new string"}'`

To delete a string: `curl -sXDELETE "http://localhost:8080/api/delete?id=1"`

To get all string:  `curl -s http://localhost:8080/api`


