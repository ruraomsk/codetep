package project

import (
	"encoding/json"
	"encoding/xml"
	"strconv"
)

//Drivers хранение всех драйверов
type Drivers struct {
	Drivers map[string]DriverXML
}

//ToString вывод в символьном виде
func (d *Drivers) ToString() string {
	result := ""
	for name, drv := range d.Drivers {
		result += name + "=" + drv.ToString() + "\n"
	}
	return result
}

//DriverXML Заголовок описания драйвера
type DriverXML struct {
	HeadDriver  xml.Name `xml:"driver" json:"driver"`
	Name        string   `xml:"name,attr" json:"name"`
	Description string   `xml:"description,attr" json:"description"`
	Code        string   `xml:"code,attr" json:"code"`
	LenData     string   `xml:"lenData,attr" json:"lenData"`
	LenInit     string   `xml:"lenInit,attr" json:"lenInit"`
	Header      string   `xml:"header,attr" json:"header"`
	Signals     Signals  `xml:"signals" json:"signals"`
	Inits       Inits    `xml:"init" json:"init"`
	MapSignals  map[string]Signal
}

//ToJSON вывод в формате Json
func (d *DriverXML) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

//ToString вывод в символьном виде
func (d *DriverXML) ToString() string {
	result := d.Name + ":" + d.Description + "\t" + d.Code + "\t" + d.LenData + "\t" + d.LenInit + "\t" + d.Header
	result += "\n" + d.Signals.ToString() + "\n" + d.Inits.ToString()
	return result
}

//Signals описание входов и выходов драйвера буфер обмена
type Signals struct {
	Signals []Signal `xml:"signal" json:"signal"`
}

//ToString вывод в символьном виде
func (s *Signals) ToString() string {
	result := ""
	for _, sig := range s.Signals {
		result += sig.ToString() + "\n"
	}
	return result
}

//Signal описание одного входы/выхода
type Signal struct {
	Name    string `xml:"name,attr" json:"name"`
	Format  string `xml:"format,attr" json:"format"`
	Mode    string `xml:"mode,attr" json:"mode"`
	Address string `xml:"address,attr" json:"address"`
	Value   string `xml:"value,attr,omitempty" json:"value"`
}

//ToString вывод в символьном виде
func (s *Signal) ToString() string {
	return s.Name + ":" + s.Format + "\t" + s.Mode + "\t" + s.Address + "\t" + s.Value
}

//Inits описание буфера инициализации драйвера
type Inits struct {
	Type  string   `xml:"type,attr" json:"type"`
	Inits []Signal `xml:"signal" json:"signal"`
}

//ToString вывод в символьном виде
func (i *Inits) ToString() string {
	result := "type=" + i.Type + "\n"
	for _, sig := range i.Inits {
		result += sig.ToString() + "\n"
	}
	return result
}

//MakeDriverTable build code for main hrader for device
func (dev *Device) MakeDriverTable(div Drivers) string {
	d, _ := div.Drivers[dev.Driver]
	res := make([]string, 0)
	l, _ := strconv.Atoi(d.LenInit)
	for i := 0; i < l; i++ {
		res = append(res, "0")
	}
	for _, i := range d.Inits.Inits {
		a, _ := strconv.Atoi(i.Address)
		res[a] = i.Value
		for _, ii := range dev.Inits {
			if ii.Name == i.Name {
				res[a] = ii.Value
				break
			}
		}
	}
	rez := ""
	rez = "#include <" + d.Header + ">\n"
	rez += "static char buf_" + dev.Name + "[" + d.LenData + "];\t//" + dev.Name + "\n"
	rez += "static " + d.Inits.Type + " ini_" + dev.Name + "={"
	for _, r := range res {
		rez += r + ","
	}
	rez += "};\n"
	rez += "#pragma pack(push,1)\n"
	rez += "static table_drv table_" + dev.Name + "={0,0,&ini_" + dev.Name + ",buf_" + dev.Name + ",0,0};\n"
	rez += "#pragma pop\n"
	rez += "#pragma pack(push,1)\n"
	rez += "static DriverRegister def_buf_" + dev.Name + "[]={\n"
	for _, def := range dev.Defs {
		s := d.MapSignals[def.DriverName]
		rez += "\t{&" + def.Name + "," + s.Format + "," + s.Address + "},\n"
	}
	rez += "\t{NULL,0,0},\n};\n"
	rez += "#pragma pop\n"
	// rez += "static char temp_" + dev.Name + "[256];\n"
	return rez
}
