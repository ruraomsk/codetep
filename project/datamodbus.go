package project

import (
	"encoding/json"
	"strconv"
)

//TableXML верхушка таблицы
type TableXML struct {
	Format FormatXML `xml:"format" json:"format"`
	Recods Records   `xml:"records" json:"records"`
}

//ToString вывод содержимого в строку
func (t *TableXML) ToString() string {
	result := t.Format.ToString()
	result += t.Recods.ToString()
	return result + "\n"
}

//Records описание раздела Records
type Records struct {
	RecordListXML []RecordXML `xml:"record" json:"record"`
}

//ToString вывод содержимого в строку
func (r *Records) ToString() string {
	result := ""
	for _, rd := range r.RecordListXML {
		result += rd.ToString() + "\n"

	}
	return result + "\n"
}

// FormatXML сбор форматов
type FormatXML struct {
	MaxRecors int       `xml:"maxRecords,attr" json:"maxrecord"`
	FieldsXML FieldsXML `xml:"fields" json:"fields"`
}

//ToString вывод содержимого в строку
func (f *FormatXML) ToString() string {
	result := strconv.Itoa(f.MaxRecors) + "\n"
	result += f.FieldsXML.ToString()
	return result + "\n"
}

// FieldsXML описывает массив полей регистов
type FieldsXML struct {
	FieldListXML []FieldXML `xml:"field" json:"field"`
	Ret          int        `xml:"ret,attr"`
}

//ToString вывод содержимого в строку
func (f *FieldsXML) ToString() string {
	result := strconv.Itoa(f.Ret) + "\n"
	for _, fd := range f.FieldListXML {
		result += fd.ToString()
	}
	return result
}

// FieldXML описывает отдельное поле
type FieldXML struct {
	Name           string `xml:"name,attr" json:"name"`
	Description    string `xml:"description,attr" json:"description"`
	Type           string `xml:"type,attr" json:"type"`
	Validators     string `xml:"validators,omitempty"`
	SelectionValue string `xml:"selectionValue,omitempty"`
	DefaultValue   string `xml:"defaultValue,omitempty" json:"defaultValue"`
}

//ToString вывод содержимого в строку
func (f *FieldXML) ToString() string {
	return f.Name + "\t\t" + f.Description + "\t\t\t" + f.Type + "\t" + f.DefaultValue + "\n"
}

// RecordXML описывает отдельное поле
type RecordXML struct {
	ValueListXML []ValueXML `xml:"value" jsons:"value"`
}

//ToString вывод содержимого в строку
func (r *RecordXML) ToString() string {
	result := ""
	for _, v := range r.ValueListXML {
		result += v.ToString()
	}
	return result + "\n"
}

// ValueXML наполнение записи
type ValueXML struct {
	Name  string `xml:"name,attr" json:"name"`
	Value string `xml:",chardata" json:"value"`
}

//ToString вывод содержимого в строку
func (v *ValueXML) ToString() string {
	return v.Name + "\t" + v.Value + "\t"
}

//Registers место хранения всех регистров данного МодБаса
type Registers struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MapRegs     map[string]Register
}

//ToString вывод содержимого в строку
func (r *Registers) ToString() string {
	result := ""
	for _, reg := range r.MapRegs {
		result += reg.ToString()
	}
	return result
}

//ToJSON вывод в JSON
func (r Registers) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
