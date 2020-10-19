package main

import (
	"context"
	"fmt"
	"time"
	"reflect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type retrival  struct {
	id int `bson:"id"`
	Title   string `bson:"title"`
	Participants []Partis `bson:"participants"`
	StartTime Date `bson:"starttime"`
	EndTime   Date `bson:"endtime"`
	Timestamp time.Time `bson:"timestamp"`
}

func ScheduleMeet(meeting Meeting) int {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://12345:12345@cluster0.pexud.mongodb.net/First?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		return -1
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return -1
	}
	defer client.Disconnect(ctx)
	dbs := client.Database("Custom").Collection("Meetings")
	if CheckBooking(meeting) {
		Meetings, err := dbs.InsertOne(ctx, meeting)
		if err != nil {
			return -1
		}
		fmt.Println(Meetings)
		return 1
	}
	// opts := options.Find().SetSort(bson.D{{"Timestamp", -1}})
	// sortCursor, err := dbs.Find(ctx, bson.M{}, opts)
	// var output []bson.M
	// if err = sortCursor.All(ctx, &output); err != nil {
	// return -1
	// }
	// what := output[len(output)-1]
	// panic, check := what["id"].(int32)
	// fmt.Println(panic, check, what)
	// return int(panic)
	return -1
}

func CheckBooking(meeting Meeting) bool {
	clientOptions := options.Client().ApplyURI("mongodb+srv://12345:12345@cluster0.pexud.mongodb.net/First?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		return false
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return false
	}
	defer client.Disconnect(ctx)
	emails := []bson.A{}
	dbs := client.Database("Custom").Collection("Meetings")
	for _, s := range meeting.Participants {
		emails = append(emails, bson.A{s.Email})
	}
	fmt.Println(emails)
	fmt.Println(meeting.StartTime)
	cur, err := dbs.Find(ctx,
		bson.D{{"participants.email",
			bson.D{{"$in", emails}}},
			{"starttime.year", meeting.StartTime.Year},
			{"starttime.month", meeting.StartTime.Month},
			{"starttime.day", meeting.StartTime.Day},
			{"endtime.hour",
				bson.D{{"$gte",
					meeting.StartTime.Hour}}},
			{"endtime.minutes",
				bson.D{{"$gte",
					meeting.StartTime.Minutes}}},
		})
	if err != nil {
		return false
	}
	defer cur.Close(ctx)
	var out []bson.M
	cur.All(ctx, &out)
	fmt.Println(out)
	if len(out) == 0 {
		return false
	}
	for _, s := range out {
		fmt.Println(s["participants"])
	}
	return true
}

func FindMeeting(s int) []bson.M{
	clientOptions := options.Client().ApplyURI("mongodb+srv://12345:12345@cluster0.pexud.mongodb.net/First?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		return nil
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil
	}
	defer client.Disconnect(ctx)
	dbs := client.Database("Custom").Collection("Meetings")
	cur, err := dbs.Find(ctx, bson.D{{"id", s}})
	if err != nil {
		return nil
	}
	defer cur.Close(ctx)
	var out []bson.M
	var retrive retrival
	cur.All(ctx, &out)
	cur.Decode(&retrive)
	for _, s := range out {
		fmt.Println(s["id"])
		fmt.Println(s["title"])
		news := reflect.TypeOf(s["participants"])
		fmt.Println(news)
		fmt.Println(s["starttime"])
		fmt.Println(s["endtime"])
		fmt.Println(s["timestamp"])
	}
	return out
}

func getLatestId() int {
	clientOptions := options.Client().ApplyURI("mongodb+srv://12345:12345@cluster0.pexud.mongodb.net/First?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		return -1
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return -1
	}
	defer client.Disconnect(ctx)
	dbs := client.Database("Custom").Collection("Meetings")
	opts := options.Find().SetSort(bson.D{{"Timestamp", -1}})
	sortCursor, err := dbs.Find(ctx, bson.M{}, opts)
	var output []bson.M
	if err = sortCursor.All(ctx, &output); err != nil {
		return -1
	}
	what := output[len(output)-1]
	panic, check := what["id"].(int32)
	fmt.Println(panic, check, what)
	return int(panic)
}

func FindDatedMeeting(d1 Date, d2 Date) int {
	clientOptions := options.Client().ApplyURI("mongodb+srv://12345:12345@cluster0.pexud.mongodb.net/First?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		return -1
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return -1
	}
	defer client.Disconnect(ctx)
	dbs := client.Database("Custom").Collection("Meetings")
	cur, err := dbs.Find(ctx, bson.D{
		{"starttime.year", d1.Year},
		{"starttime.month", d1.Month},
		{"starttime.day", d1.Day},
		{"starttime.hour",
			bson.D{{"$lte",
				d1.Hour}}},
		{"starttime.minutes",
			bson.D{{"$lte",
				d1.Minutes}}},
		{"endtime.year", d2.Year},
		{"endtime.month", d2.Month},
		{"endtime.day", d2.Day},
		{"endtime.hour",
			bson.D{{"$lte",
				d2.Hour}}},
		{"endtime.minutes",
			bson.D{{"$lte",
				d2.Minutes}}}})
	if err != nil {
		return -1
	}
	defer cur.Close(ctx)
	var out []bson.M
	cur.All(ctx, &out)
	for _, s := range  out{
		fmt.Println("API")
		fmt.Println(s)
	}
	return 1
}
