package main

import (
	"fmt"
	"slices"
	"strings"
)

type ItemPlace struct {
	Place string
	Items []string
}

type Location struct {
	Name        string
	Description string
	ItemPlaces  []ItemPlace
	Exits       []string
	Accessible  bool
	Tasks       []string
}

type Player struct {
	Inventory []string
	Wearing   []string
}

type Game struct {
	Locations  map[string]*Location
	Player     *Player
	CurrentLoc string
}

var game *Game

func initGame() {
	game = &Game{
		Locations:  make(map[string]*Location),
		Player:     &Player{Inventory: []string{}, Wearing: []string{}},
		CurrentLoc: "кухня",
	}

	game.Locations["кухня"] = &Location{
		Name:        "кухня",
		Description: "ты находишься на кухне",
		ItemPlaces:  []ItemPlace{{Place: "на столе", Items: []string{"чай"}}},
		Exits:       []string{"коридор"},
		Accessible:  true,
		Tasks:       []string{"собрать рюкзак", "идти в универ"},
	}

	game.Locations["коридор"] = &Location{
		Name:        "коридор",
		Description: "ничего интересного",
		ItemPlaces:  []ItemPlace{},
		Exits:       []string{"кухня", "комната", "улица"},
		Accessible:  true,
	}

	game.Locations["комната"] = &Location{
		Name:        "комната",
		Description: "ты в своей комнате",
		ItemPlaces: []ItemPlace{
			{Place: "на столе", Items: []string{"ключи", "конспекты"}},
			{Place: "на стуле", Items: []string{"рюкзак"}},
		},
		Exits:      []string{"коридор"},
		Accessible: true,
	}

	game.Locations["улица"] = &Location{
		Name:        "улица",
		Description: "на улице весна",
		ItemPlaces:  []ItemPlace{},
		Exits:       []string{"домой"},
		Accessible:  false,
	}
}

func joinWithSymbol(items []string, symb string) string {
	if len(items) == 1 {
		return items[0]
	}
	return fmt.Sprintf("%s %s %s", strings.Join(items[:len(items)-1], ", "), symb, items[len(items)-1])
}

func look() string {
	loc := game.Locations[game.CurrentLoc]
	var parts []string
	if loc.Name != "комната" && loc.Description != "" {
		parts = append(parts, loc.Description)
	}
	var items []string
	for _, ip := range loc.ItemPlaces {
		if len(ip.Items) > 0 {
			items = append(items, fmt.Sprintf("%s: %s", ip.Place, strings.Join(ip.Items, ", ")))
		}
	}
	if len(items) > 0 {
		parts = append(parts, strings.Join(items, ", "))
	} else if len(loc.ItemPlaces) > 0 {
		parts = append(parts, "пустая комната")
	}
	if len(loc.Tasks) > 0 {
		parts = append(parts, "надо "+joinWithSymbol(loc.Tasks, "и"))
	}
	return strings.Join(parts, ", ") + ". можно пройти - " + strings.Join(loc.Exits, ", ")
}

func goTo(place string) string {
	loc := game.Locations[game.CurrentLoc]
	if !slices.Contains(loc.Exits, place) {
		return "нет пути в " + place
	}
	target := game.Locations[place]
	if !target.Accessible {
		return "дверь закрыта"
	}
	game.CurrentLoc = place
	var entry string
	if place == "кухня" {
		entry = "кухня, ничего интересного"
	} else {
		entry = target.Description
	}
	return entry + ". можно пройти - " + strings.Join(target.Exits, ", ")
}

func updateTasks() {
	for _, loc := range game.Locations {
		var newTasks []string
		for _, task := range loc.Tasks {
			if strings.HasPrefix(task, "собрать ") {
				reqs := strings.Split(strings.TrimPrefix(task, "собрать "), " и ")
				if !requirementsMet(reqs) {
					newTasks = append(newTasks, task)
				}
			} else {
				newTasks = append(newTasks, task)
			}
		}
		loc.Tasks = newTasks
	}
}

func requirementsMet(reqs []string) bool {
	for _, req := range reqs {
		if !slices.Contains(game.Player.Inventory, req) && !slices.Contains(game.Player.Wearing, req) {
			return false
		}
	}
	return true
}

func wear(item string) string {
	loc := game.Locations[game.CurrentLoc]
	for i, ip := range loc.ItemPlaces {
		if idx := slices.Index(ip.Items, item); idx != -1 {
			loc.ItemPlaces[i].Items = slices.Delete(ip.Items, idx, idx+1)
			game.Player.Wearing = append(game.Player.Wearing, item)
			updateTasks()
			return fmt.Sprintf("вы надели: %s", item)
		}
	}
	return "нет такого"
}

func take(item string) string {
	if !slices.Contains(game.Player.Wearing, "рюкзак") {
		return "некуда класть"
	}
	loc := game.Locations[game.CurrentLoc]
	for i, ip := range loc.ItemPlaces {
		if idx := slices.Index(ip.Items, item); idx != -1 {
			loc.ItemPlaces[i].Items = slices.Delete(ip.Items, idx, idx+1)
			game.Player.Inventory = append(game.Player.Inventory, item)
			updateTasks()
			return fmt.Sprintf("предмет добавлен в инвентарь: %s", item)
		}
	}
	return "нет такого"
}

func apply(item, target string) string {
	if !slices.Contains(game.Player.Inventory, item) {
		return fmt.Sprintf("нет предмета в инвентаре - %s", item)
	}
	if item == "ключи" && target == "дверь" {
		loc := game.Locations["улица"]
		loc.Accessible = true
		return "дверь открыта"
	}
	return "не к чему применить"
}

func handleCommand(input string) string {
	parts := strings.Split(input, " ")
	switch parts[0] {
	case "осмотреться":
		return look()
	case "идти":
		return goTo(parts[1])
	case "надеть":
		return wear(parts[1])
	case "взять":
		return take(parts[1])
	case "применить":
		return apply(parts[1], parts[2])
	default:
		return "неизвестная команда"
	}
}

func main() {
	initGame()
}
