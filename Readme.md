go test -failfast -count=1 -v ./... -coverpkg=./... -coverprofile=coverpkg.out
go tool cover -html=coverpkg.out -o coverage.html