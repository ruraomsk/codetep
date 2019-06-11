package project

//ModelXML заголовок
type ModelXML struct {
	ModifC ModifC `xml:"modif_C"`
	ModifH ModifH `xml:"modif_H"`
	Name   string
}

//ModifC section for main c
type ModifC struct {
	Attachs []Attach `xml:"attach"`
}

//ModifH section for include file
type ModifH struct {
	Actions []Action `xml:"act"`
}

//Attach one operator for main c file
type Attach struct {
	ID   string `xml:"id,attr"`
	File string `xml:"file,attr"`
}

//Action jne jperator for header file
type Action struct {
	Name string `xml:"name,attr"`
}

//ToString return string
func (m *ModelXML) ToString() string {
	res := "Model:" + m.Name + "\nModificator for C\n"
	for _, attach := range m.ModifC.Attachs {
		res += attach.ToString() + "\n"
	}
	res += "Modificator for H\n"
	for _, action := range m.ModifH.Actions {
		res += action.ToString() + "\n"
	}
	return res
}

//ToString return string
func (a *Attach) ToString() string {
	return "id:" + a.ID + " file=" + a.File
}

//ToString return string
func (a *Action) ToString() string {
	return "name=" + a.Name
}
