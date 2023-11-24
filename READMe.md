# FIO_APP

## Описание

Сервис для сбора, обогащения и хранения данных о пользователях

## Структура 

```text
│   docker-compose.yml        
│   Dockerfile
│   go.mod
│   go.sum
│   READMe.md
├───api
│       api.json
│
├───cmd
│   └───app
│           main.go
│
├───logs
│       all.log
│
└───pkg
    ├───adapters
    │   ├───apis
    │   │       apis.go
    │   │
    │   └───producer
    │           producer.go
    │
    ├───dtos
    │       person.go
    │
    ├───errs
    │       errs.go
    │
    ├───logging
    │       logging.go
    │
    ├───models
    │       person.go
    │
    ├───ports
    │   ├───consumer
    │   │       consumer.go
    │   │
    │   ├───graph
    │   │       fio_type.go
    │   │       mutation.go
    │   │       query.go
    │   │       server.go
    │   │
    │   └───rest
    │           handlers.go
    │           router.go
    │           server.go
    │
    ├───repo
    │       cache.go
    │       permanent.go
    │       repo.go
    │
    ├───service
    │   │   service.go
    │   │
    │   └───utils
    │           validation.go
    │
    └───storage
        └───database
            ├───postgres
            │       storage.go
            │
            └───redis
                    storage.go

```

## Бизнес-логика

ФИО читается из FIO (топик кафки), проверяется на валидность (отсутствие спец.
символов, обязательно наличие фамилии и имени). Затем выполняются три запроса к
внешнему API, которые по имени определяют возраст, пол и национальность. В
структуре ФИО обновляются значения, и она отправляется в БД. Если ФИО
некорректно, структура отправляется в топик FIO_FAILED. При добавлении ФИО
через REST или GraphQL все поля должны быть заполнены сразу.

Через REST и GraphQL поддерживаются все CRUD-операции, также можно выбрать
список пользователей с фильтрами по всем полям (фамилия, имя, отчество,
возраст, пол, национальность) и пагинацией.

Также для ускоренного получения ФИО реализован кеш, в который на время попадают
все новые или только что обновлённые ФИО.

## Используемые технологии

* go 1.21
* Apache Kafka
* PostgreSQL
* Redis
* GORM
* Gin Web Framework
* GraphQL

## Запуск приложения

Для запуска приложения необходимо задать переменные окружения в корне проекта в файле `.env`

```text
HOST=
PORT=

GRAPHQL_HOST=
GRAPHQL_PORT=

# db
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_HOST=
DB_PORT=

#redis
REDIS_HOST=
REDIS_PORT=
REDIS_PASSWORD=
REDIS_DB=

#kafka
ADDRESS=
```

и выполнить следующую команду:

```shell
docker-compose up --build
```

## Формат REST-запросов

### Добавление пользователя

* Метод: `POST`
* Эндпоинт: `http://localhost:8080/people`
* Формат тела запроса:

```json
{
    "name": "Ivan", 
    "surname": "Ivanov", 
    "patronymic": "Ivanovich",
    "age": 39,
    "gender": "male",
    "nation": "RU"
}
```

* Формат ответа:

```json
{}
```

### Получение списка пользователей

* Метод: `GET`
* Эндпоинт: `http://localhost:8080/people`
* Формат запроса:

```json
{
    "offset": 0,
    "limit": 1,
    "gender": "",
    "nation": ""
}
```

* Формат ответа:

```json
{
    "data": [
        {
            "id": 1,
            "name": "Ivan",
            "surname": "Ivanov",
            "patronymic": "Ivanovich",
            "age": 39,
            "gender": "male",
            "nation": "RU"
        }
    ]
}
```

### Обновление пользователя

* Метод: `PUT`
* Эндпоинт: `http://localhost:8080/people/:id`
* Формат запроса:

```json
{
    "name": "Ivan",
    "surname": "Ivanov",
    "patronymic": "Ivanovich",
    "age": 40,
    "gender": "male",
    "nation": "RU"
}
```

* Формат ответа:

```json
{}
```

### Удаление пользователя

* Метод: `DELETE`
* Эндпоинт: `http://localhost:8080/people/:id`
* Формат ответа:

```json
{}
```

### Отправка сообщения в топик кафки

* Метод: `POST`
* Эндпоинт: `http://localhost:8080/kafka/produce`
* Формат запроса:

```json
{
  "name": "Arsen",
  "surname": "Yarullin",
  "patronymic": "Rustemovich"
}
```

* Формат ответа:

```json
{}
```

## Формат GraphQL-запросов

***Все запросы выполняются по адресу `http://localhost:8081`***

### Добавление пользователя

* Формат запроса:

```text
mutation AddFio {
    addFio(
        name: "Ivan"
        surname: "Ivanov"
        patronymic: "Ivanovich"
        age: 20
        gender: "male"
        nation: "RU"
    )
}
```

* Формат ответа:

```json
{
  "data": {
    "addFio": true
  }
}
```

### Получение пользователя

* Формат запроса:

```text
query GetFioById {
    getFioById(id: 2) {
        id
        name
        surname
        patronymic
        age
        gender
        nation
    }
}
```

* Формат ответа:

```json
{
  "data": {
    "getFioByID": {
      "age": 54,
      "gender": "male",
      "id": null,
      "name": "Ivan",
      "nationality": "HR",
      "patronymic": "",
      "surname": "Ivankov"
    }
  }
}
```

### Получение списка пользователей

* Формат запроса:

```text
query GetFios {
    getFios(offset: 0, limit: 1) {
        age
        gender
        id
        name
        nation
        patronymic
        surname
    }
}
```

* Формат ответа:

```json
{
    "data": {
        "getFios": [
            {
              "age": 54,
              "gender": "male",
              "id": null,
              "name": "Ivan",
              "nationality": "HR",
              "patronymic": "",
              "surname": "Ivankov"
            }
        ]
    }
}
```

### Обновление пользователя

* Формат запроса:

```text
mutation UpdateFio {
    updateFio(
        gender: "male"
        nation: "RU"
        id: 2
        name: "Ivan"
        surname: "Ivankov"
        patronymic: "Romanovich"
        age: 40
    )
}

```

* Формат ответа:

```json
{
    "data": {
        "updateFio": true
    }
}
```

### Удаление пользователя

* Формат запроса:

```text
mutation DeleteFio {
    deleteFio(id: 2)
}
```

* Формат ответа:

```json
{
    "data": {
        "deleteFio": true
    }
}
```



