package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func HandleGameConfig() {
	file, err := ioutil.ReadFile("config/game.json")
	if err != nil {
		log.Fatal(err)
	}

	data := GameConfig{}
	
	_ = json.Unmarshal([]byte(file), &data)

	window_height = data.WindowHeight
	window_width  = data.WindowWidth
}