# Установка protoc на Windows

## Быстрая установка

### Вариант 1: Автоматическая установка (PowerShell)

```powershell
# Запустите от имени администратора
make install-protoc
```

После установки добавьте `C:\protoc\bin` в PATH:
1. Откройте "Переменные среды" (Environment Variables)
2. Добавьте `C:\protoc\bin` в переменную PATH
3. Перезапустите терминал

### Вариант 2: Через Chocolatey (если установлен)

```powershell
choco install protoc
```

### Вариант 3: Через Scoop (если установлен)

```powershell
scoop install protobuf
```

### Вариант 4: Ручная установка

1. Перейдите на https://github.com/protocolbuffers/protobuf/releases
2. Скачайте `protoc-XX.X-win64.zip` (последняя версия, например v25.1)
3. Распакуйте архив в `C:\protoc`
4. Добавьте `C:\protoc\bin` в переменную PATH
5. Перезапустите терминал

### Проверка установки

```powershell
protoc --version
```

Должно вывести версию, например: `libprotoc 25.1`

## После установки protoc

1. Убедитесь, что Go плагины установлены:
```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

2. Сгенерируйте proto файлы:
```powershell
make gen
```

## Альтернатива: Использование Docker для генерации

Если не хотите устанавливать protoc локально, можно использовать Docker:

```powershell
docker run --rm -v ${PWD}:/workspace -w /workspace znly/protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/**/*.proto
```


