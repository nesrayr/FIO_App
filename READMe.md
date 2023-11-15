# FIO_APP

## Описание

Сервис для сбора, обогащения и хранения данных о пользователях

## Структура 

```text
│   docker-compose.yml
│   Dockerfile
│   go.mod
│   go.sum
│   gqlgen.yml
│   READMe.md
│
├───api
│       api.json
│
├───cmd
│   └───app
│           main.go
│
├───graph
│   │   generated.go
│   │   resolver.go
│   │   schema.graphqls
│   │   schema.resolvers.go
│   │
│   └───model
│           models_gen.go
│
└───pkg
    ├───dtos
    │       person.go
    │
    ├───handlers
    │       handlers.go
    │
    ├───kafka
    │       consumer.go
    │       producer.go
    │       utils.go
    │
    ├───router
    │       router.go
    │
    └───storage
        ├───database
        │   ├───postgres
        │   │       storage.go
        │   │
        │   └───redis
        │           storage.go
        │
        ├───models
        │       person.go
        │
        └───person
                storage.go

```

## Бизнес-логика

ФИО читается из FIO (топик кафки), проверяется на валидность (отсутствие спец.
символов, обязательно наличие фамилии и имени). Затем выполняются три запроса к
внешнему API, которые по имени определяют возраст, пол и национальность. В
структуре ФИО обновляются значения, и она отправляется в БД. Если ФИО
некорректно, структура отправляется в топик FIO_FAILED. При добавлении ФИО
через REST или GraphQL все поля должны быть заполнены сразу.

## Используемые технологии

* go 1.21
* Apache Kafka
* PostgreSQL
* Redis
* Gin Web Framework
* GraphQL

## Запуск приложения

Для запуска приложения необходимо задать переменные окружения в корне проекта в файле `.env`
и выполнить следующую команду:

```shell
docker-compose up --build
```

## OpenAPI

```json
{
  "openapi": "3.0.3",
  "info": {
    "title": "FIO_App",
    "version": "1.0.0"
  },
  "paths": {
    "/people": {
      "get": {
        "summary": "Get people",
        "parameters": [
          {
            "name": "limit",
            "in": "query",
            "description": "Number of people to show",
            "required": false,
            "schema": {
              "type": "integer",
              "format": "int32"
            }
          },
          {
            "name": "offset",
            "in": "query",
            "description": "Number of people to skip",
            "required": false,
            "schema": {
              "type": "integer",
              "format": "int32"
            }
          },
          {
            "name": "nationality",
            "in": "query",
            "description": "Filter by nationality",
            "required": false,
            "schema": {
              "type": "string"
            }
          },
          {
            "name": "gender",
            "in": "query",
            "description": "Filter by gender",
            "required": false,
            "schema": {
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "ok",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/PersonDTO"
                }
              }
            }
          },
          "400": {
            "description": "bad request",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/BadRequestResponse"
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create repo",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/PersonDTO"
              }
            }
          },
          "required": true
        },
        "responses": {
          "200": {
            "description": "ok",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object"
                }
              }
            }
          },
          "400": {
            "description": "bad request",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/BadRequestResponse"
                }
              }
            }
          }
        }
      }
    },
    "/people{people_id}": {
      "patch": {
        "summary": "Edit repo by id",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/PersonDTO"
              }
            }
          }
        },
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "type": "object"
                }
              }
            },
            "description": "ok"
          },
          "400": {
            "description": "bad request",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/BadRequestResponse"
                }
              }
            }
          }
        }
      },
      "delete": {
        "summary": "Delete repo by id",
        "responses": {
          "200": {
            "description": "ok",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object"
                }
              }
            }
          },
          "400": {
            "description": "bad request",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/BadRequestResponse"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "PersonDTO": {
        "required": [
          "name",
          "surname",
          "patronymic",
          "age",
          "gender",
          "nationality"
        ],
        "type": "object",
        "properties": {
          "name": {
            "type": "string"
          },
          "surname": {
            "type": "string"
          },
          "patronymic": {
            "type": "string"
          },
          "age": {
            "type": "integer",
            "format": "int32"
          },
          "gender": {
            "type": "string"
          },
          "nationality": {
            "type": "string"
          }
        }
      },
      "BadRequestResponse": {
        "type": "object"
      }
    }
  }
}
```

