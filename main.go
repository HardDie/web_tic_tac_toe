package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type mytest struct {
	Array [][]string
}

type request struct {
	Type string
	Name string
	line int
	row  int
}

type response struct {
	Line  int
	Row   int
	Value string
}

func parseRequest(r io.Reader) (*request, error) {
	req := request{}
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&req)
	if err != nil {
		return nil, err
	}
	fmt.Sscanf(req.Name, "%d_%d", &req.line, &req.row)
	return &req, nil
}

func main() {
	const width = 3
	const height = 3
	value := mytest{}

	for line := 0; line < height; line++ {
		tmp := make([]string, 0)
		for row := 0; row < width; row++ {
			tmp = append(tmp, " ")
		}
		value.Array = append(value.Array, tmp)
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
			log.Fatal(err)
		}

		resp := []response{}

		switch data.Type {
		case "Set":
			if value.Array[data.line][data.row] == " " {
				value.Array[data.line][data.row] = currentPlayer
				switch currentPlayer {
				case "X":
					currentPlayer = "O"
				case "O":
					currentPlayer = "X"
				}
			}

			resp = append(resp, response{data.line, data.row, currentPlayer})

		case "Reset":
			for line, line_s := range value.Array {
				for row, _ := range line_s {
					value.Array[line][row] = " "
					resp = append(resp, response{line, row, " "})
				}
			}
			fmt.Println("Reset type!")
		}

		w.Header().Set("Content-Type", "application/json")
		responseJson, err := json.Marshal(resp)
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
