## Project layout and architect

- Layout: https://github.com/golang-standards/project-layout
- Architect: [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## Framework

- HTTP Web Framework: https://github.com/gin-gonic/gin
- Database migrations: https://github.com/golang-migrate/migrate
- Go Struct and Field validation: https://github.com/go-playground/validator
- ORM: https://gorm.io/
- Log: https://github.com/sirupsen/logrus
- JWT: https://github.com/golang-jwt/jwt
- AWS SDK: [v1](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/welcome.html) and [v2](https://aws.github.io/aws-sdk-go-v2/docs/getting-started/)
- Test: https://github.com/stretchr/testify
- Mapstructure: https://github.com/mitchellh/mapstructure
- Sendgird: https://github.com/sendgrid/sendgrid-go
- Google Cloud Client: https://github.com/googleapis/google-cloud-go
- Struct Data Fake Generator: https://github.com/bxcodec/faker
- JSON Schema: https://github.com/xeipuuv/gojsonschema

## Requirements:

- [Golang](https://go.dev/)
- [Makefile](http://gnu.org/licenses/gpl.html)
- [Docker](https://docs.docker.com/get-started/)

## Launch project:

- cp .env.example .env
- make setup-database
- make run-migration
- go run main.go
