# Описание настроек (`.env`)

В отличие от Gps2Yandex, который настраивался через `settings.json`, bus2map настраивается через переменные окружения
(`.env`-файл рядом с бинарником, пример — `cmd/.env`).

```dotenv
# Допустимые значения: debug, info, warn, error
LOG_LEVEL=debug
# Записывать логи в файл (./logs/current.log, ротация по 10 МБ / 30 дней)
LOG_TO_FILE=false

WIALON_IPS_ENABLED=true
WIALON_IPS_LISTEN_ADDR=:20332

EGTS_ENABLED=true
EGTS_LISTEN_ADDR=:30332

YANDEX_ENABLED=true
YANDEX_URL=url
YANDEX_CLID=clid
YANDEX_COMPRESS=true

TWOGIS_ENABLED=true
TWOGIS_URL=url
TWOGIS_CLID=clid
TWOGIS_COMPRESS=false

GRPC_LISTEN_ADDR=:9090
GRPC_REFLECTION=true
```

Где:

|Переменная|Описание|
|---|---|
|**Логирование**||
|`LOG_LEVEL`|Уровень логирования: `debug`, `info`, `warn`, `error`|
|`LOG_TO_FILE`|Дублировать логи в файл `./logs/current.log` (ротация: до 10 МБ на файл, хранение 30 дней, старые файлы архивируются)|
|**Приём координат — Wialon IPS**||
|`WIALON_IPS_ENABLED`|Принимать ретрансляцию от "Wialon" по протоколу Wialon IPS|
|`WIALON_IPS_LISTEN_ADDR`|Адрес и порт для приёма, например `:20332` (слушать на всех интерфейсах) или `192.168.1.10:20332`|
|**Приём координат — ЕГТС/EGTS** *(новое)*||
|`EGTS_ENABLED`|Принимать данные напрямую от трекеров/бортового оборудования по протоколу ЕГТС|
|`EGTS_LISTEN_ADDR`|Адрес и порт для приёма ЕГТС, например `:30332`|
|**Отправка — "Яндекс.Карты"**||
|`YANDEX_ENABLED`|Отправлять обогащённые данные в "Яндекс.Карты"|
|`YANDEX_CLID`|Идентификатор участника программы, выдаётся Яндексом|
|`YANDEX_URL`|Endpoint, выдаётся партнёрам Яндекса|
|`YANDEX_COMPRESS`|Сжимать (gzip) тело запроса при отправке|
|**Отправка — "2ГИС"** *(новое)*||
|`TWOGIS_ENABLED`|Отправлять обогащённые данные в "2ГИС"|
|`TWOGIS_CLID`|Идентификатор участника программы, выдаётся 2ГИС|
|`TWOGIS_URL`|Endpoint, выдаётся партнёрам 2ГИС|
|`TWOGIS_COMPRESS`|Сжимать (gzip) тело запроса при отправке|
|**gRPC API** *(новое)*||
|`GRPC_LISTEN_ADDR`|Адрес и порт gRPC-сервера, например `:9090` — см. [grpc_api.md](./grpc_api.md)|
|`GRPC_REFLECTION`|Включить [gRPC reflection](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md) — удобно для ручных вызовов из `grpcurl`/Postman без `.proto`-файла на клиенте, в проде обычно выключается|

**Примечание:** директория со справочниками (`route.txt`, `transport.txt`, `schedule.txt`) в `.env` не задаётся —
пути к файлам справочников зашиты в бинарнике (`cmd/datasource/route.txt`, `transport.txt`, `schedule.txt` — файлы
кладутся рядом с исполняемым файлом, в подкаталог `datasource/`). Это отличается от Gps2Yandex, где директория со
справочниками указывалась настройкой `Catalogs.Directory`.

**Внимание:** ни одна секция (`WialonIPS`/`EGTS`/`Yandex`/`TwoGIS`) не обязана быть единственной активной — можно
включить приём одновременно с Wialon IPS и ЕГТС, и одновременно слать в "Яндекс.Карты" и "2ГИС". Секции, где
`*_ENABLED=false`, можно не заполнять содержательными значениями (кроме самого флага) — сервис их пропустит.
