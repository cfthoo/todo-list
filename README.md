<br />
<h1 align="center">todo-list</h3>

<p align="center">

  <br/>

This is a simple todo-list APIs created with Go. It allow users to login with Google , Facebook and Github.
User can perform CRUD with the todo-list APIs such as CREATE task , LIST tasks ,GET task by Id ,UPDATE and DELETE task.

## Usage

### Docker Compose

This will get postgres and todo-list image from docker hub and run the container.
There is a .env file which uses to configure the DB and oauth2 info.
Please edit the .env if you would like to provide your own configuration

```bash
$ docker-compose up
```

You shall see the postgres is running , the database is also migrated successfully.
And the http server is listening on port 8080
![Alt text](https://github.com/cfthoo/image/blob/main/up.png)

### Locally

To run locally , please make sure you have the postgres db running.
Go to the root folder of this project and run

```bash
$ go run cmd/main.go
```

### Building Image

You use use below command to build the image
Go to the root folder of this project and run

```bash
$ docker build -t todo-list .
```

or

If you have `make` in your machine , you can run with

```bash
$ make build
```

## Testing

The APIs required you to log in before you can call it.
There will be a token returned once you have log in successfully.

Step for testing

1. Navigate to localhost:8080 once you have start the application with docker-compose up
   ![Alt text](https://github.com/cfthoo/image/blob/main/a.png)

2. Click on any of the auth provider to login. For example you shall see the google login when u select google.
   ![Alt text](https://github.com/cfthoo/image/blob/main/2.png)

3. After successfully login , you shall be prompted authorized and a token.Copy the token for next step.
   ![Alt text](https://github.com/cfthoo/image/blob/main/3.png)

4. You can use curl or postman to call the api

Example of curl

```bash
$ curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODI5NTEwMzR9.5DrTnLAAEqaDAjYspEFEaJW25CqE-EpGRxcOY-CaaCY" http://localhost:8080/todolist
```

Example of postman

a.You will get error message "not authorized" if you call the api without token
![Alt text](https://github.com/cfthoo/image/blob/main/4a.png)

b.Go to Authorization tab and select Bearer Token , then paste in the token in step 3.
![Alt text](https://github.com/cfthoo/image/blob/main/4.png)

5. Call the api and get response.
   ![Alt text](https://github.com/cfthoo/image/blob/main/5.png)

### Unit test

You can run unit test locally with the command below

```bash
$ go test ./... -cover
```

or

If you have `make` in your machine , you can run with

```bash
$ make test
```

## API documentation

**1. Create task for a todolist**  
This Create method creates a task  
PATH: {url}/todolist  
METHOD: POST  
REQUEST PAYLOAD:

```json
{
  "name": "my first task"
}
```

RETURN PAYLOAD:

```json
{
  "id": 1,
  "name": "my first task",
  "created_by": "105301550950520990207",
  "created_at": "2023-05-01T03:16:57.837083Z",
  "modified_at": "2023-05-01T03:16:57.837083Z"
}
```

**2. Get tasks for a todolist**  
This Get method returns all task under a specific user.  
PATH: {url}/todolist  
METHOD: GET  
RETURN PAYLOAD:

```json
{
  "id": 1,
  "name": "my first task",
  "created_by": "105301550950520990207",
  "created_at": "2023-05-01T03:16:57.837083Z",
  "modified_at": "2023-05-01T03:16:57.837083Z"
}
```

**3. Get task by id**  
This Get method returns a task by id.  
PATH {url}/todolist{id}  
METHOD: GET  
RETURN PAYLOAD:

```json
{
  "id": 1,
  "name": "my first task",
  "created_by": "105301550950520990207",
  "created_at": "2023-05-01T03:16:57.837083Z",
  "modified_at": "2023-05-01T03:16:57.837083Z"
}
```

**4. Update task for a todolist**  
This Update method updates a task  
PATH: {url}/todolist  
METHOD: PUT  
REQUEST PAYLOAD:

```json
    {
        "id": 1
        "name": "Updating my first task"
    }
```

RETURN PAYLOAD:

```json
{
  "id": 1,
  "name": "Updating my first task",
  "created_by": "105301550950520990207",
  "created_at": "2023-05-01T03:16:57.837083Z",
  "modified_at": "2023-05-01T03:20:25.747081Z"
}
```

**5. Delete task by id**  
This Delete method delete a task by id.  
PATH {url}/todolist{id}  
METHOD: DELETE  
RETURN PAYLOAD:

```json
{}
```
