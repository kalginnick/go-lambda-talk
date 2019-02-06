# ETL-pipeline

## Предустановки

- [Golang](https://golang.org/doc/install)
- [Make](https://www.gnu.org/software/make/)
- [SAM](https://docs.aws.amazon.com/en_us/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Структура проекта

```bash
├── cmd
│   ├── extract                         <-- Лямбда получения прайса от поставщика
│   │   ├── extract-dev.yml             <-- Конфигурация для локального тестирования
│   │   ├── main.go                     <-- Исходный код функции получения
│   │   └── scheduled-event.json        <-- Тестовое событие
│   ├── load                            <-- Лямбда выгрузки прайса на сервер интернет-магазина
│   │   ├── load-dev.yml                <-- Конфигурация для локального тестирования
│   │   ├── main.go                     <-- Исходный код функции выгрузки готового прайса
│   │   └── s3-event.json               <-- Тестовое событие
│   └── transform                       <-- Лямбда преобразования формата
│       ├── main.go                     <-- Исходный код функции преобразования
│       ├── s3-event.json               <-- Тестовое событие
│       └── transform-dev.yml           <-- Конфигурация для локального тестирования
├── go.mod                              <-- Зависимости для сборки
├── go.sum
├── Makefile                            <-- Конфигурация сборки make
├── pkg
│   └── client                          <-- Пакет для работы с удалёнными хранилищами
│       ├── ftp.go
│       └── s3.go
├── README.md                           <-- Этот файл
├── slides.pdf                          <-- Презентация
└── template.yml                        <-- Конфигурация развертывания в AWS
```

## Сборка

```bash
make build
```

## Тестирование и запуск

Для локального тестирования нужно собрать функцию и запустить её, используя sam, например:

```bash
make build
cd cmd/extract
sam local invoke -t extract-dev.yml -e scheduled-event.json
```

## Деплой

Для деплоя в облаке нужно создать аккаунт в [AWS](https://portal.aws.amazon.com/billing/signup#/start)
После регистрации в [панели администратора](https://console.aws.amazon.com/iam/home#/security_credentials)
лучше создать отдельные `Access key ID` и `Secret access key` для программного доступа через API и поместить их в файл
`~/.aws/credentials`:

```properties
[default]
aws_access_key_id = YOUR_ACCESS_KEY
aws_secret_access_key = YOUR_SECRET_KEY
```

Регион по умолчанию, в котором будет происходить деплой, задать в `~/.aws/config`

```properties
[default]
region = eu-central-1
```

Эти значения будут использованы по умолчанию при работе утилит `sam` и `aws-cli`.

Процессу деплоя нужно передать через переменные окружения несколько значений, это:

- `DEPLOYMENT_NAME` - имя, идентифицирующее наш набор функций.
  Будет использовано для создания единицы развёртывания CloudFormation Stack и S3 Bucket с бинарниками функций.
- `IMPORT_SOURCE_BUCKET_NAME` - имя S3 Bucket, использующегося для хранения исходных файлов импортf.
- `IMPORT_RESULT_BUCKET_NAME` - имя S3 Bucket, использующегося для хранения файлов? преобразованных в ходе импорта.
- `LOAD_URL` - URL FTP-сервера, на который будут загружены файлы, полученные в результате импорта.
- `LOAD_USER` - логин для доступа к FTP-серверу.
- `LOAD_PSWD` - пароль для доступа к FTP-серверу.

```bash
export DEPLOYMENT_NAME=my-lambda-stack-name \
       IMPORT_SOURCE_BUCKET_NAME=my-unique-source-bucket-name \
       IMPORT_RESULT_BUCKET_NAME=my-unique-result-bucket-name \
       LOAD_URL=load-ftp-url \
       LOAD_USER=load-ftp-user \
       LOAD_PSWD=load-ftp-password

make deploy
```

Если при деплое произошла ошибка, то причину её возникновения можно найти в консоли [CloudFormation Stack](https://console.aws.amazon.com/cloudformation/home)

Созданные функции можно увидеть в разделе [AWS Lambda](https://console.aws.amazon.com/lambda/home) и даже запустить вручную.

Мониторинг и логи работы функций доступны в [CloudWatch](https://console.aws.amazon.com/cloudwatch/home)
