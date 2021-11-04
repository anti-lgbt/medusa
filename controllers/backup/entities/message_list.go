package entities

type MessageList struct {
	ID       int64  `json:"id"`
	UID      string `json:"uid"`
	Personal bool   `json:"personal"`
}
