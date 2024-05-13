# comp-club-system

Тестовое задание в YADRO

### Требования к системе

1. Go 1.22
2. Docker

### Запуск

#### Локальный запуск
##### Linux
```shell
go build -o task ./cmd
./task test_data.txt
```

##### Windows
```shell
go build -o task.exe ./cmd
task.exe test_data.txt
```

##### В Docker-контейнере

TODO

##### Запуск тестов

```shell
go test ./test -v
```