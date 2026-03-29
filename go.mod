module github.com/[[ .OrgName ]]/[[ .ProjectName ]]

go [[ .GoVersion ]]

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0
	github.com/kamalyes/go-config v0.13.1
	github.com/kamalyes/go-i18n v0.0.0-20260320070028-fea3ef2009b5
	github.com/kamalyes/go-logger v0.4.6
	github.com/kamalyes/go-pbmo v0.0.0-20260423155936-2c4d27513f74
	github.com/kamalyes/go-rpc-gateway v1.3.8
	github.com/kamalyes/go-sqlbuilder v0.3.6
	github.com/kamalyes/go-toolbox v0.11.88
	[[ .ShareProtoModule ]] v0.0.0-20260420141641-c095e3547996
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
	gorm.io/gorm v1.31.1
)

replace [[ .ShareProtoModule ]] => ../[[ .ObjectPrefix ]]-share-proto
