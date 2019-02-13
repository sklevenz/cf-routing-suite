# cf-routing-suite (cfrs)
A suite containing servers and clients for challenging the cloud foundry routing tier

## cfrs-server

The server can run locally on a given port or it can be pushed to cloud foundry.

### simulator vs. mongodb

The server can run in simulation mode or it can connect to a mongodb database. 
The simulator is for testing and rapid development. 
If the server runs locally then it can connect to a local running mongodb (localhost with fixed port).
If the server is pushed to cloud foundry then it can connect to mongodb backing service via service binding. 

The server is controlled by following environment variables:

| Variable | Port |
|----------|------|
| PORT | port the service is listening to |
| MODE | the mode can be simulator or mongodb |




### db mode


#### TODO

- ~~add version~~
- ~~add unit test~~
- ~~add version to cf push process~~
- ~~rename server.go~~
- ~~add go releaser~~
- add ci
- add blue/green deployment
- ~~add mongodb persistence local~~~
- add mongodb persistence cf service
- ~~replace flags by env only~~


### API

- https://cfrs-server-skl.cfapps.sap.hana.ondemand.com/

## cfrs-client

A cf client app written in go. 

# Script Support

Shel scripts are located in script directory. Simply call ````/script/script-name.sh PARAMETER(s)````

| Script | Description |
|--------|-------------|
| ```build.sh``` | build all go binaries |
| ```cf-push-server.sh [s|db]``` | push server to cloud foundry (s: use database simulator, db: connect to mongodb)|
| ```delete-tag.sh``` | delete a tag from github (local and remote) |
| ```make-release.sh``` | Make a new release by [goreleaser](https://goreleaser.com), set a version tag and push release to github |
| ```run-client.sh``` | run the client |
| ```run-mongo.sh``` | run a mongo db for local testing |
| ```run-server.sh [s|db] [8080]``` | run the server locally (s | use database simulator, db: connect ot mongodb
| ```test.sh``` | run all go tests (uses local mongodb if available |

