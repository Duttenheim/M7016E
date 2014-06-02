package profiling

import (
	"github.com/beevik/ntp"
	"time"
	"fmt"
)

//------------------------------------------------------------------------------
/**
*/
type Profiler struct {
	lastFrame	time.Time
	diff		int64
	Name		string
}

// store profilers in global map, this will NOT be exported so it wont be available from outside this package
var profilers 	map[string]*Profiler
var NTPAddress 	string

//------------------------------------------------------------------------------
/**
*/
func Setup() {
	profilers 	= make(map[string]*Profiler)
	NTPAddress 	= "ntp.ltu.se"
}

//------------------------------------------------------------------------------
/**
*/
func RegisterProfiler(name string) (*Profiler, error) {
	if _, ok := profilers[name]; !ok {
		profiler := new(Profiler)
		profiler.Name = name
		profilers[name] = profiler
		return profiler, nil
	}
	return nil, fmt.Errorf("Profiler with name '%s' already registered!", name)
}

//------------------------------------------------------------------------------
/**
*/
func UnregisterProfiler(name string) error {
	if _, ok := profilers[name]; !ok {
		delete(profilers, name)
		return nil
	}	
	return fmt.Errorf("Profiler with name '%s' is not registered!", name)
}

//------------------------------------------------------------------------------
/**
*/
func GetProfiler(name string) (*Profiler, error) {
	if val, ok := profilers[name]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("No profiler named '%s' registered!", name)
}

//------------------------------------------------------------------------------
/**
	Get time from LTU NTP server
*/
func (profiler *Profiler) Start() error {
	var err error
	profiler.lastFrame, err = ntp.Time(NTPAddress)
	if err != nil {
		return err
	}
	profiler.diff = -1
	return nil
}

//------------------------------------------------------------------------------
/**
	Calculates the difference in time since start
*/
func (profiler *Profiler) Stop() error {	
	time, err := ntp.Time(NTPAddress)
	if err != nil {
		return err
	}	
	difference := time.Sub(profiler.lastFrame)
	profiler.diff = difference.Nanoseconds()
	return nil
}

//------------------------------------------------------------------------------
/**
	Get time difference in nanoseconds
*/
func (profiler *Profiler) GetNanoseconds() int64 {
	return profiler.diff
}

//------------------------------------------------------------------------------
/**
	Get time difference in microseconds
*/
func (profiler *Profiler) GetMicroseconds() int64 {
	return profiler.diff / 1000
}

//------------------------------------------------------------------------------
/**
	Get time difference in milliseconds
*/
func (profiler *Profiler) GetMilliseconds() int64 {
	return profiler.GetMicroseconds() / 1000
}

//------------------------------------------------------------------------------
/**
*/
func (profiler *Profiler) Print() {
	fmt.Printf("Timer '%s'\n[\n    Time in nanoseconds: %dns\n    Time in microseconds: %dus\n    Time in milliseconds: %dms\n]\n",
		profiler.Name,
		profiler.GetNanoseconds(), 
		profiler.GetMicroseconds(),
		profiler.GetMilliseconds())
}