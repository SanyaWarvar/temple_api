
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
    "profile_pic_url": string,
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

**3. GET {{base_url}}/user/find**

Поиск людей по имени и фамилии.

Принимает json
{
    "search_string": string,
    "page": int | null
    
}.
Параметр page строго больше нуля.
Пример: {
	"search_string": "Иванов Иван",
	"page": 1
}

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|[{"first_name": string, "second_name":string, "username": string, "profile_picture": url}], длина массива не более 50|
|400|Некорректный json|{"message": string}|
***

**4. PUT {{base_url}}/users_posts/:id**

Обновляет пост.

Принимает json формата: 
{
    "body": string
}

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|{"details": "success"}|
|400|Некорректные данные|{"message": string}|
|409|Поста не существует|{"message": string}|
|409|Пост принадлежит другому пользователю|{"message": string}|
***

**5. GET {{base_url}}/user/friends/:page**

Возвращает список всех друзей пользователя

Принимает json формата: 
{
    "username": string
}

Возвращает:


|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|{friends:[{fullname: string, username: string, profile_pic: url}]}|
|400|Некорректные данные|{"message": string}|
***

**6. GET {{base_url}}/user/follows/:page**

Возвращает список всех подписок (на пользователей) пользователя

Принимает json формата: 
{
    "username": string
}

Возвращает:


|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|{follows:[{fullname: string, username: string, profile_pic: url}]}|
|400|Некорректные данные|{"message": string}|
***

**7. GET {{base_url}}/user/subs/:page**

Возвращает список всех подписчиков пользователя

Принимает json формата: 
{
    "username": string
}

Возвращает:


|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|{subs:[{fullname: string, username: string, profile_pic: url}]}|
|400|Некорректные данные|{"message": string}|
***

**8. PUT {{base_url}}/user/profile_pic**

Обновляет фотографию профиля пользователя

Принимает form data формата:

profile_pic: file

файл должен быть с расширением .gif, .png, .jpg или .svg

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|201|Все хорошо|null|
|400|Некорректные данные|{"message": string}|
***

**9. POST {{base_url}}/user/friends/:username**

Добавить пользователя в друзья (отправить приглашение)

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|204|Все хорошо|null|
|400|Некорректные данные|{"message": string}|
|400|Пользователь уже в друзьях|{"message": string}|
|400|Приглашение отправлено самому себе|{"message": string}|

***

**10. PUT {{base_url}}/user/friends/:username**

Подтвердить дружбу (принять приглашение)

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|204|Все хорошо|null|
|400|Некорректные данные|{"message": string}|
|400|Пользователь уже в друзьях|{"message": string}|

***

**11. DELETE {{base_url}}/user/friends/:username**

Удалить пользователя из друзей

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|204|Все хорошо|null|
|400|Некорректные данные|{"message": string}|
|400|Пользователь не в друзьях|{"message": string}|

***
### Posts Group

**1. GET {{base_url}}/users_posts/:id**

Возвращает пост по его id.

Схема поста для вывода:
{
    "id": uuid,
    "author_id": uuid,
    "body": string,
    "last_update": time (rfc3339),
    "edited": bool,
    "likes_count": int,
    "liked_by_me": bool
}

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|Искомый пост (схема выше)|
|400|Пост не найден или некорректный id|{"message": string}|
***

**2. POST {{base_url}}/users_posts**

Создает новый пост.

Пост:
{
    "id": uuid,
    "author_id": uuid,
    "body": string,
    "last_update": time (rfc3339),
    "edited": bool
}

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|201|Все хорошо|{"post_id": uuid}|
|400|Некорректные данные|{"message": string}|
***

**3. DELETE {{base_url}}/users_posts/:id**

Удаляет пост. 

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|{"details": "success"}|
|400|Некорректные данные|{"message": string}|
|409|Пост уже удален|{"message": string}|
|409|Пост принадлежит другому пользователю|{"message": string}|
***

