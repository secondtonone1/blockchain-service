package dbtable

type BaseNode struct {
	ID                            int64  `gorm:"primary_key;AUTO_INCREMENT:number;unique"`
	Httpserveraddress             string `gorm:"type:varchar(100)"`
	HttpvalidateHost              bool
	AccessControlAlloworigin      string `gorm:"type:varchar(100)"`
	AccessControlAllowheaders     string `gorm:"type:varchar(100)"`
	AccessControlMaxage           string `gorm:"type:varchar(100)"`
	AccessControlAllowcredentials bool
	WasmRuntime                   string `gorm:"type:varchar(100)"`
	ContractsConsole              bool
	AgentName                     string `gorm:"type:varchar(100)"`
	MaxTransactionTime            int32
	P2Ppeeraddress                string `gorm:"type:varchar(100)"`
	Project                       Project
	ProjectID                     int64
	Timestamp                     int64
	InitialPrivatekey             string
	OperationRes                  int8 //0 成功, 1失败
	NodeType                      int8 //0 FullNode , 1 ProduceNode , 2
}
