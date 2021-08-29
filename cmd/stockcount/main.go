package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"

	tc "github.com/0dayfall/ctw/tweet/recentcount"
)

type StockFile struct {
	Stock []struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"data"`
}

func main() {
	//READ FILE WITH STOCKS
	jsonFile, err := os.Open(path.Join(path.Dir(os.Args[1]), "stocks.json"))
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	var stocks StockFile
	err = json.Unmarshal(byteValue, &stocks)
	if err != nil {
		log.Fatal(err)
	}

	//WRITE RESULTS TO A FILE
	f, err := os.Create(path.Join(path.Dir(os.Args[2]), "stocks.csv"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)

	//LOOP THEM THROUGH
	for _, stock := range stocks.Stock {

		//COUNT EACH ONE
		recentCount := tc.GetRecentCount(stock.Symbol+" lang:en", "day")
		log.Println(recentCount)
		countDatas := recentCount.Data
		log.Println(countDatas)
		for _, countData := range countDatas {
			var record []string
			record = append(record, stock.Symbol)
			record = append(record, countData.Start.String())
			record = append(record, strconv.Itoa(countData.TweetCount))

			log.Println(record)
			err := w.Write(record)
			if err != nil {
				log.Fatalln(err)
			}
		}
		w.Flush()
	}
}
