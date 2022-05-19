package Model

type Administrator struct {
	AdminId   string `gorm:"PRIMARY KEY"`
	AdminName string `gorm:"NOT NULL"`
	Function  string `gorm:"DEFAULT '管理日常事务'"`
}

func GetAdministratorById(ADId string) *Administrator {
	var admin Administrator
	result := GlobalConn.Where(&Administrator{AdminId: ADId}).First(&admin)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return &admin
}
