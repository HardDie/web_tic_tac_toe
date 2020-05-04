package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func dumpResponse(resp *http.Response) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response packet:")
	fmt.Println(string(dump))
	fmt.Println("END")
}

type mytest struct {
	Array [][]string
}

type request struct {
	Name string
	line int
	row  int
}

func parseRequest(r io.Reader) (*request, error) {
	req := request{}
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&req)
	if err != nil {
		return nil, err
	}
	fmt.Sscanf(req.Name, "%d_%d", &req.line, &req.row)
	fmt.Println(req.line, req.row)
	return &req, nil
}

func main() {
	value := mytest{
		Array: [][]string{
			[]string{" ", " ", " "},
			[]string{" ", " ", " "},
			[]string{" ", " ", " "},
		},
	}

	currentPlayer := "X"

	tmplPageIndex := template.New("index")
	tmplPageIndex, err := tmplPageIndex.Parse(pageIndexTmpl)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(pageStyle))
	})
	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(pageScript))
	})

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Get request")
		defer r.Body.Close()

		data, err := parseRequest(r.Body)
		if err != nil {
			panic(err)
		}

		if value.Array[data.line][data.row] == " " {
			value.Array[data.line][data.row] = currentPlayer
			switch currentPlayer {
			case "X":
				currentPlayer = "O"
			case "O":
				currentPlayer = "X"
			}
		}

		w.Header().Set("Content-Type", "application/json")
		response := []interface{} {
			map[string]interface{} {
				"Line": data.line,
				"Row": data.row,
				"Value": value.Array[data.line][data.row],
			},
		}
		responseJson, err := json.Marshal(response)
		fmt.Println(string(responseJson))
		w.Write([]byte(responseJson))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmplPageIndex.Execute(w, value)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
