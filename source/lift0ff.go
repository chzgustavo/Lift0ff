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

func execComand(comando *exec.Cmd) error {
	c := comando
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}

func hostname() (string, error) {
	st := ""
	file, err := os.Open("/proc/sys/kernel/hostname")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		st = scanner.Text()
	}
	return st, nil
}

func rtcInfo(dato string) (string, error) {
	rtcInfo := ""
	file, err := os.Open("/proc/driver/rtc")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), dato) {
			_, err := fmt.Sscanf(scanner.Text(), dato+" : %s", &rtcInfo)
			if err != nil {
				return "", err
			}
			break
		}
	}
	return rtcInfo, nil
}

func cpuInfo(dato string) (string, error) {
	sysInfo := ""
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), dato) {
			if dato == "model name" {
				noSpaceString := strings.ReplaceAll(scanner.Text(), "\t", "")
				sysInfo = strings.TrimPrefix(noSpaceString, dato+": ")
				break
			}
			_, err := fmt.Sscanf(scanner.Text(), dato+" : %s", &sysInfo)
			if err != nil {
				return "", err
			}
			break
		}
	}
	return sysInfo, nil
}

func kernelVersion() (string, error) {
	kv := ""
	file, err := os.Open("/proc/version")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_, err := fmt.Sscanf(scanner.Text(), "Linux version %s", &kv)
		if err != nil {
			return "", err
		}
		break
	}
	return kv, nil
}

func versionSo() (string, error) {
	st := ""
	dato := "PRETTY_NAME"
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "", err
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
	return st, err
}

func uptime() (float64, error) {
	var up float64 = 0
	file, err := os.Open("/proc/uptime")
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_, err := fmt.Sscanf(scanner.Text(), "%f", &up)
		if err != nil {
			return 0, err
		}
		break
	}
	return up, nil
}

func who() (string, error) {
	f1, f2 := "", ""
	out, err := exec.Command("who", "-b").Output()
	if err != nil {
		return "", err
	}
	_, error := fmt.Sscanf(string(out), " arranque del sistema %s %s", &f1, &f2)
	if error != nil {
		return "", error
	}
	return f1 + " " + f2, nil
}

func memInfo(mem string) (float64, error) {
	var m float64
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), mem) {
			_, err := fmt.Sscanf(scanner.Text(), mem+" %f", &m)
			if err != nil {
				return 0, err
			}
			break
		}
	}
	return m, nil
}

func memDisk() ([]string, error) {
	m1, m2, m3, m4, m5, m6 := "", "", "", "", "", ""
	mDisk := make([]string, 0)
	cmd := exec.Command("df", "-h")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return mDisk, err
	}
	if err := cmd.Start(); err != nil {
		return mDisk, err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		_, error := fmt.Sscanf(scanner.Text(), "%s %s %s %s %s %s", &m1, &m2, &m3, &m4, &m5, &m6)
		if error != nil {
			return mDisk, err
		}
		if m6 == "/" {
			mDisk = append(mDisk, m2, m3, m4, m5)
			break
		}
	}
	return mDisk, nil
}

