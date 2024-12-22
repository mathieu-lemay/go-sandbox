test:
    go test -coverprofile=.coverage ./...

bench:
    go test -bench . ./...

cov-render:
    go tool cover -html .coverage -o coverage.html
    xdg-open coverage.html
