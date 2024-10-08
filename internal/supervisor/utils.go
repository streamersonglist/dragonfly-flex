/**
 * Copied from https://github.com/fly-apps/postgres-flex/blob/master/internal/supervisor/utils.go
 */
package supervisor

import (
	"fmt"
	"os"
)

func fatalOnErr(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(i ...interface{}) {
	fmt.Fprint(os.Stderr, "hivemind: ")
	fmt.Fprintln(os.Stderr, i...)
	panic(i)
}
