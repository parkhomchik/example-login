### Задача: ###

Реализовать простую страницу входа пользователя. Покажите некоторые страницы, если Логин является успешной и установить сессию и куки. 
Внедрение механизма выхода из системы (просто удалить сессию / куки). 
Для CSRF, используйте https://github.com/codahale/charlie


### Задача сделать ###

#### Заказчик: ####

1. Просто используйте обычный JavaScript. 
Если вам необходимо использовать структуру, вы можете использовать JQuery.


#### Сервер: ####

1. Используйте идти
1. Нет необходимости в базе данных .. Просто нет пользователей и проводных пароль
1. Реализация CSRF токен, используя простую библиотеку https://github.com/codahale/charlie~~HEAD=dobj

### Как подать заявление на этом испытании ###

Многие люди слепо просить эту работу, не читая требования. 
Я хочу прочитать требования, и ответить обратно с "golangbuild" в качестве первого слова в строке темы Вашего ответа. 
Если это слово отсутствует в передней, я знаю, что вы не читали требования.

### Практические результаты ###

1. Ссылка на страницу, так что я могу проверить демо
1. Один .go файл, который имеет thecode, который реализует сервер.


Главная страница (Кнопка Логин если пользователь не залогинен, если залогинен то пишем приветсвие с его именем)
    - Логин 
        Удачный логин - Главная страница
        Не удачный логин - ОШИБКА
