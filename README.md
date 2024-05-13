# comp-club-system

Тестовое задание в YADRO

### Требования к системе

1. Go 1.22
2. Docker

### Запуск

#### Локальный запуск
Вместо ``<filename>`` нужно писать имя файла с входными данными, который находится в папке, из которой запускается программа. 
Для примера можно использовать файл ``test_file.txt``.

##### Linux
```shell
go build -o task ./cmd
./task <filename>
```

##### Windows
```shell
go build -o task.exe ./cmd
task.exe <filename>
```

##### В Docker-контейнере

```shell
docker build -t comp-club-system .   
docker run -v "$(pwd)"/<filename>:/app/<filename> comp-club-system <filename>
```

##### Запуск тестов

```shell
go test ./test -v
```