package arch

import (
	"container/list"
	"fmt"
)

// MsgQueue is a FIFO queue for buffering messages.
type MsgQueue struct {
	list *list.List
}

// NewMsgQueue creates a new empty message queue.
func NewMsgQueue() *MsgQueue {
	return &MsgQueue{
		list: list.New(),
	}
}

// Push pushes a message into the queue.
func (q *MsgQueue) Push(m *Message) {
	q.list.PushBack(m)
}

// Pull pulls a message out of the queue.
func (q *MsgQueue) Pull() *Message {
	if q.list.Len() == 0 {
		return nil
	}
	return q.list.Remove(q.list.Front()).(*Message)
}

// Len returns the length of the queue.
func (q *MsgQueue) Len() int {
	return q.list.Len()
}

// Receive returns the next message.
func (q *MsgQueue) Receive() (*Message, error) {
	m := q.Pull()
	if m == nil {
		return nil, fmt.Errorf("empty queue")
	}
	return m, nil
}
