//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico
package main

import (
	"os"
	"strings"

	"github.com/magiconair/properties"
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

	driversPath := CreateFolder(PathJoin(Papp.DataPath, ".metadata", "drivers"))
	logsPath := CreateFolder(PathJoin(Papp.DataPath, ".metadata", "logs"))
	corePrefsPath := CreateFolder(PathJoin(Papp.DataPath, ".metadata", ".plugins", "org.eclipse.core.runtime", ".settings"))
	corePrefsFile := PathJoin(corePrefsPath, "org.jkiss.dbeaver.core.prefs")

	defaultProps := properties.NewProperties()
	_, _, _ = defaultProps.Set("dialog.default.folder", formatPath(Papp.DataPath))
	_, _, _ = defaultProps.Set("logs.debug.location", formatPath(PathJoin(logsPath, "dbeaver-debug.log")))
	_, _, _ = defaultProps.Set("qm.logDirectory", formatPath(logsPath))
	_, _, _ = defaultProps.Set("ui.auto.update.check", "false")
	_, _, _ = defaultProps.Set("ui.drivers.home", formatPath(driversPath))

	if !Exists(corePrefsFile) {
		Log.Info("Creating default props...")
		if err := WriteToFile(corePrefsFile, defaultProps.String()); err != nil {
			Log.Error("Cannot write default props to org.jkiss.dbeaver.core.prefs: ", err)
		}
	} else {
		Log.Info("Loading org.jkiss.dbeaver.core.prefs file...")
		corePrefsProps, err := properties.LoadFile(corePrefsFile, properties.UTF8)
		if err != nil {
			Log.Error("Cannot load org.jkiss.dbeaver.core.prefs file: ", err)
		}
		corePrefsProps.Merge(defaultProps)
		Log.Info("Wrting to org.jkiss.dbeaver.core.prefs...")
		if err := WriteToFile(corePrefsFile, corePrefsProps.String()); err != nil {
			Log.Error("Cannot write to org.jkiss.dbeaver.core.prefs: ", err)
		}
	}

	Launch(os.Args[1:])
}

func formatPath(path string) string {
	return strings.Replace(strings.Replace(path, `/`, `\`, -1), `\`, `\\`, -1)
}
