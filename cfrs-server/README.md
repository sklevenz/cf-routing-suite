# cfrs-server

A a web server written in go. 

## Push it to cloud foundry:

```
cf push
curl -sS https://cfrs-server-skl.cfapps.sap.hana.ondemand.com/
```

## Build Locally

### Build Server

```
pushd server
  go test
  go build
  
  # alternatively
  go run cfrs-server.go
  curl -sS http://localhost:8080/ 
popd
```

### TODO

- ~~add version~~
- ~~add unit test~~
- add version to cf push process
- ~~rename server.go~~
- add go releaser
- add ci
- add blue/green deployment
- add mongodb persistence local
- add mongodb persistence service



## API

- https://cfrs-server-skl.cfapps.sap.hana.ondemand.com/
