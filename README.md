[![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/golang-migrate/migrate/CI/master)](https://github.com/golang-migrate/migrate/actions/workflows/ci.yaml?query=branch%3Amaster)
[![GoDoc](https://pkg.go.dev/badge/github.com/golang-migrate/migrate)](https://pkg.go.dev/github.com/golang-migrate/migrate/v4)
[![Coverage Status](https://img.shields.io/coveralls/github/golang-migrate/migrate/master.svg)](https://coveralls.io/github/golang-migrate/migrate?branch=master)
[![packagecloud.io](https://img.shields.io/badge/deb-packagecloud.io-844fec.svg)](https://packagecloud.io/golang-migrate/migrate?filter=debs)
[![Docker Pulls](https://img.shields.io/docker/pulls/migrate/migrate.svg)](https://hub.docker.com/r/migrate/migrate/)
![Supported Go Versions](https://img.shields.io/badge/Go-1.16%2C%201.17-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/golang-migrate/migrate.svg)](https://github.com/golang-migrate/migrate/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/golang-migrate/migrate)](https://goreportcard.com/report/github.com/golang-migrate/migrate)
## Versions

Version | Skills | Done?
--------|------------|------
**[v1.0.1](https://github.com/ghkdqhrbals/simplebank/tree/1.1v)** | Postresq, migration, Testing_enviroments, Sqlc, Git-Workflow | :white_check_mark: |
**[v1.0.2](https://github.com/ghkdqhrbals/simplebank/tree/1.2v)** | __Gin__, __Viper__, __Gomock__, Postresq, migration, Testing_enviroments, Sqlc, Git-Workflow | :white_check_mark: |
**[v1.0.3](https://github.com/ghkdqhrbals/simplebank/tree/1.2v)** | __Bcrypt__, Gin, Viper, Gomock, Postresq, migration, Testing_enviroments, Sqlc, Git-Workflow | :white_check_mark: |

## Update
__Bcrypt로 사용자 PW 저장(Blowfish encryption algorithm)__([Detail](https://github.com/ghkdqhrbals/simplebank/wiki/ghkdqhrbals:bcrypt))
* [Bcrypt](https://www.usenix.org/legacy/event/usenix99/provos/provos.pdf)

## Update History
* __Gin으로 RPC 통신 추가 ([Details](https://github.com/ghkdqhrbals/simplebank/wiki/ghkdqhrbals:gin))__
1. Set router, routes
2. Set various handler
3. Get http request
4. Use custom validator to check if it is a valid request.
5. Binding JSON to STRUCT(request)
6. Access Local Database -> Execute transactions -> Get results(all process can handle with error)
7. Response

* __Viper으로 configuration 자동설정 ([Details](https://github.com/ghkdqhrbals/simplebank/wiki/ghkdqhrbals:viper))__
1. Set /app.env
2. Set /util/config.go
3. import configurations in /main.go

* __Gomock으로 서비스 레이어의 테스트에서 DB 의존성을 제거 ([Details](https://github.com/ghkdqhrbals/simplebank/wiki/ghkdqhrbals:mockdb))__
1. Use sqlc interface with all query functions to interface
2. Edit /.bash_profile for PATH to go/bin(to using mockgen)
3. Execute mockgen to generate mock functions
4. __Set APIs for testing(TestGetAccountAPI)__

__명령어는 Makefile에 정의__
__Work in VScode and Extensions below__
* [Go Coverage Viewer](https://marketplace.visualstudio.com/items?itemName=soren.go-coverage-viewer)
* [Go Extension Pack](https://marketplace.visualstudio.com/items?itemName=doggy8088.go-extension-pack)
* [Go Test Explorer](https://marketplace.visualstudio.com/items?itemName=premparihar.gotestexplorer)
* [Git Extension Pack](https://marketplace.visualstudio.com/items?itemName=donjayamanne.git-extension-pack)