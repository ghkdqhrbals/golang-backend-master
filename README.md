[![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/golang-migrate/migrate/CI/master)](https://github.com/golang-migrate/migrate/actions/workflows/ci.yaml?query=branch%3Amaster)
[![GoDoc](https://pkg.go.dev/badge/github.com/golang-migrate/migrate)](https://pkg.go.dev/github.com/golang-migrate/migrate/v4)
[![Coverage Status](https://img.shields.io/coveralls/github/golang-migrate/migrate/master.svg)](https://coveralls.io/github/golang-migrate/migrate?branch=master)
[![packagecloud.io](https://img.shields.io/badge/deb-packagecloud.io-844fec.svg)](https://packagecloud.io/golang-migrate/migrate?filter=debs)
[![Docker Pulls](https://img.shields.io/docker/pulls/migrate/migrate.svg)](https://hub.docker.com/r/migrate/migrate/)
![Supported Go Versions](https://img.shields.io/badge/Go-1.16%2C%201.17-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/golang-migrate/migrate.svg)](https://github.com/golang-migrate/migrate/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/golang-migrate/migrate)](https://goreportcard.com/report/github.com/golang-migrate/migrate)
# Updates
* gomock 라이브러리를 이용하여 서비스 레이어의 테스트에서 DB 의존성을 제거함.
즉, fake DB(in memory)를 통해 테스트함.

### sqlc로 쿼리문 인터페이스 생성
```bash
type Querier interface {
    AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error)
    CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
    CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
    CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
    DeleteAccount(ctx context.Context, id int64) error
    GetAccount(ctx context.Context, id int64) (Account, error)
    GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
    GetEntry(ctx context.Context, id int64) (Entry, error)
    GetTransfer(ctx context.Context, id int64) (Transfer, error)
    ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
    ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
    ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
    UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
}
```
sqlc의 emit_interface=true를 통해 query의 account.go, transfer.sql, entry.sql의 인터페이스를 생성.


### mockgen으로 기본함수생성
```bash
type Store interface {
	Querier // 인터페이스
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}
```
이후 simplebank/ghkdqhrbals/db/sqlc/store.go에 Store 인터페이스 정의
이 인터페이스는 의존성을 제거한 db테스트에 사용



```bash
mockgen -destination db/mock/store.go github.com/ghkdqhrbals/simplebank/db/sqlc Store
```
mockgen을 통해 앞서 정의한 Store 인터페이스를 받아오고, 가상으로 실행하는 함수를 자동으로 정의.


### mock DB
```bash
func (server *Server) getAccount(ctx *gin.Context) {
	~
	account, err := server.store.GetAccount(ctx, req.ID)
	~
}
type Querier interface {
	~
	GetAccount(ctx context.Context, id int64) (Accounts, error)
	~
}
```
store의 GetAccount는 다음과 같음. 즉, Accounts 구조와 에러를 반환.

```bash
func (q *Queries) GetAccount(ctx context.Context, id int64) (Accounts, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Accounts
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}
```



Unit Test
```bash
{
    name:      "OK",
    accountID: account.ID,
    buildStubs: func(store *mockdb.MockStore) {
        store.EXPECT().
            // mock으로 테스트할 function 이름.
            //
            GetAccount(gomock.Any(), gomock.Eq(account.ID)). // account.ID를 argument로 받는 store.GetAccount함수가 call 되길 예상함.
            Times(1).                                        // 위의 함수가 call 되는 횟수 예상
            Return(account, nil)                             // test function의 return 값 예상, 실제 GetAccount에서 account.ID로 검색했을 때 account구조체와 nil에러를 반환받기를 원함.
    },
    checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
        require.Equal(t, http.StatusOK, recorder.Code)
        requireBodyMatchAccount(t, recorder.Body, account)
    },

},
```


```bash
go test -run "path/function name" -v(detaily describe) -cover(coverage) //test
```




# Information

mock은 전체 테스트: coverage = 100% && 행동 관찰

stubs는 특정 기능부분 테스트: coverage <= 100% && 상태 관찰

의미는 이러하지만 사실 이 두 가지 테스트 방법은 따로 구분되지 않음.
ex) mock 또한 < 100% 가능

# Gin 사용법

### Microservices 제어흐름
__Request -> Route Parser -> [Optional Middleware] -> Route Handler -> [Optional Middleware] -> Response__

### Gin router 생성 -> 루트 정의 및 해당 루트 엑세스 시 handler 설정
```bash
router := gin.Default() // router 생성

// router.[GET or POST or etc](url string,HANDLER)
// Route handler 정의
router.GET("/accounts/:id",server.getAccount)
```
" localhost:8080/accounts/의 뒤에 int64로 GET이 들어오면, id로 tagging하고 server의 getAccount를 실행하여라 " 라는 의미임.


```bash
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
```

getAccount가 http로 오면 서버에서 다음과 같은 형식에 저장할 것이라고 명시한다.

이때, 기본적인 validate를 통해 조건을 정할 수 있다.(ID는 1이상이고, 비어서는 안되며, ctx.ShouldBindUri을 사용하여 req에 저장할 것이기에 uri: 명시)

__이는 기본적인 validate방식이며, gin에서 custom 가능하다.__
* [Basic validation in Golang](https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/#basic-validation-using-gin)
* [Writing custom validation with reflection Module in Golang](https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme)

```bash
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

    v.RegisterValidation("currency", validCurrency) // Register Custom Validator {tag, validator.Func}
}

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// fieldLevel.Field() to get the value of the field
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	WON = "WON"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, WON:
		return true
	}
	return false
}
```
__currency를 확인하는 custom validator 생성__

즉, Oneof, min과 같은 basic validator들과 같이 __currency__라는 validator을 설정한다는 뜻

이러한 currency는 validCurrency라는 조건을 핸들링하고 true시 keep going

```bash
type createAccountRequest struct {
	~
	Currency string ` ~ , currency"` // currency custom validator 추가
}
```

```bash
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil { // 우리는 이전 tagging된 id를 getAccountRequest ID에 바인딩.
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.ID) // DB 쿼리 트랜젝션 실행
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}
```
Unmarshelling(gin에서 지원) => ShouldBindUri

이는 다양한 방법 존재. __ShouldBindJSON, ShouldBindUri, ShouldBindQuery, etc__
* [ShouldBindWith](https://pkg.go.dev/github.com/gin-gonic/gin#section-readme)

server의 getAccount는 server.store.GetAccount로 DB 쿼리 및 기타 ERROR 전송





# migrate

__Database migrations written in Go. Use as [CLI](#cli-usage) or import as [library](#use-in-your-go-project).__

* Migrate reads migrations from [sources](#migration-sources)
   and applies them in correct order to a [database](#databases).
* Drivers are "dumb", migrate glues everything together and makes sure the logic is bulletproof.
   (Keeps the drivers lightweight, too.)
* Database drivers don't assume things or try to correct user input. When in doubt, fail.

Forked from [mattes/migrate](https://github.com/mattes/migrate)

## Databases

Database drivers run migrations. [Add a new database?](database/driver.go)

* [PostgreSQL](database/postgres)
* [PGX](database/pgx)
* [Redshift](database/redshift)
* [Ql](database/ql)
* [Cassandra](database/cassandra)
* [SQLite](database/sqlite)
* [SQLite3](database/sqlite3) ([todo #165](https://github.com/mattes/migrate/issues/165))
* [SQLCipher](database/sqlcipher)
* [MySQL/ MariaDB](database/mysql)
* [Neo4j](database/neo4j)
* [MongoDB](database/mongodb)
* [CrateDB](database/crate) ([todo #170](https://github.com/mattes/migrate/issues/170))
* [Shell](database/shell) ([todo #171](https://github.com/mattes/migrate/issues/171))
* [Google Cloud Spanner](database/spanner)
* [CockroachDB](database/cockroachdb)
* [ClickHouse](database/clickhouse)
* [Firebird](database/firebird)
* [MS SQL Server](database/sqlserver)

### Database URLs

Database connection strings are specified via URLs. The URL format is driver dependent but generally has the form: `dbdriver://username:password@host:port/dbname?param1=true&param2=false`

Any [reserved URL characters](https://en.wikipedia.org/wiki/Percent-encoding#Percent-encoding_reserved_characters) need to be escaped. Note, the `%` character also [needs to be escaped](https://en.wikipedia.org/wiki/Percent-encoding#Percent-encoding_the_percent_character)

Explicitly, the following characters need to be escaped:
`!`, `#`, `$`, `%`, `&`, `'`, `(`, `)`, `*`, `+`, `,`, `/`, `:`, `;`, `=`, `?`, `@`, `[`, `]`

It's easiest to always run the URL parts of your DB connection URL (e.g. username, password, etc) through an URL encoder. See the example Python snippets below:

```bash
$ python3 -c 'import urllib.parse; print(urllib.parse.quote(input("String to encode: "), ""))'
String to encode: FAKEpassword!#$%&'()*+,/:;=?@[]
FAKEpassword%21%23%24%25%26%27%28%29%2A%2B%2C%2F%3A%3B%3D%3F%40%5B%5D
$ python2 -c 'import urllib; print urllib.quote(raw_input("String to encode: "), "")'
String to encode: FAKEpassword!#$%&'()*+,/:;=?@[]
FAKEpassword%21%23%24%25%26%27%28%29%2A%2B%2C%2F%3A%3B%3D%3F%40%5B%5D
$
```





## Link
[Basic validation in Golang](https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/#basic-validation-using-gin)
[Writing custom validation with reflection Module in Golang](https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme)






## Migration Sources

Source drivers read migrations from local or remote sources. [Add a new source?](source/driver.go)

* [Filesystem](source/file) - read from filesystem
* [io/fs](source/iofs) - read from a Go [io/fs](https://pkg.go.dev/io/fs#FS)
* [Go-Bindata](source/go_bindata) - read from embedded binary data ([jteeuwen/go-bindata](https://github.com/jteeuwen/go-bindata))
* [pkger](source/pkger) - read from embedded binary data ([markbates/pkger](https://github.com/markbates/pkger))
* [GitHub](source/github) - read from remote GitHub repositories
* [GitHub Enterprise](source/github_ee) - read from remote GitHub Enterprise repositories
* [Bitbucket](source/bitbucket) - read from remote Bitbucket repositories
* [Gitlab](source/gitlab) - read from remote Gitlab repositories
* [AWS S3](source/aws_s3) - read from Amazon Web Services S3
* [Google Cloud Storage](source/google_cloud_storage) - read from Google Cloud Platform Storage

## CLI usage

* Simple wrapper around this library.
* Handles ctrl+c (SIGINT) gracefully.
* No config search paths, no config files, no magic ENV var injections.

__[CLI Documentation](cmd/migrate)__

### Basic usage

```bash
$ migrate -source file://path/to/migrations -database postgres://localhost:5432/database up 2
```

### Docker usage

```bash
$ docker run -v {{ migration dir }}:/migrations --network host migrate/migrate
    -path=/migrations/ -database postgres://localhost:5432/database up 2
```

## Use in your Go project

* API is stable and frozen for this release (v3 & v4).
* Uses [Go modules](https://golang.org/cmd/go/#hdr-Modules__module_versions__and_more) to manage dependencies.
* To help prevent database corruptions, it supports graceful stops via `GracefulStop chan bool`.
* Bring your own logger.
* Uses `io.Reader` streams internally for low memory overhead.
* Thread-safe and no goroutine leaks.

__[Go Documentation](https://godoc.org/github.com/golang-migrate/migrate)__

```go
import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/github"
)

func main() {
    m, err := migrate.New(
        "github://mattes:personal-access-token@mattes/migrate_test",
        "postgres://localhost:5432/database?sslmode=enable")
    m.Steps(2)
}
```

Want to use an existing database client?

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    db, err := sql.Open("postgres", "postgres://localhost:5432/database?sslmode=enable")
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    m, err := migrate.NewWithDatabaseInstance(
        "file:///migrations",
        "postgres", driver)
    m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
}
```

## Getting started

Go to [getting started](GETTING_STARTED.md)

## Tutorials

* [CockroachDB](database/cockroachdb/TUTORIAL.md)
* [PostgreSQL](database/postgres/TUTORIAL.md)

(more tutorials to come)

## Migration files

Each migration has an up and down migration. [Why?](FAQ.md#why-two-separate-files-up-and-down-for-a-migration)

```bash
1481574547_create_users_table.up.sql
1481574547_create_users_table.down.sql
```

[Best practices: How to write migrations.](MIGRATIONS.md)

## Versions

Version | Supported? | Import | Notes
--------|------------|--------|------
**master** | :white_check_mark: | `import "github.com/golang-migrate/migrate/v4"` | New features and bug fixes arrive here first |
**v4** | :white_check_mark: | `import "github.com/golang-migrate/migrate/v4"` | Used for stable releases |
**v3** | :x: | `import "github.com/golang-migrate/migrate"` (with package manager) or `import "gopkg.in/golang-migrate/migrate.v3"` (not recommended) | **DO NOT USE** - No longer supported |

## Development and Contributing

Yes, please! [`Makefile`](Makefile) is your friend,
read the [development guide](CONTRIBUTING.md).

Also have a look at the [FAQ](FAQ.md).

---

Looking for alternatives? [https://awesome-go.com/#database](https://awesome-go.com/#database).
