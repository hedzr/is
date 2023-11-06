// Copyright © 2023 Hedzr Yeh.

package color

// import (
// 	"fmt"
// 	"strings"
// 	"sync"
//
// 	"github.com/hedzr/env/log"
// )
//
// type lxS struct {
// 	log.Logger
// }
//
// var onceLxInitializer sync.Once //nolint:gochecknoglobals //no
// var lx *lxS                     //nolint:gochecknoglobals //no
//
// const extrasLogSkip = 4
//
// // LazyInit initials local lx instance properly.
// // While you are storing a log.Logger copy locally, DO NOT put these
// // codes into func init() since it's too early to get a unconfigurec
// // log.Logger.
// //
// // The best time to call LazyInit is at cmdr.Root(...).AddGlobalPreAction(...)
// //
// //	root := cmdr.Root(appName, version).
// //			AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
// //				logx.LazyInit()
// //			})
// //
// // For this usage, refer to: https://github.com/hedzr/bgo/blob/master/cli/bgo/cmd/root_cmd.go#L31
// func LazyInit() { lazyInit() }
//
// func lazyInit() *lxS {
// 	onceLxInitializer.Do(func() {
// 		lx = &lxS{
// 			log.Skip(extrasLogSkip),
// 		}
// 	})
// 	return lx
// }
//
// //nolint:lll //no
// func _internalLogTo(tofn func(sb strings.Builder, ln bool), format string, args ...interface{}) { //nolint:goprintffuncname //so what
// 	var sb strings.Builder
// 	sb.WriteString(fmt.Sprintf(format, args...))
// 	tofn(sb, strings.HasSuffix(sb.String(), "\n"))
// }
//
// func init() {
// 	lazyInit()
// }
