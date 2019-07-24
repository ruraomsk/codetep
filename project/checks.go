package project

import "strconv"

type subVariable struct {
	SubName  string   //Имя подсистемы где первый объявлена переменная
	Variable Variable //Собственно переменная
}

//VerifyAllVariables check all Variables into project
func (pr *Project) VerifyAllVariables() string {
	result := ""
	// TODO: Вначале проверяем переменные все они с одним именем должны быть одного типа
	vars := make(map[string]subVariable, 10000)
	for _, sub := range pr.Subsystems {
		for name, v := range sub.Variables {
			sv, ok := vars[name]
			if ok {
				//Переменная уже встречалась
				// fmt.Println(sv.SubName, sv.Variable.Name)
				if sv.Variable.Format != v.Format || sv.Variable.Size != v.Size {
					result += "Error! subsytem " + sv.SubName + " var " + sv.Variable.Name + " not equal var " + v.Name + " into " + sub.Name + "\n"
				}
			} else {
				//Переменная попалась первый раз
				vars[v.Name] = subVariable{SubName: sub.Name, Variable: v}
			}
		}
	}
	// TODO: Затем проверяем переменные модбаса
	for _, sub := range pr.Subsystems {
		for _, mod := range sub.Modbuses {
			for name, m := range mod.Registers {
				sv, ok := sub.Variables[name]
				if ok {
					//Переменная есть
					format := "1"
					if m.Type > 1 {
						format = strconv.Itoa(m.Format)
					}
					if len(sv.Size) == 0 {
						sv.Size = "1"
					}
					if sv.Format != format || sv.Size != strconv.Itoa(m.Size) {
						result += "Error! subsytem " + sub.Name + " var " + sv.Name + " not equal modbus var " + m.Name + " into " + mod.Name + "\n"
					}
				} else {
					//Переменная попалась первый раз
					result += "Error! modbus var " + m.Name + " into " + mod.Name + " haven't variable\n"
				}
			}
		}
	}
	// TODO: Потом проверяем все назначения на устройства со стороны переменных
	for _, sub := range pr.Subsystems {
		for _, dev := range sub.RealDevices {
			drv, ok := pr.DefDrivers.Drivers[dev.Driver]
			if !ok {
				continue
			}
			for _, d := range dev.Defs {
				v, ok := sub.Variables[d.Name]
				if !ok {
					//Переменная нет
					result += "Error! Subsystem " + sub.Name + " def var " + d.Name + " into " + dev.Name + " haven't variable\n"
					continue
				}
				// Теперь проверим совпадение форматов переменной и пина драйвера
				if drv.ReadPINFormat(d.DriverName) != v.Format {
					result += "Error! Subsystem " + sub.Name + " def  var " + d.Name + " into " + dev.Name + " wrong format\n"
					continue

				}

			}
		}
	}

	return result
}

//VerifyAllDevices check all Devices into project
func (pr *Project) VerifyAllDevices() string {
	result := ""
	// TODO: Вначале проверяем наличие всех драйверов
	for _, sub := range pr.Subsystems {
		for _, dev := range sub.RealDevices {
			drv, ok := pr.DefDrivers.Drivers[dev.Driver]
			if !ok {
				result += "Error! subsystem " + sub.Name + " device " + dev.Name + " not found driver " + dev.Driver + "\n"
				continue
			}
			// TODO: Потом проверяем все назначения на устройства со стороны драйвера
			for _, d := range dev.Defs {
				if !drv.FindPIN(d.DriverName) {
					result += "Error! subsystem " + sub.Name + " device " + dev.Name + " driver " + dev.Driver
					result += " not found " + d.DriverName
					continue
				}
			}
		}
	}

	return result

}

//FindPIN find pin on driver
func (ds *DriverXML) FindPIN(name string) bool {
	for _, s := range ds.Signals.Signals {
		if s.Name == name {
			return true
		}
	}
	for _, i := range ds.Inits.Inits {
		if i.Name == name {
			return true
		}
	}
	return false
}

//ReadPINFormat return format pin var
func (ds *DriverXML) ReadPINFormat(name string) string {
	result := ""
	for _, s := range ds.Signals.Signals {
		if s.Name == name {
			return s.Format
		}
	}
	for _, i := range ds.Inits.Inits {
		if i.Name == name {
			return i.Format
		}
	}
	return result
}
