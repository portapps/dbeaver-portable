//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"
	"strings"

	"github.com/magiconair/properties"
	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
)

var (
	app *portapps.App
)

func init() {
	var err error

	// Init app
	if app, err = portapps.New("dbeaver-portable", "DBeaver"); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "dbeaver.exe")
	app.Args = []string{
		"-data",
		app.DataPath,
		"-vm",
		utl.PathJoin(app.AppPath, "jre", "bin", "javaw.exe"),
	}

	driversPath := utl.CreateFolder(app.DataPath, ".metadata", "drivers")
	logsPath := utl.CreateFolder(app.DataPath, ".metadata", "logs")
	corePrefsPath := utl.CreateFolder(app.DataPath, ".metadata", ".plugins", "org.eclipse.core.runtime", ".settings")
	corePrefsFile := utl.PathJoin(corePrefsPath, "org.jkiss.dbeaver.core.prefs")

	defaultProps := properties.NewProperties()
	_, _, _ = defaultProps.Set("dialog.default.folder", formatPath(app.DataPath))
	_, _, _ = defaultProps.Set("logs.debug.location", formatPath(utl.PathJoin(logsPath, "dbeaver-debug.log")))
	_, _, _ = defaultProps.Set("qm.logDirectory", formatPath(logsPath))
	_, _, _ = defaultProps.Set("ui.auto.update.check", "false")
	_, _, _ = defaultProps.Set("ui.drivers.home", formatPath(driversPath))

	if !utl.Exists(corePrefsFile) {
		log.Info().Msg("Creating default props...")
		if err := utl.WriteToFile(corePrefsFile, defaultProps.String()); err != nil {
			log.Error().Err(err).Msg("Cannot write default props to org.jkiss.dbeaver.core.prefs")
		}
	} else {
		log.Info().Msg("Loading org.jkiss.dbeaver.core.prefs file...")
		corePrefsProps, err := properties.LoadFile(corePrefsFile, properties.UTF8)
		if err != nil {
			log.Error().Err(err).Msg("Cannot load org.jkiss.dbeaver.core.prefs file")
		}
		corePrefsProps.Merge(defaultProps)
		log.Info().Msg("Writing to org.jkiss.dbeaver.core.prefs")
		if err := utl.WriteToFile(corePrefsFile, corePrefsProps.String()); err != nil {
			log.Error().Err(err).Msg("Cannot write to org.jkiss.dbeaver.core.prefs")
		}
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}

func formatPath(path string) string {
	return strings.Replace(strings.Replace(path, `/`, `\`, -1), `\`, `\\`, -1)
}
