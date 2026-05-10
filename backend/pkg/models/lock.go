package models

type Lock struct {
    ID        string `json:"id"`
    Status    string `json:"status"` // e.g., "locked" or "unlocked"
    OwnerID   string `json:"owner_id"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}