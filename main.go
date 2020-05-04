package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"web_test/game"
)

type request struct {
	Type string
	Name string
	line int
	row  int
}

type responseBody struct {
	Status string
	Step   struct {
		Line   int
		Row    int
		Player string
	}
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
	gg := game.New()
	currentPlayer := game.PlayerX

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

		respBody := responseBody{}

		switch data.Type {
		case "Set":
			pos := data.line*game.FieldWidth + data.row
			if flag, _ := gg.CheckWin(); flag {
				respBody.Status = "Game done"
				fmt.Println("Game done!")
				break
			}
			if err := gg.MakeStep(currentPlayer, pos); err != nil {
				respBody.Status = "Field busy"
				fmt.Println("Field busy")
				break
			}

			playerStr := ""
			switch currentPlayer {
			case game.PlayerX:
				currentPlayer = game.PlayerO
				playerStr = "X"
			case game.PlayerO:
				currentPlayer = game.PlayerX
				playerStr = "O"
			}

			respBody.Status = "Step made"
			respBody.Step.Line = data.line
			respBody.Step.Row = data.row
			respBody.Step.Player = playerStr

		case "Reset":
			gg.Draw()
			gg.Reset()
			currentPlayer = game.PlayerX
			fmt.Println("Reset type!")
		}

		w.Header().Set("Content-Type", "application/json")
		responseJson, err := json.Marshal(respBody)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(responseJson))
		w.Write([]byte(responseJson))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmplPageIndex.Execute(w, gg.GameToSlice())
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
