package dbtable

import (
	"fmt"

	"lbaas/dbservice/dbmanager"
)

func UsrBaseTblExist() bool {

	if !dbmanager.GetDBMgrInst().GetOrmDB().HasTable(&UsrBase{}) {
		fmt.Println("usrbase table not exist")
		return false
	}
	fmt.Println("usrbase table has existed")
	return true
}

func UsrBaseTblCreate() error {

	if !dbmanager.GetDBMgrInst().GetOrmDB().HasTable(&UsrBase{}) {
		if err := dbmanager.GetDBMgrInst().GetOrmDB().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&UsrBase{}).Error; err != nil {
			return err
		}

	}
	return nil
}

type UsrBase struct {
	Name   string `gorm:"primary_key;type:varchar(100);not null;unique"`
	Passwd string `type:varchar(100);not null`
}

func UsrBaseTblInsert(usrbase *UsrBase) error {

	if err := dbmanager.GetDBMgrInst().GetOrmDB().Create(usrbase).Error; err != nil {
		fmt.Println("usr base insertfailed ")
		return err
	}
	fmt.Println("usr base tbl insert success")
	return nil
}

func UsrBaseTblFind(usrbase *UsrBase) error {
	res := dbmanager.GetDBMgrInst().GetOrmDB().Where(usrbase).First(usrbase)
	if res.Error != nil {
		fmt.Println("usr base tbl find failed! ", res.Error)
		return res.Error
	}
	fmt.Println("usr data find success!")
	fmt.Println("usr name is : ", usrbase.Name)
	fmt.Println("usr passwd is: ", usrbase.Passwd)
	return nil
}

func UsrBaseTblDel(usrbase *UsrBase) error {
	if err := dbmanager.GetDBMgrInst().GetOrmDB().Where(usrbase).Delete(UsrBase{}).Error; err != nil {
		fmt.Println("usrbase table delete failed")
		return err
	}
	fmt.Println("usrbase table delete success")
	return nil
}

func UsrBaseUnscoped(usrbase *UsrBase) error {
	if err := dbmanager.GetDBMgrInst().GetOrmDB().Unscoped().Where(usrbase).Find(&usrbase).Error; err != nil {
		fmt.Println("usr base unscoped find failed")
		return err
	}
	fmt.Println("usr base unscoped find success")
	fmt.Println("usr name is : ", usrbase.Name)
	fmt.Println("usr passwd is: ", usrbase.Passwd)
	return nil
}

func UsrBaseTblUpdate(before *UsrBase, after *UsrBase) error {
	if err := dbmanager.GetDBMgrInst().GetOrmDB().Model(before).Updates(after).Error; err != nil {
		fmt.Println("usr base update failed")
		return err
	}

	fmt.Println("usr base update  success")

	return nil
}
