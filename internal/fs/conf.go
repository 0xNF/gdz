package fs

type Conf struct {
	Verbose  bool
	OSType   string
	Versions []string
}

func NewConfig() Conf {
	conf := Conf{
		Verbose:  true,
		Versions: []string{"3.8"},
	}
	return conf
}
