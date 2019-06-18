package project

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

//MakeMaster make the master header for project
func (p *Project) MakeMaster(prPath string) error {
	for _, s := range p.Subs {
		path := prPath + "/" + s.Path + "/src/master.h"
		path = RepairPath(path)
		err := os.Remove(path)
		if err != nil {
			err = fmt.Errorf("Error! Remove file " + path + " " + err.Error())
		}

		sw, err := os.Create(path)
		if err != nil {
			err = fmt.Errorf("Error! Create file " + path + " " + err.Error())
		}

		defer sw.Close()
		sw.WriteString("#ifndef " + strings.ToUpper(s.Name) + "_H\n")
		sw.WriteString("#defile " + strings.ToUpper(s.Name) + "_H\n")
		sub := p.Subsystems[s.Name]
		model := p.Models[sub.Model.Name]
		sw.WriteString("// Подсистема  " + s.Name + ":" + s.Description + "\n")
		simul := "0"
		if p.General.Simul == "on" {
			simul = "1"
		}
		sw.WriteString("static char SimulOn=" + simul + ";\n")
		sw.WriteString("static short CodeSub=" + s.ID + ";\n")
		sw.WriteString("static char SimulIP[]=\"" + p.General.IP + "\\0\";\n")
		sw.WriteString("static int SimulPort=" + p.General.Port + ";\n")
		sw.WriteString("static int StepCycle=" + s.Step + ";\t //Время цикла мс\n")
		sw.WriteString("float takt,taktScheme=0,taktSS=0;\n")
		sw.WriteString("#define SIZE_BUFFER " + strconv.Itoa(sub.SizeBuffer) + "\n")
		sw.WriteString("static char BUFFER[SIZE_BUFFER];\n")
		if !model.ActionContain("one_ip") {
			sw.WriteString("#include <fp8/UDPTransport.h>\n")
			sw.WriteString("SetupUDP setUDP ={\"" + s.Main + "\\0\",5432,\"" + s.Second + "\\0\",5432,BUFFER,sizeof(BUFFER),};\n")
			sw.WriteString("int master=1,nomer=1;\n")
		}
		for _, v := range sub.Variables {
			sw.WriteString("#define " + v.Name + "\tBUFFER[" + strconv.Itoa(v.Address) + "]\t// " + v.Description + "\n")
			sw.WriteString("#define id" + v.Name + "\t" + strconv.Itoa(v.ID) + "\t// " + v.Description + "\n")
		}
		sw.WriteString("#pragma pack(push,1)\n")
		sw.WriteString("static VarCtrl allVariables[]={ \t\t\t //Описание всех переменных\n")
		for _, v := range sub.Variables {
			sw.WriteString("\t " + strconv.Itoa(v.ID) + "\t," + v.Format + "\t," + v.Size + "\t,&" + v.Name + "},\t//" + v.Description + "\n")
		}
		sw.WriteString("\t{-1,0,NULL},\n}\n")
		sw.WriteString("static char NameSaveFile[]=\"" + sub.NameSaveFile + "\\0\"; //Имя файла сохранения переменных\n")
		sw.WriteString("#pragma pop\n")
		sw.WriteString("static VarSaveCtrl saveVariables[]={\t//Id переменных для сохранения\n")
		for _, sv := range sub.MapSaves {
			v := sub.Variables[sv.Name]
			sw.WriteString("\t{" + strconv.Itoa(v.ID) + ",\"" + v.Name + "\\0\"},\t//" + v.Description + "\n")
		}
		sw.WriteString("\t{0,NULL}\n};\n")
		modStr := "static ModbusDevice modbuses[]={\n"
		for _, mb := range sub.Modbuses {
			coil := "#pragma pack(push,1)\nstatic ModbusRegister coil_" + mb.Name + "[]={\n"
			di := "#pragma pack(push,1)\nstatic ModbusRegister di_" + mb.Name + "[]={\n"
			ir := "#pragma pack(push,1)\nstatic ModbusRegister ir_" + mb.Name + "[]={\n"
			hr := "#pragma pack(push,1)\nstatic ModbusRegister hr_" + mb.Name + "[]={\n"
			for _, r := range mb.Registers {
				str := "\t{&" + r.Name + "," + strconv.Itoa(r.Format) + "," + strconv.Itoa(r.Address) + "},\t//" + r.Description + "\n"
				switch r.Type {
				case 0:
					coil += str
				case 1:
					di += str
				case 2:
					ir += str
				case 3:
					hr += str
				}
			}
			coil += "\t{NULL,0,0},\n}\n"
			di += "\t{NULL,0,0},\n}\n"
			ir += "\t{NULL,0,0},\n}\n"
			hr += "\t{NULL,0,0},\n}\n"
			sw.WriteString(coil)
			sw.WriteString("#pragma pop\n")
			sw.WriteString(di)
			sw.WriteString("#pragma pop\n")
			sw.WriteString(ir)
			sw.WriteString("#pragma pop\n")
			sw.WriteString(hr)
			sw.WriteString("#pragma pop\n")
			modStr += "\t{"
			if mb.IsMaster() {
				sw.WriteString("static char " + mb.Name + "_ip1[]={\"" + mb.IP1 + "\\0\"};\n")
				sw.WriteString("static char " + mb.Name + "_ip2[]={\"" + mb.IP2 + "\\0\"};\n")
				modStr += "1"
			} else {
				modStr += "0"
			}
			modStr += "," + mb.Port + ",&coil_" + mb.Name + "[0],&di_" + mb.Name + "[0],&di_" + mb.Name + "[0],&hr_" + mb.Name + "[0]"
			modStr += ",NULL"
			if mb.IsMaster() {
				modStr += "," + mb.Name + "_ip1"
				modStr += "," + mb.Name + "_ip2"
				modStr += "," + mb.Step
			} else {
				modStr += ",NULL,NULL,0"
			}
			modStr += "},\t//" + mb.Description + "\n"
		}
		sw.WriteString("#pragma pack(push,1)\n")
		sw.WriteString(modStr)
		sw.WriteString("\t{0,-1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,0},\n};\n")
		sw.WriteString("#pragma pop\n")
		for _, dev := range sub.RealDevices {
			sw.WriteString(dev.MakeDriverTable(p.DefDrivers))
		}
		sw.WriteString("#pragma pack(push,1)\n")
		sw.WriteString("static Drive drivers[]={\n")
		for _, dev := range sub.RealDevices {
			drv := p.DefDrivers.Drivers[dev.Driver]
			sw.WriteString("\t{" + drv.Code + ",0x" + dev.Slot + "," + drv.LenData + ",def_buf_" + dev.Name + ",&table_" + dev.Name + "},\t//" + dev.Description + "\n")
		}
		sw.WriteString("\t{0,0,0,0,NULL,NULL},\n};\n")
		sw.WriteString("#pragma pop\n")

		sw.Close()
	}
	return nil
}
