package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "strconv"
)

type MemStat struct {
	cache int64
	rss int64
	mapped_file int64
	pgpgin int64
	pgpgout int64
	pgfault int64
	pgmajfault int64
	inactive_anon int64
	active_anon int64
	inactive_file int64
	active_file int64
	unevictable int64
	hierarchical_memory_limit int64 //NULL for total
}

type CpuStat struct {
	user int64
	system int64
}

func read(path string) string {
    bs, err := ioutil.ReadFile(path)
    if err != nil {
		fmt.Println("Error : requested file doesn't exist.")
        return ""
    }	
    return (string(bs))
}

func findId() string {	
	var id,name [100]string								//Not dynamic !!
	var indexId, indexName [200]int
	var index int = 0									//index used to go through id and name
	var keepFilling bool = true
	var containerName string = "polo417/docpal_old"		//name to change
    var idFound string
    
    str := read("ps")
 
    /*--Displaying all the interesting characters of the file--*/
    for i := 0; i<len(strings.Split(str, " "))-1; i++ {
		
		if ( (strings.Split(str, " ")[i]!="") && (strings.Split(str, " ")[i]!="About") && (strings.Split(str, " ")[i]!="CONTAINER") && (strings.Split(str, " ")[i]!="ID") && (strings.Split(str, " ")[i]!="IMAGE") && (strings.Split(str, " ")[i]!="COMMAND") && (strings.Split(str, " ")[i]!="CREATED") && (strings.Split(str, " ")[i]!="STATUS") && (strings.Split(str, " ")[i]!="PORTS") ){
			if ( (strings.Split(strings.Split(str, " ")[i],"\n")[0]=="") || (i==136) ){
				indexId[index]=i
				indexName[index]=i+3
				index++
			} else {
			}
		}		
    }
    
    /*--retrieving the IDs and the names of the containers from the file--*/
    index=0
    for (keepFilling==true) {
		if (indexId[index]!=0) {
			id[index]=strings.Split(strings.Split(str, " ")[indexId[index]],"\n")[1]
			name[index]=strings.Split(str, " ")[indexName[index]]
			index++
		} else {
			keepFilling=false
		}
	}
	
	index=0
	for ( index < len(name) ) {
		if ( (name[index]==containerName) || (name[index]==containerName+":latest") ) {
			idFound = id[index]
			break
		} else {
			index++
		}
	}
	return idFound
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
	fmt.Println("Memory Statistics : ", memory_stat)
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
	fmt.Println("Memory Statistics (total) : ", memory_stat_total)
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
    fmt.Println("CPU Statistics : ", cpu_stat)
    return cpu_stat
}

func main() {
	getMemInfo()
	getMemTotInfo()
	getCpuInfo()
}
	




