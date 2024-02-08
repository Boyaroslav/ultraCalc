package main

import (
	"fmt"
	"net/http"
	"os"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

type Page struct {
	Title string
	Body  []byte
}

type ToDoItem struct {
	body       string
	answer     int
	status     int
	err        error
	inprogress int
}

const NotDone = 0
const Done = 1
const Error = -1

var ToDo []ToDoItem

func backend() {
	main_page_start, _ := os.ReadFile("index.html")
	main_page_end, _ := os.ReadFile("index_end.html")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(main_page_start))
		for _, e := range ToDo {
			fmt.Fprintf(w, e.body)
			if e.status == NotDone {
				fmt.Fprintf(w, "  wait...")
			}
			if e.status == Done {

				fmt.Fprintf(w, "="+string(e.answer)+"   Done!")
			}
			if e.status == Error {
				fmt.Fprintf(w, "  error  "+fmt.Sprint(e.err))
			}
			fmt.Fprintf(w, "<br>")
		}
		fmt.Fprintf(w, string(main_page_end))

	})
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ToDo = append(ToDo, ToDoItem{string(r.FormValue("text")), 0, 0, nil, 0})
		}
		fmt.Fprintf(w, "okay")

	})

	http.ListenAndServe(":8080", nil)
}

func Calculate() {

	for {
		for i, _ := range ToDo {
			if ToDo[i].inprogress == 0 {
				ToDo[i].inprogress = 1
				go func() {
					r := <-StartCalculating(ToDo)
					if r.err != nil {
						ToDo[i].status = Error
						ToDo[i].err = r.err
					} else {
						ToDo[i].status = Done
						ToDo[i].answer = r.answer
					}

				}()
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
	go Calculate()
	backend()
}
