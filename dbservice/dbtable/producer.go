package dbtable

type EOSProducer struct {
	ID         int64  `gorm:"primary_key;AUTO_INCREMENT:number;unique"`
	name       string `gorm:"type:varchar(100);not null;unique"`
	Privatekey string
	url        string
	location   int32
	BaseNode   BaseNode
	BaseNodeID int64
}
