package project

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

//LoadProject загружает заголовок конкретного проекта
func LoadProject(path string) (Project, error) {
	var pr Project
	pr.Path = path
	namefile := path + "/main.xml"
	buf, err := ioutil.ReadFile(namefile)
	if err != nil {
		fmt.Println("Error! " + err.Error())
		return pr, err
	}
	err = xml.Unmarshal(buf, &pr)
	// fmt.Println(pr.ToString())
	pr.Subsystems = make(map[string]*Subsystem)
	for _, sub := range pr.Subs {
		subt, err := pr.LoadSubsystem(sub.Name)
		if err != nil {
			return pr, err
		}
		pr.Subsystems[sub.Name] = subt
	}
	return pr, err
}

//LoadSubsystem загружает подсистему проекта
func (p *Project) LoadSubsystem(name string) (*Subsystem, error) {
	for _, sub := range p.Subs {
		// fmt.Println(sub.Name)
		if sub.Name == name {
			subb := new(Subsystem)
			subb.Name = name
			namefile := RepairPath(p.Path + "/" + sub.Path + "/" + sub.File + ".xml")
			buf, err := ioutil.ReadFile(namefile)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			err = xml.Unmarshal(buf, &subb)
			if err != nil {
				fmt.Println("Error! " + err.Error())
				return nil, err
			}
			//Load Variables section
			namefile = RepairPath(p.Path + "/" + sub.Path + "/" + subb.VariableFile.XML + ".xml")
			buf, err = ioutil.ReadFile(namefile)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			subb.Variables = make(map[string]Variable)
			err = xml.Unmarshal(buf, &subb.Vars)
			fullSize := 0
			id := 1
			sort.Slice(subb.Vars.ListVariable, func(i, j int) bool { return subb.Vars.ListVariable[i].Name < subb.Vars.ListVariable[j].Name })
			for i := 0; i < len(subb.Vars.ListVariable); i++ {
				v := subb.Vars.ListVariable[i]
				if len(v.Size) == 0 {
					v.Size = "1"
				}
				v.ID = id
				v.Address = fullSize
				fullSize += v.FullSize()
				id++
				subb.Variables[v.Name] = v
				subb.Vars.ListVariable[i] = v

			}
			subb.SizeBuffer = fullSize
			//Load Saves section
			namefile = RepairPath(p.Path + "/" + sub.Path + "/" + subb.Saves.XML + ".xml")
			buf, err = ioutil.ReadFile(namefile)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			var saves = new(Saved)
			subb.MapSaves = make(map[string]Save)
			err = xml.Unmarshal(buf, &saves)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			subb.NameSaveFile = saves.NameFile
			for _, s := range saves.Saves {
				subb.MapSaves[s.Name] = s
			}
			//Load Devices section
			namefile = RepairPath(p.Path + "/" + sub.Path + "/" + subb.Devices.XML + ".xml")
			buf, err = ioutil.ReadFile(namefile)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			subb.RealDevices = make(map[string]Device)
			devXML := new(DevicesXML)
			err = xml.Unmarshal(buf, &devXML)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			for _, dev := range devXML.Devices {
				subb.RealDevices[dev.Name] = dev
			}
			//Load Assign Device section
			namefile = RepairPath(p.Path + "/" + sub.Path + "/" + devXML.XML + ".xml")
			subb.LoadAssign(namefile)
			//Load Modbus section
			for i, m := range subb.Modbuses {
				if strings.Contains(m.XMLModbus, ".xml") {
					namefile = RepairPath(p.Path + "/" + sub.Path + "/" + m.XMLModbus)

				} else {
					namefile = RepairPath(p.Path + "/" + sub.Path + "/" + m.XMLModbus + ".xml")

				}
				table, err := LoadModbusTable(namefile)
				if err != nil {
					return nil, err
				}
				m.Registers = table.GetRegisters(m.Name, m.Description).MapRegs
				subb.Modbuses[i] = m
			}
			if subb.Initsig.XML != "" {
				namefile = RepairPath(p.Path + "/" + sub.Path + "/" + subb.Initsig.XML + ".xml")
				buf, err = ioutil.ReadFile(namefile)
				if err != nil {
					fmt.Println(err.Error())
					return nil, err
				}
				ini := new(IniSignal)
				err = xml.Unmarshal(buf, &ini)
				if err != nil {
					fmt.Println(err.Error())
					return nil, err
				}

			}
			return subb, err
		}
	}
	return nil, fmt.Errorf("Error! Нет такой подсистемы %s", name)
}

//AppendNewVariables добавить переменные из внутреннего
func (s *Subsystem) AppendNewVariables(vars map[string]string) {
	id := 30000
	for name, st := range vars {
		v := new(Variable)
		v.Name = name
		v.Address = s.SizeBuffer
		if strings.Contains(st, ".b=0") {
			v.Format = "1"
		} else if strings.Contains(st, ".i=0") {
			v.Format = "3"
		} else if strings.Contains(st, ".f=0") {
			v.Format = "8"
		} else if strings.Contains(st, ".l=0") {
			v.Format = "11"
		}
		v.Size = "1"
		v.Description = "Внутренняя переменная " + name
		v.ID = id
		id++
		s.SizeBuffer += v.FullSize()
		s.Variables[v.Name] = *v
		s.Vars.ListVariable = append(s.Vars.ListVariable, *v)
	}
}
