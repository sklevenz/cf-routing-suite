# cf-routing-suite (cfrs)
A suite containing servers and clients for challenging the cloud foundry routing tier

## cfrs-server

A a web server written in go. 

### Push it to cloud foundry:

```
cf push
curl -sS https://cfrs-server-skl.cfapps.sap.hana.ondemand.com/
```

### Build Locally

#### Build Server

```
pushd server
  go test
  go build
  
  # alternatively
  go run cfrs-server.go
  curl -sS http://localhost:8080/ 
popd
```

#### TODO

- ~~add version~~
- ~~add unit test~~
- ~~add version to cf push process~~
- ~~rename server.go~~
- ~~add go releaser~~
- add ci
- add blue/green deployment
- ~~add mongodb persistence local~~~
- add mongodb persistence service



### API

- https://cfrs-server-skl.cfapps.sap.hana.ondemand.com/

## cfrs-client

A cf client app written in go. 

# Script Support

Shel scripts are located in script directory. Simply call ````/script/script-name.sh PARAMETER(s)````

| Script | Description |
|--------|-------------|
| ```build.sh``` | build all go binaries |
| ```cf-push-server.sh``` | push server to cloud foundry |
| ```delete-tag.sh``` | delete a tag from github (local and remote) |
| ```make-release.sh``` | Make a new release by [goreleaser](https://goreleaser.com), set a version tag and push release to github |
| ```run-client.sh``` | run the client |
| ```run-mongo.sh``` | run a mongo db for local testing |
| ```run-server.sh``` | run the server |
| ```test.sh``` | run all go tests |

