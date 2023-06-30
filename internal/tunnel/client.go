package tunnel

type Client interface {
	Connect() (Connection, error)
}
