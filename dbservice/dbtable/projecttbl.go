package dbtable

import (
	"blockchain-service/dbservice/dbmanager"
	"errors"
	"fmt"
)

type Project struct {
	ID        int64  `gorm:"primary_key;AUTO_INCREMENT:number;unique"`
	Name      string `gorm:"type:varchar(100);not null;unique"`
	Desc      string `gorm:"type:varchar(100)"`
	Protype   int8   `gorm:"not null"` //project 类型 EOS, ENTHUM
	Runstate  int8   `gorm:"not null"` //0 未初始化创建 1可更新
	UsrBase   UsrBase
	UsrBaseID int64 `gorm:"not null"` //usr id

}

func ProjectTblExist() bool {

	if !dbmanager.GetDBMgrInst().GetOrmDB().HasTable(&Project{}) {
		fmt.Println("project table not exist")
		return false
	}
	fmt.Println("project table has existed")
	return true
}

func ProjectTblCreate() error {

	if !dbmanager.GetDBMgrInst().GetOrmDB().HasTable(&Project{}) {
		if err := dbmanager.GetDBMgrInst().GetOrmDB().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Project{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func ProjectTblInsert(project *Project) error {

	if err := dbmanager.GetDBMgrInst().GetOrmDB().Create(project).Error; err != nil {
		fmt.Println("project insert failed ")
		return err
	}
	fmt.Println("project tbl insert success")
	return nil
}

func ProjectTblFind(project *Project) error {
	res := dbmanager.GetDBMgrInst().GetOrmDB().Where(project).First(project)
	if res.Error != nil {
		fmt.Println("project tbl find failed! ", res.Error)
		return res.Error
	}
	fmt.Println("project data find success!")
	fmt.Println("project name is : ", project.Name)
	return nil
}

func ProjectTblDel(project *Project) error {
	if err := dbmanager.GetDBMgrInst().GetOrmDB().Where(project).Delete(Project{}).Error; err != nil {
		fmt.Println("Project table delete failed")
		return err
	}
	fmt.Println("Project table delete success")
	return nil
}

func ProjectTblFindByUsr(usr *UsrBase, project *Project) error {

	dbmanager.GetDBMgrInst().GetOrmDB().Model(usr).Related(project)
	if project.ID == 0 {
		fmt.Println("Can't find project by usr")
		return errors.New("Can't find project by usr")
	}
	fmt.Println("find project by usr success")
	return nil
}

func ProjectTblUpdate(before *Project, after *Project) error {
	if before.ID == 0 {
		return errors.New("primary key is empty")
	}

	if err := dbmanager.GetDBMgrInst().GetOrmDB().Model(before).Updates(after).Error; err != nil {
		fmt.Println("Project update failed")
		return err
	}

	fmt.Println("Project update  success")

	return nil
}