**4. PUT {{base_url}}/users_posts/:id**

Обновляет пост.

Принимает json формата: 
{
    "body": string
}

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|{"details": "success"}|
|400|Некорректные данные|{"message": string}|
|409|Поста не существует|{"message": string}|
|409|Пост принадлежит другому пользователю|{"message": string}|
***

**5. PUT {{base_url}}/users_posts/like/:id**

Ставит/снимает лайк к посту.

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|204|Все хорошо|null|
|400|Некорректные данные|{"message": string}|
|409|Поста не существует|{"message": "post not found"}|
***

**6. GET {{base_url}}/user/posts/:username**

Получить посты пользователя

Принимает json формата:

{
page: int
}

Возвращает массив (не более 50) из постов:
{
    "id": uuid,
    "author_id": uuid,
    "body": string,
    "last_update": time (rfc3339),
    "edited": bool,
    "likes_count": int,
    "liked_by_me": bool
}

|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|[]posts|
|400|Некорректные данные|{"message": string}|
***

### Chats Group

**1. GET {{base_url}}/chats/**

Получить все чаты пользователя (не более 50) и сообщения из каждого чата (по 25 на чат)

Принимает json формата:

{
page: int
}

Возвращает:
chat:
{
    "id": uuid,
    "members": []string,
    "messages": []message
}

message:
{
	"id":        uuid, 
	"body" :     string,    
	"author_id":  uuid, 
	"chat_id" :   uuid, 
	"created_at": time,
	"readed":    bool,      
	"edited" :   bool,      
	'reply_to" :  uuid 
}

|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|[]chats|
|400|Некорректные данные|{"message": string}|
***

**2. POST {{base_url}}/chats/**

Создать новый чат с другим пользователем

Принимает json формата:

{
	username: string
}

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|201|Все хорошо|{"chat_id": uuid}|
|400|Некорректные данные|{"message": string}|
|400|Чат уже создан|{"message": string}|
***

**3. GET {{base_url}}/chats/messages/:chat_id**

Получить все сообщения из чата.

Принимает json формата:

{
	page: int
}

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|[]messages|
|400|Некорректные данные|{"message": string}|
***

**4. POST {{base_url}}/chats/messages/:chat_id**

Отправить сообщение в чат.

Принимает json message

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|{"message_id": uuid}|
|400|Некорректные данные|{"message": string}|
***

**5. PUT {{base_url}}/chats/messages/read/:message_id**

Пометить сообщение прочитанным.

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|204|Все хорошо|null|
|400|Некорректные данные|{"message": string}|
***

**6. PUT {{base_url}}/chats/messages/**

Изменить сообщение. 

Принимает json формата:

{
	"id": uuid,
	"body": string
}

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|204|Все хорошо|null|
|400|Некорректные данные|{"message": string}|
***

**7. PUT {{base_url}}/chats/messages/:message_id**

Удалить сообщение

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|204|Все хорошо|null|
|400|Некорректные данные|{"message": string}|
***

### Rofls group

**1. POST {{base_url}}/rofls**

Выложить рофланыч (короткое видео).

Принимает form data:

{
"file": file,
"title": string
}

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|201|Все хорошо|{"id": uuid}|
|400|Некорректные данные|{"message": string}|
***

**2. GET {{base_url}}/rofls/:id**

Получить рофланыч по id

Возвращает rofl json формата:

{
	"id": uuid,
	"author_username": string,
	"created_at": time,
	"title" : string,
	"body": string
}

Поле body содержит url для получения файла.

|Статус код|Причины|Тело ответа|
|-|--------|------|
|200|Все хорошо|rofl|
|400|Некорректные данные|{"message": string}|
***

**3. DELETE {{base_url}}/rofls/:id**

Удалить рофланыч

Возвращает:

|Статус код|Причины|Тело ответа|
|-|--------|------|
|204|Все хорошо|null|
|400|Некорректные данные|{"message": string}|
***
