# calculator-with-authorization

## GUI

Если переменныем окружения не были изменены, графический интерфейс доступен по адресу ```http://localhost:8080```

У пользователя есть возможность зарегистрироваться или войти в уже существующий аккаунт при помощи логина и пароля. После этого происходит переадресация на страницу самого калькулятора

- Поддерживаются операции сложения, вычитания, умножения и деления, а также выражения в скобках
- Выражение может вводиться как с пробелами между числом и операндом, так и без
- Калькулятор принимает на вход положительные целые числа

## Зависимости

- Go версии ```1.23``` или новее
- Дополнительные библиотеки (указаны в ```go.mod```)

## Переменные окружения

| Переменная                             | Описание                                            | Значение по умолчанию |
|----------------------------------------|-----------------------------------------------------|-----------------------|
| ```server.HTTP_PORT```                 | Порт для запуска HTTP сервера                       | 8080                  |
| ```server.HTTP_HOST```                 | Хост для запуска HTTP сервера оркестратора          | localhost             |
| ```server.GRPC_PORT```                 | Порт для запуска gPRC сервера                       | 50051                 |
| ```server.GRPC_HOST```                 | Хост для запуска gRPC сервера                       | localhost             |
| ```duration.TIME_ADDITION_MS```        | Время выполнения операции сложения в миллисекундах  | 100                   |
|```duration.TIME_SUBTRACTION_MS```      | Время выполнения операции вычитания в миллисекундах | 100                   |
| ```duration.TIME_MULTIPLICATIONS_MS``` | Время выполнения операции умножения в миллисекундах | 100                   |
| ```duration.TIME_DIVISIONS_MS```       | Время выполнения операции деления в миллисекундах   | 100                   |
| ```DATABASE_PATH```                    | Путь к базе данных                                  |./db/calc.db           |
| ```worker.COMPUTING_POWER```           | Количество горутин, выполняющих вычисления          | 5                     |
| ```jwt.secret_key```                   | Используется для создания цифровой подписи токена   | your_secret_key_here  |
| ```jwt.token_duration```               | Время жизни токена                                  | 24                    |

## Тестирование

```bash
go test -cover .\internal\agent\

# ok        coverage: 65.8% of statements
```

```bash
go test -cover .\internal\auth\

# ok        coverage: 89.5% of statements
```
```bash
go test -cover .\internal\orchestrator\transport\http\handlers\

# ok        coverage: 65.3% of statements
```
```bash
go test -cover .\internal\orchestrator\transport\grpc\

# ok        coverage: 19.2% of statements
```
```bash
go test -cover .\internal\orchestrator\service\ 

# ok        coverage: 88.6% of statements
```