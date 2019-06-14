package main

import (
	"fmt"
	"rura/codetep/project"
)

func main() {

	fmt.Println("Начало работы")
	prPath := "/home/rura/dataSimul/pr"
	pr, err := project.LoadProject(prPath)
	if err != nil {
		fmt.Println("Найдены ошибки " + err.Error())
		return
	}
	defDrivers, err := project.LoadAllDrivers(prPath + "/settings/default")
	if err != nil {
		fmt.Println("Найдены ошибки " + err.Error())
		return
	}
	defModels, err := project.LoadAllModels(prPath + "/settings/models")
	// fmt.Println(pr.ToString())
	fmt.Println(defDrivers.ToString())
	for _, model := range defModels {
		fmt.Println(model.ToString())
	}
	// TODO: Глобальная прооверка на правильность данных
	fmt.Println(pr.VerifyAllDevices(defDrivers))

	fmt.Println(pr.VerifyAllVariables(defDrivers))
	// TODO: Написать кодогенератор
	// TODO: Open
	fmt.Println("Конец работы")

}
