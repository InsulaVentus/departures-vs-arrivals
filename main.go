package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const urlFormat = "https://avinor.no/Api/Flights/Airport/%s?direction=%s&start=%s&end=%s&language=en"

type flight struct {
	ToAirportIATA   string `json:"ToAirport"`
	ToAirportName   string `json:"ToAirportName"`
	FromAirportIATA string `json:"FromAirport"`
	FromAirportName string `json:"FromAirportName"`
	AirlineName     string `json:"AirlineName"`
	FlightId        string `json:"FlightId"`
	When            int64  `json:"Timestamp"`
}

func (d *flight) String() string {
	return fmt.Sprintf("%s - %s (%s) - %s %s", time.Unix(d.When, 0).Format("15:04"), d.ToAirportName, d.ToAirportIATA, d.AirlineName, d.FlightId)
}

func main() {
	fmt.Println("Comparing departures vs arrivals")

	departuresFrom, err := time.Parse(time.RFC3339, "2019-10-04T05:00:00+01:00")
	if err != nil {
		log.Fatal(err)
	}

	departuresTo, err := time.Parse(time.RFC3339, "2019-10-04T10:00:00+01:00")
	if err != nil {
		log.Fatal(err)
	}

	urlDepartures := CreateDepartureUrl("OSL", departuresFrom, departuresTo)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	rawDepartures, err := GetRawContent(urlDepartures, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	departures, err := Parse(rawDepartures)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d Departures from OSL between %s and %s:\n", len(departures), departuresFrom.Format(time.Stamp), departuresTo.Format(time.Stamp))
	for _, value := range departures {
		fmt.Printf("%s", value.String())
	}

	departureMap := CreateDepartureMap(departures)

	arrivalsFrom, err := time.Parse(time.RFC3339, "2019-10-06T13:00:00+01:00")
	if err != nil {
		log.Fatal(err)
	}

	arrivalsTo, err := time.Parse(time.RFC3339, "2019-10-06T19:00:00+01:00")
	if err != nil {
		log.Fatal(err)
	}

	urlArrivals := CreateArrivalUrl("OSL", arrivalsFrom, arrivalsTo)

	rawArrivals, err := GetRawContent(urlArrivals, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	arrivals, err := Parse(rawArrivals)
	if err != nil {
		log.Fatal(err)
	}

	arrivalMap := CreateArrivalMap(arrivals)

	fmt.Println()

	for airport, arrivalsFromAirport := range arrivalMap {

		departuresFromAirport, ok := departureMap[airport]
		if ok {
			fmt.Printf("%s:\n", airport)
			fmt.Printf("Departures: %v\n", departuresFromAirport)
			fmt.Printf("Arrivals: %v\n", arrivalsFromAirport)
			fmt.Println()
		}

	}
}

func CreateDepartureUrl(departureAirport string, from time.Time, to time.Time) string {
	return fmt.Sprintf(urlFormat, departureAirport, "d", from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))
}

func CreateArrivalUrl(arrivalAirport string, from time.Time, to time.Time) string {
	return fmt.Sprintf(urlFormat, arrivalAirport, "a", from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))
}

func GetRawContent(url string, client *http.Client) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func CreateDepartureMap(flights []flight) map[string][]flight {
	var m = make(map[string][]flight)

	for _, f := range flights {
		m[f.ToAirportIATA] = append(m[f.ToAirportIATA], f)
	}
	return m
}

func CreateArrivalMap(flights []flight) map[string][]flight {
	var m = make(map[string][]flight)

	for _, f := range flights {
		m[f.FromAirportIATA] = append(m[f.FromAirportIATA], f)
	}
	return m
}


func Parse(bs []byte) ([]flight, error) {

	var flights = struct {
		Flights []flight `json:"Flights"`
	}{}

	err := json.Unmarshal(bs, &flights)
	if err != nil {
		return nil, err
	}

	return flights.Flights, nil
}
