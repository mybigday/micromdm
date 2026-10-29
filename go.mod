module github.com/micromdm/micromdm

require (
	github.com/RobotsAndPencils/buford v0.12.0
	github.com/boltdb/bolt v1.3.1
	github.com/fullsailor/pkcs7 v0.0.0-20180824154052-36585635cb64
	github.com/garyburd/go-oauth v0.0.0-20180319155456-bca2e7f09a17
	github.com/go-kit/kit v0.7.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.6.2
	github.com/groob/finalizer v0.0.0-20170707115354-4c2ed49aabda
	github.com/groob/plist v0.0.0-20180203051248-dd56909aee38
	github.com/micromdm/go4 v0.0.0-20191221011012-654e10aaab18
	github.com/micromdm/scep v1.0.1-0.20181014170139-9be65e185499
	github.com/pkg/errors v0.8.0
	github.com/pressly/goose v2.3.0+incompatible
	github.com/satori/go.uuid v1.2.0
	golang.org/x/crypto v0.58.0
	golang.org/x/net v0.61.0
)

require (
	github.com/go-logfmt/logfmt v0.3.0 // indirect
	github.com/go-sql-driver/mysql v1.4.0 // indirect
	github.com/go-stack/stack v1.7.0 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515 // indirect
	github.com/lib/pq v1.0.0 // indirect
	github.com/mattn/go-sqlite3 v1.9.0 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/appengine v1.2.0 // indirect
	google.golang.org/protobuf v1.23.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/fullsailor/pkcs7 => github.com/groob/pkcs7 v0.0.0-20180824154052-36585635cb64

go 1.26.0
