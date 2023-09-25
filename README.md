# Go Deploy

Go Deploy is a Go project that automates the execution of shell commands to perform deployment tasks.

## Features

- Automated execution of shell commands
- Configuration of environment variables through the .env file

## Installation

To install Go Deploy, follow these steps:

1. Clone the repository to your local machine:

```shell
git clone https://github.com/cc-hearts/go-deploy.git
```

2. Navigate to the project directory:

```shell
cd go-deploy
```

3. Build and install the project:

```shell
go install

```

## Usage

Create a .env file in the root directory of your project and configure the environment variables according to your requirements. For example:

```go
SSH_HOST = "" // ip:port
SSH_USER = ""
SSH_PASSWORD = ""
SSH_MYSQL_HOST = "" // user:password@tcp(ip:port)/database

PORT = "" // :port
```

Run Go Deploy:

> Go Deploy will automatically load the environment variables from the `.env` file and execute the configured deployment command.

```shell
   go run main.go
```

## Build

If you are using Linux, you can use the following command:

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o <binary_name>  main.go
```

## License

[MIT](./LICENSE)
