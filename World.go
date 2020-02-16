package multi_threading

type World struct {
	ProcessesList []*Process
	Associates    map[string]workFunction
	Nl            NetworkLayer
}
