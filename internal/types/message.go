package types

import "time"

type MessageCategory string

const (
    Emergency MessageCategory = "EMERGENCY"
    General   MessageCategory = "GENERAL"
    Help      MessageCategory = "HELP"
)

type Message struct {
    ID        string          `json:"id"`
    Content   string          `json:"content"`
    Category  MessageCategory `json:"category"`
    Sender    User           `json:"sender"`
    Timestamp time.Time      `json:"timestamp"`
}

type User struct {
    ID       string `json:"id"`
    Nickname string `json:"nickname"`
    NodeID   string `json:"nodeId"`
}