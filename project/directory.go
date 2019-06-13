package project

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/R358/xmldoc"
	"github.com/clbanning/mxj"
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
		drvs.Drivers[drv.Name] = *drv
		// drv.SaveXML(path + "_" + file.Name())
	}

	return *drvs, nil
}

//LoadDriverTable загрузка таблицы описания драйвера
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

//LoadAllModels load all models fro path dir
func LoadAllModels(path string) (map[string]ModelXML, error) {
	Result := make(map[string]ModelXML)
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		err = fmt.Errorf("Error! Ошибка чтения каталога %s %s ", path, err.Error())
		return Result, err
	}
	path += "/"
	for _, file := range dirs {
		if file.IsDir() {
			continue
		}
		npath := path + file.Name()
		model, err := LoadModel(npath)
		if err != nil {
			err = fmt.Errorf("Error Ошибка загрузки модули %s %s", npath, err.Error())
			return Result, err
		}
		model.Name = strings.Replace(file.Name(), ".xml", "", 1)
		Result[model.Name] = *model
	}

	return Result, nil

}

//LoadModel load from XML one model
func LoadModel(namefile string) (*ModelXML, error) {
	t := new(ModelXML)
	buf, err := ioutil.ReadFile(namefile)
	if err != nil {
		err = fmt.Errorf("Error! " + err.Error())
		return nil, err
	}
	err = xml.Unmarshal(buf, &t)
	// fmt.Println(t.ToString())

	return t, err

}

//SaveXML сохраняет в XML
func (t *DriverXML) SaveXML(path string) error {
	result, err := xml.Marshal(t)
	if err != nil {
		fmt.Println("Error !" + err.Error())
		return err
	}
	result, err = mxj.BeautifyXml(result, "", "\t")
	if err != nil {
		fmt.Println("Error !" + err.Error())
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error !" + err.Error())
		return err
	}
	defer file.Close()
	_, err = file.Write(result)
	return err

}
func textValue(node xmldoc.XDNode) string {
	if ele, ok := node.(*xmldoc.XDElement); ok {
		b := &bytes.Buffer{}
		ele.TraverseChildren(func(child xmldoc.XDNode) (stop bool) {
			if child.GetType() == xmldoc.CDataType {
				b.Write(child.(*xmldoc.XDCData).Data)
			}
			return stop
		})
		return b.String()
	}
	return ""
}

//LoadAssign загружает для устройства таблицу назначений
func (s *Subsystem) LoadAssign(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	doc, err := xmldoc.Parse(file)
	for _, child := range doc.Root.Children {
		if child.GetName().String() == "" {
			continue
		}
		if child.GetType() == xmldoc.ElementType {
			dev := child.GetName().String()
			rd, ok := s.RealDevices[dev]
			// rd.Defs = make([]Def, 100)
			if !ok {
				return fmt.Errorf("Error нет устройсва " + dev + " в подсистеме " + s.Name)
			}
			if newChild, ok := child.(*xmldoc.XDElement); ok {
				for _, ass := range newChild.Children {
					if ass.GetName().String() == "" {
						continue
					}
					assign, _ := ass.(*xmldoc.XDElement)
					name := xmldoc.XDName{LocalName: "name"}
					def := Def{Name: assign.Attributes[name], DriverName: textValue(ass)}
					rd.Defs = append(rd.Defs, def)
				}

			}
			s.RealDevices[dev] = rd
		}
	}
	return nil

}
