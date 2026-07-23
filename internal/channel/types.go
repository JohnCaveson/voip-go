package channel

type ChannelInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsDefault bool   `json:"is_default"`
}
