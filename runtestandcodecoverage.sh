go test -cover ./message/
go test -cover ./server/
go test -cover ./client/
go test -cover ./utils/

go test -bench=. ./message/
go test -bench=. ./client/
go test -bench=. ./server/
go test -bench=. ./utils/
