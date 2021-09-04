package main

import (
	"encoding/json"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"os"
	"flag"
	"strings"
	"net/url"
	"strconv"
)


type result struct {
	Callsign 	string	`json:"callsign"`
	City 		string	`json:"city"`
	Country 	string	`json:"country"`
	FirstName 	string	`json:"fname"`
	Id			int		`json:"id"`
	Remarks 	string	`json:"remarks"`
	State 		string	`json:"state"`
	Surname 	string	`json:"surname"`
}

type results struct {
	Count int 
	Results []result
}



func main() {
	countryFlag := flag.String("country", "United States", "Country")
	stateFlag := flag.String("state", "", "State")
	cityFlag := flag.String("city", "", "City")
	flag.Parse()

	country := strings.Title(strings.ToLower(*countryFlag))
	state := strings.Title(strings.ToLower(*stateFlag))
	city := strings.Title(strings.ToLower(*cityFlag))
	

	fmt.Println("County:", country)
	fmt.Println("State:", state)
	fmt.Println("City:", city)
	
	params := url.Values{}
	if country != "" { params.Add("country", country) }
	if state != "" { params.Add("state", state) }
    if city != "" { params.Add("city", city) }

	url := "https://www.radioid.net/api/dmr/user/?"+params.Encode()

	fmt.Println(url)

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	//fmt.Println(body)

	results1 := results{}
	jsonErr := json.Unmarshal(body, &results1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//fmt.Println(results1.Results[0])


	csvFile, err := os.Create("./data.csv")

	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()


	//Call Type,Contacts Alias,City,Province,Country,Call ID
	//Private Call,$CALLSIGN $FIRSTNAME,Macon,Georgia,United States,3113622

	writer := csv.NewWriter(csvFile)
	var header []string
	header = append(header, "Call Type")
	header = append(header, "Contacts Alias")
	header = append(header, "City")
	header = append(header, "Province")
	header = append(header, "Country")
	header = append(header, "Call ID")
	writer.Write((header))

	for _, result := range results1.Results {
		var row []string
		row = append(row, "Private Call")
		row = append(row, result.Callsign + " " + result.FirstName)
		row = append(row, result.City)
		row = append(row, result.State)
		row = append(row, result.Country)
		row = append(row, strconv.Itoa(result.Id))
		writer.Write(row)
	}

// remember to flush!
writer.Flush()
}