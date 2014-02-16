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


func main() {
	
	var id,name [100]string								//Not dynamic !!
	var indexId, indexName [200]int
	var index int = 0									//index used to go through id and name
	var keepFilling bool = true
	
	var containerName string = "polo417/docpal_old"		//name to change
	var path string
	
	var memory_stat, memory_stat_total MemStat
	var cpu_stat CpuStat
	var memBuffer [13]int64
	var cpuBuffer [2]int64
	
    /*--Reading the file called "ps"--*/
    bs, err := ioutil.ReadFile("ps")
    if err != nil {
		fmt.Println("Error : requested file doesn't exist.")
        return
    }	
    
    /*--retrieving all the characters of the file--*/
    str := string(bs)
    //fmt.Println(str)
    
    /*--Displaying all the interesting characters of the file--*/
    for i := 0; i<len(strings.Split(str, " "))-1; i++ {
		//fmt.Println(i)
		
		if ( (strings.Split(str, " ")[i]!="") && (strings.Split(str, " ")[i]!="About") && (strings.Split(str, " ")[i]!="CONTAINER") && (strings.Split(str, " ")[i]!="ID") && (strings.Split(str, " ")[i]!="IMAGE") && (strings.Split(str, " ")[i]!="COMMAND") && (strings.Split(str, " ")[i]!="CREATED") && (strings.Split(str, " ")[i]!="STATUS") && (strings.Split(str, " ")[i]!="PORTS") ){
			if ( (strings.Split(strings.Split(str, " ")[i],"\n")[0]=="") || (i==136) ){
				//fmt.Println("i=", i, " : ", strings.Split(strings.Split(str, " ")[i],"\n")[1]," +++ ", len(strings.Split(strings.Split(str, " ")[i],"\n")[1]))
				indexId[index]=i
				indexName[index]=i+3
				index++
			} else {
				//fmt.Println("i=", i, " : ", strings.Split(str, " ")[i]," +++ ", len(strings.Split(str, " ")[i]))
			}
		}		
    }
    
    /*--retrieving the IDs and the names of the containers from the file--*/
    index = 0
    for (keepFilling==true) {
		if (indexId[index]!=0) {
			id[index]=strings.Split(strings.Split(str, " ")[indexId[index]],"\n")[1]
			name[index]=strings.Split(str, " ")[indexName[index]]
			index++
		} else {
			keepFilling=false
		}
	}
    
/*    fmt.Println("PIKAAAAAAAAAAAAAAAAA")
	fmt.Println(indexId)
	fmt.Println(indexName)
	fmt.Println(id)
	fmt.Println(name)
	fmt.Println(keepFilling)	*/
	
	/*--Reading the right memory metric file--*/
	index=0
	for ( index < len(name) ) {
		if ( (name[index]==containerName) || (name[index]==containerName+":latest") ) {
			//fmt.Println("Found ! Index = ", index)
			path="/sys/fs/cgroup/memory/lxc/"+id[index]+"/memory.stat"
			break
		} else {
			index++
		}
	}
	//fmt.Println(containerName)
	//fmt.Println(path)
	bs, err = ioutil.ReadFile(path)
    if err != nil {
        return
    }
    str = string(bs)
    fmt.Println(str)	
	fmt.Println("PIKAAAAAAAAAAAAAAAAAAAAAA")
	
	/*--Assigning values to the variables used to store the memory metrics--*/
	//fmt.Println(strings.Split(str,"\n"))
	for i:=0; i<13; i++ {
		//fmt.Println(strings.Split(strings.Split(str,"\n")[i]," ")[1])
		charizard, _ := strconv.ParseInt(strings.Split(strings.Split(str,"\n")[i]," ")[1], 0, 64)
		//fmt.Println("CHAAAAAAAAAAAAAAAAAA", charizard)
		memBuffer[i]=charizard;
	}
	memory_stat = MemStat {memBuffer[0],memBuffer[1],memBuffer[2],memBuffer[3],memBuffer[4],memBuffer[5],memBuffer[6],memBuffer[7],memBuffer[8],memBuffer[9],memBuffer[10],memBuffer[11],memBuffer[12]}
	for i:=0; i<13; i++ {
		//fmt.Println(strings.Split(strings.Split(str,"\n")[i]," ")[1])
		charizard, _ := strconv.ParseInt(strings.Split(strings.Split(str,"\n")[i]," ")[1], 0, 64)
		//fmt.Println("CHAAAAAAAAAAAAAAAAAA", charizard)
		memBuffer[i]=charizard;
	}
	memory_stat_total = MemStat {memBuffer[0],memBuffer[1],memBuffer[2],memBuffer[3],memBuffer[4],memBuffer[5],memBuffer[6],memBuffer[7],memBuffer[8],memBuffer[9],memBuffer[10],memBuffer[11],0}
	
	//fmt.Println(memBuffer)
	fmt.Println("Memory Statistics : ", memory_stat)
	fmt.Println("Memory Statistics (total) : ", memory_stat_total)
	fmt.Println("\n")
	
	
	/*--Reading the right CPU metric file--*/
	//fmt.Println(index)
	path="/sys/fs/cgroup/cpuacct/lxc/"+id[index]+"/cpuacct.stat"
	//fmt.Println(path)
	bs, err = ioutil.ReadFile(path)
    if err != nil {
        return
    }
    str = string(bs)
    fmt.Println(str)
    
	for i:=0; i<2; i++ {
		//fmt.Println(strings.Split(strings.Split(str,"\n")[i]," ")[1])
		charizard, _ := strconv.ParseInt(strings.Split(strings.Split(str,"\n")[i]," ")[1], 0, 64)
		//fmt.Println("CHAAAAAAAAAAAAAAAAAA", charizard)
		cpuBuffer[i]=charizard;
	}
	cpu_stat = CpuStat {cpuBuffer[0],cpuBuffer[1]}
    
	fmt.Println("PIKAAAAAAAAAAAAAAAAAAAAAA")	
    
    //fmt.Println(cpuBuffer)
    fmt.Println("CPU Statistics : ", cpu_stat)
    
}

