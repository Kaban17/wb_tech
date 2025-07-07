package models

type DeliveryInfo struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	ZipCode uint   `json:"zip_code"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}
type PaymentInfo struct {
	Transaction  string `json:"transaction"`
	RequestID    uint   `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    uint   `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost uint   `json:"delivery_cost"`
	GoodsTotal   uint   `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}
type ItemInfo struct {
	ChartID     int    `json:"chart_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        int    `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
type ItemsInfo []ItemInfo
type Order struct {
	ID                string       `json:"order_uid"`
	TrackNumber       string       `json:"track_number"`
	Entry             string       `json:"entry"`
	Delivery          DeliveryInfo `json:"delivery"`
	Payment           PaymentInfo  `json:"payment"`
	Items             ItemsInfo    `json:"items"`
	Locale            string       `json:"locale"`
	InternalSignature string       `json:"internal_signature"`
	CustomerID        string       `json:"customer_id"`
	DeliveryService   string       `json:"delivery_service"`
	ShardKey          int          `json:"shard_key"`
	SMID              int          `json:"sm_id"`
	DataCreated       string       `json:"data_created"`
	OofShard          int          `json:"oof_shard"`
}
