package main

import (
	"fmt"
	"rura/codetep/project"
)

func main() {

	fmt.Println("Начало работы")
	prPath := "/home/rura/dataSimul/prnew"
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
	// fmt.Println(pr.ToString())
	// fmt.Println(defDrivers.ToString())
	for _, model := range pr.Models {
		fmt.Println(model.ToString())
	}
	// TODO: Глобальная прооверка на правильность данных
	// fmt.Println(pr.VerifyAllDevices())
	// fmt.Println(pr.VerifyAllVariables())
	err = pr.MakeMaster(prPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Конец работы")
}
