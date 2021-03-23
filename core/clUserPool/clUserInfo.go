package clUserPool

import (
	"clws-framework/core/clCommon"
)



// 设置用户登录状态
//@param _isLogin 是否登录
func (this *ClNetUserInfo) SetLogin(_isLogin bool) {
	this.IsLogin = _isLogin
}



// 设置扩展参数
//@param _key 参数key
//@param _val 参数数值
func (this *ClNetUserInfo) SetParam(_key string, _val string) {
	this.ParamLock.Lock()
	defer this.ParamLock.Unlock()

	this.Params[_key] = _val
}


// 获取Int32扩展参数
func (this *ClNetUserInfo) GetInt32(_key string, _def int32) int32 {
	this.ParamLock.RLock()
	defer this.ParamLock.RUnlock()

	v, exist := this.Params[_key]
	if !exist {
		return _def
	}
	return clCommon.Int32(v, _def)
}


// 获取Uint32扩展参数
func (this *ClNetUserInfo) GetUint32(_key string, _def uint32) uint32 {
	this.ParamLock.RLock()
	defer this.ParamLock.RUnlock()

	v, exist := this.Params[_key]
	if !exist {
		return _def
	}
	return clCommon.Uint32(v, _def)
}


// 获取int64扩展参数
func (this *ClNetUserInfo) GetInt64(_key string, _def int64) int64 {
	this.ParamLock.RLock()
	defer this.ParamLock.RUnlock()

	v, exist := this.Params[_key]
	if !exist {
		return _def
	}

	return clCommon.Int64(v, _def)
}


// 获取int64扩展参数
func (this *ClNetUserInfo) GetUInt64(_key string, _def uint64) uint64 {
	this.ParamLock.RLock()
	defer this.ParamLock.RUnlock()

	v, exist := this.Params[_key]
	if !exist {
		return _def
	}
	return clCommon.Uint64(v, _def)
}


// 获取float32扩展参数
func (this *ClNetUserInfo) GetFloat32(_key string, _def float32) float32 {
	this.ParamLock.RLock()
	defer this.ParamLock.RUnlock()

	v, exist := this.Params[_key]
	if !exist {
		return _def
	}
	return clCommon.Float32(v, _def)
}


// 获取float64扩展参数
func (this *ClNetUserInfo) GetFloat64(_key string, _def float64) float64 {
	this.ParamLock.RLock()
	defer this.ParamLock.RUnlock()

	v, exist := this.Params[_key]
	if !exist {
		return _def
	}

	return clCommon.Float64(v, _def)
}


// 获取bool
func (this *ClNetUserInfo) GetBool(_key string, _def bool) bool {
	this.ParamLock.RLock()
	defer this.ParamLock.RUnlock()

	v, exist := this.Params[_key]
	if !exist {
		return _def
	}
	return clCommon.Bool(v)
}



