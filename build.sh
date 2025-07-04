#!/bin/bash

# for d in *; do
# 	if [ -d $d ]; then
# 	fi
# done

# find . -type f -iname 'go.mod' -print0 | xargs -0I% echo "pushd \$(dirname %)>/dev/null && pwd && go mod tidy && popd >/dev/null; echo;echo;echo" | sh

set -e

MAIN_BUILD_PKG=(".")
MAIN_APPS=(_examples)
APPS=(blocks colors color-tool prompt pipe small)

# for an in ${APPS[*]}; do
# 	echo "$an //"
# done

# LDFLAGS = -s -w -X 'github.com/hedzr/cmdr/v2/conf.Buildstamp=2024-10-25T18:09:06+08:00' -X 'github.com/hedzr/cmdr/v2/conf.GIT_HASH=580ca50' -X 'github.com/hedzr/cmdr/v2/conf.GitSummary=580ca50-dirty' -X 'github.com/hedzr/cmdr/v2/conf.GitDesc=580ca50 upgrade deps' -X 'github.com/hedzr/cmdr/v2/conf.BuilderComments=' -X 'github.com/hedzr/cmdr/v2/conf.GoVersion=go version go1.23.7 darwin/arm64' -X 'github.com/hedzr/cmdr/v2/conf.Version=0.5.1'

extract-app-version() {
	local ok=0
	local DEFAULT_DOC_NAME="${DEFAULT_DOC_NAME}"
	for dn in "${DEFAULT_DOC_NAME}" \
		slog/doc.go _examples/small/doc.go examples/small/doc.go \
		cli/cmdr-cli/consts/def.go cli/cmdr/consts/def.go cli/cmdr-cli/doc.go cli/cmdr/doc.go \
		"${1:-doc.go}"; do
		if [[ $ok -eq 0 ]]; then
			local docfn="${dn}"
			if [ -f "$docfn" ]; then
				dbg "checking ${docfn}..."
				APPNAME="$(grep -E "appName[ \t]+=[ \t]+" ${docfn} | grep -Eo "\\\".+\\\"")"
				VERSION="$(grep -E "version[ \t]+=[ \t]+" ${docfn} | grep -Eo "[0-9.]+")"
				if [[ "$APPNAME" != "" && "$VERSION" != "" ]]; then
					# $ foo=${string#"$prefix"}
					# $ foo=${foo%"$suffix"}
					VERSION="v${VERSION#v}"
					let ok++
				fi
			fi
		fi
	done
	if [[ $ok -eq 0 ]]; then
		false
	fi
}

publish() {
	extract-app-version
	local version="$VERSION"
	tip "version extracted: $version"
	if is_git_dirty; then
		err ".git repo is NOT clean, cannot publish right now"
	elif [ ! "$version" = v1.* ]; then
		err "BAD version ($version) extracted, cannot publish right now"
	else
		git tag $version
		git push origin --all && git push origin --tags
		# sleepx 5
		# git tag lite/$version
		# git push origin --all && git push origin --tags
	fi
}

