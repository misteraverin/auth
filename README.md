# Домашнее задание. Сервис аутенфикации

<aside>
Представьте, что вы приходите на работу и бизнес дает вам задачку с 0 написать сервис аутенфикации. Задание специально не разжевано, все функицональные и нефункциональные требования вы должны узнать сами у нас.
</aside>

## Задача от бизнеса

Сервис auth - сервис аутентификации. Должен реализовывать следующие методы:

- login: пользователь передает логин/пароль через BasicAuth (см. метод http.Request.BasicAuth), если логин / пароль правильные, то ему возвращается ответ 200 и в хедере JWT-токена. Время жизни токена 60 мин.
- verify: возвращает 200, если access-токен валиден. Обновляет токен.

### Про JWT

Если хочется приблизить сервис к реальной системе, то пароли в сервисе можно лучше хранить в шифрованном виде с использованием библиотек bcrypt(https://pkg.go.dev/golang.org/x/crypto/bcrypt) и pdkdf2(https://pkg.go.dev/golang.org/x/crypto/pbkdf2), т.к. они поддерживают соль.

### Примеры CURL-запросов

```
curl -X POST \
  http://localhost:8080/verify \
  -H "Authorization: Bearer <access_token>" \
  -d "{}"

curl -X POST \
  http://localhost:8080/login \
  -u "login:password"
```

## Basic версия сервиса

- хранение в памяти
- REST
- Docker, Docker-Compose

## Advanced версия сервиса

- шифрование токенов
- grpc
- хранение в PostgreSQL
- мониторинг: Grafana, Prometheus
- (из Basic версии) Docker, Docker-Compose
