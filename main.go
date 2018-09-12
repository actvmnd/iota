package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/antonholmquist/jason"
)

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	binance, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=IOTABTC")
	if err != nil {
		log.Println(err)
		return
	}
	defer binance.Body.Close()
	bodyBinance, err := ioutil.ReadAll(binance.Body)
	if err != nil {
		log.Println(err)
		return
	}
	vBinance, err := jason.NewObjectFromBytes([]byte(bodyBinance))
	if err != nil {
		log.Println(err)
		return
	}
	iotabtc, err := vBinance.GetString("price")
	if err != nil {
		log.Println(err)
		return
	}
	blockchainINFO, err := http.Get("https://blockchain.info/ticker")
	if err != nil {
		log.Println(err)
		return
	}
	defer blockchainINFO.Body.Close()
	bodyBlockchainINFO, err := ioutil.ReadAll(blockchainINFO.Body)
	if err != nil {
		log.Println(err)
		return
	}
	vBlockchainINFO, err := jason.NewObjectFromBytes([]byte(bodyBlockchainINFO))
	if err != nil {
		log.Println(err)
		return
	}
	btcinr, err := vBlockchainINFO.GetFloat64("INR", "last")
	if err != nil {
		log.Println(err)
		return
	}
	iotabtcFloat, err := strconv.ParseFloat(iotabtc, 64)
	if err != nil {
		log.Println(err)
		return
	}
	iota := iotabtcFloat * btcinr
	query := r.URL.Query().Get("q")
	if query != "" {
		q, err := strconv.ParseFloat(query, 64)
		if err != nil {
			log.Println(err)
			return
		}
		r := fmt.Sprintf("%.2f", iota*q)
		fmt.Fprintf(w, r)
	} else {
		r := fmt.Sprintf("%.2f", iota)
		fmt.Fprintf(w, r)
	}
}

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", handler)
	log.Printf("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}