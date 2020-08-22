package main

import (
	"bufio"
	"container/list"
	"fmt"
	"math"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/shirou/gopsutil/cpu"
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

func networkData() *list.List {
	n1, n2, n3, n4, n5, n6, n7, n8 := "", "", "", "", "", "", "", ""
	listaRed := list.New()
	red := make([]string, 9)
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan() == true; i++ {
		if i > 1 {
			ndata := strings.TrimSpace(scanner.Text())
			space := regexp.MustCompile(`\s+`)
			noSpaces := space.ReplaceAllString(ndata, " ")
			netData := strings.Split(noSpaces, " ")

			interfaz := netData[0:1]
			download := netData[1:5]
			upload := netData[9:13]

			intfz := strings.Join(interfaz, "")
			intfz = intfz[:len(intfz)-1]
			downl := strings.Join(download, " ")
			upl := strings.Join(upload, " ")

			_, er1 := fmt.Sscanf(downl, "%s %s %s %s", &n1, &n2, &n3, &n4)
			if er1 != nil {
				fmt.Println(er1)
			}
			_, er2 := fmt.Sscanf(upl, "%s %s %s %s", &n5, &n6, &n7, &n8)
			if er2 != nil {
				fmt.Println(er2)
			}
			red[0], red[1], red[2], red[3], red[4], red[5], red[6], red[7], red[8] = intfz, n1, n2, n3, n4, n5, n6, n7, n8
			listaRed.PushBack(red)
			red = []string{"", "", "", "", "", "", "", "", ""}
		}
	}
	return listaRed
}

func help() {
	fmt.Println("\nUsage:	./liftoff [OPTIONS] COMMAND")
	fmt.Println("\nOptions:")
	fmt.Println("	-i, --info			Muestra información del sistema")
	fmt.Println("	-v, --version			Imprime la version de liftoff")
	fmt.Println("\nCommands:")
	fmt.Println("	mem				Muestra información de memoria Disco")
	fmt.Println("	port				Muestra puertos TCP y UDP (LISTEN, ESTABLISHED)")
	fmt.Println("	processor			Muestra los porcentajes del procesador")
	fmt.Println("	proce			    	Muestra los procesos del sistema de manera estatica")
	fmt.Println("	proct			    	Muestra procesos del sistema actualizados cada 5 segundos en tiempo real. Finalizar con Ctrl + c")
	fmt.Println("	ram				Muestra información de memoria RAM")
	fmt.Println("	red				Muestra datos recibidos y transmitidos en la Red\n")
}
func info() {
	fmt.Println("\nInformacion del Sistema\n")
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
	fmt.Printf("	 Inicio sistema: %s %s\n\n", f1, f2)
}

func version() {
	fmt.Println()
	fmt.Println(" _        _    __   _              __    __  ")
	fmt.Println("| |      (_)  / _| | |            / _|  / _| ")
	fmt.Println("| |       _  | |_  | |_    ____  | |_  | |_  ")
	fmt.Println("| |      | | |  _| | __|  / _  ) |  _| |  _| ")
	fmt.Println("| |____  | | | |   | |_  | (_) | | |   | |   ")
	fmt.Println("|______| |_| |_|   |___| (____/  |_|   |_|   ")
	fmt.Println("                                            Version 1.0.0 beta\n")
}

func mem() {
	m1, m2, m3, m4 := memDisk()
	fmt.Println("\nMemoria Disco\n")
	fmt.Println("Memory:")
	fmt.Println("    Memory Total:", m1)
	fmt.Println("     Memory Used:", m2)
	fmt.Println("Memory Available:", m3)
	fmt.Printf(" percentage used: %s\n\n", m4)
}

func port() {
	execComand(exec.Command("sudo", "lsof", "-i", "-P"))
	fmt.Println()
}

func processor() {
	fmt.Println()
	for i := 0; i < 5; i++ {
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgHiRed).SprintfFunc()

		tblp := table.New("CPU", "Porcentaje")
		tblp.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		pCpu, _ := cpu.Percent(time.Second, false)
		pCore, _ := cpu.Percent(1500*time.Millisecond, true)
		tblp.AddRow("Cpu-Total", fmt.Sprintf("%.2f %%", pCpu[0]))
		for i, core := range pCore {
			tblp.AddRow("Cpu"+strconv.Itoa(i), fmt.Sprintf("%.2f %%", core))
		}
		tblp.Print()
		fmt.Println()
	}
}

func proct() {
	execComand(exec.Command("top", "-d", "5"))
	fmt.Println()
}

func proce() {
	execComand(exec.Command("ps", "aux"))
	fmt.Println()
}

func ram() {
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
	fmt.Println("      Swap Free:", math.Round(memSwapFree*0.000001*100)/100, "GiB", "\n")
}

func red() {
	fmt.Println("\nDatos de Red\n")
	fmt.Println("Datos Recibidos:")

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgHiRed).SprintfFunc()

	tbld := table.New("Interface", "bytes", "packets", "errs", "drop")
	tbld.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	tblu := table.New("Interface", "bytes", "packets", "errs", "drop")
	tblu.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	listaRed := networkData()
	for temp := listaRed.Front(); temp != nil; temp = temp.Next() {
		s := reflect.ValueOf(temp.Value)
		tbld.AddRow(s.Index(0), s.Index(1), s.Index(2), s.Index(3), s.Index(4))
		tblu.AddRow(s.Index(0), s.Index(5), s.Index(6), s.Index(7), s.Index(8))
	}
	tbld.Print()
	fmt.Println("\nDatos Transmitidos:")
	tblu.Print()
	fmt.Println()
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
		help()
	case "-i", "--info":
		info()
	case "-v", "--version":
		version()
	case "mem":
		mem()
	case "port":
		port()
	case "processor":
		processor()
	case "proct":
		proct()
	case "proce":
		proce()
	case "ram":
		ram()
	case "red":
		red()
	default:
		fmt.Println("\n	Ejecute ./liftoff --help, para ver los comandos validos\n")
	}
}
