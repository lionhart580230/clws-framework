package clConfig

import (
	"time"
	"sync"
)

type sectionType map[string] string

type Config struct {
	config map[string] sectionType
	fileName string
	autoLoad time.Duration
	lock sync.RWMutex
	stringMap map[/*section:key*/ string] autoLoadString
	int64Map map[/*section:key*/ string] autoLoadInt64
	int32Map map[/*section:key*/ string] autoLoadInt32
	uint64Map map[/*section:key*/ string] autoLoadUint64
	uint32Map map[/*section:key*/ string] autoLoadUint32
	float32Map map[/*section:key*/ string] autoLoadFloat32
	float64Map map[/*section:key*/ string] autoLoadFloat64
	boolMap map[/*section:key*/ string] autoLoadBool
	sectionMap map[/*section*/ string] autoLoadSection
	ArrMap map[/*section*/ string] autoLoadArr
}

type autoLoadString struct {
	section string
	key string
	def string
	value *string
}

type autoLoadInt64 struct {
	section string
	key string
	def int64
	value *int64
}

type autoLoadInt32 struct {
	section string
	key string
	def int32
	value *int32
}

type autoLoadUint64 struct {
	section string
	key string
	def uint64
	value *uint64
}

type autoLoadUint32 struct {
	section string
	key string
	def uint32
	value *uint32
}

type autoLoadBool struct {
	section string
	key string
	def bool
	value *bool
}

type autoLoadFloat32 struct {
	section string
	key string
	def float32
	value *float32
}

type autoLoadFloat64 struct {
	section string
	key string
	def float64
	value *float64
}

type autoLoadSection struct {
	section string
	value *map[string] string
}

type autoLoadArr struct {
	section string
	value *[] string
}