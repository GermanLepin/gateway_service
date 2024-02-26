# Information about the gateway service

I have designed the microservices interaction, and there are three services:
- gateway-service is an entry point into applications;
- payment-service is the main service where I wrote the main handlers and logic;
- bank-api is an imitation of an external banking service with which we communicate using the REST API. You can check this out on the scheme at the very bottom.

What I already did:
1. user creation handler `v1/user/create`
2. user logging handler `v1/user/login`
3. user deletion handler `v1/user/delete/{uuid}`
4. user fetching handler `v1/user/fetch/{uuid}`
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

type User struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Password     string
	Email        string
	Phone        int
	UserType     string
	JWTToken     string
	RefreshToken string
}


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

    http://localhost:9000/v1/user/create


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


