package manager

type Options struct {
	StoragePath string
}

type Option func(*Options)

func WithStoragePath(path string) Option {
	return func(o *Options) {
		o.StoragePath = path
	}
}
