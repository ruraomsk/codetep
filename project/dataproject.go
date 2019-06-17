package project

import (
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"
)

//Project описание одного проекта системы
type Project struct {
	General    General `xml:"general" json:"general"`
	Subs       []Sub   `xml:"subs" json:"subs"`
	Path       string
	Subsystems map[string]*Subsystem
	DefDrivers Drivers
	Models     map[string]ModelXML
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
		result += "Saves:\n"
		for _, s := range subt.MapSaves {
			result += s.ToString() + "\n"
		}
		result += "Devices:\n"
		for _, rd := range subt.RealDevices {
			result += rd.ToString() + "\n"
		}
	}

	return result
}

//Subsystem описание подсистемы
type Subsystem struct {
	Name         string
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
	MapSaves     map[string]Save
	RealDevices  map[string]Device
	SizeBuffer   int
	NameSaveFile string
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
	XML string `xml:"xml,attr" json:"xml"`
}

//Modbus описание интерфейса мModBus подситемы
type Modbus struct {
	Name        string `xml:"name,attr" json:"name"`
	Description string `xml:"description,attr" json:"desription"`
	Type        string `xml:"type,attr" json:"type"`
	Port        string `xml:"port,attr,omitempty" json:"port"`
	// Slave       string `xml:"slave,attr,omitempty" json:"slave"`
	Step      string `xml:"step,attr,omitempty" json:"step"`
	IP1       string `xml:"ip1,attr,omitempty"`
	IP2       string `xml:"ip2,attr,omitempty"`
	XMLModbus string `xml:"xml,attr,omitempty" json:"xml"`
	Registers map[string]Register
}

//ToString возвращает в символьном виде
func (m *Modbus) ToString() string {
	result := "\t" + m.Name + "\t:" + m.Description + "\t\t\t" + m.Type + "\t" + m.Port + "\t" + m.Step + "\t" + m.XMLModbus + "\n"
	for _, reg := range m.Registers {
		result += reg.ToString()
	}
	return result
}

//IsMaster return true if modbus is master mode
func (m *Modbus) IsMaster() bool {
	if strings.ToLower(m.Type) == "master" {
		return true
	}
	return false
}

//DevicesXML struct
type DevicesXML struct {
	DevicesHead xml.Name `xml:"devices"`
	XML         string   `xml:"xml,attr"`
	Devices     []Device `xml:"device"`
}

//Device Перечень драйверов
type Device struct {
	Name        string `xml:"name,attr" json:"name"`
	Description string `xml:"description,attr" json:"desription"`
	Driver      string `xml:"driver,attr" json:"driver"`
	Slot        string `xml:"slot,attr" json:"slot"`
	Defs        []Def
	Inits       []Init
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

//Init одна строка настройки драйвера
type Init struct {
	Name  string `xml:"name,attr" json:"name"`
	Value string `xml:"value,attr" json:"value"`
}

//Saved сохранениеи переменных на внешний носитель
type Saved struct {
	Sav      xml.Name `xml:"saves"`
	NameFile string   `xml:"name,attr"`
	Saves    []Save   `xml:"save" json:"save"`
}

//ToString возвращает в символьном виде
func (s *Saved) ToString() string {
	result := "Saved " + "\n"
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
	ID          int
	Address     int
}

//FullSize возвращает размер в байтах вместе с байтом достоверности
func (v *Variable) FullSize() int {
	size := 0
	format, _ := strconv.Atoi(v.Format)
	if format == 1 {
		size = 2
	} else if format <= 4 {
		size = 3
	} else if format <= 9 {
		size = 5
	} else if format <= 15 {
		size = 9
	} else if format == 18 {
		size = 2
	}
	res, _ := strconv.Atoi(v.Size)
	return size * res
}

//ToString возвращает в символьном виде
func (v *Variable) ToString() string {
	return "\t" + v.Name + "\t:" + v.Description + "\t" + v.Format + "\t" + v.Size + "\n"
}
