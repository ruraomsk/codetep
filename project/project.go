package project

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
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
			var vars = new(Variables)
			subb.Variables = make(map[string]Variable)
			err = xml.Unmarshal(buf, &vars)
			for _, v := range vars.ListVariable {
				if len(v.Size) == 0 {
					v.Size = "1"
				}
				subb.Variables[v.Name] = v
			}
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

			return subb, err
		}
	}
	return nil, fmt.Errorf("Error! Нет такой подсистемы %s", name)
}
