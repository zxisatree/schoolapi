# SchoolAPI

This project is an API for a school's mailing list. It is written in Golang, using Gin and Gorm, with MySQL as the intended database.

# Quick start

Before getting started, ensure that you have [Go](https://go.dev/dl/) and [MySQL](https://www.mysql.com/downloads/) installed. Then, create a database in MySQL for the API to use.

To start a local instance of the API server:

1. Clone the repository to your local machine
2. Create a `.env` file in the root directory with the following values:
   - `PORT=3700` or any unused port number on your local machin
   - `MYSQL_DB_URL="<username>:<password>@tcp(127.0.0.1:<mysql_port>)/<mysql_db_name>?charset=utf8mb4&parseTime=True&loc=Local"` where `<username>` is your username for MySQL, `<password>` is the password for that account, `<mysql_port>` is the port that MySQL is running on (default 3306) and `<mysql_db_name>` is the name of the database you created earlier.
3. Navigate (`cd`) to the root directory of the cloned repository
4. Run `go run migrate/migrate.go` to create the necessary tables in the database.
5. EITHER run `go run seed/seedData.go` to seed the database with sample data OR insert data into the database manually. The schema of the tables can be found via `DESC <table_name>` where `<table_name>` is either `students`, `teachers` or `teachers_students`.
6. Run `go build`, then run the executable called `schoolapi` generated in the root directory to start the server on the port chosen in step 2.

# Tests

Navigate to the `tests` directory, copy the `.env` file from your root directory(created in step 2 of Quick Start), then run `go test` to run all the unit tests.
