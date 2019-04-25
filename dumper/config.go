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
}
