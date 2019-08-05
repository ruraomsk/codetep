package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"rura/codetep/project"
)

func pressAnyKey() {
	fmt.Print("\npress Enter key ....")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return
	}
}
func main() {

	fmt.Println("Начало работы...")
	prPath := ""
	if len(os.Args) == 1 {
		if runtime.GOOS == "linux" {
			prPath = "/home/rura/dataSimul/pr"
		} else {
			prPath = "d:/md/pti/pr"

		}
	} else {
		prPath = os.Args[1]
	}
	pr, err := project.LoadProject(prPath)
	if err != nil {
		fmt.Println("Найдены ошибки " + err.Error())
		pressAnyKey()
		return
	}
	pr.DefDrivers, err = project.LoadAllDrivers(prPath + "/settings/default")
	if err != nil {
		fmt.Println("Найдены ошибки " + err.Error())
		pressAnyKey()
		return
	}
	pr.Models, err = project.LoadAllModels(prPath + "/settings/models")
	if err != nil {
		fmt.Println("Найдены ошибки " + err.Error())
		pressAnyKey()
		return
	}
	result := pr.VerifyAllVariables()
	if result != "" {
		fmt.Println("Найдены ошибки проверке имен " + result)
		pressAnyKey()
		return
	}
	result = pr.VerifyAllDevices()
	if result != "" {
		fmt.Println("Найдены ошибки проверке устройств " + result)
		pressAnyKey()
		return
	}

	err = pr.MakeMaster(prPath)
	if err != nil {
		fmt.Println(err.Error())
		pressAnyKey()
		return
	}
	err = pr.MakeMainC(prPath)
	if err != nil {
		fmt.Println(err.Error())
		pressAnyKey()
		return
	}
	fmt.Println("Конец работы")
}
