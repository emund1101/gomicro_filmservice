package main

import (
	"films/utils/orm"
	"fmt"
	"gorm.io/gorm"
)

var db *gorm.DB

type seat struct {
	Id     int
	Hall   int
	Line   int
	No     int
	Status int
}

//初始化前
func init() {
	db = orm.Initconf("database_films", "test")

}

func main() {
	//db.Exec("truncate seat")
	data := make([]seat, 0, 280)
	for a := 1; a < 5; a++ {
		for i := 1; i <= 3; i++ {
			for no := 1; no < 11; no++ {
				temp := seat{
					Hall:   a,
					Line:   i,
					No:     no,
					Status: 1,
				}
				data = append(data, temp)
			}
		}
	}

	db.Table("seat").Create(&data)
	s := fmt.Sprintf("%+v ", data)
	fmt.Println(s)
	//fmt.Printf("%+v", data) //打印完整结构体
}
