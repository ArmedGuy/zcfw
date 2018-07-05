package registry

type Registry interface {
	Watch() error
}
