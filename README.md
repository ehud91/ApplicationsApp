# ApplicationsApp

This project is an API application handle an applications:

Application Api:

- Add new Application
- Get all Applications

## Getting Started

### Requirements

- [go](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [React.js](https://reactjs.org/)

### Clone

To get started you can simply clone this repository using git:

```shell
$ git https://github.com/ehud91/ApplicationsApp.git
$ cd ApplicationsApp
```

### Application Installation

- Clone the repo by using git clone.
- Add the following application configurations to your local environment variables in your machine

```bash
# APP_SERVICE_PORT
# PG_URL
# PG_HOST
# PG_PASSWORD
# PG_PORT
# PG_DBNAME
# PG_USER
```

- Run the following 

```bash
# install dependencies
go get

# run Server
go run main.go
```

### Database Installation

- Install [PostgreSQL](https://www.guru99.com/download-install-postgresql.html). You can find how to install the PostgreSQL in here.

### Create the Database and tables

- Please run the 2 scripts under -> /ApplicationsApp/sql/InitializeDb.sql

### Test the application

- Go to your [Postman](https://www.postman.com/downloads/) and try all the services URI's 

### API Endpoints

Contains 2 main controllers api's:

- http://localhost:8123/addnewapplication - Add new Application

```json
{
    "app_name": "app 2",
    "app_key": "0cde0459-a6ba-4f47-bc28-7893eafe112f"
}
```

- http://localhost:8123/getapplication - get application by name / key

```json
{
	"name": "Contacts Application",
	"key": "f5847be3-6ec5-4b24-b61a-93a487c7b2fe"
}
```
