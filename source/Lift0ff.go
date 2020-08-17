package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

func execComand(com *exec.Cmd) {
	c := com
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()

	if err != nil {
		fmt.Println("Error execution: ", err)
		fmt.Println("Error")
	} /*else {
		fmt.Println("Ejecucion correcta")
	}*/
}

func singleData(path string) string {
	st := ""
	file, err := os.Open(path)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		st = scanner.Text()
	}
	return st
}

func simpleData(path string, dato string) string {
	st := ""
	file, err := os.Open(path)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), dato) {
			st = scanner.Text()
			//st = removeSpace(st)
			noSpaceString := strings.ReplaceAll(st, "\t", "")
			st = strings.TrimPrefix(noSpaceString, dato+": ")
			break
		}
	}
	return st
}

func kernelVersion(path string) (string, string) {
	kv1, kv2 := "", ""
	file, err := os.Open(path)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_, err := fmt.Sscanf(scanner.Text(), "%s version %s", &kv1, &kv2)
		if err != nil {
			//panic(err)
			fmt.Println(err)
		}
		break
	}
	return kv1, kv2
}

func versionSo(path string, dato string) string {
	st := ""
	file, err := os.Open(path)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), dato) {
			st = scanner.Text()
			st = strings.TrimPrefix(st, dato+"=")
			st = st[1 : len(st)-1]
			break
		}
	}
	return st
}

func tiempoActivoSo(path string) float64 {
	var tA float64 = 0.0

	file, err := os.Open(path)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		_, err := fmt.Sscanf(scanner.Text(), "%f", &tA)
		if err != nil {
			panic(err)
		}
		break
	}
	return tA
}

func fechaInicioSistema(com *exec.Cmd) (string, string) {
	f1, f2 := "", ""
	out, err := com.Output()
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	_, error := fmt.Sscanf(string(out), " arranque del sistema %s %s", &f1, &f2)
	if error != nil {
		//panic(error)
		fmt.Println(error)
	}
	return f1, f2
}

func memInfo(mem string) float64 {
	var m1 float64
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), mem) {
			_, err := fmt.Sscanf(scanner.Text(), mem+" %f", &m1)
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}
	return m1
}

func memDisk() (string, string, string, string) {
	m1, m2, m3, m4, m5, m6 := "", "", "", "", "", ""
	cmd := exec.Command("df", "-h")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		_, error := fmt.Sscanf(scanner.Text(), "%s %s %s %s %s %s", &m1, &m2, &m3, &m4, &m5, &m6)
		if error != nil {
			//panic(error)
			fmt.Println(error)
		}
		if m6 == "/" {
			break
		}
	}
	return m2, m3, m4, m5
}

