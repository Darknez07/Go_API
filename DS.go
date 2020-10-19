package main
import(
	"time"
)

// type rsvp int
// var rsvp map[string]int;
type (
	//Date Type
	Date struct{
		Month int
		Year int
		Day int
		Hour int
		Minutes int
		Seconds int
	}
	//Partis short for Participants
	Partis struct{
		Name string
		Email string
		RSVP int
	}
	// Meeting is The main DS
	Meeting struct{
		//Id for just
		Id int
		Title string
		Participants []Partis
		StartTime Date
		EndTime Date
		Timestamp time.Time
	}
)