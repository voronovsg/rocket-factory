module github.com/voronovsg/rocket-factory/iam

go 1.26.1

replace (
	github.com/voronovsg/rocket-factory/platform => ../platform
	github.com/voronovsg/rocket-factory/shared => ../shared
)

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/caarlos0/env/v11 v11.4.1
	github.com/gomodule/redigo v1.9.3
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.10.0
	github.com/joho/godotenv v1.5.1
	github.com/pkg/errors v0.9.1
	github.com/voronovsg/rocket-factory/platform v0.0.0-00010101000000-000000000000
	github.com/voronovsg/rocket-factory/shared v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.28.0
	golang.org/x/crypto v0.52.0
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pressly/goose/v3 v3.27.1 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260420184626-e10c466a9529 // indirect
)
