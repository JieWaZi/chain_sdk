package client

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
