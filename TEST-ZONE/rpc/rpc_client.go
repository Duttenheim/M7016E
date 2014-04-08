
/**
* ArithClient
 */

package main

import (
	"net/rpc"
	"fmt"
	"log"
	"os"
)

type CpuStat struct {
	User int64
	System int64
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "server")
		os.Exit(1)
	}
	serverAddress := os.Args[1]

	client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	
	
	// Synchronous call
	args := 0
	var cpu CpuStat
	var mem MemStat
	var memTot MemStat
	
	err = client.Call("CpuStat.SendCpu", args, &cpu)
	fmt.Println("CPU Statistics : ", cpu)
	
	err = client.Call("MemStat.SendMem", args, &mem)
	fmt.Println("Mem Statistics : ", mem)
	
	err = client.Call("MemStat.SendMemTot", args, &memTot)
	fmt.Println("Mem Statistics (total) : ", memTot)
	
}
