# Микросервис для работы с балансом: Авито Job

## Запуск приложения
```
git clone https://github.com/AndreyChemurov/microservice_job_avito.git
cd microservice_job_avito/
[sudo] docker-compose up
```

## Информация
Приложение работает на порту ```8000```. </br>
Реализовано доп. задание №1.

**Запросы** </br>

- Зачисление средств.
- Списание средств.
- Перевод средств от пользователя к пользователю.
- Получение баланса пользователя.

**Сущности и атрибуты** </br>

***USER_JOB*** </br>
- user_id: уникальный идентификатор.

***BALANCE_JOB*** </br>
- balance_id: уникальный идентификатор.
- user_id: id пользователя. (один-к-одному)
- amount: баланс пользователя.

## Примеры запросов/ответов

### Через Postman
![Screenshot_20200910_144742](https://user-images.githubusercontent.com/58785926/92725402-c845e080-f374-11ea-859d-dc68d0dcfb5c.png) </br>
![Screenshot_20200910_144803](https://user-images.githubusercontent.com/58785926/92725415-cda32b00-f374-11ea-8865-2d00bf2a2612.png) </br>
![Screenshot_20200910_144825](https://user-images.githubusercontent.com/58785926/92725426-d0058500-f374-11ea-86e2-2c24dcf163b1.png) </br>
![Screenshot_20200910_144847](https://user-images.githubusercontent.com/58785926/92725431-d1cf4880-f374-11ea-9102-e90cdb4eb2fa.png) </br>
### Через curl
- /balance (текущий баланс) </br>
```curl http://localhost:8000/balance?id=someuser``` </br> </br>
```curl 'http://localhost:8000/balance?id=someuser&currency=USD'``` </br> </br>
```curl --header "Content-Type: application/x-www-form-urlencoded" --data "id=someuser" --request POST http://localhost:8000/balance``` </br> </br>
```curl --header "Content-Type: application/x-www-form-urlencoded" --data "id=someuser&currency=USD" --request POST http://localhost:8000/balance``` </br> </br>
```curl --header "Content-Type: application/json" --data '{"id": "someuser"}' --request POST http://localhost:8000/balance``` </br> </br>
```curl --header "Content-Type: application/json" --data '{"id": "someuser", "currency": "USD"}' --request POST http://localhost:8000/balance``` </br> </br>
- /increase (начислить средства) </br> 
```curl 'http://localhost:8000/increase?id=someuser&money=100'``` </br> </br>
```curl --header "Content-Type: application/x-www-form-urlencoded" --data "id=someuser&money=100" --request POST http://localhost:8000/increase``` </br> </br>
```curl --header "Content-Type: application/json" --data '{"id": "someuser", "money": "100"}' --request POST http://localhost:8000/increase``` </br> </br>
- /decrease (списать средства) </br>
```curl 'http://localhost:8000/decrease?id=someuser&money=100'``` </br> </br>
```curl --header "Content-Type: application/x-www-form-urlencoded" --data "id=someuser&money=100" --request POST http://localhost:8000/decrease``` </br> </br>
```curl --header "Content-Type: application/json" --data '{"id": "someuser", "money": "100"}' --request POST http://localhost:8000/decrease``` </br> </br>
- /remittance (перевод средств) </br>
```curl 'http://localhost:8000/remittance?from=someuser&to=someuser2&money=100'``` </br> </br>
```curl --header "Content-Type: application/x-www-form-urlencoded" --data "from=someuser&to=someuser2&money=100" --request POST http://localhost:8000/remittance``` </br> </br>
```curl --header "Content-Type: application/json" --data '{"from": "someuser", "to": "someuser2", "money": "100"}' --request POST http://localhost:8000/remittance``` </br>

![Screenshot_20200910_142645](https://user-images.githubusercontent.com/58785926/92724450-4f925480-f373-11ea-9478-6b5f9bb248d3.png)
![Screenshot_20200910_142743](https://user-images.githubusercontent.com/58785926/92724477-5751f900-f373-11ea-9d2d-1f5aeacd6971.png)
![Screenshot_20200910_142813](https://user-images.githubusercontent.com/58785926/92724486-59b45300-f373-11ea-9e14-5c329e2cbf66.png)
![Screenshot_20200910_142844](https://user-images.githubusercontent.com/58785926/92724491-5b7e1680-f373-11ea-9216-bcccb1b1bfd0.png)

## Запуск тестов
Из-под директории **microservice_job_avito/** </br>
Тесты для main пакета: </br>
```
[sudo] docker-compose run --rm web go test . [-v [-cover]] && [sudo] docker-compose stop db
```
Тесты для database пакета: </br>
```
[sudo] docker-compose run --rm web go test ./internal/database/ [-v [-cover]] && [sudo] docker-compose stop db
```

## Результаты тестов
Покрытие main - 55.6% </br>
Покрытие database - 79.2% </br>
</br>
![Screenshot_20200913_210845](https://user-images.githubusercontent.com/58785926/93025294-6c45bb00-f605-11ea-9c84-bb0e1e488911.png)
![Screenshot_20200913_210910](https://user-images.githubusercontent.com/58785926/93025295-6e0f7e80-f605-11ea-9569-116d21c39876.png)

## Проблемы и решения
### Проблема №1
Т.к. нет никакой информации о пользователях, то новый пользователь будет создаваться каждый раз, когда указывается новый уникальный id при вызове "increase": http://localhost:8000/increase?id=some-id&money=100 </br>

### Проблема №2
Некоторые запросы на списание денег с баланса пользователя могут приходить синхронно, таким образом будет гонка между транзакциями. Поэтому присутствует условие CHECK (amount >= 0.0), которое не позволит списать средства с баланса пользователя, если у него их недостаточно.

### Проблема №3
Если баланс = **n**, а списание = **n + m**, то средства не списываются с баланса.
