package buildtags

// IsK8sBuild tests if go build with -tags=k8s
func IsK8sBuild() bool { return k8sEnabled }

// IsIstioBuild tests if go build with -tags=istio
func IsIstioBuild() bool { return istioEnabled }

// IsDockerBuild tests if go build with -tags=docker-build
func IsDockerBuild() bool { return dockerEnabled }
