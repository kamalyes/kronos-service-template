module github.com/<% .OrgName %>/<% .ProjectName %>

go <% .GoVersion %>

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0
	github.com/kamalyes/go-bizerror v0.0.0-20260814173213-9f68830eaadf
	github.com/kamalyes/go-config v0.23.0
	github.com/kamalyes/go-i18n v0.2.8
	github.com/kamalyes/go-logger v0.7.1
	github.com/kamalyes/go-pbmo v0.2.1
	github.com/kamalyes/go-rpc-gateway v1.9.6
	github.com/kamalyes/go-sqlbuilder v0.7.0
	github.com/kamalyes/go-toolbox v0.16.3
	<% .ShareProtoModule %> v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
	gorm.io/gorm v1.31.1
)

// 协议中心依赖的版本由 scripts/init.sh 解析并写入（见该脚本注释）：
// 上面那个 v0.0.0-00010101000000-000000000000 是占位，没解析过是跑不起来的
//
// [EN] The contract-hub version is resolved and written by scripts/init.sh; the
// placeholder above does not resolve on its own.

// 本地联调协议中心时再打开它：指向同级目录的协议仓库
// 模板默认不带，因为生成出来就带着一条指向不存在目录的 replace，
// 别人没有那个兄弟目录时 go mod tidy 直接失败
//
// [EN] Uncomment while developing against a local contract hub. Shipping it by default
// would break `go mod tidy` for anyone without that sibling directory.
// replace <% .ShareProtoModule %> => ../<% .ObjectPrefix %>-share-proto
