package base

type File interface {
	Bash()
}

type Dir interface {
	Cd()
}

func CP(f File) bool {
	f.Bash()
	return false
}
