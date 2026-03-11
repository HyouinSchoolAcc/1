package web

// Config holds server configuration.
type Config struct {
	RootDir         string
	TemplatesDir    string
	StaticDir       string
	PresetBaseDir   string
	BackupRootDir   string
	RetentionDays   int
	ServePort       string
	EnableDevReload bool
}
