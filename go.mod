module go-rpc

go 1.15

require (
	github.com/filecoin-project/go-jsonrpc v0.1.2-0.20201008195726-68c6a2704e49
	github.com/filecoin-project/lotus v1.1.2
	go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738
	modernc.org/mathutil v1.1.1
)

replace go.etcd.io/bbolt v1.3.4 => github.com/coreos/bbolt v1.3.4

//replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/filecoin-project/filecoin-ffi => github.com/filecoin-project/filecoin-ffi v0.30.4-0.20200910194244-f640612a1a1f
