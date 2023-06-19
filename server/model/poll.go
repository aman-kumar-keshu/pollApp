package model

type Poll struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Topic     string  `json:"topic"`
	Src       string  `json:"src"`
	Upvotes   int     `json:"upvotes"`
	Downvotes int     `json:"downvotes"`
}

type PollCollection struct {
	Polls []Poll `json:"items"`
}
