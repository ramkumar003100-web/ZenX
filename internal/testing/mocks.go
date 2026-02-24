package testing

type Generator interface {
	Generate(interfacePath string) error
}
