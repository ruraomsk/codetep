package project

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

//MakeMaster make the master header for project
func (p *Project) MakeMaster(prPath string) error {
	for _, s := range p.Subs {
		path := prPath + "/" + s.Path + "/src/master.h"
		path = RepairPath(path)
		err := os.Remove(path)
		// if err != nil {
		// 	err = fmt.Errorf("Error! Remove file " + path + " " + err.Error())
		// 	return err
		// }

		sw, err := os.Create(path)
		if err != nil {
			err = fmt.Errorf("Error! Create file " + path + " " + err.Error())
			return err
		}

		defer sw.Close()
		sw.WriteString("#ifndef " + strings.ToUpper(s.Name) + "_H\n")
		sw.WriteString("#define " + strings.ToUpper(s.Name) + "_H\n")
		sub := p.Subsystems[s.Name]
		model := p.Models[sub.Model.Name]
		sw.WriteString("// Подсистема  " + s.Name + ":" + s.Description + "\n")
		simul := "0"
		if p.Simul == "on" {
			simul = "1"
		}
		sts, tvars, err := p.LoadShema(s)
		if err != nil {
			return err
		}
		sub.AppendNewVariables(tvars)
		sw.WriteString("static char SimulOn=" + simul + ";\n")
		sw.WriteString("static short CodeSub=" + s.ID + ";\n")
		sw.WriteString("static char SimulIP[]=\"" + p.IP + "\\0\";\n")
		sw.WriteString("static int SimulPort=" + p.Port + ";\n")
		sw.WriteString("static int StepCycle=" + s.Step + ";\t //Время цикла мс\n")
		sw.WriteString("float takt,taktScheme=0,taktSS=0;\n")
		sw.WriteString("#define SIZE_BUFFER " + strconv.Itoa(sub.SizeBuffer) + "\n")
		sw.WriteString("static char BUFFER[SIZE_BUFFER];\n")
		if !model.ActionContain("one_ip") {
			sw.WriteString("#include <fp8/UDPTransport.h>\n")
			sw.WriteString("SetupUDP setUDP ={\"" + s.Main + "\\0\",5432,\"" + s.Second + "\\0\",5432,BUFFER,sizeof(BUFFER),};\n")
			sw.WriteString("int master=1,nomer=1;\n")
		}
		sort.Slice(sub.Vars.ListVariable, func(i, j int) bool { return sub.Vars.ListVariable[i].Name < sub.Vars.ListVariable[j].Name })
		for _, v := range sub.Vars.ListVariable {
			sw.WriteString("#define " + v.Name + "\tBUFFER[" + strconv.Itoa(v.Address) + "]\t// " + v.Description + "\n")
			sw.WriteString("#define id" + v.Name + "\t" + strconv.Itoa(v.ID) + "\t// " + v.Description + "\n")
		}
		sw.WriteString("#pragma pack(push,1)\n")
		sw.WriteString("static VarCtrl allVariables[]={ \t\t\t //Описание всех переменных\n")
		for _, v := range sub.Vars.ListVariable {
			sw.WriteString("\t{" + strconv.Itoa(v.ID) + "\t," + v.Format + "\t," + v.Size + "\t,&" + v.Name + "},\t//" + v.Description + "\n")
		}
		sw.WriteString("\t{-1,0,NULL},\n};\n")
		sw.WriteString("static char NameSaveFile[]=\"" + sub.NameSaveFile + "\\0\"; //Имя файла сохранения переменных\n")
		sw.WriteString("#pragma pop\n")
		sw.WriteString("static VarSaveCtrl saveVariables[]={\t//Id переменных для сохранения\n")
		saves := make([]Save, 0)
		for _, sv := range sub.MapSaves {
			saves = append(saves, sv)
		}
		sort.Slice(saves, func(i, j int) bool { return saves[i].Name < saves[j].Name })
		for _, sv := range saves {
			v := sub.Variables[sv.Name]
			sw.WriteString("\t{" + strconv.Itoa(v.ID) + ",\"" + v.Name + "\\0\"},\t//" + v.Description + "\n")
		}
		sw.WriteString("\t{0,NULL}\n};\n")
		modStr := "static ModbusDevice modbuses[]={\n"
		sort.Slice(sub.Modbuses, func(i, j int) bool { return sub.Modbuses[i].Port < sub.Modbuses[j].Port })
		for _, mb := range sub.Modbuses {
			coil := "#pragma pack(push,1)\nstatic ModbusRegister coil_" + mb.Name + "[]={\n"
			di := "#pragma pack(push,1)\nstatic ModbusRegister di_" + mb.Name + "[]={\n"
			ir := "#pragma pack(push,1)\nstatic ModbusRegister ir_" + mb.Name + "[]={\n"
			hr := "#pragma pack(push,1)\nstatic ModbusRegister hr_" + mb.Name + "[]={\n"
			regs := make([]Register, 0)
			for _, r := range mb.Registers {
				regs = append(regs, r)
			}

			sort.Slice(regs, func(i, j int) bool { return regs[i].Address < regs[j].Address })

			for _, r := range regs {
				format := "0"
				switch r.Type {
				case 0:
					format = "1"
				case 1:
					format = "1"
				case 2:
					format = strconv.Itoa(r.Format)
				case 3:
					format = strconv.Itoa(r.Format)
				}

				str := "\t{&" + r.Name + "," + format + "," + strconv.Itoa(r.Address) + "},\t//" + r.Description + "\n"
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
			coil += "\t{NULL,0,0},\n};\n"
			di += "\t{NULL,0,0},\n};\n"
			ir += "\t{NULL,0,0},\n};\n"
			hr += "\t{NULL,0,0},\n};\n"
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
			modStr += "," + mb.Port + ",&coil_" + mb.Name + "[0],&di_" + mb.Name + "[0],&ir_" + mb.Name + "[0],&hr_" + mb.Name + "[0]"
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
		realDevice := make([]Device, 0, len(sub.RealDevices))
		for _, dev := range sub.RealDevices {
			realDevice = append(realDevice, dev)
		}
		sort.Slice(realDevice, func(i, j int) bool { return realDevice[i].Name < realDevice[j].Name })
		for _, dev := range realDevice {
			sw.WriteString(dev.MakeDriverTable(p.DefDrivers))
		}
		sw.WriteString("#pragma pack(push,1)\n")
		sw.WriteString("static Driver drivers[]={\n")

		for _, dev := range realDevice {
			drv := p.DefDrivers.Drivers[dev.Driver]
			sw.WriteString("\t{" + drv.Code + ",0x" + dev.Slot + "," + drv.LenInit + "," + drv.LenData + ",def_buf_" + dev.Name + ",&table_" + dev.Name + "},\t//" + dev.Description + "\n")
		}
		sw.WriteString("\t{0,0,0,0,NULL,NULL},\n};\n")
		sw.WriteString("#pragma pop\n")
		sw.WriteString("void InitSetConst(void){\t//Инициализация переменных для хранения\n")
		sliceSaves := make([]Save, 0, len(sub.MapSaves))
		for _, sv := range sub.MapSaves {
			sliceSaves = append(sliceSaves, sv)
		}
		sort.Slice(sliceSaves, func(i, j int) bool { return sliceSaves[i].Name < sliceSaves[j].Name })

		for _, sv := range sliceSaves {
			v := sub.Variables[sv.Name]
			fname := v.getFunctionSet()
			sw.WriteString("\t" + fname + "(" + strconv.Itoa(v.ID) + "," + sv.Value + ");\n")
		}
		for _, ii := range sub.IniSignal.Isignals {
			v := sub.Variables[ii.Name]
			fname := v.getFunctionSet()
			sw.WriteString("\t" + fname + "(" + strconv.Itoa(v.ID) + "," + ii.Value + ");\n")
		}
		sw.WriteString("}\n")
		if model.ActionContain("add_vchs") {
			s, err := p.LoadCodePart("add_VCHS")
			if err != nil {
				return err
			}
			sw.WriteString(s)
		}
		ss := sub.MakeDelayHeader()
		// fmt.Println(sub.Name, ss)
		sw.WriteString(ss)
		// sts, tvars, err := p.LoadShema(s)
		// if err != nil {
		// 	return err
		// }
		for _, s := range sts {
			// if strings.Contains(s, "void Scheme()") {
			// 	sw.WriteString("void ZeroVar() {\n")
			// 	for name, tv := range tvars {
			// 		sw.WriteString("\t" + name + tv + "\n")
			// 	}
			// 	sw.WriteString("}\n")
			// }
			sw.WriteString(s + "\n")
		}
		ss = sub.MakeMainCycleFunc()
		sw.WriteString(ss)
		sw.WriteString("#endif")
		sw.Close()
	}
	return nil
}

//MakeMainC maker maic C file for all subsystems of project
func (p *Project) MakeMainC(prPath string) error {
	for _, s := range p.Subs {
		sub := p.Subsystems[s.Name]
		model := p.Models[sub.Model.Name]
		path := prPath + "/" + s.Path + "/src/mainfp.c"
		path = RepairPath(path)
		err := os.Remove(path)
		// if err != nil {
		// 	err = fmt.Errorf("Error! Remove file " + path + " " + err.Error())
		// 	return err
		// }

		sw, err := os.Create(path)
		if err != nil {
			err = fmt.Errorf("Error! Create file " + path + " " + err.Error())
			return err
		}

		defer sw.Close()

		ipath := RepairPath(prPath + "/settings/src-FP/mainFP.c")
		lpath := RepairPath(prPath + "/settings/src-FP/")
		file, err := os.Open(ipath)
		if err != nil {
			err = fmt.Errorf("Error! Opening file " + ipath + " " + err.Error())
			return err
		}
		defer file.Close()

		sReader := bufio.NewScanner(file)
		for sReader.Scan() {
			line := sReader.Text()
			if !strings.Contains(line, "%attach_") {
				sw.WriteString(line + "\n")
				continue
			}
			us := strings.Split(line, "_")
			if len(us) != 2 {
				continue
			}
			nameFile := model.AttachPath(us[1])
			if nameFile == "" {
				// err = fmt.Errorf("Error! В подсистеме " + sub.Name + " при генерации main.c нет вставки " + us[1])
				// return err
				continue
			}
			pp := lpath + nameFile + ".c"
			buf, err := ioutil.ReadFile(pp)
			if err != nil {
				err = fmt.Errorf("Error! В подсистеме " + sub.Name + " при генерации main.c нет " + pp + "!")
				return err
			}
			ssline := string(buf)
			ssline = strings.Replace(ssline, "%NetBlKey%", "id"+sub.Netblkey.Name, 1)

			sw.WriteString(ssline + "\n")
		}
		file.Close()
		sw.Close()
	}
	return nil
}
