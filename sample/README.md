# Sample Project
REST API sample used sqlw

The samaple app has three apis below
* GET  http://localhost:8080/items
* GET  http://localhost:8080/items/{item_id}
* POST http://localhost:8080/items

## Running the sample app
1. Executes docker-compose command to launch mysql
2. Run the go command to launch sample app

```go
$ cd mysql
$ docker-compose up -d
$ cd ../sample
$ go run main.go
```

## Layer and package structure of the sample app
Sample app has three layers like port, usecase, domain.
* port
  * includes http request handlers and more
* usecase
  * includes business logics
* domain
  * includes domain models
  
```
main.go
├─ port
│  ├─ handler
│  └─ repository(implementation)
├─ usecase
│  ├─ interactor
│  └─ repository(interface)
└─ domain
   └─ model

```

Each layers dependency is below
```
port -> usecase -> domain
```

## API Specifications

### Get all items
http://localhost:8080/items

#### HTTP Method
GET

#### URL Params
None

#### Response

```
[
  {
    "id": "48c0c6c2-268c-11eb-850f-acde48001122",
    "name": "bar"
  },
  {
    "id": "4c3e0fda-268c-11eb-850f-acde48001122",
    "name": "foo"
  }
]
```

### Get an item
http://localhost:8080/items/{item_id}

#### HTTP Method
GET

#### URL Params
item_id string

#### Response

```
{
  "id": "48c0c6c2-268c-11eb-850f-acde48001122",
  "name": "bar"
}
```

### Create an item
http://localhost:8080/items

#### HTTP Method
POST

#### Request Body
name string

#### Response

```
{
  "id": "48c0c6c2-268c-11eb-850f-acde48001122",
  "name": "baz"
}
```
