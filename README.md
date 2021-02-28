# Beego simple API example.

#### The app shows user registration and login process, email verification via url sent by MailGun.
#### It has example of using middleware in beego: getting user info form token, checking verified email and role "admin"

## 1: Initial Setup
```bash
git clone git@github.com:Artyem-Y/beego-crud-api.git
```

## 2: Install PostgreSQL:
```bash
sudo apt-get update
sudo apt-get install postgresql
```

#### create db with linux command:
```bash

sudo -u postgres createdb crud-api
```

#### Enter to the database:
```bash
sudo -u postgres psql -d crud-api
```
#### Create user of the db from postgres command line:
```bash
CREATE USER admin WITH PASSWORD '111';

\q
```

## 3: Getting Started

#### Copy .env-example tp .env and add your credentials

```bash
DB_SERVER = 127.0.0.1
DB_PORT = 5432
DB_NAME = crud-api
DB_USER = admin
DB_USER_PASS = 111

REFRESH_SECRET=
ACCESS_SECRET=

# add app url for local version
APP_URL=http://localhost:8080

# it's necessary to register mailgun account, get api key and create domain for getting mails
MAILGUN_API_KEY=
MAILGUN_DOMAIN=

# email for users notification
NOTIFICATION_EMAIL=
```

#### Inside beego-crud-api/db directory create "dbconf.yml" file with such content:
```bash
development:
    driver: postgres
    open: user=admin dbname=crud-api password=111 host=127.0.0.1 port=5432 sslmode=disable
```

#### Install Goose migration tool:
https://bitbucket.org/liamstask/goose/src/master/
```bash
go get bitbucket.org/liamstask/goose/cmd/goose
```
#### Install Bee command-line tool:
https://github.com/beego/bee
```bash

#### To install all dependencies of a Golang project or golang projects recursively with
#### the go get command, change directory into the project root and simply run::
```bash
go get .
```
#### Inside the project root directory migrate the database:
```bash
goose up
```
#### To delete tables, repeat below command for each table:
```bash
goose down
```
#### If we have several DB environment is "dbconf.yml", we can write needed, e.g.:
```bash
goose -env local up
```

#### Inside the project directory run app with terminal command:
```bash
bee run -downdoc=true -gendoc=true

or 

bee run
```
