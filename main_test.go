package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {

	input := []byte(`{"Flights":[{"Id":"sk4675-osl-alc-20191004","AirlineName":"SAS","AirportNameToOrFrom":"Alicante","AirportToOrFrom":"ALC","CodeShares":[],"FlightId":"SK4675","FromAirport":"OSL","FromAirportName":"Oslo","Gate":null,"Belt":null,"GateOrBelt":null,"GateOrBeltStatus":null,"IsDeparture":true,"IsOld":false,"ScheduleChanged":false,"ScheduledTime":"06:00","ScheduledTimeFull":"201910040600","Status":null,"ToAirport":"ALC","ToAirportName":"Alicante","Date":"20191004","CheckInZones":null,"Timestamp":1570161600}]}`)

	want := map[string][]flight{
		"ALC": {{
			ToAirportIATA: "ALC",
			ToAirportName: "Alicante",
			AirlineName:   "SAS",
			FlightId:      "SK4675",
			When:          1570161600,
		}},
	}

	got, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}
