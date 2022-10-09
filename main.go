package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

type User struct {
	Uid      int
	Name     string
	Phone    string
	Email    string
	Password string
}

func queryMultiRow() []User {
	rows, err := db.Query("select uid, name, phone, email from `user` where uid > ?", 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil
	}
	defer rows.Close()
	users := []User{}
	var u User
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

func ExportCsv(filePath string, data [][]string) {
	fp, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("创建文件["+filePath+"]句柄失败，%v", err)
		return
	}
	defer fp.Close()
	fp.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(fp)
	w.WriteAll(data)
	w.Flush()
}

func main() {
	filename := "./exportUsers.csv"
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
	ExportCsv(filename, column)
}
