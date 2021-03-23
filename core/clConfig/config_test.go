package clConfig

import (
	"testing"
	"fmt"
)

func TestConfig_GetBoolConfig(T *testing.T) {
	myConfig := New("temp.ini", 1)
	if myConfig == nil {
		fmt.Printf("载入失败!\n")
		T.Error("载入失败!")
		return
	}
	var strConfig = ""
	var float32Config = float32(0.0)
	var float64Config = float64(0.0)
	var int32Config = int32(0)
	var int64Config = int64(0)
	var boolConfig = false
	var section = make(map[string] string)
	myConfig.GetStrConfig("config", "teststring", "", &strConfig)
	myConfig.GetFloat32Config("config", "testfloat", 0, &float32Config)
	myConfig.GetFloat64Config("config", "testfloat", 0, &float64Config)
	myConfig.GetInt32Config("config", "testint", 0, &int32Config)
	myConfig.GetInt64Config("config", "testint", 0, &int64Config)
	myConfig.GetBoolConfig("config", "testbool", false, &boolConfig)
	myConfig.GetFullSection("config", &section )


	fmt.Printf("strconfig: %v\n", strConfig)
	fmt.Printf("floatconfig: %v - %v\n", float32Config, float64Config)
	fmt.Printf("intconfig: %v - %v\n", int32Config, int64Config)
	fmt.Printf("boolconfig: %v\n", boolConfig)
	fmt.Printf("sectionconfig: %v\n", section)
}