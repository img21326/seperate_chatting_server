package pubsub

type PubSubInterface interface {
	Publish([]byte) error
	Subscribe() []byte
}
