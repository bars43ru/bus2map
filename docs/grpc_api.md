# gRPC API

Это функциональность, которой не было в Gps2Yandex: bus2map поднимает gRPC-сервер (адрес — `GRPC_LISTEN_ADDR`
в [.env](./settings.md)), который решает две независимые задачи — **приём** координат от внешних систем и
**раздачу** уже обогащённых данных подписчикам. Использовать можно любую из них по отдельности или обе сразу.

Описание сервиса — `api/proto/bustracking.proto`, сервис `BusTrackingService`.

## StreamGPSData — отправка координат в bus2map

Клиентский потоковый (client-streaming) метод: внешняя система открывает поток и пишет в него сырые GPS-точки,
bus2map обрабатывает их так же, как координаты, полученные по Wialon IPS или ЕГТС — то есть ищет транспортное
средство по `uid` в `transport.txt`, текущее расписание в `schedule.txt` и маршрут в `route.txt`.

Пригодится, если:
* у транспорта уже есть собственная телематика/бортовой компьютер с доступом в интернет, и нет смысла заводить его
  через "Wialon" или отдельный ЕГТС-трекер — систему можно подключить к bus2map напрямую;
* нужно протестировать/эмулировать сигнал без реального GPS-трекера.

```proto
rpc StreamGPSData(stream GPSData) returns (StreamGPSDataResponse);

message GPSData {
  string uid = 1;       // тот же uid, что и в transport.txt
  double latitude = 2;
  double longitude = 3;
  uint32 speed = 4;
  uint32 course = 5;
  google.protobuf.Timestamp time = 6;
}
```

`uid` должен совпадать со значением из первой колонки `transport.txt` — иначе, как и в случае с Wialon IPS/ЕГТС,
координата будет отброшена с предупреждением в логах ("not found UID in transport").

Ответ (`StreamGPSDataResponse`) сервис отправляет один раз — при закрытии клиентом потока (пустое сообщение,
подтверждающее, что все данные приняты).

## StreamBusTrackingInfo — подписка на обогащённые данные

Серверный потоковый (server-streaming) метод: клиент отправляет пустой запрос и получает бесконечный поток
`BusTrackingInfo` — то есть ровно те же данные, которые сервис параллельно отправляет в "Яндекс.Карты"/"2ГИС", но в
структурированном виде и без похода во внешние карты.

Пригодится для собственных интеграций — внутренних карт/дашбордов, алертинга (например, "автобус сошёл с
расписания"), аналитики — без необходимости заново разбирать протоколы Wialon IPS/ЕГТС и повторно писать логику
сопоставления со справочниками.

```proto
rpc StreamBusTrackingInfo(StreamBusDataRequest) returns (stream BusTrackingInfo);

message BusTrackingInfo {
  GPSData gps_data = 1;    // координаты, скорость, курс, время
  Route route = 2;         // number / yandex / two_gis — из route.txt
  Transport transport = 3; // uuid / state_number / type — из transport.txt
  Schedule schedule = 4;   // number / state_number / from / to — активная запись из schedule.txt
}
```

Подписчиков может быть несколько одновременно — каждый получает одинаковый поток событий. В поток попадают только
координаты, которые сервис смог полностью сопоставить со справочниками (транспорт → расписание → маршрут); координаты
с ошибкой сопоставления (см. [file_format.md](./file_format.md#как-эти-файлы-связаны-между-собой)) в
`StreamBusTrackingInfo` не передаются.

## Проверка вручную (без написания клиента)

Если включена рефлексия (`GRPC_REFLECTION=true`), сервис можно опросить [`grpcurl`](https://github.com/fullstorydev/grpcurl)
без `.proto`-файла на клиенте:

```bash
# список сервисов
grpcurl -plaintext localhost:9090 list

# подписаться на поток обогащённых данных
grpcurl -plaintext -d '{}' localhost:9090 BusTrackingService/StreamBusTrackingInfo
```

**Внимание:** в продакшене `GRPC_REFLECTION` обычно стоит выключать — рефлексия раскрывает структуру API кому угодно,
у кого есть сетевой доступ к порту.
