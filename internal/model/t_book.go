package model

type TBook struct {
	Id      int32  `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name    string `gorm:"column:name;default:NULL"`
	Price   string `gorm:"column:price;default:NULL"`
	Author  string `gorm:"column:author;default:NULL"`
	Sales   int32  `gorm:"column:sales;default:NULL"`
	Stock   int32  `gorm:"column:stock;default:NULL"`
	ImgPath string `gorm:"column:img_path;default:NULL"`
}

func (t *TBook) TableName() string {
	return "t_book"
}
