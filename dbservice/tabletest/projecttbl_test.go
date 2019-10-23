package tabletest

import (
	"blockchain-service/dbservice/dbmanager"
	"blockchain-service/dbservice/dbtable"
	"testing"
)

func Test_ProjectTbl(t *testing.T) {
	t.Log("第一个测试通过了")
	dbmgr := dbmanager.GetDBMgrInst()
	if dbmgr == nil {
		t.Error("dbmgr init failed!")
		return
	}
	defer dbmgr.ReleaseDB()
	dbtable.ProjectTblCreate()
	pProject := &dbtable.Project{Name: "lbaaspro1", Desc: "", UsrBaseID: 1}
	insertres := dbtable.ProjectTblInsert(pProject)
	if insertres != nil {
		t.Error("project insert failed!")
		return
	}
	t.Log("project insert success!")

	findres := dbtable.ProjectTblFind(pProject)
	if findres != nil {
		t.Error("project find failed!")
		return
	}
	t.Log("project find success!")

	pProjectR := &dbtable.Project{}
	rfindres := dbtable.ProjectTblFindByUsr(&dbtable.UsrBase{ID: 1}, pProjectR)
	if rfindres != nil {
		t.Error("project relate find failed!")
		return
	}

	t.Log("project relate find success!")
	pProjectR.Name = "zackweb5.0"
	updateres := dbtable.ProjectTblUpdate(pProject, pProjectR)
	if updateres != nil {
		t.Error("project update failed ")
		return
	}
	t.Log("project update success!")

	delres := dbtable.ProjectTblDel(pProject)
	if delres != nil {
		t.Error("project del failed!")
		return
	}
	t.Log("project delete success!")

}
