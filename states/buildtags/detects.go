package buildtags

// IsK8sBuild tests if go build with -tags=k8s
//
//goland:noinspection GoBoolExpressions
func IsK8sBuild() bool { return k8sEnabled }

// IsIstioBuild tests if go build with -tags=istio
//
//goland:noinspection GoBoolExpressions
func IsIstioBuild() bool { return istioEnabled }

// IsDockerBuild tests if go build with -tags=docker-build
//
//goland:noinspection GoBoolExpressions
func IsDockerBuild() bool { return dockerEnabled }
