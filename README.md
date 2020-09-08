# Микросервис для работы с балансом: Авито Job

## Запуск приложения
```
git clone https://github.com/AndreyChemurov/microservice_job_avito.git
cd microservice_job_avito/
[sudo] docker-compose up
```

## Информация
Приложение работает на порту ```8000```.

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
- /balance </br>

- /increase </br>

- /decrease </br>

- /remittance </br>

### Через curl
- /balance </br>
```curl http://localhost:8000/balance?id=someuser``` </br>
```curl 'http://localhost:8000/balance?id=someuser&currency=USD'``` </br>
```curl --header "Content-Type: application/x-www-form-urlencoded" --data "id=someuser" --request POST http://localhost:8000/balance``` </br>
```curl --header "Content-Type: application/x-www-form-urlencoded" --data "id=someuser&currency=USD" --request POST http://localhost:8000/balance``` </br>
```curl --header "Content-Type: application/json" --data '{"id": "someuser"}' --request POST http://localhost:8000/balance``` </br>
```curl --header "Content-Type: application/json" --data '{"id": "someuser", "currency": "USD"}' --request POST http://localhost:8000/balance``` </br>
- /increase </br>
```curl 'http://localhost:8000/increase?id=someuser&money=100'``` </br>
```curl --header "Content-Type: application/x-www-form-urlencoded" --data "id=someuser&money=100" --request POST http://localhost:8000/increase``` </br>
```curl --header "Content-Type: application/json" --data '{"id": "someuser", "money": "100"}' --request POST http://localhost:8000/increase``` </br>
- /decrease </br>
```curl 'http://localhost:8000/decrease?id=someuser&money=100'``` </br>
```curl --header "Content-Type: application/x-www-form-urlencoded" --data "id=someuser&money=100" --request POST http://localhost:8000/decrease``` </br>
```curl --header "Content-Type: application/json" --data '{"id": "someuser", "money": "100"}' --request POST http://localhost:8000/decrease``` </br>
- /remittance </br>
```curl 'http://localhost:8000/remittance?from=someuser&to=someuser2&money=100'``` </br>
```curl --header "Content-Type: application/x-www-form-urlencoded" --data "from=someuser&to=someuser2&money=100" --request POST http://localhost:8000/remittance``` </br>
```curl --header "Content-Type: application/json" --data '{"from": "someuser", "to": "someuser2", "money": "100"}' --request POST http://localhost:8000/remittance``` </br>
### Через браузер
- /balance

- /increase

- /decrease

- /remittance

## Запуск тестов
```
pass
```

## Результаты тестов
pass

## Проблемы и решения
### Проблема №1
Т.к. нет никакой информации о пользователях, то новый пользователь будет создаваться каждый раз, когда указывается новый уникальный id при вызове "increase": http://localhost:8000/increase?id=some-id&money=100 </br>

### Проблема №2
Некоторые запросы на списание денег с баланса пользователя могут приходить синхронно, таким образом будет гонка между транзакциями. Поэтому присутствует условие CHECK (amount >= 0.0), которое не позволит списать средства с баланса пользователя, если у него их недостаточно.

### Проблема №3
Если баланс = **n**, а списание = **n + m**, то средства не списываются с баланса.
