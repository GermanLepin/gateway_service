# Information about the project X

I have designed the microservices interaction, and there are two services:
- gateway service is an entry point into applications;
- authentication-service is the service checking authorization and authentication;


What I already did:
1. user creation handler `v1/api/user/create`
2. user logging handler `v1/api/user/login`
3. user deletion handler `v1/api/user/delete/{uuid}`
4. user fetching handler `v1/api/user/fetch/{uuid}`
5. start the database and the servers using only one command
6. migrations

In the `project` repository can be found all files with basic commands and the main `docker-compose.yml`

To start all services in Docker, you need to clone this repository to your local computer:
```
git clone git@github.com:GermanLepin/payment_service.git
```

Go to `project`
```
cd project/
```

Let's start all services and databases with the command:
```
make up_build
```

// TODO Swagger
# gateway service API

Implemented a creation method. Accepts a user name, a user last name, a user phone number, a user email, and a user password.

| Key              | Data type | Description         | Example
|------------------|-----------|---------------------|--------------------- |
| first_name       | string    | a user first name   | John                 |
| last_name        | string    | a user last name    | Smith                |
| password         | string    | a user password     | 1234qwer             |
| email            | string    | a user email        | john@gmail.com       |
| phone            | int       | a user phone number | 4912345678901        |
| user_type        | string    | a user type         | admin/user           |

POST method.

    https://localhost:9999/v1/api/user/create

*Add to the request body (JSON format):*
```
{
	"first_name":"Jonn",
	"last_name": "Smith",
	"password":"1234qwer",
	"email":"john@gmail.com",
	"phone": 4912345678901,
	"user_type":"user"
}
```



Implemented a login method. Accepts a user email, and a user password.

| Key              | Data type | Description         | Example
|------------------|-----------|---------------------|--------------------- |
| email            | string    | a user email        | john@gmail.com       |
| password         | string    | a user password     | 1234qwer             |

POST method.

    https://localhost:9999/v1/api/user/create

*Add to the request body (JSON format):*
```
{
	"email": "john@gmail.com",
	"password":"1234qwer"
}
```

*Responce from the login request  (JSON format):*
```
{
	"session_id": "01095789-72f9-46c4-b3fb-e74c0b47d85b",
	"is_bloked": false,
	"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZ21haWwuY29tIiwiZXhwIjoxNzA5MDY5NDUyLCJ1c2VyX2lkIjoiNmQ5YmUyOWUtYWI4Yi00NmVjLTg4ZDctNzgwYzE4MDM3MGE2In0.zsUzuGor3x1EtYAZ9rFN919VGtNLdBlyxl_Agti0Xqk",
	"access_token_expires_at": "2024-02-27T21:30:52.136613036Z",
	"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZ21haWwuY29tIiwiZXhwIjoxNzA5MzI3NzUyLCJ1c2VyX2lkIjoiNmQ5YmUyOWUtYWI4Yi00NmVjLTg4ZDctNzgwYzE4MDM3MGE2In0.Cp6V9wM1CTER33Itac0bNgfPKrlVdgXhZ765TQmoK9Y",
	"refresh_token_expires_at": "2024-03-01T21:15:52.136651535Z",
	"user_id": "6d9be29e-ab8b-46ec-88d7-780c180370a6"
}
```

