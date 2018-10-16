// +build linux

package pomotimer

import "github.com/erikdubbelboer/gspt"

func setProcessName(name string) {
	gspt.SetProcTitle(name)
}
