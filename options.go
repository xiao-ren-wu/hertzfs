package hertzfs

type staticFSConf struct {
	BasePath     string
	CacheControl int
}

type FSOption func(conf *staticFSConf)

func BasePath(basePath string) FSOption {
	return func(conf *staticFSConf) {
		conf.BasePath = basePath
	}
}

func CacheControl(maxCache int) FSOption {
	return func(conf *staticFSConf) {
		conf.CacheControl = maxCache
	}
}
