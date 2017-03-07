package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-toast/toast"
)

var apiKey string
var lastTransactionSell int
var lastTransactionBuy int

type apiString struct {
	ApiKey string      `json:"apiKey"`
	Items  []watchItem `json:"watchItem"`
}

type watchItem struct {
	Id        int
	CheckBuy  bool
	BuyPrice  int
	CheckSell bool
	SellPrice int
}

type transaction struct {
	Id        int
	Item_id   int
	Price     int
	Quantity  int
	Created   time.Time
	Purchased time.Time
}

type objet struct {
	Id           int64
	Chat_link    string
	Name         string
	Icon         string
	Description  string
	Type         string `json:"type"`
	Rarity       string
	Level        int64
	Vendor_value int64
	Default_skin int64
	Flags        []string
	GameType     []string `json:"gametype"`
	Restrictions []string
	Details      []string
}

func main() {
	apiKey = getApiKey("./config.json")
	lastTransactionSell = 0
	lastTransactionBuy = 0
	doEvery(5 * time.Second)

}

func doEvery(d time.Duration) {
	for x := range time.Tick(d) {
		checkTransaction(x)
	}
}

func checkTransaction(t time.Time) {
	var transactionSell []transaction
	var transactionBuy []transaction
	var item objet
	urlItem := "https://api.guildwars2.com/v2/items/"
	urlSell := "https://api.guildwars2.com/v2/commerce/transactions/history/sells?access_token=" + apiKey
	urlBuy := "https://api.guildwars2.com/v2/commerce/transactions/history/buys?access_token=" + apiKey
	getJson(urlSell, &transactionSell)
	getJson(urlBuy, &transactionBuy)

	/*if lastTransactionSell == 0 || transactionSell[0].Id == lastTransactionSell {
		lastTransactionSell = transactionSell[0].Id
	} else {
		lastTransactionSell = transactionSell[0].Id
		getJson(urlItem+strconv.Itoa(transactionSell[0].Item_id), &item)
		doANotif(transactionSell[0].Quantity, transactionSell[0].Price, item.Name, "sold")
	}*/
	fmt.Println(urlBuy)
	if transactionSell[0].Id != lastTransactionSell && lastTransactionSell != 0 {
		lastTransactionSell = transactionSell[0].Id
		getJson(urlItem+strconv.Itoa(transactionSell[0].Item_id), &item)
		doANotif(transactionSell[0].Quantity, transactionSell[0].Price, item.Name, "sold")
	} else {
		lastTransactionSell = transactionSell[0].Id
		getJson(urlItem+strconv.Itoa(transactionSell[0].Item_id), &item)
		doANotif(transactionSell[0].Quantity, transactionSell[0].Price, item.Name, "sold")
	}

	/*if lastTransactionBuy == 0 || transactionBuy[0].Id == lastTransactionBuy+1 {
		lastTransactionBuy = transactionBuy[0].Id
		fmt.Println("transaction.Id")
		fmt.Println(transactionBuy[0].Id)
		fmt.Println("lasttransaction")
		fmt.Println(lastTransactionBuy)
		getJson(urlItem+strconv.Itoa(transactionBuy[0].Item_id), &item)
		doANotif(transactionBuy[0].Quantity, transactionBuy[0].Price, item.Name, "bought")
	} else {
		lastTransactionBuy = transactionBuy[0].Id
		getJson(urlItem+strconv.Itoa(transactionBuy[0].Item_id), &item)
		doANotif(transactionBuy[0].Quantity, transactionBuy[0].Price, item.Name, "bought")
	}*/

	if transactionBuy[0].Id != lastTransactionBuy && lastTransactionBuy != 0 {
		lastTransactionBuy = transactionBuy[0].Id
		getJson(urlItem+strconv.Itoa(transactionBuy[0].Item_id), &item)
		doANotif(transactionBuy[0].Quantity, transactionBuy[0].Price, item.Name, "bought")
	} else {
		lastTransactionBuy = transactionBuy[0].Id
	}

}

func doANotif(quantity int, price int, name string, transaction string) {

	notification := toast.Notification{
		AppID:   "Personal Transaction Notifier",
		Title:   "A wild transaction appears",
		Message: strconv.Itoa(quantity) + " " + name + " " + transaction + " for " + strconv.Itoa(price) + " coppers.",
		//Icon:    ".\\gw2Icon.png",

		Audio: toast.Default,
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func getApiKey(url string) string {
	file, err := ioutil.ReadFile(url)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var jsontype apiString
	json.Unmarshal(file, &jsontype)
	return jsontype.ApiKey
}
