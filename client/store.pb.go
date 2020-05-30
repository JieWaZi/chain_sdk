package client

type PutStateRequest struct {
	Keys  []string `json:"keys,omitempty"`
	Value string   `json:"value,omitempty"`
}

type Vehicle struct {
	ObjectType string  `json:"objectType"`
	CreateTime int64   `json:"createTime"`
	Id         string  `json:"id"`
	Brand      string  `json:"brand"`
	Price      float64 `json:"price"`
	OwnerId    string  `json:"ownerId"`
	Status     int32   `json:"status"`
	UserId     string  `json:"userId"`
}

type GetStateRequest struct {
	Keys []string `json:"keys,omitempty"`
}

type GetHistoryRequest struct {
	Keys []string `json:"keys,omitempty"`
}

type GetHistoryResponse struct {
	KeyModifications []*KeyModification `json:"key_modifications,omitempty"`
}

type KeyModification struct {
	TxId      string `json:"tx_id,omitempty"`
	Value     string `json:"value,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	IsDelete  bool   `json:"is_delete,omitempty"`
}
