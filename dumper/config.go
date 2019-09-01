package dumper

// BaseConfig Base dump command configuration.
type BaseConfig struct {
	DumpPath     string
	CrtInfo      FileInfo
	KeyInfo      FileInfo
	DomainSubDir bool
	Clean        bool
	Watch        bool
	Hook         string
	Version      string
}

// FileInfo File information.
type FileInfo struct {
	Name string
	Ext  string
}
