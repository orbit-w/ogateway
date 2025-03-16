package gateway

type Stopper interface {
	Stop() error
}
