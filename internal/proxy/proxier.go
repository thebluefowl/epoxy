package proxy

type Proxier interface {
	Do([]byte) ([]byte, error)
}
