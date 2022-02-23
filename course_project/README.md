# Курсовой проект [TFS: Golang 2021](https://fintech.tinkoff.ru/study/fintech/golang/). Торговый робот

Робот торгует на демо-платформе [kraken-demo](https://futures.kraken.com/ru.html)
Для использования робота необходимо зарегистрироваться на платформе, и получить API ключи с полным доступом

## Предварительно:

Робот написан на языке Golang и для его работы необходмо установить Golang версии 1.17

В качестве хранения заявок и работы телеграм бота необходимо установить PostgreSQL12

После установки PostgreSQL необходимо создать базу данных, и в ней запустить [скрипт](init_table.sql) инициализирующий все таблицы

## Настройка робота

Для запуска робота необходимо указать указать параметры в [файле](resources/application.properties)

```
log.level - уровень логирования
apikey.public - public ключ для доступа к API биржы
apikey.private - private ключ для доступа к API биржы
telegram.token - token для телеграм бота
postgresql.host - адресс для базы данных
postgresql.sslmode - ssl mode
```

## Запуск робота

билдим проект

```shell
go build github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang
```

запускаем робота

```shell
./golang
```

## Работа робота

Для работы робота используется REST API

пример использования есть в [файле](src/golang/example.http)

Endpoints:

* POST Start -- запуск робота
* Post Stop -- остановить робота
* POST Configure -- конфигурация робота, ожидает в теле json

```json
{
  "product_id": "PI_XBTUSD", // Валюта для торгов 
  "stop_val": 10, // Сколько готовы потерять
  "profit_val": 20 // Сколько хотим получить
}
```

