package providers

type Provider interface {
	GetVersions(moduleName, artifactName string) ([]string, error)
}
