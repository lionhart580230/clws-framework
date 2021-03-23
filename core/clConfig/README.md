# INIT配置文件读取库

### 引入库文件
~~~
import "ciaolan/config"
~~~


### 生成config对象
~~~
myConfig := config.New("temp.ini", 300)
~~~
如上代码, 即可得到config实例化对象, 该对象绑定temp.ini的配置文件。
后面的300是配置文件重载时间, 如果=0则不启动重载机制。
300为每300秒检查一次配置文件的变动，并重新载入。直接作用于变量


~~~
var floatConfig = float32(0)

myConfig.GetFloat32Config("config", "testFloat32", 0, &floatConfig)

fmt.Printf("获取到配置值: %0.2f\n", floatConfig)
~~~

如上代码, 可以轻松从init中读取数据, 由于是传入变量指针, 所以当设置重载时间的时候
配置文件有任何变动，都会实时更新变量, 直接作用于程序,不需要再编写额外的代码

~~~

var sectionConfig = make(map[string] string)

myConfig.GetSectionConfig("config", &sectionConfig)

~~~

如上代码将区段[config]到区段结尾的所有key=value结构全部载入到sectionConfig中
可以使用for进行遍历


### 所有接口

~~~

GetStrConfig        获取字符串型的配置

GetFloat32Config    获取float32型的配置

GetFloat64Config    获取float64型的配置

GetInt32Config      获取int32型的配置

GetUint32Config     获取uint32型的配置

GetInt64Config      获取int64型的配置

GetUint64Config     获取uint64型的配置

GetBoolConfig       获取bool类型的配置

GetFullSection      获取整个区段的配置

~~~


## 配置文件结构

~~~
[config]

# 测试字符串型
teststring=FUCKFUCK
# 测试浮点数
testfloat=10.05
# 测试有符号整数
testint=-10
# 测试无符号整数
testuint=10
# 测试bool
testbool=false

// 我也是一个注释

/*

 这
 里
 到
 这
 里
 都
 是
 一
 个
 注
 释

*/

[config2]
# 另一个区段的开始

~~~