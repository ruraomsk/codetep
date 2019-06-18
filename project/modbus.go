package project

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/clbanning/mxj"
)

//Register Собственно регистр МодБас
type Register struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	Format      int    `json:"format"`
	Address     int    `json:"address"`
	Size        int    `json:"size"`
	UID         int    `json:"uId"`
}

//ToString возвращает в символьном виде
func (r *Register) ToString() string {
	return fmt.Sprintf("%s\t%s\t%d\t%d\t%d\t%d\t%d \n", r.Name, r.Description, r.Type, r.Format, r.Address, r.Size, r.UID)
}

//LoadModbusTable загрузка таблицы описания модбаса
func LoadModbusTable(namefile string) (*TableXML, error) {
	namefile = RepairPath(namefile)

	t := new(TableXML)
	buf, err := ioutil.ReadFile(namefile)
	if err != nil {
		fmt.Println("Error! " + err.Error())
		return nil, err
	}
	err = xml.Unmarshal(buf, &t)
	// fmt.Println(pr.ToString())
	return t, err

}

//NewTAbleXML создает новую таблицу для описания ModBus
func NewTAbleXML() *TableXML {
	table := new(TableXML)
	name := FieldXML{Name: "name", Description: "Имя", Type: "S"}
	description := FieldXML{Name: "description", Description: "Описание", Type: "S"}
	ftype := FieldXML{Name: "type", Description: "Тип", Type: "I", DefaultValue: "0"}
	format := FieldXML{Name: "format", Description: "Формат", Type: "I", DefaultValue: "2"}
	address := FieldXML{Name: "address", Description: "Десятичный адрес регистра", Type: "I", DefaultValue: "0"}
	size := FieldXML{Name: "size", Description: "Размер", Type: "I", DefaultValue: "1"}
	UID := FieldXML{Name: "unitId", Description: "Идентификатор устройства", Type: "I", DefaultValue: "1"}
	table.Format.MaxRecors = 2147483647
	table.Format.FieldsXML.FieldListXML = []FieldXML{name, description, ftype, format, address, size, UID}
	table.Recods.RecordListXML = make([]RecordXML, 0)
	return table
}

//UpdateRegisters записывает новые значения регистров в таблицу настройки модбас
func (t *TableXML) UpdateRegisters(regs Registers) {
	t.Recods.RecordListXML = make([]RecordXML, 0)
	for _, reg := range regs.MapRegs {
		rec := new(RecordXML)
		// rec.ValueListXML=make([]ValueXML,0)
		name := ValueXML{Name: "name", Value: reg.Name}
		description := ValueXML{Name: "description", Value: reg.Description}
		ftype := ValueXML{Name: "type", Value: string(reg.Type)}
		format := ValueXML{Name: "format", Value: string(reg.Format)}
		address := ValueXML{Name: "address", Value: string(reg.Address)}
		size := ValueXML{Name: "size", Value: string(reg.Size)}
		UID := ValueXML{Name: "unitId", Value: string(reg.UID)}
		rec.ValueListXML = []ValueXML{name, description, ftype, format, address, size, UID}
		t.Recods.RecordListXML = append(t.Recods.RecordListXML, *rec)
	}
}

//SaveTable сохраняет таблицу в формате XML
func (t *TableXML) SaveTable(path string) error {
	path = RepairPath(path)
	result, err := xml.Marshal(t)
	if err != nil {
		fmt.Println("Error !" + err.Error())
		return err
	}
	result = bytes.Replace(result, []byte("<TableXML"), []byte("<table"), -1)
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

//GetRegisters создает новый пул регистров из таблицы
func (t *TableXML) GetRegisters(name, description string) *Registers {
	regs := new(Registers)
	regs.MapRegs = make(map[string]Register)
	regs.Name = name
	regs.Description = description
	for _, rec := range t.Recods.RecordListXML {
		reg := new(Register)
		reg.Size = 1
		reg.Format = 2
		reg.Type = 0
		reg.Address = 0
		for _, val := range rec.ValueListXML {
			switch val.Name {
			case "name":
				reg.Name = val.Value
				break
			case "description":
				reg.Description = val.Value
				break
			case "type":
				reg.Type, _ = strconv.Atoi(val.Value)
				break
			case "format":
				reg.Format, _ = strconv.Atoi(val.Value)
				break
			case "address":
				reg.Address, _ = strconv.Atoi(val.Value)
				break
			case "size":
				reg.Size, _ = strconv.Atoi(val.Value)
				break
			case "unitId":
				reg.UID, _ = strconv.Atoi(val.Value)
				break
			}
		}
		regs.MapRegs[reg.Name] = *reg
	}
	return regs
}
