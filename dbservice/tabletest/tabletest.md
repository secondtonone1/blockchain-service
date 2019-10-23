## table api
### 判断表存在
``` golang
func ProjectTblExist() bool {

	if !dbmanager.GetDBMgrInst().GetOrmDB().HasTable(&Project{}) {
		fmt.Println("project table not exist")
		return false
	}
	fmt.Println("project table has existed")
	return true
}
```
### 创建表
``` golang
func ProjectTblCreate() error {

	if !dbmanager.GetDBMgrInst().GetOrmDB().HasTable(&Project{}) {
		if err := dbmanager.GetDBMgrInst().GetOrmDB().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Project{}).Error; err != nil {
			return err
		}
	}
	return nil
}
```
### 插入表
``` golang
func ProjectTblInsert(project *Project) error {

	if err := dbmanager.GetDBMgrInst().GetOrmDB().Create(project).Error; err != nil {
		fmt.Println("project insert failed ")
		return err
	}
	fmt.Println("project tbl insert success")
	return nil
}
```
### 单表查询
``` golang
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
```

### 联表查询

``` golang
func ProjectTblFindByUsr(usr *UsrBase, project *Project) error {

	dbmanager.GetDBMgrInst().GetOrmDB().Model(usr).Related(project)
	if project.ID == 0 {
		fmt.Println("Can't find project by usr")
		return errors.New("Can't find project by usr")
	}
	fmt.Println("find project by usr success")
	return nil
}
```

### 删除表

``` golang
func ProjectTblDel(project *Project) error {
	if err := dbmanager.GetDBMgrInst().GetOrmDB().Where(project).Delete(Project{}).Error; err != nil {
		fmt.Println("Project table delete failed")
		return err
	}
	fmt.Println("Project table delete success")
	return nil
}
```

### 更新表
``` golang
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
```
## 表结构设计
代码层面通过gorm设定外键关联，而这只是代码层的外键，并没有真正作用于表结构。
所以为了安全，需要将表之间的外键也设置一下。而且保证主键自增且非0.

## 测试代码如何运行
在tabletest目录下
执行 go test -v 可测试所有test文件。
或者 go test 指定测试文件