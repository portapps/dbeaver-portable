//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/magiconair/properties"
	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/files"
	"github.com/portapps/portapps/v3/pkg/log"
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
	if err := os.MkdirAll(app.DataPath, 0o755); err != nil {
		log.Fatal().Err(err).Msg("Cannot create data path")
	}
	app.Process = filepath.Join(app.AppPath, "dbeaver.exe")
	app.Args = []string{
		"-data",
		app.DataPath,
		"-vm",
		filepath.Join(app.AppPath, "jre", "bin", "javaw.exe"),
	}

	driversPath := filepath.Join(app.DataPath, ".metadata", "drivers")
	logsPath := filepath.Join(app.DataPath, ".metadata", "logs")
	corePrefsPath := filepath.Join(app.DataPath, ".metadata", ".plugins", "org.eclipse.core.runtime", ".settings")
	for _, dir := range []string{driversPath, logsPath, corePrefsPath} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Fatal().Err(err).Msgf("Cannot create directory %s", dir)
		}
	}
	corePrefsFile := filepath.Join(corePrefsPath, "org.jkiss.dbeaver.core.prefs")

	defaultProps := properties.NewProperties()
	_, _, _ = defaultProps.Set("dialog.default.folder", formatPath(app.DataPath))
	_, _, _ = defaultProps.Set("logs.debug.location", formatPath(filepath.Join(logsPath, "dbeaver-debug.log")))
	_, _, _ = defaultProps.Set("qm.logDirectory", formatPath(logsPath))
	_, _, _ = defaultProps.Set("ui.auto.update.check", "false")
	_, _, _ = defaultProps.Set("ui.drivers.home", formatPath(driversPath))

	if !files.Exists(corePrefsFile) {
		log.Info().Msg("Creating default props...")
		if err := os.WriteFile(corePrefsFile, []byte(defaultProps.String()), 0o644); err != nil {
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
		if err := os.WriteFile(corePrefsFile, []byte(corePrefsProps.String()), 0o644); err != nil {
			log.Error().Err(err).Msg("Cannot write to org.jkiss.dbeaver.core.prefs")
		}
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}

func formatPath(path string) string {
	return strings.ReplaceAll(filepath.FromSlash(path), `\`, `\\`)
}
