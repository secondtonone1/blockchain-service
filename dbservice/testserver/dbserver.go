package main

import (
	"fmt"
	"lbaas/dbservice/dbmanager"
	"lbaas/dbservice/dbtable"
)

func main() {
	dbmgr := dbmanager.GetDBMgrInst()
	if dbmgr == nil {
		fmt.Println("dbmgr init failed!")
	}
	defer dbmgr.ReleaseDB()
	dbtable.UsrBaseTblExist()
	dbtable.UsrBaseTblCreate()
	dbtable.UsrBaseTblExist()

	dbtable.UsrBaseTblInsert(&dbtable.UsrBase{
		Name:   "lemon",
		Passwd: "lemon123",
	})
	dbtable.UsrBaseTblUpdate(&dbtable.UsrBase{Name: "lemon"}, &dbtable.UsrBase{Passwd: "lemon456"})
	dbtable.UsrBaseTblFind(&dbtable.UsrBase{Name: "lemon"})
	dbtable.UsrBaseTblDel(&dbtable.UsrBase{Name: "lemon"})
	dbtable.UsrBaseTblFind(&dbtable.UsrBase{Name: "lemon"})
	//dbtable.UsrBaseUnscoped(&dbtable.UsrBase{Name: "lemon"})

}