build-main-platforms() {
	local all=1
	(($#)) && all=0
	local appname="${1:-colors}"
	local binname="${1:-colors}"
	(($#)) && shift

	extract-app-version
	local version="$VERSION"

	for GOOS in darwin linux windows freebsd openbsd netbsd plan9; do
		local suffix=''
		[[ "$GOOS" = "windows" ]] && suffix='.exe'
		for GOARCH in amd64 arm64 riscv64 mips64; do
			if go tool dist list | grep -qE "$GOOS/$GOARCH"; then
				if [ "$all" = "1" ]; then
					for appname in ${APPS[@]}; do
						binname="${appname}"
						tip "--- build '$binname' for $GOOS/$GOARCH ---"
						GOOS=$GOOS GOARCH=$GOARCH go build "$@" \
							-o "./bin/${binname}-${version}-${GOOS}_${GOARCH}${suffix}" \
							./${MAIN_APPS[0]}/$appname #./_examples/$appname/
					done
				else
					tip "--- build for $GOOS/$GOARCH ---"
					GOOS=$GOOS GOARCH=$GOARCH go build "$@" \
						-o "./bin/${binname}-${version}-${GOOS}_${GOARCH}${suffix}" \
						./${MAIN_APPS[0]}/$appname #./_examples/$appname/
				fi
			fi
		done
	done
	ls -la $LS_OPT ./bin/${binname}*
}

build-all-platforms() {
	local W_PKG=github.com/hedzr/cmdr/v2/conf
	local TIMESTAMP="$(date -u '+%Y-%m-%dT%H:%M:%S')"
	local GOVERSION="$(go version)"
	local GIT_VERSION="$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")"
	local GIT_REVISION="$(git rev-parse --short HEAD)"
	local GIT_SUMMARY="$(git describe --tags --dirty --always)"
	local GIT_DESC="$(git log --oneline -1)"
	local GIT_HASH="$(git rev-parse HEAD)"
	local GOBUILD_TAGS=-tags="hzstudio sec antonal"
	extract-app-version
	tip "APPNAME = $APPNAME, VERSION = $VERSION"
	local LDFLAGS="-s -w -X '${W_PKG}.Buildstamp=$TIMESTAMP' \
		-X '${W_PKG}.GIT_HASH=$GIT_REVISION' \
		-X '${W_PKG}.GitSummary=$GIT_SUMMARY' \
		-X '${W_PKG}.GitDesc=$GIT_DESC' \
		-X '${W_PKG}.BuilderComments=$BUILDER_COMMENT' \
		-X '${W_PKG}.GoVersion=$GOVERSION' \
		-X '${W_PKG}.Version=$VERSION' \
		-X '${W_PKG}.AppName=$APPNAME'"
	local CGO_ENABLED=0
	local GOBIN="./bin" GOOS GOARCH
	local mbp ma an

	for GOOS in $(go tool dist list | awk -F'/' '{print $1}' | sort -u); do
		if [[ "$GOOS" != "aix" && "$GOOS" != "android" && "$GOOS" != "illumos" && "$GOOS" != "ios" ]]; then
			for mbp in ${MAIN_BUILD_PKG[*]}; do
				for ma in ${MAIN_APPS[*]}; do
					for an in ${APPS[*]}; do
						local ANAME="${mbp}/${ma}/${an}"
						echo -e "\n\n" && tip "BUILDING FOR ${ANAME} / $GOOS ...\n"
						if [ -d "$ANAME" ]; then
							for GOARCH in $(go tool dist list | grep -E "^$GOOS" | awk -F'/' '{print $2}' | sort -u); do
								SUFFIX="-${GOOS}_${GOARCH}"
								[[ $GOOS = windows ]] && SUFFIX="${SUFFIX}.exe"
								echo "  >> building ${ANAME} - ${GOARCH} ..."
								go build -trimpath -gcflags=all='-l -B' \
									"${GOBUILD_TAGS}" -ldflags "${LDFLAGS}" \
									-o ${GOBIN}/${an}${SUFFIX} \
									${ANAME}/ || exit
							done
							ls -la $LS_OPT $GOBIN/${an}_${GOOS}*
						fi
						# return
					done
				done
			done
		fi
	done
}

test-all-platforms() {
	local func="${1:-}"
	[ -d ./logs ] || mkdir -pv ./logs
	# for GOOS in $(go tool dist list | awk -F'/' '{print $1}' | sort -u); do
	for GOOS in darwin linux windows freebsd openbsd netbsd plan9; do
		for GOARCH in amd64 arm64 riscv64 mips64; do
			go tool dist list | grep -qE "$GOOS/$GOARCH" &&
				echo && echo && tip "TESTING FOR $GOOS/$GOARCH ...\n" &&
				go test ./... -v -race -cover -coverprofile=./logs/coverage-${GOOS}_${GOARCH}.txt \
					-test.run=^TestDirTimestamps$ -covermode=atomic -test.short -vet=off 2>&1 || exit
		done
	done
}

cov() {
	for GOOS in darwin linux windows; do
		go test -v -race -coverprofile=coverage-$GOOS.txt ./...
		go tool cover -html=coverage-$GOOS.txt -o cover-$GOOS.html
	done
	open cover-darwin.html
}

bf1() {
	# https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
	# Or: go tool dist list
	# Missed: posix
	for GOOS in $(go tool dist list | awk -F'/' '{print $1}' | sort -u); do
		echo -e "\n\nTESTING FOR $GOOS ...\n"
		go test -v -race -coverprofile=coverage-${GOOS}_${GOARCH}.txt \
			-test.run=^TestDirTimestamps$ ./dir/ || exit
	done
}

fmt() {
	echo fmt...
	gofmt -l -w -s .
}

lint() {
	echo lint...
	golint ./...
}

cyclo() {
	echo cyclo...
	gocyclo -top 10 .
}

all-ops() {
	fmt && lint && cyclo
}

test-all() { test-all-platforms "$@"; }

all() { build-all-platforms "$@"; }

main() { build-main-platforms "$@"; }

mk-ver() {
	if extract-app-version "$@"; then
		echo "version info extracted: APPNAME=$APPNAME, VERSION=$VERSION"
	fi
}

# if [[ $# -eq 0 ]]; then
# 	cmd=cov
# else
# 	cmd=${1:-cov} && shift
# fi
# $cmd "$@"

sleepx() { tip "sleeping..." && (($#)) && \sleep "$@"; }

######### SIMPLE BASH.SH FOOTER BEGIN #########

# The better consice way to get baseDir, ie. $CD, is:
#       CD=$(cd `dirname "$0"`;pwd)
# It will open a sub-shell to print the folder name of the running shell-script.

is_darwin() { [[ $OSTYPE == darwin* ]]; }
is_darwin_sillicon() { is_darwin && [[ $(uname_mach) == arm64 ]]; }
is_linux() { [[ $OSTYPE == linux* ]]; }
is_freebsd() { [[ $OSTYPE == freebsd* ]]; }
is_win() { in_wsl; }
in_wsl() { [[ "$(uname -r)" == *windows_standard* ]]; }

is_git_clean() { git diff-index --quiet "$@" HEAD -- 2>/dev/null; }
is_git_dirty() {
	if is_git_clean "$@"; then
		false
	else
		true
	fi
}

dbg() { ((DEBUG)) && printf ">>> \e[0;38;2;133;133;133m$@\e[0m\n" || :; }
tip() { printf "\e[0;38;2;133;133;133m>>> $@\e[0m\n"; }
err() { printf "\e[0;33;1;133;133;133m>>> $@\e[0m\n" 1>&2; }
fn_exists() { LC_ALL=C type $1 2>/dev/null | grep -qE '(shell function)|(a function)'; }
CD="$(cd $(dirname "$0") && pwd)" && BASH_SH_VERSION=v20241021 && DEBUG=${DEBUG:-0} && PROVISIONING=${PROVISIONING:-0}
LS_OPT="--color"
is_darwin && LS_OPT="-G"
(($#)) && {
	dbg "$# arg(s) | CD = $CD"
	check_entry() {
		local prefix="$1" cmd="$2" && shift && shift
		if fn_exists "${prefix}_${cmd}_entry"; then
			eval "${prefix}_${cmd}_entry" "$@"
		elif fn_exists "${cmd}_entry"; then
			eval "${cmd}_entry" "$@"
		else
			prefix="${prefix}_${cmd}"
			if fn_exists $prefix; then
				eval $prefix "$@"
			elif fn_exists ${prefix//_/-}; then
				eval ${prefix//_/-} "$@"
			elif fn_exists $cmd; then
				eval $cmd "$@"
			elif fn_exists ${cmd//_/-}; then
				eval ${cmd//_/-} "$@"
			else
				err "command not found: $cmd $@"
				return 1
			fi
		fi
	}
	check_entry "boot" "$@"
} || { dbg "empty: $# | CD = $CD"; }
######### SIMPLE BASH.SH FOOTER END #########
