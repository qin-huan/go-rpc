module go-rpc

go 1.16

require (
	github.com/echoingtech/go-micro v1.18.8
	github.com/echoingtech/hakot v1.10.37
	github.com/echoingtech/messager v0.0.0-20211110085909-64e57a1b851a
	github.com/echoingtech/pgw v0.2.0
	github.com/echoingtech/uc v1.2.0
	github.com/filecoin-project/go-jsonrpc v0.1.2-0.20201008195726-68c6a2704e49
	github.com/filecoin-project/lotus v1.1.2
	github.com/go-redis/redis/v8 v8.11.3
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-xorm/xorm v0.7.9
	github.com/gocql/gocql v0.0.0-20211015133455-b225f9b53fa1
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.1.2
	github.com/gorilla/websocket v1.4.2
	github.com/spf13/cobra v1.2.1
	github.com/xuri/excelize/v2 v2.4.1
	github.com/yanyiwu/gojieba v1.1.2
	go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/protobuf v1.27.1
	modernc.org/mathutil v1.1.1
)

replace go.etcd.io/bbolt v1.3.4 => github.com/coreos/bbolt v1.3.4

//replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/filecoin-project/filecoin-ffi => github.com/filecoin-project/filecoin-ffi v0.30.4-0.20200910194244-f640612a1a1f
