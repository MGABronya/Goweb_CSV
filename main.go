// @Title  main
// @Description  从数据库中导出一个CSV文件
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-09 10:35
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User				定义了用户相关信息
type User struct {
	Uid      int
	Name     string
	Phone    string
	Email    string
	Password string
}

var db *gorm.DB

// @title    main
// @description   用于将数据库中的数据导出为csv文件
// @auth      MGAronya（张健）             2022-10-08 19:26
// @param     void
// @return    void
func main() {
	// TODO 设置导出的文件名
	filename := "./exportUsers.csv"

	// TODO 从数据库中获取数据
	users := queryMultiRow()
	column := [][]string{{"手机号", "用户UID", "Email", "用户名"}}
	for _, u := range users {
		str := []string{}
		str = append(str, u.Phone)
		str = append(str, strconv.Itoa(u.Uid))
		str = append(str, u.Email)
		str = append(str, u.Name)
		column = append(column, str)
	}

	// TODO 导出csv文件
	ExportCsv(filename, column)
}

// @title    queryMultiRow
// @description   用于获取数据库中的数据
// @auth      MGAronya（张健）             2022-10-08 19:26
// @param     void
// @return    []User			获取的所有的用户信息
func queryMultiRow() []User {
	rows, err := db.Query("select uid, name, phone, email from `user` where uid > ?", 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil
	}
	// TODO 关闭rows， 释放持有的数据库连接
	defer rows.Close()
	users := []User{}
	var u User

	// TODO 循环读取结果集中的数据
	for rows.Next() {
		err := rows.Scan(&u.Uid, &u.Name, &u.Phone, &u.Email)
		users = append(users, u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil
		}
	}
	return users
}


// @title    ExportCsv
// @description   用于获取数据库中的数据
// @auth      MGAronya（张健）             2022-10-08 19:26
// @param     filePath string, data [][]string		导出csv文件
// @return    void
func ExportCsv(filePath string, data [][]string) {
	// TODO 创建文件句柄
	fp, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("创建文件["+filePath+"]句柄失败，%v", err)
		return
	}
	defer fp.Close()

	// TODO 写入UTF-8 BOM
	fp.WriteString("\xEF\xBB\xBF")

	// TODO 创建一个新的文件流
	w := csv.NewWriter(fp)
	w.WriteAll(data)
	w.Flush()
}

// @title    InitDB
// @description   从配置文件中读取数据库相关信息后，完成数据库初始化
// @auth      MGAronya（张健）             2022-9-16 10:07
// @param     void        void         没有入参
// @return    db        *gorm.DB         将返回一个初始化后的数据库指针
func InitDB() *gorm.DB {
	//driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc),
	)
	var err error
	db, err = gorm.Open(mysql.Open(args), &gorm.Config{})
	// TODO  如果未能连接到数据库，终止程序并返回错误信息
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}
}