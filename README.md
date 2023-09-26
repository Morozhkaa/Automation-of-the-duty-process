# Автоматизация процесса дежурств.

Задача: разработать приложение, которое по ___linkedin oncall__ REST API_ создаст команды/сотрудников команд и их дежурства согласно описанию, приведенному в config.yaml файле.


### Usage:

1. Предполагается, что локально у вас запущено linkedin oncall приложение. Если нет, то вот предпринятые мной шаги для этого:

- Склонировала следующий репозиторий:

```
git clone -b v2.0.1 https://github.com/linkedin/oncall.git
```

- В docker-compose файле скачанного репозитория заменила образ mysql на mysql:8.1.0. И в config/config.docker.yaml файле установила значение:
```
auth:
  debug: True
```

2. Склонируйте текущий репозиторий и в config.yaml файле опишите желаемый график дежурств в следующем формате: 
```
---
teams:
  - name: "k8s SRE"
    scheduling_timezone: "Europe/Moscow"
    email: "k8s@sre-course.ru"
    slack_channel: "#k8s-team"
    users:
      - name: "o.ivanov"
        full_name: "Oleg Ivanov"
        phone_number: "+1 111-111-1111"
        email: "o.ivanov@sre-course.ru"
        duty:
          - date: "02/10/2023"
            role: "primary"
          - date: "03/10/2023"
            role: "secondary"
          - date: "04/10/2023"
            role: "primary"
          - date: "05/10/2023"
            role: "secondary"
          - date: "06/10/2023"
            role: "primary"
      - name: "d.petrov"
        full_name: "Dmitriy Petrov"
        phone_number: "+1 211-111-1111"
        email: "d.petrov@sre-course.ru"
        duty:
          - date: "02/10/2023"
            role: "secondary"
          - date: "03/10/2023"
            role: "primary"
          - date: "04/10/2023"
            role: "secondary"
          - date: "05/10/2023"
            role: "primary"
          - date: "06/10/2023"
            role: "secondary"
```

3. Запустите наше приложение-скрипт, которое по переданному описанию добавит записи о дежурствах в onCall:
```
go run main.go
```

### Comments:
Не стала заворачивать в докер, тк пришлось бы дублировать код oncall репозитория, чтобы запустить их в одной network.