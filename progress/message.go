package progress

type ProgressMessages struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Number string `json:"number"`
	Text   string `json:"text"`
}


