package main

import (
	// "context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	// "reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func JsonOut(message string) string {
	var raw map[string]interface{}
	json.Unmarshal([]byte(message), &raw)
	out, _ := json.Marshal(raw)
	return string(out)
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func checkforerr(p Partis) bool {
	if p.Email == "error" {
		return true
	}
	return false
}

func CheckEmail(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func CheckErr(d Date) bool {
	if d.Day == 0 {
		return true
	} else {
		return false
	}
}

func toDate(s string) Date {
	dp := Date{}
	p := strings.Split(s, "/")
	if len(p) == 3 {
		dp.Day, _ = strconv.Atoi(p[0])
		dp.Month, _ = strconv.Atoi(p[1])
		yearTime := strings.Split(p[2], " ")
		if len(yearTime) > 1 {
			dp.Year, _ = strconv.Atoi(yearTime[0])
			onlyTime := strings.Split(yearTime[1], ":")
			if len(onlyTime) > 1 {
				dp.Hour, _ = strconv.Atoi(onlyTime[0])
				dp.Minutes, _ = strconv.Atoi(onlyTime[1])
			} else {
				dp.Day = 0
			}
		} else {
			dp.Day = 0
		}
		dp.Seconds = 0
	} else {
		dp.Day = 0
	}
	return dp
}

func Checktime(d1 Date, d2 Date) bool {
	if d1.Year != d2.Year {
		return true
	} else if d1.Month != d2.Month {
		return true
	} else if d1.Day != d2.Day {
		return true
	} else if d2.Hour <= d1.Hour {
		return true
	} else if d2.Minutes <= d1.Minutes {
		return true
	}
	return false
}

type test_struct struct {
	Id           int
	Title        string
	Participants []struct {
		Name  string
		Email string
		RSVP  string
	}
	StartTime string
	EndTime   string
	Timestamp time.Time
}

var x = getLatestId() + 1

func meetings(rw http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		decoder := json.NewDecoder(request.Body)
		errors := false
		var t test_struct
		err := decoder.Decode(&t)
		var meet Meeting
		var p Partis
		pump := []Partis{}
		if err != nil {
			panic(err)
		}
		t.Timestamp = time.Now()
		meet.Id = x
		meet.Title = t.Title
		for _, s := range t.Participants {
			p = Partis{}
			p.Name = s.Name
			if !CheckEmail(s.Email) {
				fmt.Fprintf(rw, JsonOut(`{"Message":"Email Not correct"}`))
				fmt.Fprintf(rw, "\n")
				errors = true
				break
			}
			p.Email = s.Email
			p.RSVP = rsvp[s.RSVP]
			pump = append(pump, p)
		}
		meet.Participants = pump
		meet.StartTime = toDate(t.StartTime)
		if CheckErr(meet.StartTime) && !errors {
			fmt.Fprintf(rw, JsonOut(`{"Message": "Date format Wrong" }`))
			fmt.Fprintf(rw, "\n")
			errors = true
		}
		meet.EndTime = toDate(t.EndTime)
		if CheckErr(meet.EndTime) && !errors {
			fmt.Fprintf(rw, JsonOut(`{"Message": "Date format Wrong" }`))
			fmt.Fprintf(rw, "\n")
			errors = true
		}
		meet.Timestamp = t.Timestamp
		meets := &meet
		if Checktime(meet.StartTime, meet.EndTime) && !errors {
			fmt.Fprintf(rw, JsonOut(`{"Message":"This meeting cannot be set because may be longer than 1 day or illformatted"}`))
			errors = true
		}
		b, _ := json.Marshal(&meets)
		if !errors {
			x = ScheduleMeet(meet)
			if x == -1 {
				fmt.Fprintf(rw, JsonOut(`{"Message":"The meeting already exist or participant is not free"}`))
			} else {
				fmt.Fprintf(rw, string(b))
			}
		}
		x++
	case "GET":
		ends := Date{}
		starts := Date{}
		end := "Start"
		start := "Start"
		num, err := strconv.Atoi(strings.Split(request.URL.Path, "/")[2])
		start = request.URL.Query().Get("start")
		end = request.URL.Query().Get("end")
		ends = toDate(end)
		starts = toDate(start)
		var final string
		if err == nil {
			x := FindMeeting(num)
			if x != nil {
				for _, q := range x {
					// fmt.Println(q)
					for ll, s := range q {
						// fmt.Println(reflect.TypeOf(s))
						// fmt.Println(f)
						switch t := s.(type) {
						case primitive.M:
							final+=`, "`+ll+`" : [{`
							for label, nw := range s.(primitive.M){
								// fmt.Printf(label)
								// fmt.Println(label, nw)
								// fmt.Println(reflect.TypeOf(nw))
								// fmt.Println(label)
								final+=`"`+label+`" : `+ strconv.Itoa(int(nw.(int32)))+`, `
							}
							final = final[:len(final) - 2]
							final+=`}], `
							final = final[:len(final) - 2]
						case primitive.A:
							final += `"Participants" : [`
							for last, w := range s.(primitive.A) {
								final += `{ `
								for mj, y := range w.(primitive.M) {
									fmt.Println(reflect.TypeOf(mj), "KJ")
									switch h := y.(type) {
									case int32:
										f := y.(int32)
										final += `"RSVP" : ` + `"`+revrsvp[int(f)]+`"` + `, `
									case string:
										f := y.(string)
										final += `"`+mj +`" : "` +f +`", `
									default:
										fmt.Println(h)
									}
								}
								final = final[:len(final) - 2]
								if last == len(s.(primitive.A))-1 {
									final += `} `
								} else {
									final += `}, `
								}
							}
							final += `] `
						case int32:
							f := s.(int32)
							final += `"Id" : ` +strconv.Itoa(int(f)) +`, `
						case string:
							final += `"Title" : "`+s.(string)+`", `
						case primitive.ObjectID:
							continue
						case primitive.DateTime:
							fmt.Println(s)
						default:
							fmt.Println(t)
						}
					}
				}
				fmt.Fprintf(rw, JsonOut(`{`+final+`}`))
			}else if x == nil{
			fmt.Fprintf(rw, JsonOut(`{"Message":"The Format of Request Not correct"}`))
			}
		} else {
			if starts.Day != 0 && ends.Day != 0 {
				FindDatedMeeting(starts, ends)
			} else {
				tt:= allMeetings()
				if tt != nil {
					for _, q := range tt {
						final = ""
						// fmt.Println(q)
						for ll, s := range q {
							// fmt.Println(reflect.TypeOf(s))
							// fmt.Println(f)
							switch t := s.(type) {
							case primitive.M:
								final+=`, "`+ll+`" : [{`
								for label, nw := range s.(primitive.M){
									// fmt.Printf(label)
									// fmt.Println(label, nw)
									// fmt.Println(reflect.TypeOf(nw))
									// fmt.Println(label)
									final+=`"`+label+`" : `+ strconv.Itoa(int(nw.(int32)))+`, `
								}
								final = final[:len(final) - 2]
								final+=`}], `
								final = final[:len(final) - 2]
							case primitive.A:
								fmt.Println("Idhar")
								final += `"Participants" : [`
								for last, w := range s.(primitive.A) {
									final += `{ `
									for mj, y := range w.(primitive.M) {
										fmt.Println(reflect.TypeOf(mj), "KJ")
										switch h := y.(type) {
										case int32:
											f := y.(int32)
											final += `"RSVP" : ` + `"`+revrsvp[int(f)]+`"` + `, `
										case string:
											f := y.(string)
											final += `"`+mj +`" : "` +f +`", `
										default:
											fmt.Println(h)
										}
									}
									final = final[:len(final) - 2]
									if last == len(s.(primitive.A))-1 {
										final += `} `
									} else {
										final += `}, `
									}
								}
								final += `] `
							case int32:
								f := s.(int32)
								fmt.Println(strconv.Itoa(int(f)))
								final += `"Id" : ` +strconv.Itoa(int(f)) +`, `
							case string:
								fmt.Println(s.(string))
								final += `"Title" : "`+s.(string)+`", `
							case primitive.ObjectID:
								continue
							case primitive.DateTime:
								fmt.Println(s)
							default:
								fmt.Println("Default")
								fmt.Println(t)
							}
						}
						fmt.Fprintf(rw, JsonOut(`{`+final+`}`))
					}
					fmt.Println(final)
					var raw map[string] interface{}
					fmt.Println(json.Unmarshal([]byte(final), &raw))
					fmt.Fprintf(rw, JsonOut(`{`+final+`}`))
				}else if tt == nil{
				fmt.Fprintf(rw, JsonOut(`{"Message":"The Format of Request Not correct"}`))
				}
			}
		}
	}
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
}

func handler() {
	http.HandleFunc("/meetings", meetings)
	http.HandleFunc("/meetings/", meetings)
	log.Fatal(http.ListenAndServe(":801", nil))
}
