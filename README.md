# avito-backend-trainee-assignment-2023

## Запуск

Приложения запускается через docker-compose на localhost:8080

```sh
docker-compose up -d
```

## База данных

В директории assets/postgres лежит init.sql файл, который инициализирует начальное состояние базы данных.

## Примеры запросов и ответов

В директории assets/swagger описал swagger.yml

# Создание сегмента
```sh
curl --request POST \
  --url http://localhost:8080/api/segment/add \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "TEST_SEGMENT"
}'
```

# Удаление сегмента
```sh
curl --request POST \
  --url http://localhost:8080/api/segment/delete \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "vbcvbcb"
}'
```

# Получение сегментов пользователя
```sh
curl --request POST \
  --url http://localhost:8080/api/user/segments \
  --header 'Content-Type: application/json' \
  --data '{
	"user_id": 123
}'
```

# Получение ссылки на csv-отчет с историей сегментов пользователей за период
```sh
curl --request POST \
  --url http://localhost:8080/api/user/segments-history \
  --header 'Content-Type: application/json' \
  --data '{
	"period": "2023-09"
}'
```

# Изменение сегментов пользователя
```sh
curl --request POST \
  --url http://localhost:8080/api/user/change-segments \
  --header 'Content-Type: application/json' \
  --data '{
	"add_segments": [
		{
			"name": "test1",
			"ttl": "2023-08-26 21:51:05"
		},
		{
			"name": "test3"
		}
	],
	"delete_segments": ["test2"],
	"user_id": 123
}'
```

## Вопросы и их решения

По доп. заданиям:
1. Сам бы отчет сохранял в S3-хранилище и давал ссылку на соответствующий объект, если бы это был продакшен. Сделал проще - сохраняю отчет на сервере и отдаю на него ссылку
2. Было несколько мыслей, как это лучше реализовать. Задание на крон, триггер в базе данных, отдельная горутина. Выбрал последний вариант.