package main
import (
	"time"
	"fmt"
)

var rsvp map[string]int;
var revrsvp map[int]string;
func main() {
	rsvp = make(map[string]int)
	rsvp["Yes"]= 1
	rsvp["No"] = 2
	rsvp["Maybe"] = 3
	rsvp["Not Answered"] = 4
	revrsvp = make(map[int]string)
	revrsvp = reverseMap(rsvp)
	pump :=[]Partis{};
	q :=Partis{
		Name: "One",
		Email: "Hello123@gmail.com",
		RSVP: 1,
	};
	pump = append(pump, q)
	meeting := Meeting{
		Id:           0,
		Title:        "Title 1",
		Participants: pump,
		StartTime:    Date{
			Month:   10,
			Year:    2020,
			Day:     18,
			Hour:    11,
			Minutes: 24,
			Seconds: 11,
		},
		EndTime:      Date{
			Month:   10,
			Year:    2020,
			Day:     18,
			Hour:    12,
			Minutes: 40,
			Seconds: 12,
		},
		Timestamp:  time.Now(),
	};
	fmt.Println(meeting);
	handler()
}

func reverseMap(m map[string]int) map[int]string {
    n := make(map[int]string)
    for k, v := range m {
        n[v] = k
    }
    return n
}
