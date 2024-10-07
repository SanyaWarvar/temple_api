# Temple Union API

base_url = https://temple-api.onrender.com

## Api Docs

# Auth group

**1. POST {{base_url}}/auth/register**

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
   
Статус коды:

200 статус, json {"details": "ok"}

Если статус не 200, то json {"message": string}

400 статус:
  
  1. В принимамемом json нет всех необходимых полей
  2. username, password невалидны

409 статус:
  
   1. email или username уже заняты.

**2. POST {{base_url}}/auth/send_code**

Отправляет на почту юзера код подтверждения.

   Принимает json формата
   {
    "username": string,
    "email": string,
    "password": string
}
   
Статус коды:

200 статус, json {
		"exp_time":       string,
		"next_code_time": string,
	}

exp_time - время, по прошествии которого, отправленный код перестанет быть действительным

next_code_time - время, по прошествии которого, можно будет запросить следующий код

Пример ответа:

{
    "exp_time": "15m0s",
    "next_code_time": "1m0s"
}


Если статус не 200, то json {"message}": string}

400 статус:
  
  1. В принимамемом json нет всех необходимых полей
  2. username, password невалидны

401 статус:

   1. неверная пара логин, пароль.
   2. email не принадлежит этому пользователю

409 статус:
  
   1. Этот пользователь уже подтвердил email

**3. POST {{base_url}}/auth/confirm_email**

Подтверждает почту пользователя.

   Принимает json формата
   {
    "email": string,
    "code": string
   }
   
Статус коды:

200 статус, json {
		"details": "ok"
	}

exp_time - время, по прошествии которого, отправленный код перестанет быть действительным
next_code_time - время, по прошествии которого, можно будет запросить следующий код

Пример ответа:

{
    "exp_time": "15m0s",
    "next_code_time": "1m0s"
}


Если статус не 200, то json {"message}": string}

400 статус:
  
  1. В принимамемом json нет всех необходимых полей
  2. Полученный email не зарегистрирован

409 статус:
  
   1. Этот пользователь уже подтвердил email

**4. POST {{base_url}}/auth/sign_in**

Создает пару access, refresh токенов при помощи логина, пароля и почты.

   Принимает json формата
   {
    "username": string,
    "email": string,
    "password": string
   }
   
Статус коды:

200 статус, json 
{
	"access_token": string,
   "refresh_token": string
}

Если статус не 200, то json {"message}": string}

400 статус:
  
  1. В принимамемом json нет всех необходимых полей
  2. Полученный email не зарегистрирован

401 статус:
   1. email не принадлежит пользователю

403 статус:
   1. email не подтвержден

**5. POST {{base_url}}/auth/refresh_token**

Создает пару access, refresh токенов при помощи пары ранее выданных токенов.

   Принимает json формата
  {
	"access_token": string,
   "refresh_token": string
 }
   
Статус коды:

200 статус, json 
{
	"access_token": string,
   "refresh_token": string
}

Если статус не 200, то json {"message}": string}

400 статус:
  
  1. В принимамемом json нет всех необходимых полей
  2. Некорректная подпись access токена
  3. Некорректный refresh токен

401 статус:
   1. Время действия refresh токена истекло
   2. refresh токен не найден (не был выдан сервером, либо уже был использован)

