package dbtable

type Task struct {
	ID        int64 `gorm:"primary_key;AUTO_INCREMENT:number;unique"`
	Project   Project
	ProjectID int64
	OP        int8 //0未执行，1执行
}
