package clDebug

import "fmt"

func Info(_str string, _args ...interface{}) {
	if _args != nil && len(_args) > 0 {
		_str = fmt.Sprintf(_str, _args...)
	}
	fmt.Printf("[INFO] %v\n", _str)
}

func Debug(_str string, _args ...interface{}) {
	if _args != nil && len(_args) > 0 {
		_str = fmt.Sprintf(_str, _args...)
	}
	fmt.Printf("[DEBUG] %v\n", _str)
}

func Err(_str string, _args ...interface{}) {
	if _args != nil && len(_args) > 0 {
		_str = fmt.Sprintf(_str, _args...)
	}
	fmt.Printf("[ERR] %v\n", _str)
}

