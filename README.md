# Instructions to run locally

## Directly with go
1. Ensure you have go installed. This project uses go version `1.24.2`
1. At the root of the project, run `go mod download`
1. Run `go run cmd/main.go`
1. Server will be live at `localhost:8080`

## Docker
1. Ensure docker daemon is running.
1. Run `docker build -t receipt-processor:latest .`
1. Run `docker run -it -p 8080:8080 receipt-processorf:latest`
1. Server will be live at `localhost:8080`

---

### Running tests
Run `go test ./server/` to run all the tests

### To see test coverage:
1. Run `go test -coverprofile=testcoverage ./server/`
1. Run `go tool cover -html=testcoverage -o testcoverage.html`
1. Open the generated `testcoverage.html` in desired browser
