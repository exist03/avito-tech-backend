# avito-tech-backend

Проект представляет собой сервис для централизации работы с экспериментами(тесты новых продуктов, тесты интерфейса,
скидочные и многие другие)

# Инструкция по запуску

```shell
#Склонировать репозиторий и перейти в рабочую директорию
https://github.com/exist03/avito-tech-backend.git
cd avito-tech-backend
#Запуск сервиса
make docker
#Swagger документация
localhost:8080/swagger/index.html
#Запуск тестов
make tests
```

# Запросы обрабатываемые сервисом

`GET /user/get/{user_id}`<br/>
`GET /user/user/get_history/`<br/>
`POST /segment body{"id":1, "name":"someName"}`<br/>
`DELETE /segment/{id}`<br/>
`PATCH /user/update body{"segments_add": [
{
"segment": {
"id": 15,
"name": "someName"
},
"ttl": "2024-03-12T13:37:27+00:00"
}
],
"segments_del": [
{
"segment": {
"id": 48,
"name": "anotherName"
},
"ttl": "2024-03-12T13:37:27+00:00"
}
],
"user_id": 0}`<br/>

# Примеры запросов

```shell
#GET
curl -X GET localhost:8080/api/service/user/get/1
#GET
curl -X GET 'localhost:8080/api/service/user/get_history?user_id=1&start=1693063271&end=1693209561'
#POST
curl -X POST localhost:8080/api/service/segment -H "User-role: admin" -H "Content-Type: application/json" -d '{"id":20, "name":"somename", "percent": 35, "ttl": "2024-03-12T13:37:27+00:00"}'
#PATCH
curl -X PATCH localhost:8080/api/service/user/update -H "Content-Type: application/json" -d '{"user_id":1, "segments_add":[{"id":1,"name":"somename","ttl":"2024-03-12T13:37:27+00:00"}]}'
#DELETE
curl -X DELETE localhost:8080/api/service/segment/123
```

# Возникшие проблемы

- Методы создания и удаления сегментов принимают название.
    - Логично, что при попытке создания мы не должны создавать ранее существующий сегмент и аналогично с удалением. Это
      можно решить если название сегмента сделать PK, но строка первичный ключ это плохая практика. Соответсвенно примем
      за PK id, и примем на вход id и название.
- Сегменты не должны добавляться/удаляться лишними пользователями.
  - В подобных запросах для валидации необходимо передавать в заголовке ```User-role: admin```
