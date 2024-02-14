package config

type Args struct {
	StartAddr string
	BaseAddr  string
}

type GetArgsBuilder interface {
	SetStart(string) GetArgsBuilder
	SetBase(string) GetArgsBuilder
	Build() *Args
}
type ConcreteGetArgsBuilder struct {
	args *Args
}

func NewGetArgsBuilder() *ConcreteGetArgsBuilder {
	return &ConcreteGetArgsBuilder{args: &Args{}}
}

func (cgab *ConcreteGetArgsBuilder) SetStart(startAddr string) GetArgsBuilder {
	cgab.args.StartAddr = startAddr
	return cgab
}

func (cgab *ConcreteGetArgsBuilder) SetBase(baseAddr string) GetArgsBuilder {
	cgab.args.BaseAddr = baseAddr
	return cgab
}
func (cgab *ConcreteGetArgsBuilder) Build() *Args {
	return cgab.args
}
