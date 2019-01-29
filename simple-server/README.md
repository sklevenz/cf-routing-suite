# simple-server

A a simple web server written in go. 

## Push it to cloud foundry:

```
cf push
curl -sS https://simple-server-skl.cfapps.sap.hana.ondemand.com/
```

## Build Locally

### Build Simple Server

```
pushd simple-server
  # go test
  go build
  go install
  
  # alternatively
  go run server.go
  curl -sS http://localhost:8080/ 
popd
```

### TODO

- ~~add version~~
- ~~add unit test~~
- add version to cf push process
- rename server.go
- add go releaser
- add ci
- add blue/green deployment
- add mongodb persistence local
- add mongodb persistence service



## API

- https://simple-server-skl.cfapps.sap.hana.ondemand.com/
- https://simple-server-skl.cfapps.sap.hana.ondemand.com/reset/
- https://simple-server-skl.cfapps.sap.hana.ondemand.com/count/
