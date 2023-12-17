# Welcome to the Riddles API

## A curated collection of the best riddles out there

A pure-Go API using the database/sql package.

### Features:

- RESTful architecture
- HATEOAS-driven navigation
- Error handling
- Comprehensive non-intrusive logging
- IP-based access control for selected methods

### Available methods

| Operation                    | URI                                          | Method | Status                          | Status Code             | Availability |
| ---------------------------- | -------------------------------------------- | ------ | ------------------------------- | ----------------------- | ------------ |
| All riddles                  | /api/riddles         | GET    | OK<br>Internal Server Error     | 200<br>500              | public       |
| Random riddle                | /api/riddles/random  | GET    | Success<br>Internal Server Error| 200<br>500              | public       |
| Specific riddle              | /api/riddles/{id}    | GET    | OK<br>Bad Request<br>Not found  | 200<br>400<br>404       | public       |
| Post riddle                  | /api/riddles         | POST   | Success<br>Bad Request<br>Internal Server Error| 201<br>400<br>500 | public |
| Delete riddle                | /api/riddles/{id}    | DELETE | OK<br>Bad Request<br>Not Found<br>Internal Server Error<br>Forbidden | 200<br>400<br>404<br>500<br>403 | restricted |
| Update riddle                | /api/riddles/{id}    | PATCH  | OK<br>Bad Request<br>Not Found<br>Internal Server Error<br>Forbidden | 200<br>400<br>404<br>500<br>403 | restricted |

#### Request body example for Update Riddle:

```json
{
  "riddle": "What am I?",
  "solution": "riddle",
  "username": "mr smith", //optional
  "user_email": "mr@smith.com" //optional
}
```


### Special Methods

| Operation                        | URI                                            | Method | Status                          | Status Code                   | Availability |
| -------------------------------- | ---------------------------------------------- | ------ | ------------------------------- | ----------------------------- | ------------ |
| Riddle with DALLE generated image| /api/riddles/images/{id}| GET    | OK<br>Bad Request<br>Not Found<br>Internal Server Error<br>Forbidden | 200<br>400<br>404<br>500<br>403 | restricted   |

#### Description:
Generates an image with the riddle as a prompt, and returns the riddle together with the image URL. In the backend, it retrieves the image and stores it in an image table.

#### Optional request body for additions to the prompt:
```json
{"style": "impressionist"}