# polling-app

## Environment setup

You need to have [Go](https://golang.org/),
[Node.js](https://nodejs.org/),
installed on your computer.

Verify the tools by running the following commands:

```sh
go version
npm --version
node --version
```

## Start in Local mode

Start a local PostgreSQL database on `localhost:5432`.
The database will be populated with test records from the
[init-db.sql](init-db.sql) file.

Navigate to the `server` folder and start the back end:

```sh
cd server
go run server.go
```
The back end will serve on http://localhost:8080.

Navigate to the `webapp` folder, install dependencies,
and start the front end development server by running:

```sh
cd webapp
npm install
npm start
```
The application will be available on http://localhost:3000.
