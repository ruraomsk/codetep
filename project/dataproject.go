package project

import (
	"encoding/json"
)

//Project описание одного проекта системы
type Project struct {
	General    General `xml:"general" json:"general"`
	Subs       []Sub   `xml:"subs" json:"subs"`
	Path       string
	Subsystems map[string]*Subsystem
}

//ToJSON вывод в JSOM
func (p *Project) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

//General описание параматров всего проекта
type General struct {
	Name        string `xml:"name,attr" json:"name"`
	Description string `xml:"description,attr" json:"desription"`
	DefDrv      string `xml:"defdrv,attr" json:"defdrv"`
	Simul       string `xml:"simul,attr" json:"simul"`
	IP          string `xml:"ip,attr" json:"ip"`
	Port        string `xml:"port,attr" json:"port"`
}

//Sub описание одной подсистемы
type Sub struct {
	Name        string `xml:"name,attr" json:"name"`
	Path        string `xml:"path,attr" json:"path"`
	Description string `xml:"description,attr" json:"desription"`
	File        string `xml:"file,attr" json:"file"`
	ID          string `xml:"id,attr" json:"id"`
	Step        string `xml:"step,attr" json:"step"`
	Main        string `xml:"main,attr" json:"main"`
	Second      string `xml:"second,attr" json:"second"`
}

//ToString вывод в строку
func (p *Project) ToString() string {
	result := p.General.Name + ":" + p.General.Description + " =" + "\n"
	for _, sub := range p.Subs {
		result += sub.Name + "\t" + sub.Description + "\t" + sub.Path + "\t" + sub.File + "\t" + sub.Path + "\t" + sub.Step + "\t" + sub.Path + "\n"
		subt := p.Subsystems[sub.Name]
		result += subt.ToString()
		result += "Variables:\n"
		for _, v := range subt.Variables {
			result += v.ToString() + "\n"
		}
	}

	return result
}

//Subsystem описание подсистемы
type Subsystem struct {
	Model        Model        `xml:"model" json:"model"`
	Netblkey     Netblkey     `xml:"netblkey" json:"netblkey"`
	Result       Result       `xml:"result" json:"result"`
	Devices      Devices      `xml:"devices" json:"devices"`
	Saves        Saves        `xml:"saves" json:"saves"`
	VariableFile VariableFile `xml:"variable" json:"vars"`
	Key          Key          `xml:"key" json:"key"`
	Initsig      Initsig      `xml:"initsig" json:"initsig"`
	Modbuses     []Modbus     `xml:"modbus" json:"modbus"`
	Delay        Delay        `xml:"delay" json:"delay"`
	Variables    map[string]Variable
}

//ToString подсистему в строку
func (s *Subsystem) ToString() string {
	result := "Model:" + s.Model.Name + " Netblkey:" + s.Netblkey.Name + " result:" + s.Result.Path + " key:" + s.Key.Name + " delay:" + s.Delay.Time + "\n"
	result += "Modbuses:\n"
	for _, mb := range s.Modbuses {
		result += mb.ToString() + "\n"
	}
	return result
}

//ToJSON вывод в формате JSON
func (s *Subsystem) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

//Key Переменная по которой
type Key struct {
	Name string `xml:"name,attr" json:"name"`
}

//VariableFile Имя файла с переменными
type VariableFile struct {
	XML string `xml:"xml,attr" json:"xml"`
}

//Initsig начальные значения сигналов
type Initsig struct {
	XML string `xml:"xml,attr" json:"xml"`
}

//Model описывает модель процессора
type Model struct {
	Name string `xml:"name,attr" json:"name"`
}

//Netblkey опысывает ключ блокировки сетевых обменов
type Netblkey struct {
	Name string `xml:"name,attr" json:"name"`
}

//Result путь куда ложить результат генерации программ
type Result struct {
	Path string `xml:"path,attr" json:"path"`
}

//Delay Задержка начала основоного цикла
type Delay struct {
	Time string `xml:"time,attr" json:"time"`
}

//Saves Путь к перечню сохраненнных переменных
type Saves struct {
	XML string `xml:"xml,attr" json:"xml"`
}

//Devices Путь описание списка всех устройств
type Devices struct {
	XML string `xml:"xml" json:"xml"`
}

//Modbus описание интерфейса мModBus подситемы
type Modbus struct {
	Name        string `xml:"name,attr" json:"name"`
	Description string `xml:"description,attr" json:"desription"`
	Type        string `xml:"type,attr" json:"type"`
	Port        string `xml:"port,attr,omitempty" json:"port"`
	Slave       string `xml:"slave,attr,omitempty" json:"slave"`
	Step        string `xml:"step,attr,omitempty" json:"step"`
	XMLModbus   string `xml:"xml,attr,omitempty" json:"xml"`
}

//ToString возвращает в символьном виде
func (m *Modbus) ToString() string {
	return "\t" + m.Name + "\t:" + m.Description + "\t\t\t" + m.Type + "\t" + m.Port + "\t" + m.Slave + "\t" + m.Step + "\t" + m.XMLModbus
}

//Device Перечень драйверов и закрепление переменных на них
type Device struct {
	Name        string `xml:"name,attr" json:"name"`
	Description string `xml:"description,attr" json:"desription"`
	Driver      string `xml:"driver,attr" json:"driver"`
	Slot        string `xml:"slot,attr" json:"slot"`
	Defs        []Def  `xml:"def" json:"def"`
}

//ToString возвращает в символьном виде
func (d *Device) ToString() string {
	result := "Device " + d.Name + ":" + d.Description + " driver=" + d.Driver + ":" + d.Slot + "\n"
	for _, def := range d.Defs {
		result += def.ToString()
	}
	return result + "\n"
}

//Def одна строка свящи переменной с имем на драйвере
type Def struct {
	Name       string `xml:"name,attr" json:"name"`
	DriverName string `xml:",chardata" json:"drivername"`
}

//ToString возвращает в символьном виде
func (d *Def) ToString() string {
	return "<<\t" + d.Name + "\t\t:\t" + d.DriverName + "\n"
}

//Saved сохранениеи переменных на внешний носитель
type Saved struct {
	NameFile string `xml:"name,attr" json:"namefile"`
	Saves    []Save `xml:"save" json:"save"`
}

//ToString возвращает в символьном виде
func (s *Saved) ToString() string {
	result := "Saved " + s.NameFile + "\n"
	for _, sav := range s.Saves {
		result += sav.ToString()
	}
	return result + "\n"
}

// Save описание одной переменной хранения
type Save struct {
	Value string `xml:"value,attr" json:"value"`
	Name  string `xml:",chardata" json:"name"`
}

//ToString возвращает в символьном виде
func (s *Save) ToString() string {
	return ">>\t" + s.Name + "\t\t:\t" + s.Value + "\n"
}

//Variables описание всех переменных модели
type Variables struct {
	ListVariable []Variable `xml:"var" json:"var"`
}

//ToString возвращает в символьном виде
func (v *Variables) ToString() string {
	result := "Variables \n"
	for _, vv := range v.ListVariable {
		result += vv.ToString()
	}
	return result + "\n"
}

//Variable собственно описание переменной
type Variable struct {
	Name        string `xml:"name,attr" json:"name"`
	Description string `xml:"description,attr" json:"desription"`
	Format      string `xml:"format,attr" json:"format"`
	Size        string `xml:"size,attr,omitempty" json:"size"`
}

//ToString возвращает в символьном виде
func (v *Variable) ToString() string {
	return "\t" + v.Name + "\t:" + v.Description + "\t" + v.Format + "\t" + v.Size + "\n"
}