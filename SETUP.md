# Скрипт для установки protoc на Windows

# Вариант 1: Через Scoop (рекомендуется)
# Если у вас установлен Scoop:
scoop install protobuf

# Вариант 2: Через Chocolatey
# Если у вас установлен Chocolatey:
choco install protoc

# Вариант 3: Ручная установка
# 1. Скачайте protoc с https://github.com/protocolbuffers/protobuf/releases
# 2. Распакуйте в C:\protoc
# 3. Добавьте C:\protoc\bin в PATH

# Вариант 4: Используйте Docker (не требует установки protoc)
make gen-docker


