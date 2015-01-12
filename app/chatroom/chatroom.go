package chatroom

type Event struct {
	Type      string
	User      string
	Timestamp int
	Text      string
}

type Subscription struct {
	Archive []Event
	New     <-chan Event
}
