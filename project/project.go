package project

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
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
			namefile := p.Path + "/" + sub.Path + "/" + sub.File + ".xml"
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
			namefile = p.Path + "/" + sub.Path + "/" + subb.VariableFile.XML + ".xml"
			buf, err = ioutil.ReadFile(namefile)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			var vars = new(Variables)
			subb.Variables = make(map[string]Variable)
			err = xml.Unmarshal(buf, &vars)
			for _, v := range vars.ListVariable {
				subb.Variables[v.Name] = v
			}
			return subb, err
		}
	}
	return nil, fmt.Errorf("Error! Нет такой подсистемы %s", name)
}
