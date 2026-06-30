package voice

import "strings"

type Transcript struct {
	Text    string
	IsFinal bool
}

func (t Transcript) FinalText() (string, bool) {
	text := strings.TrimSpace(t.Text)
	return text, t.IsFinal && text != ""
}

type UserUtterance struct {
	Text string
}

type AssistantResponse struct {
	Text string
}
