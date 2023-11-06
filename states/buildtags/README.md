# buildtags

As building with `hedzr/log/`, the tags are:

```bash
go build -tags=docker,k8s,istio ./...
```

The usages will be:

```go
if buildtags.IsDockerBuild() {
	//
}
if buildtags.IsK8sBuild() {
//
}
if buildtags.IsIstioBuild() {
//
}
```
