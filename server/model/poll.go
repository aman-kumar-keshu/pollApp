package model

type Poll struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Topic     string  `json:"topic"`
	Src       string  `json:"src"`
	Upvotes   int     `json:"upvotes"`
	Downvotes int     `json:"downvotes"`
	Options 	[]Option  `json:"options"`
}

type PollRequest struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Topic     string  `json:"topic"`
	Src       string  `json:"src"`
	Upvotes   int     `json:"upvotes"`
	Downvotes int     `json:"downvotes"`
	Options 	[]string  `json:"options"`
}

type PollCollection struct {
	Polls []Poll `json:"items"`
}

type Option struct {
    ID	         int  `json:"id"`
    Option       string
	Votes		 int  
    PollId 	 	 int
}