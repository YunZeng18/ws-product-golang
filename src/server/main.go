package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"encoding/json"
)


type Counters struct {
    Views  int //field name that starts with capital letter are public, so the json package can access it
    Clicks int 
}

type EventCounter struct {
    sync.Mutex
    counters map[string]Counters
}
func NewEventCounter() *EventCounter {
	e := new(EventCounter)
	e.counters = map[string]Counters {}
	return e
}	
var (
	content = []string{"sports", "entertainment", "business", "education"}
	stats = NewEventCounter()
)

func (e * EventCounter) addView(key string) {
	e.Lock()
	v := e.counters[key]
	v.Views++
	e.counters[key] = v
	defer e.Unlock()
}

func (e * EventCounter) addClick(key string) {
	e.Lock()
	v := e.counters[key]
	v.Clicks++
	e.counters[key] = v
	defer e.Unlock()
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	event_time := fmt.Sprintf("%s:%s",content[rand.Intn(len(content))],time.Now().Format("2006-01-02 15:04"))//why this date is used for formatting though? ask the Go devs XD
	stats.addView(event_time)
	err := processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		processClick(event_time)
	}
}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func processClick(event_time string) error {
	stats.addClick(event_time)
	return nil
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	statsJS, err := json.Marshal(stats.counters)
	fmt.Println(err)
	fmt.Fprint(w, string(statsJS))
}

func isAllowed() bool {
	return true
}

func uploadCounters() error {
	return nil
}
func init() {
	
  }
func main() {
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/stats/", statsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
	
}
