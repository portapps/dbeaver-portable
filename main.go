//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico
package main

import (
	"os"

	. "github.com/portapps/portapps"
)

func init() {
	Papp.ID = "dbeaver-portable"
	Papp.Name = "DBeaver"
	Init()
}

func main() {
	Papp.AppPath = AppPathJoin("app")
	Papp.DataPath = CreateFolder(AppPathJoin("data"))
	Papp.Process = PathJoin(Papp.AppPath, "dbeaver.exe")
	Papp.Args = []string{
		"-data",
		Papp.DataPath,
		"-vm",
		PathJoin(Papp.AppPath, "jre", "bin", "javaw.exe"),
	}
	Papp.WorkingDir = Papp.AppPath

	Launch(os.Args[1:])
}
