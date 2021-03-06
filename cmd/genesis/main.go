package main

import (
	"runtime"
	"time"

	"github.com/whiteblock/genesis-cli/pkg/cmds"
	"github.com/whiteblock/genesis-cli/pkg/message"
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"
)

var (
	buildTime  string
	commitHash string
)

func main() {
	updateChan := service.CheckForUpdates(commitHash)
	timeChan := time.After(100 * time.Millisecond)
	cmds.Execute(buildTime, commitHash)
	select {
	case hasUpdate := <-updateChan:
		if hasUpdate {
			if runtime.GOOS == "windows" {
				util.Print(message.WindowsUpdateAvailable)
			} else {
				util.Print(message.UpdateAvailable)
			}

		}

	case <-timeChan:
		// do nothing
	}
}
