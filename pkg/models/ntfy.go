package models

type Notification struct {
	Topic    string    `json:"topic"`
	Message  string    `json:"message"`
	Title    string    `json:"title"`
	Tags     []string  `json:"tags"`
	Priority int       `json:"priority"`
	Actions  []*Action `json:"actions"`
	Click    string    `json:"click"`
	Attach   string    `json:"attach"`
	Markdown bool      `json:"markdown"`
	Icon     string    `json:"icon"`
	Filename string    `json:"filename"`
	Delay    string    `json:"delay"`
	Email    string    `json:"email"`
	Call     string    `json:"call"`
}

type Action struct {
	Action string `json:"action"`
	Label  string `json:"label"`
	Url    string `json:"url"`
	Clear  bool   `json:"clear"`
}
