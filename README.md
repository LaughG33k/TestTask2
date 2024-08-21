Тестовое задание

В связи с добавление аунтификации были добавлены новые ручки. смотреть swagger

Так же что бы можео было редактировать news нужно создать админа.
insert into users (login, password, permission) values ('somelogin', crypt('somepass', gen_salt('md5')), 'admin')

Для того что бы не кидало ошибку о авторизации, получаем jwt админа через /auth/login, и устанавливаем jwt в заголовок 
Header{
    Authorization: ["Bearer someJwtToken"]
}


Так же буду рад получить фид бэк по коду и, что стоило бы улучшить
