package main

import (
	"fmt"
	"rura/codetep/project"
)

func main() {

	fmt.Println("Начало работы")
	pr, err := project.LoadProject("/home/rura/dataSimul/pr")
	if err != nil {
		fmt.Println("Найдены ошибки " + err.Error())
		return
	}
	defDrivers, err := project.LoadAllDrivers("/home/rura/dataSimul/pr/settings/default")
	if err != nil {
		fmt.Println("Найдены ошибки " + err.Error())
		return
	}
	fmt.Println(pr.ToString())
	fmt.Println(defDrivers.ToString())
	fmt.Println("Конец работы")
}
