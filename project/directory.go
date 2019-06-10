package project

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

//LoadAllDrivers загружает все драйвера
func LoadAllDrivers(path string) (Drivers, error) {
	drvs := new(Drivers)
	drvs.Drivers = make(map[string]DriverXML)
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		err = fmt.Errorf("Error! Ошибка чтения каталога %s %s ", path, err.Error())
		return *drvs, err
	}
	path += "/"
	for _, file := range dirs {
		if file.IsDir() {
			continue
		}
		npath := path + file.Name()
		drv, err := LoadDriverTable(npath)
		if err != nil {
			err = fmt.Errorf("Error Ошибка загрузки драйвера  %s %s", npath, err.Error())
			return *drvs, err
		}
		drvs.Drivers[drv.HeadDriver.Name] = *drv
	}

	return *drvs, nil
}

//LoadDriverTable загрузка таблицы описания модбаса
func LoadDriverTable(namefile string) (*DriverXML, error) {
	t := new(DriverXML)
	buf, err := ioutil.ReadFile(namefile)
	if err != nil {
		err = fmt.Errorf("Error! " + err.Error())
		return nil, err
	}
	err = xml.Unmarshal(buf, &t)
	// fmt.Println(t.ToString())
	return t, err

}
