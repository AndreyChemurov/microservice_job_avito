# Микросервис для работы с балансом: Авито Job

## Запуск приложения
```
git clone https://github.com/AndreyChemurov/microservice_job_avito.git
cd microservice_job_avito/
[sudo] docker-compose up
```
## Примеры запросов/ответов

### Через Postman
> /balance
pass

> /increase
pass

> /decrease
pass

> /remittance
pass

### Через curl
> /balance
pass

> /increase
pass

> /decrease
pass

> /remittance
pass

### Через браузер
> /balance
pass

> /increase
pass

> /decrease
pass

> /remittance
pass

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
Некоторые запросы на списание денег с баланса пользователя могут приходить синхронно, таким образом будет гонка между транзакциями... **Solution is coming :)**

### Проблема №3
Если баланс = **n**, а списание = **n + m**, то средства не списываются с баланса.

## Информация

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
