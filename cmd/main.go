package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Boyaroslav/ultraCalc/pkg/agent"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

type Page struct {
	Title string
	Body  []byte
}

type ToDoItem struct {
	body       string
	answer     int64
	status     int32
	err        error
	inprogress int32
	duration   time.Duration
	date       time.Time
}

var ToDoMx *sync.Mutex
var number int64

const NotDone = 0
const Done = 1
const Error = -1

var ToDo []ToDoItem

func backend() {
	main_page_start, _ := os.ReadFile("index.html")
	main_page_end, _ := os.ReadFile("index_end.html")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(main_page_start))
		num := strconv.Itoa(int(number))

		fmt.Fprintf(w, "<p>"+num+"</p>") // число запущенных агентов для вычисления
		fmt.Fprintf(w, "<br>")
		for _, e := range ToDo {
			fmt.Fprintf(w, e.body)
			if e.status == NotDone {
				fmt.Fprint(w, "  wait...")
			}
			if e.status == Done {

				fmt.Fprint(w, "=", e.answer, "   Done!  Время выполнения - "+e.duration.String())
			}
			if e.status == Error {
				fmt.Fprint(w, "  error  ", fmt.Sprint(e.err))
			}
			fmt.Fprintf(w, "     "+e.date.Format("2006-01-02 15:04:05"))
			fmt.Fprintf(w, "<br>")
		}
		fmt.Fprintf(w, string(main_page_end))

	})
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ToDo = append(ToDo, ToDoItem{string(r.FormValue("text")), 0, 0, nil, 0, 0, time.Now()})
		}
		fmt.Fprintf(w, "okay")

	})

	http.ListenAndServe(":8080", nil)
}
func remove(slice []ToDoItem, s int) []ToDoItem {
	return append(slice[:s], slice[s+1:]...)
}

func Calculate(numberofagents int64) {
	number = 0

	for {
		if number < numberofagents {
			for i, _ := range ToDo {
				if ToDo[i].inprogress == 0 {
					go func() {
						if ToDo[i].inprogress == 1 {
							return
						}
						ToDoMx.Lock()
						atomic.AddInt32(&(ToDo[i].inprogress), 1)
						ToDo[i].inprogress = 1
						atomic.AddInt64(&number, 1)
						start := time.Now()
						res := make(chan agent.Result, 1)
						agent.StartCalculating(ToDo[i].body, res)
						r := <-res

						if r.Err != nil {
							ToDo[i].status = Error
							ToDo[i].err = r.Err
						} else {
							atomic.AddInt32(&(ToDo[i].status), Done)
							atomic.AddInt64(&(ToDo[i].answer), int64(r.Answer))
						}
						duration := time.Since(start)
						ToDo[i].duration = duration
						ToDoMx.Unlock()
						atomic.AddInt64(&number, -1)

					}()
				}
			}
		}
	}
}

func main() {
	//var collection *mongo.Collection
	//var ctx = context.TODO()
	//func init() {
	//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	//client, err := mongo.Connect(ctx, clientOptions)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//err = client.Ping(ctx, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//collection = client.Database("Calculator").Collection("tasks")
	//}
	ToDoMx = &sync.Mutex{}
	numberofagents, err := strconv.Atoi(os.Getenv("NUMBER_OF_AGENTS"))
	if err != nil {
		fmt.Println("Cant get NUMBER_OF_AGENTS variable from env. Default is 10.")
		numberofagents = 10
	}
	go Calculate(int64(numberofagents))
	backend()
}
