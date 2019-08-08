module github.com/komuw/meli

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.12 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.7.3-0.20181221150755-2cb26cfe9cbf
	github.com/docker/docker-credential-helpers v0.6.3
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/gorilla/mux v1.7.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	golang.org/x/net v0.0.0-20190424112056-4829fb13d2c6 // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.0.0-20190804053845-51ab0e2deafa // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/genproto v0.0.0-20190801165951-fa694d86fc64 // indirect
	google.golang.org/grpc v1.22.1 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.2
	gotest.tools v2.2.0+incompatible // indirect
)

// look at: https://github.com/golang/go/issues/29376#issuecomment-449416502

// github.com/docker/engin v19.03.1
replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190725163905-fa8dd90ceb7b

go 1.13
