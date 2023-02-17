# COMPANIES API MICROSERVICE

   - This is a REST API microservice which handles companies. Is exposes REST API to create, get, patch and delete a company.
   - Postgres database is used for persistance.
   - Only authenticated users will be able to access create, patch and delete API
   - Get company API in not protected.


### Project Tree
```
.
├── Dockerfile
├── Makefile
├── README.md
├── api
│   ├── constants
│   │   └── constants.go
│   ├── controller
│   │   ├── company.go
│   │   └── login.go
│   ├── database
│   │   └── db.go
│   ├── dto
│   │   └── dto.go
│   ├── errors
│   │   └── errors.go
│   ├── logging
│   │   └── logger.go
│   ├── middleware
│   │   └── auth.go
│   ├── models
│   │   └── models.go
│   ├── repository
│   │   ├── mocks
│   │   │   └── mock_repository.go
│   │   ├── repository.go
│   │   └── repository_test.go
│   ├── router
│   │   └── router.go
│   ├── service
│   │   ├── auth.go
│   │   ├── company.go
│   │   ├── company_test.go
│   │   ├── login.go
│   │   └── mocks
│   │       └── mock_company.go
│   └── utils
│       └── utils.go
├── companies
├── db-migration
│   ├── Dockerfile
│   ├── V1__create_table_companies.sql
│   └── flyway.conf
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── e2e
│   ├── docker-compose.yml
│   └── e2e_test.go
├── go.mod
├── go.sum
└── main.go
```

### Requirements:
1. go 1.20+ 
2. docker
3. docker compose
4. golangci-lint   [ref](https://golangci-lint.run/usage/install/)


### Format:
`$ make format`

### Lint:
`$ make lint`

### Unit test:
`$ make test-unit`

### E2E test:
`$ make test-e2e`

### To start all services (db, migration job, webserver) and to test APIs:
`$ make start-all-services`

### To stop all services (db, migration job, webserver) and to test APIs:
`$ make stop-all-services`



## **Swagger**

 1. APIs are listed in swagger endpoint `http://localhost:8080/api/company/v1/swagger/index.html#/`
 2. This microservices currently supports only static authentication which has only one user, the credentials are
 ```
 {
    "email": "admin@company.com",
    "password": "password"
}
 ``` 

## Documentation for API Endpoints


All URIs are relative to *http://127.0.0.1:8080/api/company/v1*


### /api/v1/company

#### POST
##### Summary:

create company

##### Description:

creation of new company

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| CreateCompany | body | request body | Yes | [models.Company](#models.Company) |
| authorization | header | string | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Created | [models.Company](#models.Company) |
| 400 | Bad Request | [errors.ErrorResponse](#errors.ErrorResponse) |
| 403 | Forbidden | [errors.ErrorResponse](#errors.ErrorResponse) |
| 500 | Internal Server Error | [errors.ErrorResponse](#errors.ErrorResponse) |

### /api/v1/company/:id

#### GET
##### Summary:

get company

##### Description:

get company info by ID

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Company](#models.Company) |
| 400 | Bad Request | [errors.ErrorResponse](#errors.ErrorResponse) |
| 403 | Forbidden | [errors.ErrorResponse](#errors.ErrorResponse) |
| 500 | Internal Server Error | [errors.ErrorResponse](#errors.ErrorResponse) |

#### DELETE
##### Summary:

delete a company

##### Description:

delete company by ID

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| authorization | header | string | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | string |
| 400 | Bad Request | [errors.ErrorResponse](#errors.ErrorResponse) |
| 403 | Forbidden | [errors.ErrorResponse](#errors.ErrorResponse) |
| 500 | Internal Server Error | [errors.ErrorResponse](#errors.ErrorResponse) |

#### PATCH
##### Summary:

update a company

##### Description:

update company by ID

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| updateReq | body | request body | Yes | [models.Company](#models.Company) |
| authorization | header | string | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Company](#models.Company) |
| 400 | Bad Request | [errors.ErrorResponse](#errors.ErrorResponse) |
| 403 | Forbidden | [errors.ErrorResponse](#errors.ErrorResponse) |
| 500 | Internal Server Error | [errors.ErrorResponse](#errors.ErrorResponse) |

### Models


#### errors.ErrorResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| error_code | string |  | No |
| error_message | string |  | No |
| status | integer |  | No |

#### models.Company

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| amount_of_employees | integer |  | No |
| ceated_at | string |  | No |
| description | string |  | No |
| id | string |  | No |
| name | string |  | No |
| registered | boolean |  | No |
| type | string |  | No |
| updated_at | string |  | No |








## Deployment

### PreRequisites:

1. To connect this microservice with external DB, export below env variables.
```
DB_HOST=<host>
DB_PORT=<port>
DB_USER=<user>
DB_PWD=<password>
DB_NAME=<name>
DB_SSL_MODE=<ssl mode>
```


### How to run:

1. Clone the repo

	- $ git clone https://github.com/kumareswaramoorthi/companies.git

2. Navigate to project directory 

	- $ cd  companies

3. Build the application by following command

	- $ go build -o companies main.go

4. Make sure you have a postgres database up and running and env variables are set

5. Run the application by the following command 

	- $ ./companies


Alternatively, using docker,


1. Clone the repo

	- $ git clone https://github.com/kumareswaramoorthi/companies.git

2. Navigate to project directory 

	- $ cd companies

3. Build the docker image by following command

	- $ docker build -t companies:1.0 .
	
4. Make sure you have a postgres database up and running and env variables are set

5. Run the application by the following command 

	- $ docker run -p 8080:8080 companies:1.0

