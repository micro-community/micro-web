module github.com/micro-community/micro-webui

go 1.15

require (
	github.com/go-acme/lego/v3 v3.9.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/micro/micro/v3 v3.0.1-0.20201106123606-1fc5978f576e
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
	github.com/serenize/snaker v0.0.0-20201027110005-a7ad2135616e
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/net v0.0.0-20201031054903-ff519b6c9102
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.29.0
