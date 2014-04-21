package main

import (
    "fmt"
    "strings"
    "strconv"
	"bytes"
	"log"
	"os/exec"
    "io/ioutil"
)

type MemStat struct {
	Cache int64
	Rss int64
	Mapped_file int64
	Pgpgin int64
	Pgpgout int64
	Pgfault int64
	Pgmajfault int64
	Inactive_anon int64
	Active_anon int64
	Inactive_file int64
	Active_file int64
	Unevictable int64
	Hierarchical_memory_limit int64 //NULL for total
}

type CpuStat struct {
	User int64
	System int64
}

func read(path string) string {
    bs, err := ioutil.ReadFile(path)
    if err != nil {
		fmt.Println("Error : requested file doesn't exist.")
        return ""
    }	
    return (string(bs))
}

func doPs() string {
	cmd := exec.Command("docker", "ps", "-notrunc")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return(out.String())
}

func findId() string {	
    str := string(doPs())
	str2 := strings.Split(str, "\n")
	for i := 1; i<len(str2)-1; i++ {
		if ( (strings.Fields(str2[i])[1]=="polo417/docpal_old") || (strings.Fields(str2[i])[1]=="polo417/docpal_old:latest") ) {	//container name to change
			return strings.Fields(str2[i])[0]
		}
	}
	return "[000]"
}

func getMemInfo() MemStat {
	var path,str string
	var memory_stat MemStat
	var memBuffer [13]int64

	path="/sys/fs/cgroup/memory/lxc/"+findId()+"/memory.stat"
	str = read(path)	
	for i:=0; i<13; i++ {
		charizard, _ := strconv.ParseInt(strings.Split(strings.Split(str,"\n")[i]," ")[1], 0, 64)
		memBuffer[i]=charizard;
	}
	memory_stat = MemStat {memBuffer[0],memBuffer[1],memBuffer[2],memBuffer[3],memBuffer[4],memBuffer[5],memBuffer[6],memBuffer[7],memBuffer[8],memBuffer[9],memBuffer[10],memBuffer[11],memBuffer[12]}
	//fmt.Println("Memory Statistics : ", memory_stat)
	return memory_stat
}

func getMemTotInfo() MemStat {
	var path,str string
	var memory_stat_total MemStat
	var memBuffer [13]int64

	path="/sys/fs/cgroup/memory/lxc/"+findId()+"/memory.stat"
	str = read(path)
	for i:=0; i<13; i++ {
		charizard, _ := strconv.ParseInt(strings.Split(strings.Split(str,"\n")[i]," ")[1], 0, 64)
		memBuffer[i]=charizard;
	}
	memory_stat_total = MemStat {memBuffer[0],memBuffer[1],memBuffer[2],memBuffer[3],memBuffer[4],memBuffer[5],memBuffer[6],memBuffer[7],memBuffer[8],memBuffer[9],memBuffer[10],memBuffer[11],0}
	//fmt.Println("Memory Statistics (total) : ", memory_stat_total)
	return memory_stat_total
}

func getCpuInfo() CpuStat {
	var path,str string
	var cpu_stat CpuStat
	var cpuBuffer [2]int64

	path="/sys/fs/cgroup/cpuacct/lxc/"+findId()+"/cpuacct.stat"
    str = read(path)
	for i:=0; i<2; i++ {
		charizard, _ := strconv.ParseInt(strings.Split(strings.Split(str,"\n")[i]," ")[1], 0, 64)
		cpuBuffer[i]=charizard;
	}
	cpu_stat = CpuStat {cpuBuffer[0],cpuBuffer[1]}
    //fmt.Println("CPU Statistics : ", cpu_stat)
    return cpu_stat
}

func main() {
	fmt.Println(getMemInfo().Cache, getCpuInfo().User)
}
	

