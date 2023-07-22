package telegram

import "testing"

func TestSendMessage(t *testing.T) {
	m := NewManager("", 0)

	err := m.Send("Hello World!")
	if err != nil {
		t.Fatal(err)
	}
}