func main() {
	argsConProg := os.Args
	if len(argsConProg) == 1 {
		fmt.Println("Error, ejecute  liftoff --help, para ver comando habilitados")
		os.Exit(1)
	}
	sentencia := os.Args[1:]
	switch stSentencia := strings.Join(sentencia, " "); stSentencia {
	case "-h", "--help":
		fmt.Println("\nUsage:	./liftoff [OPTIONS] COMMAND")
		fmt.Println("\nOptions:")
		fmt.Println("	-v, --version			Imprime la version de liftoff")
		fmt.Println("	-i, --info			Muestra información del sistema")
		fmt.Println("\nCommands:")
		fmt.Println("	ram				Muestra información de memoria RAM")
		fmt.Println("	mem				Muestra información de memoria Disco")
		fmt.Println("	port				Muestra puertos TCP y UDP (LISTEN, ESTABLISHED)")
		fmt.Println("	proct			    	Muestra procesos del sistema actualizados cada 5 segundos en tiempo real. Finalizar con Ctrl + c")
		fmt.Println("	proce			    	Muestra los procesos del sistema de manera estatica")
	case "-v", "--version":
		fmt.Println(" _        _    __   _              __    __  ")
		fmt.Println("| |      (_)  / _| | |            / _|  / _| ")
		fmt.Println("| |       _  | |_  | |_    ____  | |_  | |_  ")
		fmt.Println("| |      | | |  _| | __|  / _  ) |  _| |  _| ")
		fmt.Println("| |____  | | | |   | |_  | (_) | | |   | |   ")
		fmt.Println("|______| |_| |_|   |___| (____/  |_|   |_|   ")
		fmt.Println("                                            Version 1.0.0 beta")
	case "-i", "--info":
		fmt.Println("Informacion del Sistema\n")
		fmt.Println("       Fecha y Hora RTC:", simpleData("/proc/driver/rtc", "rtc_date"), simpleData("/proc/driver/rtc", "rtc_time"))
		fmt.Println("	       Hostname:", singleData("/proc/sys/kernel/hostname"))
		fmt.Println("  Fabricante Procesador:", simpleData("/proc/cpuinfo", "vendor_id"))
		fmt.Println("      Modelo Procesador:", simpleData("/proc/cpuinfo", "model name"))
		kv1, kv2 := kernelVersion("/proc/version")
		fmt.Println("      Vesion del kernel:", kv1, kv2)
		fmt.Println("	  Version de SO:", versionSo("/etc/os-release", "PRETTY_NAME"))
		tA := tiempoActivoSo("/proc/uptime")
		dias := int(tA / 86400)
		horas := int((tA) / 3600)
		minutos := int(((tA / 3600) - float64(horas)) * 60)
		segundos := int(((((tA / 3600) - float64(horas)) * 60) - float64(minutos)) * 60)
		fmt.Printf("       Tiempo activo So: %vd :%vh :%vm :%vs \n", dias, horas, minutos, segundos)
		f1, f2 := fechaInicioSistema(exec.Command("who", "-b"))
		fmt.Println("	 Inicio sistema:", f1, f2)
	case "ram":
		fmt.Println("\nMemoria RAM\n")
		fmt.Println("Memory:")
		memTotal := memInfo("MemTotal:")
		memAvailable := memInfo("MemAvailable:")
		fmt.Println("    Memory Total:", math.Round(memTotal*0.000001*100)/100, "GiB")
		fmt.Println("     Memory Used:", math.Round((memTotal-memAvailable)*0.000001*100)/100, "GiB")
		fmt.Println("Memory Available:", math.Round(memAvailable*0.001*100)/100, "MiB")
		fmt.Println("    Memory Cache:", math.Round(memInfo("Cached:")*0.001*100)/100, "MiB")
		fmt.Println("     Memory Free:", math.Round(memInfo("MemFree:")*0.001*100)/100, "MiB")
		fmt.Println("\nMemory Swap:")
		memTotalSwap := memInfo("SwapTotal:")
		memSwapFree := memInfo("SwapFree:")
		fmt.Println("     Total Swap:", math.Round(memTotalSwap*0.000001*100)/100, "GiB")
		fmt.Println("      Swap Used:", math.Round((memTotalSwap-memSwapFree)*0.000001*100)/100, "GiB")
		fmt.Println("      Swap Free:", math.Round(memSwapFree*0.000001*100)/100, "GiB")
	case "mem":
		m1, m2, m3, m4 := memDisk()
		fmt.Println("\nMemoria Disco\n")
		fmt.Println("Memory:")
		fmt.Println("    Memory Total:", m1)
		fmt.Println("     Memory Used:", m2)
		fmt.Println("Memory Available:", m3)
		fmt.Println(" percentage used:", m4)
	case "proct":
		execComand(exec.Command("top", "-d", "5"))
	case "proce":
		execComand(exec.Command("ps", "aux"))
	case "port":
		execComand(exec.Command("sudo", "lsof", "-i", "-P"))
	default:
		fmt.Println("\n	Ejecute ./liftoff --help, para ver los comandos validos\n")
	}
}
