package main

import (
	"fmt"
	"os"
	"rura/codetep/project"
)

func main() {

	fmt.Println("Начало работы...")
	prPath := ""
	if len(os.Args) == 1 {
		prPath = "/home/rura/dataSimul/prnew"
	} else {
		prPath = os.Args[1]
	}
	pr, err := project.LoadProject(prPath)
	if err != nil {
		fmt.Println("Найдены ошибки " + err.Error())
		return
	}
	pr.DefDrivers, err = project.LoadAllDrivers(prPath + "/settings/default")
	if err != nil {
		fmt.Println("Найдены ошибки " + err.Error())
		return
	}
	pr.Models, err = project.LoadAllModels(prPath + "/settings/models")
	err = pr.MakeMaster(prPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Конец работы")
}
