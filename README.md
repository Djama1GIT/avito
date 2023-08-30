# Avito Backend Trainee Test Assignment

[English](#avito-backend-trainee-test-assignment) <br>
[Russian](#russian)

## Installation and Setup

1. Install Docker and Docker Compose if they are not already installed on your system.

2. Clone the project repository:

```bash
git clone https://github.com/Djama1GIT/avito.git
cd avito
```

3. Start the project:

```bash
docker-compose up --build
```
## User Interface

After starting the project, you can access the Swagger user interface at: http://localhost:8000/swagger/index.html.<br>
In Swagger, you can view the available endpoints and their parameters, and also make requests to the API.

## Usage Examples

### Main Task

<p>Method for creating a segment. Accepts the segment's slug (name)</p>

```
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "example"}'
```
> 200: {"slug":"example"} <br>
> 400: {"message":"pq: duplicate key value violates unique constraint \"segments_pkey\""}

<p>Method for deleting a segment. Accepts the segment's slug (name)</p>

```
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "example"}'
```
> 200: {"slug":"example"} <br>
> 400: {"message":"segment with slug example does not exist"}

<p>Method for adding a user to a segment. Accepts a list of segment slugs to add to the user, a list of segment slugs to remove from the user, and the user's ID</p>

```
curl -X PATCH http://127.0.0.1:8000/api/segments/ \
-d '{"user_id": 1, "segments_to_add": ["AVITO_VOICE_MESSAGES"], "segments_to_delete": ["AVITO_PERFORMANCE_VAS", "AVITO_DISCOUNT_30"]}'
```
> 200: {"user_id":1} <br>
> 400: {"message":"error occurred while processing segment to add 'AVITO_VOICE_MESSAGES': pq: duplicate key value violates unique constraint \"user_segments_pkey\""} <br>
> 400: {"message":"error occurred while checking segment to delete existence 'AVITO_PERFORMANCE_VAS': user(1) is not in this segment"}%

<p>Method for retrieving a user's active segments. Accepts the user's ID</p>

```
curl -X GET http://127.0.0.1:8000/api/segments/ -d '{"user_id": 1}'
```
> 200: {"segments":["AVITO_VOICE_MESSAGES"],"user_id":1} <br>
> 200: {"segments":[],"user_id":1} <br>

### Additional Task 1

<p>A method of obtaining the history of a user's entry/exit from a segment over a certain period in csv format.</p>

```
curl -X GET http://127.0.0.1:8000/api/users/history/ -d '{"user_id": 1, "year_month": "2023-08"}'
```
> 200: {"report":"http://localhost:8000/files/reports/user_history_2023-08_1.csv","user_id":1} <br>
> 400: {"message":"invalid YearMonth"} <br>
> 400: {"message":"json: cannot unmarshal string into Go struct field UserHistory.user_id of type int"} <br>

### Additional Task 2

<p>Method for adding a user to a segment with the ability to set TTL</p>

```
curl -X PATCH http://127.0.0.1:8000/api/segments/ -d \
'{"user_id": 1, "segments_to_add": ["AVITO_VOICE_MESSAGES"], "segments_to_add_expiration":"2023-08-29 00:00:00", "segments_to_delete": []}'
```
> 200: {"user_id":1} <br>
> 400: {"message":"parsing time \"2023-08-31 00:00:\" as \"2006-01-02 15:04:05\": cannot parse \"\" as \"05\""} <br>

<p>Method for deleting expired segments of a user</p>
<sup>Automatically called using cron every minute</sup>

```
curl -X DELETE http://127.0.0.1:8000/api/users/expired-segments/
```
> 200: empty result <br>

### Extended Examples
##### You can use the following curls in order and get the same responses

```
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_VOICE_MESSAGES"}' -w '\n'
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_PERFORMANCE_VAS"}' -w '\n'
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_DISCOUNT_30"}' -w '\n'
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_DISCOUNT_50"}' -w '\n'
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_DELIVERY_FEATURE"}' -w '\n'
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_CLOUD_FEATURE"}' -w '\n'
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_UNDEFINED_FEATURE"}' -w '\n'
```

> {"slug":"AVITO_VOICE_MESSAGES"} <br>
> {"slug":"AVITO_PERFORMANCE_VAS"} <br>
> {"slug":"AVITO_DISCOUNT_30"} <br>
> {"slug":"AVITO_DISCOUNT_50"} <br>
> {"slug":"AVITO_DELIVERY_FEATURE"} <br>
> {"slug":"AVITO_CLOUD_FEATURE"} <br>
> {"slug":"AVITO_UNDEFINED_FEATURE"}

```
curl -X PATCH http://127.0.0.1:8000/api/segments/ \
-d '{"user_id": 1, "segments_to_add": ["AVITO_VOICE_MESSAGES", "AVITO_PERFORMANCE_VAS", "AVITO_DISCOUNT_30", "AVITO_DISCOUNT_50"], "segments_to_delete": []}'  -w '\n'

curl -X PATCH http://127.0.0.1:8000/api/segments/ \
-d '{"user_id": 1, "segments_to_add": ["AVITO_DELIVERY_FEATURE", "AVITO_CLOUD_FEATURE", "AVITO_UNDEFINED_FEATURE"], "segments_to_delete": ["AVITO_VOICE_MESSAGES", "AVITO_PERFORMANCE_VAS", "AVITO_DISCOUNT_30"]}' -w '\n'

curl -X GET http://127.0.0.1:8000/api/segments/ -d '{"user_id": 1}' -w '\n'

curl -X PATCH http://127.0.0.1:8000/api/segments/ \
-d '{"user_id": 1, "segments_to_add": [], "segments_to_delete": ["AVITO_DISCOUNT_50", "AVITO_DELIVERY_FEATURE", "AVITO_CLOUD_FEATURE", "AVITO_UNDEFINED_FEATURE"]}' -w '\n'

curl -X GET http://127.0.0.1:8000/api/segments/ -d '{"user_id": 1}' -w '\n'
```
> {"user_id":1} <br>
> {"user_id":1} <br>
> {"segments":["AVITO_CLOUD_FEATURE","AVITO_DELIVERY_FEATURE","AVITO_DISCOUNT_50","AVITO_UNDEFINED_FEATURE"],"user_id":1} <br>
> {"user_id":1} <br>
> {"segments":[],"user_id":1}

```
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_VOICE_MESSAGES"}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_PERFORMANCE_VAS"}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_DISCOUNT_30"}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_DISCOUNT_50"}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_DELIVERY_FEATURE"}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_CLOUD_FEATURE"}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_UNDEFINED_FEATURE"}' -w '\n'
```

> {"slug":"AVITO_VOICE_MESSAGES"} <br>
> {"slug":"AVITO_PERFORMANCE_VAS"} <br>
> {"slug":"AVITO_DISCOUNT_30"} <br>
> {"slug":"AVITO_DISCOUNT_50"} <br>
> {"slug":"AVITO_DELIVERY_FEATURE"} <br>
> {"slug":"AVITO_CLOUD_FEATURE"} <br>
> {"slug":"AVITO_UNDEFINED_FEATURE"} <br>

```
curl -X GET http://127.0.0.1:8000/api/users/history/ -d '{"user_id": 1, "year_month": "2023-08"}'
```

> {"report":"http://localhost:8000/files/reports/user_history_2023-08_1.csv","user_id":1}
>
> user_history_2023-08_1.csv:
> 
> 1,AVITO_VOICE_MESSAGES,добавление,2023-08-29 20:33:11 <br>
> 1,AVITO_PERFORMANCE_VAS,добавление,2023-08-29 20:33:11 <br>
> 1,AVITO_DISCOUNT_30,добавление,2023-08-29 20:33:11 <br>
> 1,AVITO_DISCOUNT_50,добавление,2023-08-29 20:33:11 <br>
> 1,AVITO_DELIVERY_FEATURE,добавление,2023-08-29 20:33:12 <br>
> 1,AVITO_CLOUD_FEATURE,добавление,2023-08-29 20:33:12 <br>
> 1,AVITO_UNDEFINED_FEATURE,добавление,2023-08-29 20:33:12 <br>
> 1,AVITO_VOICE_MESSAGES,удаление,2023-08-29 20:33:12 <br>
> 1,AVITO_PERFORMANCE_VAS,удаление,2023-08-29 20:33:12 <br>
> 1,AVITO_DISCOUNT_30,удаление,2023-08-29 20:33:12 <br>
> 1,AVITO_DISCOUNT_50,удаление,2023-08-29 20:33:12 <br>
> 1,AVITO_DELIVERY_FEATURE,удаление,2023-08-29 20:33:12 <br>
> 1,AVITO_CLOUD_FEATURE,удаление,2023-08-29 20:33:12 <br>
> 1,AVITO_UNDEFINED_FEATURE,удаление,2023-08-29 20:33:12 <br>

```
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_VOICE_MESSAGES"}' -w '\n'
curl -X PATCH http://127.0.0.1:8000/api/segments/ \
-d '{"user_id": 1, "segments_to_add": ["AVITO_VOICE_MESSAGES"], "segments_to_add_expiration":"1970-08-31 00:00:00", "segments_to_delete": []}' -w '\n'
curl -X GET http://127.0.0.1:8000/api/segments/ -d '{"user_id": 1}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/users/expired-segments/ -w '\n'
curl -X GET http://127.0.0.1:8000/api/segments/ -d '{"user_id": 1}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_VOICE_MESSAGES"}' -w '\n'
```

> {"slug":"AVITO_VOICE_MESSAGES"} <br>
> {"user_id":1} <br>
> {"segments":["AVITO_VOICE_MESSAGES"],"user_id":1} <br>
>  <br>
> {"segments":[],"user_id":1} <br>
> {"slug":"AVITO_VOICE_MESSAGES"} <br>

##### Russian
## Установка и настройка

1. Установите Docker и Docker Compose, если они еще не установлены на вашей системе.

2. Клонируйте репозиторий проекта:

```bash
git clone https://github.com/Djama1GIT/avito.git
cd avito
```

3. Запустите проект:

```bash
docker-compose up --build
```

## Пользовательский интерфейс

После запуска проекта вы сможете получить доступ к пользовательскому интерфейсу Swagger по адресу:<br>http://localhost:8000/swagger/index.html.<br>
В Swagger вы сможете просмотреть доступные ручки и их параметры, а также выполнять запросы к API.

## Примеры использования

### Основное задание

<p>Метод создания сегмента. Принимает slug (название) сегмента</p>

```
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "example"}'
```
> 200: {"slug":"example"} <br>
> 400: {"message":"pq: duplicate key value violates unique constraint \"segments_pkey\""}

<p>Метод удаления сегмента. Принимает slug (название) сегмента</p>

```
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "example"}'
```
> 200: {"slug":"example"} <br>
> 400: {"message":"segment with slug example does not exist"}

<p>Метод добавления пользователя в сегмент. Принимает список slug (названий) сегментов которые нужно добавить пользователю, список slug (названий) сегментов которые нужно удалить у пользователя, id пользователя</p>

```
curl -X PATCH http://127.0.0.1:8000/api/segments/ \
-d '{"user_id": 1, "segments_to_add": ["AVITO_VOICE_MESSAGES"], "segments_to_delete": ["AVITO_PERFORMANCE_VAS", "AVITO_DISCOUNT_30"]}'
```
> 200: {"user_id":1} <br>
> 400: {"message":"error occurred while processing segment to add 'AVITO_VOICE_MESSAGES': pq: duplicate key value violates unique constraint \"user_segments_pkey\""} <br>
> 400: {"message":"error occurred while checking segment to delete existence 'AVITO_PERFORMANCE_VAS': user(1) is not in this segment"}%

<p>Метод получения активных сегментов пользователя. Принимает на вход id пользователя</p>

```
curl -X GET http://127.0.0.1:8000/api/segments/ -d '{"user_id": 1}'
```
> 200: {"segments":["AVITO_VOICE_MESSAGES"],"user_id":1} <br>
> 200: {"segments":[],"user_id":1} <br>

### Дополнительное задание 1

<p>Метод получения истории попадания/выбывания пользователя из сегмента за определенный период в формате csv</p>

```
curl -X GET http://127.0.0.1:8000/api/users/history/ -d '{"user_id": 1, "year_month": "2023-08"}'
```
> 200: {"report":"http://localhost:8000/files/reports/user_history_2023-08_1.csv","user_id":1} <br>
> 400: {"message":"invalid YearMonth"} <br>
> 400: {"message":"json: cannot unmarshal string into Go struct field UserHistory.user_id of type int"} <br>

### Дополнительное задание 2

<p>Метод добавления пользователя в сегмент с возможностью задавать TTL</p>

```
curl -X PATCH http://127.0.0.1:8000/api/segments/ -d \
'{"user_id": 1, "segments_to_add": ["AVITO_VOICE_MESSAGES"], "segments_to_add_expiration":"2023-08-29 00:00:00", "segments_to_delete": []}'
```
> 200: {"user_id":1} <br>
> 400: {"message":"parsing time \"2023-08-31 00:00:\" as \"2006-01-02 15:04:05\": cannot parse \"\" as \"05\""} <br>

<p>Метод удаления истекших сегментов пользователя</p>
<sup>Вызывается автоматически при помощи cron каждую минуту</sup>

```
curl -X DELETE http://127.0.0.1:8000/api/users/expired-segments/
```
> 200: empty result <br>


### Расширенные примеры
##### Вы можете использовать curl'ы ниже по порядку и получить такие же ответы

```
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_VOICE_MESSAGES"}' -w '\n'
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_PERFORMANCE_VAS"}' -w '\n'
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_DISCOUNT_30"}' -w '\n'
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_DISCOUNT_50"}' -w '\n'
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_DELIVERY_FEATURE"}' -w '\n'
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_CLOUD_FEATURE"}' -w '\n'
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_UNDEFINED_FEATURE"}' -w '\n'
```

> {"slug":"AVITO_VOICE_MESSAGES"} <br>
> {"slug":"AVITO_PERFORMANCE_VAS"} <br>
> {"slug":"AVITO_DISCOUNT_30"} <br>
> {"slug":"AVITO_DISCOUNT_50"} <br>
> {"slug":"AVITO_DELIVERY_FEATURE"} <br>
> {"slug":"AVITO_CLOUD_FEATURE"} <br>
> {"slug":"AVITO_UNDEFINED_FEATURE"}

```
curl -X PATCH http://127.0.0.1:8000/api/segments/ \
-d '{"user_id": 1, "segments_to_add": ["AVITO_VOICE_MESSAGES", "AVITO_PERFORMANCE_VAS", "AVITO_DISCOUNT_30", "AVITO_DISCOUNT_50"], "segments_to_delete": []}'  -w '\n'

curl -X PATCH http://127.0.0.1:8000/api/segments/ \
-d '{"user_id": 1, "segments_to_add": ["AVITO_DELIVERY_FEATURE", "AVITO_CLOUD_FEATURE", "AVITO_UNDEFINED_FEATURE"], "segments_to_delete": ["AVITO_VOICE_MESSAGES", "AVITO_PERFORMANCE_VAS", "AVITO_DISCOUNT_30"]}' -w '\n'

curl -X GET http://127.0.0.1:8000/api/segments/ -d '{"user_id": 1}' -w '\n'

curl -X PATCH http://127.0.0.1:8000/api/segments/ \
-d '{"user_id": 1, "segments_to_add": [], "segments_to_delete": ["AVITO_DISCOUNT_50", "AVITO_DELIVERY_FEATURE", "AVITO_CLOUD_FEATURE", "AVITO_UNDEFINED_FEATURE"]}' -w '\n'

curl -X GET http://127.0.0.1:8000/api/segments/ -d '{"user_id": 1}' -w '\n'
```
> {"user_id":1} <br>
> {"user_id":1} <br>
> {"segments":["AVITO_CLOUD_FEATURE","AVITO_DELIVERY_FEATURE","AVITO_DISCOUNT_50","AVITO_UNDEFINED_FEATURE"],"user_id":1} <br>
> {"user_id":1} <br>
> {"segments":[],"user_id":1}

```
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_VOICE_MESSAGES"}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_PERFORMANCE_VAS"}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_DISCOUNT_30"}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_DISCOUNT_50"}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_DELIVERY_FEATURE"}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_CLOUD_FEATURE"}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_UNDEFINED_FEATURE"}' -w '\n'
```

> {"slug":"AVITO_VOICE_MESSAGES"} <br>
> {"slug":"AVITO_PERFORMANCE_VAS"} <br>
> {"slug":"AVITO_DISCOUNT_30"} <br>
> {"slug":"AVITO_DISCOUNT_50"} <br>
> {"slug":"AVITO_DELIVERY_FEATURE"} <br>
> {"slug":"AVITO_CLOUD_FEATURE"} <br>
> {"slug":"AVITO_UNDEFINED_FEATURE"}

```
curl -X GET http://127.0.0.1:8000/api/users/history/ -d '{"user_id": 1, "year_month": "2023-08"}'
```

> {"report":"http://localhost:8000/files/reports/user_history_2023-08_1.csv","user_id":1}
>
> user_history_2023-08_1.csv:
> 
> 1,AVITO_VOICE_MESSAGES,добавление,2023-08-29 20:33:11 <br>
> 1,AVITO_PERFORMANCE_VAS,добавление,2023-08-29 20:33:11 <br>
> 1,AVITO_DISCOUNT_30,добавление,2023-08-29 20:33:11 <br>
> 1,AVITO_DISCOUNT_50,добавление,2023-08-29 20:33:11 <br>
> 1,AVITO_DELIVERY_FEATURE,добавление,2023-08-29 20:33:12 <br>
> 1,AVITO_CLOUD_FEATURE,добавление,2023-08-29 20:33:12 <br>
> 1,AVITO_UNDEFINED_FEATURE,добавление,2023-08-29 20:33:12 <br>
> 1,AVITO_VOICE_MESSAGES,удаление,2023-08-29 20:33:12 <br>
> 1,AVITO_PERFORMANCE_VAS,удаление,2023-08-29 20:33:12 <br>
> 1,AVITO_DISCOUNT_30,удаление,2023-08-29 20:33:12 <br>
> 1,AVITO_DISCOUNT_50,удаление,2023-08-29 20:33:12 <br>
> 1,AVITO_DELIVERY_FEATURE,удаление,2023-08-29 20:33:12 <br>
> 1,AVITO_CLOUD_FEATURE,удаление,2023-08-29 20:33:12 <br>
> 1,AVITO_UNDEFINED_FEATURE,удаление,2023-08-29 20:33:12 <br>

```
curl -X POST http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_VOICE_MESSAGES"}' -w '\n'
curl -X PATCH http://127.0.0.1:8000/api/segments/ \
-d '{"user_id": 1, "segments_to_add": ["AVITO_VOICE_MESSAGES"], "segments_to_add_expiration":"1970-08-31 00:00:00", "segments_to_delete": []}' -w '\n'
curl -X GET http://127.0.0.1:8000/api/segments/ -d '{"user_id": 1}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/users/expired-segments/ -w '\n'
curl -X GET http://127.0.0.1:8000/api/segments/ -d '{"user_id": 1}' -w '\n'
curl -X DELETE http://127.0.0.1:8000/api/segments/ -d '{"slug": "AVITO_VOICE_MESSAGES"}' -w '\n'
```

> {"slug":"AVITO_VOICE_MESSAGES"} <br>
> {"user_id":1} <br>
> {"segments":["AVITO_VOICE_MESSAGES"],"user_id":1} <br>
>  <br>
> {"segments":[],"user_id":1} <br>
> {"slug":"AVITO_VOICE_MESSAGES"} <br>