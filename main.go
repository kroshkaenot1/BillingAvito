package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"strconv"
)

type Users struct {
	Id      uint64 `json:"id"`
	Balance uint64 `json:"balance"`
	Reserve uint64 `json:"reserve"`
}
type Orders struct {
	Id         uint64 `json:"id"`
	Id_user    uint64 `json:"idUser"`
	Id_service uint64 `json:"idService"`
	Cost       uint64 `json:"cost"`
}

func addMoney(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)

	idUser := params["id"]
	money := params["money"]

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/Billing")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if idUser != "" {
		update, err := db.Query(fmt.Sprintf("UPDATE `Users` SET `balance`= `balance` + %s WHERE `id`=%s", money, idUser))
		if err != nil {
			panic(err)
		}
		defer update.Close()
	} else {
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `Users` (`balance`) VALUES(%s)", money))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
	}
}
func reserveMoney(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)

	idUser := params["id_user"]
	idService := params["id_service"]
	idOrder := params["id_order"]
	cost := params["cost"]

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/Billing")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	us, err := db.Query(fmt.Sprintf("SELECT * FROM `Users` WHERE `id`=%s", idUser))
	if err != nil {
		panic(err)
	}
	defer us.Close()

	var user Users
	for us.Next() {
		err = us.Scan(&user.Id, &user.Balance, &user.Reserve)
		if err != nil {
			panic(err)
		}
	}

	b, err := strconv.ParseUint(cost, 10, 64)
	if err != nil {
		panic(err)
		return
	}
	if b > user.Balance {
		panic(err)
		return
	}
	stmt1, err := db.Query(fmt.Sprintf("UPDATE `Users` SET `balance`=`balance` - %s,`reserve`= `reserve` + %s WHERE `id`=%s", cost, cost, idUser))
	if err != nil {
		panic(err)
	}
	stmt1.Close()

	stmt, err := db.Query(fmt.Sprintf("INSERT INTO `Orders` (`id`,`id_user`,`id_service`,`cost`) VALUES (%s,%s,%s,%s)", idOrder, idUser, idService, cost))
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

}
func profit(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)

	idUser := params["id_user"]
	idService := params["id_service"]
	idOrder := params["id_order"]
	cost := params["cost"]

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/Billing")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Query(fmt.Sprintf("SELECT * FROM `Users` WHERE `id`=%s", idUser))
	if err != nil {
		panic(err)
	}
	stmt.Close()

	var user Users
	for stmt.Next() {
		err = stmt.Scan(&user.Id, &user.Balance, &user.Reserve)
		if err != nil {
			panic(err)
		}
	}

	reserve := user.Reserve

	stmt1, err := db.Query(fmt.Sprintf("SELECT * FROM `Orders` WHERE `id`=%s AND `id_user`=%s AND `id_service`=%s AND`cost`=%s", idOrder, idUser, idService, cost))
	if err != nil {
		panic(err)
	}
	stmt1.Close()

	var order Orders
	for stmt1.Next() {
		err := stmt1.Scan(&order.Id, &order.Id_user, &order.Id_service, &order.Cost)
		if err != nil {
			panic(err)
		}
	}

	reserve = reserve - order.Cost
	stmt2, err := db.Query(fmt.Sprintf("UPDATE `Users` SET `reserve` = %d WHERE `id`= %s", reserve, idUser))
	if err != nil {
		panic(err)
	}
	stmt2.Close()

	report, err := os.Create("./Report.csv")
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(report)
	var data = [][]string{
		{"Name", "Age", "Occupation"},
		{"Sally", "22", "Nurse"},
		{"Joe", "43", "Sportsman"},
		{"Louis", "39", "Author"},
	}

	err = writer.WriteAll(data)

	if err != nil {
		panic(err)
	}
}
func getBalance(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	id_user, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/Billing")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.Query(fmt.Sprintf("SELECT * FROM `Users` WHERE `id`=%d", id_user))
	if err != nil {
		panic(err)
	}
	defer res.Close()
	var user Users
	for res.Next() {
		err = res.Scan(&user.Id, &user.Balance, &user.Reserve)
		if err != nil {
			panic(err)
		}
	}
	answ, _ := json.Marshal(user)
	w.Write(answ)
}
func handleFunc() {
	http.HandleFunc("/add", addMoney)
	http.HandleFunc("/reserve", reserveMoney)
	http.HandleFunc("/profit", profit)
	http.HandleFunc("/getBalanceOfUser", getBalance)
	http.ListenAndServe(":8080", nil)
}
func main() {
	handleFunc()
}
