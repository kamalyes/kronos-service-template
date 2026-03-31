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
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
	gorm.io/gorm v1.31.1
)

// 这里故意不写协议中心的 require：由 scripts/init.sh 执行 go get 时自动加上
// 写一个占位版本（v0.0.0-000101...）反而更糟——go get 会先解析 go.mod 里已有的
// 版本约束，直接报 invalid version: unknown revision 000000000000 而拿不到 latest
//
// [EN] The contract-hub require is deliberately absent and added by scripts/init.sh via
// go get. A placeholder version would be worse: go get resolves the existing requirement
// first and fails with "invalid version: unknown revision 000000000000".

// 本地联调协议中心时再打开它：指向同级目录的协议仓库
// 模板默认不带，因为生成出来就带着一条指向不存在目录的 replace，
// 别人没有那个兄弟目录时 go mod tidy 直接失败
//
// [EN] Uncomment while developing against a local contract hub. Shipping it by default
// would break `go mod tidy` for anyone without that sibling directory.
// replace <% .ShareProtoModule %> => ../<% .ObjectPrefix %>-share-proto