func networkData() (*list.List, error) {
	n1, n2, n3, n4, n5, n6, n7, n8 := "", "", "", "", "", "", "", ""
	listaRed := list.New()
	red := make([]string, 9)
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return listaRed, err
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
				return listaRed, er1
			}
			_, er2 := fmt.Sscanf(upl, "%s %s %s %s", &n5, &n6, &n7, &n8)
			if er2 != nil {
				return listaRed, er2
			}
			red[0], red[1], red[2], red[3], red[4], red[5], red[6], red[7], red[8] = intfz, n1, n2, n3, n4, n5, n6, n7, n8
			listaRed.PushBack(red)
			red = []string{"", "", "", "", "", "", "", "", ""}
		}
	}
	return listaRed, nil
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

	fecha, err1 := rtcInfo("rtc_date")
	if err1 != nil {
		fmt.Println(err1)
	}
	hora, err2 := rtcInfo("rtc_time")
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println("       Fecha y Hora RTC:", fecha, hora)

	hostname, err3 := hostname()
	if err3 != nil {
		fmt.Println(err3)
	}
	fmt.Println("	       Hostname:", hostname)

	vendor, err4 := cpuInfo("vendor_id")
	if err4 != nil {
		fmt.Println(err4)
	}
	fmt.Println("  Fabricante Procesador:", vendor)

	modelo, err5 := cpuInfo("model name")
	if err5 != nil {
		fmt.Println(err5)
	}
	fmt.Println("      Modelo Procesador:", modelo)

	kv, err6 := kernelVersion()
	if err6 != nil {
		fmt.Println(err6)
	}
	fmt.Println("      Vesion del kernel:", kv)

	vs, err7 := versionSo()
	if err7 != nil {
		fmt.Println(err7)
	}
	fmt.Println("	  Version de SO:", vs)

	up, err8 := uptime()
	if err8 != nil {
		fmt.Println(err8)
	}
	dias := int(up / 86400)
	horas := int((up) / 3600)
	minutos := int(((up / 3600) - float64(horas)) * 60)
	segundos := int(((((up / 3600) - float64(horas)) * 60) - float64(minutos)) * 60)
	fmt.Printf("       Tiempo activo So: %vd :%vh :%vm :%vs \n", dias, horas, minutos, segundos)

	start, err9 := who()
	if err9 != nil {
		fmt.Println(err9)
	}
	fmt.Printf("	 Inicio sistema: %s \n\n", start)
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
	disk, err := memDisk()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("\nMemoria Disco\n")
		fmt.Println("Memory:")
		fmt.Println("    Memory Total:", strings.Join(disk[0:1], ""))
		fmt.Println("     Memory Used:", strings.Join(disk[1:2], ""))
		fmt.Println("Memory Available:", strings.Join(disk[2:3], ""))
		fmt.Printf(" percentage used: %s\n\n", strings.Join(disk[3:4], ""))
	}
}

func port() {
	err := execComand(exec.Command("sudo", "lsof", "-i", "-P"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
}

func processor() {
	fmt.Println()
	for i := 0; i < 5; i++ {
		titulo := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columna := color.New(color.FgHiRed).SprintfFunc()

		tblp := table.New("CPU", "Porcentaje")
		tblp.WithHeaderFormatter(titulo).WithFirstColumnFormatter(columna)
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
	err := execComand(exec.Command("top", "-d", "5"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
}

func proce() {
	err := execComand(exec.Command("ps", "aux"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
}

func ram() {
	fmt.Println("\nMemoria RAM\n")
	fmt.Println("Memory:")
	memTotal, _ := memInfo("MemTotal:")
	memAvailable, _ := memInfo("MemAvailable:")
	fmt.Println("    Memory Total:", math.Round(memTotal*0.000001*100)/100, "GiB")
	fmt.Println("     Memory Used:", math.Round((memTotal-memAvailable)*0.000001*100)/100, "GiB")
	fmt.Println("Memory Available:", math.Round(memAvailable*0.001*100)/100, "MiB")
	cache, _ := memInfo("Cached:")
	memFree, _ := memInfo("MemFree:")
	fmt.Println("    Memory Cache:", math.Round(cache*0.001*100)/100, "MiB")
	fmt.Println("     Memory Free:", math.Round(memFree*0.001*100)/100, "MiB")
	fmt.Println("\nMemory Swap:")
	memTotalSwap, _ := memInfo("SwapTotal:")
	memSwapFree, _ := memInfo("SwapFree:")
	fmt.Println("     Total Swap:", math.Round(memTotalSwap*0.000001*100)/100, "GiB")
	fmt.Println("      Swap Used:", math.Round((memTotalSwap-memSwapFree)*0.000001*100)/100, "GiB")
	fmt.Println("      Swap Free:", math.Round(memSwapFree*0.000001*100)/100, "GiB", "\n")
}

func red() {
	fmt.Println("\nDatos de Red\n")
	fmt.Println("Datos Recibidos:")

	titulo := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columna := color.New(color.FgHiRed).SprintfFunc()

	tbld := table.New("Interface", "bytes", "packets", "errs", "drop")
	tbld.WithHeaderFormatter(titulo).WithFirstColumnFormatter(columna)

	tblu := table.New("Interface", "bytes", "packets", "errs", "drop")
	tblu.WithHeaderFormatter(titulo).WithFirstColumnFormatter(columna)

	listaRed, err := networkData()
	if err != nil {
		fmt.Println(err)
	}
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
		fmt.Println("Error, ejecute  ./liftoff --help, para ver comando habilitados")
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
