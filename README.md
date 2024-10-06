# Temple Union API

base_url = https://temple-api.onrender.com

## Api Docs

# Auth group

1. POST {{base_url}}/auth/register

   Принимает json формата
   {
    "username": string,
    "email": string,
    "password": string
}

username и password должны удовлетворять "^[-a-zA-Z0-9_#$&*]+$"
длина username - 4-32
длина password - 8-32

  Возвращает 200 статус, json {"details": "ok"}

  Возвращает 400 статус:
  1. В принимамемом json нет всех необходимых полей
  2. username, password невалидны

  Возвращает 409, если email или username уже заняты.
