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
			for _, d := range dev.Defs {
				_, ok := sub.Variables[d.Name]
				if !ok {
					//Переменная попалась первый раз
					result += "Error! Subsystem " + sub.Name + " device  var " + d.Name + " into " + dev.Name + " haven't variable\n"
				}
			}
		}
	}

	return result
}

//VerifyAllDevices check all Devices into project
func (pr *Project) VerifyAllDevices() error {
	// TODO: Вначале проверяем наличие всех драйверов
	// TODO: Потом проверяем все назначения на устройства со стороны драйвера
	return nil

}
