package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

type Users struct {
	Id      uint64 `json:"id"`
	Balance uint64 `json:"balance"`
	Reserve uint64 `json:"reserve"`
}

func add_money(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)

	id_user := params["id"]
	money := params["money"]

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/Billing")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if id_user != "" {
		update, err := db.Query(fmt.Sprintf("UPDATE `Users` SET `balance`= `balance` + %s WHERE `id`=%s", money, id_user))
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
func reserve_money(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)

	id_user := params["id_user"]
	id_service := params["id_service"]
	id_order := params["id_order"]
	cost := params["cost"]
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/Billing")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	us, err := db.Query(fmt.Sprintf("SELECT * FROM `Users` WHERE `id`=%s", id_user))
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
	stmt1, err := db.Query(fmt.Sprintf("UPDATE `Users` SET `balance`=`balance` - %s,`reserve`= `reserve` + %s WHERE `id`=%s", cost, cost, id_user))
	if err != nil {
		panic(err)
	}
	stmt1.Close()
	stmt, err := db.Query(fmt.Sprintf("INSERT INTO `Orders` (`id`,`id_user`,`id_service`,`cost`) VALUES (%s,%s,%s,%s)", id_order, id_user, id_service, cost))
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

}
func profit(w http.ResponseWriter, r *http.Request) {
	id_user := r.FormValue("id")
	id_service := r.FormValue("id_service")
	id_order := r.FormValue("id_order")
	cost := r.FormValue("cost")
	fmt.Print(id_user, id_service, id_order, cost)
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
	http.HandleFunc("/add", add_money)
	http.HandleFunc("/reserve", reserve_money)
	http.HandleFunc("/profit", profit)
	http.HandleFunc("/getBalanceOfUser", getBalance)
	http.ListenAndServe(":8080", nil)
}
func main() {
	handleFunc()
}
