package message

type Message struct {
	Body   string `json:"body"`
	ConnId string `json:"id"`
	User   string `json:"user"`
}
