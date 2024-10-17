
base_url = https://temple-api.onrender.com

## Api Docs

Если статус не 2xx, то сообщение с ошибкой будет в {"message": string}.

Все ендпоинты, кроме ендпоинтов группы **Auth group**, требуют header: 'Authorization: ••••••', который будет содержать access токен!

### Auth group
***
**1. POST {{base_url}}/auth/sign_up**

Создает нового пользователя.

   Принимает json формата
   {
    "username": string,
    "email": string,
    "password": string
}

   username и password должны удовлетворять "^[-a-zA-Z0-9_#$&*]+$"
   длина username - 4-32,
   длина password - 8-32
   
Варианты ответа:

|Статус код|Причины|
|-|--------|
|201|Все хорошо|
|400|В принимамемом json нет всех необходимых полей|
|400|username, password невалидны|
|409|email или username уже заняты|
***
**2. POST {{base_url}}/auth/send_code**

Отправляет на почту юзера код подтверждения.

   Принимает json формата
   {
    "email": string,
    "password": string
    }
   
Варианты ответа:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|201|Все хорошо|{"exp_time": string, "next_code_time": string}|
|400|В принимамемом json нет всех необходимых полей|{"message": string}|
|400|password невалидty|{"message": string}|
|401|Неверный пароль|{"message": string}|
|409|email уже подтвержден|{"message": string}|

exp_time - время, по прошествии которого, отправленный код перестанет быть действительным

next_code_time - время, по прошествии которого, можно будет запросить следующий код
***
**3. POST {{base_url}}/auth/confirm_email**

Подтверждает почту пользователя.

   Принимает json формата
   {
    "email": string,
    "code": string
   }
   
Варианты ответа:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|204|Все хорошо|null|
|400|В принимамемом json нет всех необходимых полей|{"message": string}|
|401|Неверный код|{"message": string}|
|400|email не зарегистрирован|{"message": string}|
|409|email уже подтвержден|{"message": string}|
***
**4. POST {{base_url}}/auth/sign_in**

Создает пару access, refresh токенов при помощи логина, пароля и почты.

   Принимает json формата
   {
    "email": string,
    "password": string
   }
   
Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|201|Все хорошо|{"access_token": string, "refresh_token": string}|
|400|В принимамемом json нет всех необходимых полей|{"message": string}|
|400|Email не зарегистрирован|{"message": string}|
|401|Неверный пароль|{"message": string}|
|403|email не подтвержден|{"message": string}|
***
**5. POST {{base_url}}/auth/refresh_token**

Создает пару access, refresh токенов при помощи пары ранее выданных токенов.

   Принимает json формата
  {
	"access_token": string,
   "refresh_token": string
 }
   
Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|201|Все хорошо|{"access_token": string, "refresh_token": string}|
|400|В принимамемом json нет всех необходимых полей|{"message": string}|
|400|Некорректный access токен|{"message": string}|
|400|Некорректный refresh токен|{"message": string}|
|401|refresh токен не принадлежит пользователю, либо не найден (просрочен, не был выдан или уже был использован)|{"message": string}|

***
### User Group
***
**1. GET {{base_url}}/user/:username**

**!этому ендпоинту header с токеном не нужен!**

Возвращает информацию о пользователе. 

Возвращает:

200 статус и json

{
    "first_name": string | null,
    "second_name": string | null,
    "status": string | null,
    "birthday": datetime (формат RFC3339) | null,  
    "gender": string | null,
    "country": string | null,
    "city": string | null
}.

***

**2. PUT {{base_url}}/user/**

Обновляет данные о пользователе.

Принимает json
{
    "first_name": string | null,
    "second_name": string | null,
    "status": string | null,
    "birthday": datetime (формат RFC3339) | null,  
    "gender": string | null,
    "country": string | null,
    "city": string | null
}.

Ни одно из полей не является обязательным, но необходимо передать хотя бы из полей.

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|204|Все хорошо|null|
|400|Некорректный json|{"message": string}|
|401|Не передан header с токеном|{"message": string}|
|401|Срок действия токена истек|{"message": string}|
***
