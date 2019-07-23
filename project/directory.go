package project

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/R358/xmldoc"
	"github.com/clbanning/mxj"
	"golang.org/x/text/encoding/charmap"
)

//LoadAllDrivers загружает все драйвера
func LoadAllDrivers(path string) (Drivers, error) {
	path = RepairPath(path)
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
	namefile = RepairPath(namefile)

	t := new(DriverXML)
	buf, err := ioutil.ReadFile(namefile)
	if err != nil {
		err = fmt.Errorf("Error! " + err.Error())
		return nil, err
	}
	err = xml.Unmarshal(buf, &t)
	// fmt.Println(t.ToString())
	t.MapSignals = make(map[string]Signal)
	for _, s := range t.Signals.Signals {
		t.MapSignals[s.Name] = s
	}
	return t, err

}

//LoadAllModels load all models fro path dir
func LoadAllModels(path string) (map[string]ModelXML, error) {
	Result := make(map[string]ModelXML)
	path = RepairPath(path)
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
		npath := RepairPath(path + file.Name())
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

//LoadCodePart load part c code from setting
func (p *Project) LoadCodePart(nameFile string) (string, error) {
	npath := RepairPath(p.Path + "/settings/src-FP/" + nameFile + ".c")
	buf, err := ioutil.ReadFile(npath)
	if err == nil {
		return string(buf), nil
	}
	return "\n!!!*************NOT FILE**********************\n" + npath, err
}

//LoadShema Загружает хедер схемы попутно вынимая все имена временных переменных
func (p *Project) LoadShema(sub Sub) ([]string, map[string]string, error) {
	lstrings := make([]string, 0)
	tvar := make(map[string]string)
	path := RepairPath(p.Path + "/" + sub.Path + "/scheme/Scheme.h")
	file, err := os.Open(path)
	if err != nil {
		return lstrings, tvar, err
	}
	dec := charmap.Windows1251.NewDecoder()
	defer file.Close()
	sReader := bufio.NewScanner(file)
	for sReader.Scan() {
		line, err := dec.String(sReader.Text())
		if err != nil {
			return lstrings, tvar, err
		}
		ls := strings.Split(line, " ")
		if len(ls) != 2 {
			continue
		}
		if strings.Contains(ls[0], "ss") && strings.Contains(ls[1], "va") {
			ls[1] = strings.Replace(ls[1], ";", "", -1)
			var r string
			if ls[0] == "ssbool" {
				r = ".b=0;"
			} else if ls[0] == "ssfloat" {
				r = ".f=0.0;"
			} else if ls[0] == "ssint" {
				r = ".i=0;"
			} else if ls[0] == "sslong" {
				r = ".l=0;"
			} else if ls[0] == "sslong" {
				r = ".c=0;"
			}
			tvar[ls[1]] = r
			continue
		}
		lstrings = append(lstrings, line)

	}

	return lstrings, tvar, nil
}

//LoadModel load from XML one model
func LoadModel(namefile string) (*ModelXML, error) {
	t := new(ModelXML)
	namefile = RepairPath(namefile)
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
	path = RepairPath(path)
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
	path = RepairPath(path)
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
					if ass.GetName().String() == "def" {
						assign, _ := ass.(*xmldoc.XDElement)
						name := xmldoc.XDName{LocalName: "name"}
						def := Def{Name: assign.Attributes[name], DriverName: textValue(ass)}
						rd.Defs = append(rd.Defs, def)
					}
					if ass.GetName().String() == "init" {
						assign, _ := ass.(*xmldoc.XDElement)
						name := xmldoc.XDName{LocalName: "name"}
						value := xmldoc.XDName{LocalName: "value"}
						init := Init{Name: assign.Attributes[name], Value: assign.Attributes[value]}
						rd.Inits = append(rd.Inits, init)

					}
				}

			}
			s.RealDevices[dev] = rd
		}
	}
	return nil

}

//RepairPath правит обратные косые на обычные
func RepairPath(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}
