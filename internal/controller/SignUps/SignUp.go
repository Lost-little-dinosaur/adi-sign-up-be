package SignUps

import (
	"adi-sign-up-be/internal/db"
	"adi-sign-up-be/internal/logger"
	"adi-sign-up-be/internal/model/Mysql"
	"gorm.io/gorm"
	"sync"
)

var dbManage *SignUpDBManage = nil

func init() {
	logger.Info.Println("[ catalogues ]start init Table ...")
	dbManage = GetManage()
}

type SignUpDBManage struct {
	mDB     *db.MainGORM
	sDBLock sync.RWMutex
}

func (m *SignUpDBManage) getGOrmDB() *gorm.DB {
	return m.mDB.GetDB()
}

func (m *SignUpDBManage) atomicDBOperation(op func()) {
	m.sDBLock.Lock()
	op()
	m.sDBLock.Unlock()
}

func GetManage() *SignUpDBManage {
	if dbManage == nil {
		var catalogueDb = db.MustCreateGorm()
		err := catalogueDb.GetDB().AutoMigrate(&Mysql.SignUp{}, &Mysql.Member{}) //自动创建表
		if err != nil {
			logger.Error.Fatalln(err)
			return nil
		}
		dbManage = &SignUpDBManage{mDB: catalogueDb}
	}
	return dbManage
}

//以上代码是初始化数据库表以及自动创建表所需的代码，下面是具体操作数据库的代码

func AddSignUp(signUp *Mysql.SignUp) (error, int64) {
	err := GetManage().getGOrmDB().Model(&Mysql.SignUp{}).Create(signUp).Error
	var countID *int64
	countID = new(int64)
	GetManage().getGOrmDB().Model(&Mysql.SignUp{}).Count(countID)
	return err, *countID
}

func AddMember(member *Mysql.Member) (error, string) {
	return GetManage().getGOrmDB().Model(&Mysql.Member{}).Create(member).Error, member.ID
}

func GetAllSignUp() (error, []Mysql.SignUp) {
	var returnArr []Mysql.SignUp
	return GetManage().getGOrmDB().Model(&Mysql.SignUp{}).Find(&returnArr).Error, returnArr
}

func GetMemberByID(ID string) (error, Mysql.Member) {
	var returnStruct Mysql.Member
	return GetManage().getGOrmDB().Model(&Mysql.Member{}).Where("id = ?", ID).First(&returnStruct).Error, returnStruct
}

func CheckIdNumberExist(idNumber string) (error, bool) {
	err := GetManage().getGOrmDB().Model(&Mysql.Member{}).Where("id_number = ?", idNumber).First(&Mysql.Member{}).Error
	if err != nil {
		return err, false
	}
	return nil, true
}
