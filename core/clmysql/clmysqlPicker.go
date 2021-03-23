package clmysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)


type DBPointer struct {
	conn *sql.DB
	dbname string
	isconnect bool
	lastErr string
	lastSql string
	key string
	Dbname string	//数据库名字,外部调用
	lastupdate uint32
}

func (this DBPointer)IsUsefull() bool {
	return this.isconnect
}

var DBPointerPool map[string] *DBPointer
var DBPointerPoolLocker sync.Mutex


// 通过连线配置返回一个连线结构
// @param host string 连线ip
// @param user string 用户名
// @param pass string 密码
// @param dbname string 数据库名字
func NewDB(host string, user string, pass string, dbname string) (*DBPointer, error) {
	//DBPointerPoolLocker.Lock()
	//defer DBPointerPoolLocker.Unlock()
	//
	//if val, find := DBPointerPool[host+user+dbname]; find {
	//	if nil == val.conn.Ping() {
	//		DBPointerPool[host+user+dbname].lastupdate = uint32(time.Now().Unix())
	//		return val, nil
	//	} else {
	//		val.Close()
	//	}
	//}

	dbp, err := sql.Open("mysql", user+":"+pass+"@tcp("+host+")/"+dbname+"?charset=utf8")
	if err != nil {
		return nil, err
	}

	err = dbp.Ping()
	if err != nil {
		return nil, err
	}

	dbp.SetMaxOpenConns(30)
	//dbp.SetMaxOpenConns(0) // TODO 暂时增加到60
	dbp.SetMaxIdleConns(10)
	dbp.SetConnMaxLifetime(4*time.Hour)

	object := DBPointer{
		conn: dbp,
		dbname: dbname,
		isconnect: true,
		key: host+user+dbname,
		Dbname:dbname,
		lastupdate: uint32(time.Now().Unix()),
	}

	//fmt.Printf(">> %v 连接成功!!\n", user+":"+pass+"@tcp("+host+")/"+dbname+"?charset=utf8")
	//DBPointerPool[host+user+dbname] = &object
	return &object, nil
}


// 设置最大空闲连接数,0为不限制
func (this *DBPointer) SetMaxIdleConns(conns int) {
	this.conn.SetMaxIdleConns(conns)
}


// 设置数据库
func (this *DBPointer) UseDB(dbname string) {
	this.Exec("USE %v", dbname)
	this.Dbname = dbname
}

func NewDBSimple(host string, user string, pass string, dbname string) (*DBPointer) {

	dbp, err := sql.Open("mysql", user+":"+pass+"@tcp("+host+")/"+dbname+"?charset=utf8")
	if err != nil {
		return nil
	}

	err = dbp.Ping()
	if err != nil {
		return nil
	}

	dbp.SetMaxOpenConns(0)
	dbp.SetMaxIdleConns(-1)
	//dbp.SetConnMaxLifetime(10*time.Minute)

	object := DBPointer{
		conn: dbp,
		dbname: dbname,
		isconnect: true,
		key: host+user+dbname,
		Dbname:dbname,
		lastupdate: uint32(time.Now().Unix()),
	}

	return &object
}


// 查看数据库状态..
func (this *DBPointer) UpdateStatus() bool {
	if this.conn == nil {
		return false
	}
	err := this.conn.Ping()
	if err != nil {
		return false
	}
	return true
	//DBPointerPoolLocker.Lock()
	//defer DBPointerPoolLocker.Unlock()
	//
	//DBPointerPool[this.key].lastupdate = uint32(time.Now().Unix())
}

// 操作完毕，释放资源
func (this *DBPointer) Close() {
	if this.conn != nil {
		this.conn.Close()
	}
}


func (this *DBPointer) StartTrans() (*ClTranslate, error) {
	myTx, err := this.conn.Begin()
	if err != nil {
		return nil, err
	}
	return &ClTranslate {
		tx: myTx,
		DBName: this.Dbname,
	}, nil
}

/**
 * 普通的数据库查询
 * 不支持自定义主键，但返回的是slice类型，更加精简
 * @param {[type]} sqlstr string 需要查询的数据库语句
 * @param {[type]} cache  int    缓存存在的时间, 单位: 秒
 *
 * return
 * @1  结果集
 * @2  结果条数
 * @3  如果有错误发生,这里是错误内容,否则为nil
 */
func (this *DBPointer) Query(sqlstr string, args... interface{}) (*DbResult, error){

	lastSql := sqlstr
	if args != nil && len(args) != 0 {
		lastSql = fmt.Sprintf(sqlstr, args...)
	}

	this.lastSql = lastSql

	rows, err := query(lastSql, this.conn)
	if err != nil {
		this.lastSql = fmt.Sprintf(sqlstr, args...)
		this.lastErr = fmt.Sprintf("查询失败: %v", err)
		return nil, err
	}
	var result DbResult
	result.ArrResult = make([] TdbResult, 0)
	for _, val := range rows {
		result.ArrResult = append(result.ArrResult, val)
	}
	result.Length = uint32(len(result.ArrResult))

	return &result, nil
}


/**
 * 普通的数据库查询
 * 不支持自定义主键，但返回的是slice类型，更加精简
 * @param {[type]} sqlstr string 需要查询的数据库语句
 * @param {[type]} cache  int    缓存存在的时间, 单位: 秒
 *
 * return
 * @1  结果集
 * @2  结果条数
 * @3  如果有错误发生,这里是错误内容,否则为nil
 */
func (this *DBPointer) GetLastSql() string {
	return this.lastSql
}
/**
 * 普通的数据库查询
 * 不支持自定义主键，但返回的是slice类型，更加精简
 * @param {[type]} sqlstr string 需要查询的数据库语句
 * @param {[type]} cache  int    缓存存在的时间, 单位: 秒
 *
 * return
 * @1  结果集
 * @2  结果条数
 * @3  如果有错误发生,这里是错误内容,否则为nil
 */
func (this *DBPointer) QueryByKey(sqlstr string, key string, args... interface{}) (*DbResult, error){

	lastSql := sqlstr
	if args != nil && len(args) != 0 {
		lastSql = fmt.Sprintf(sqlstr, args...)
	}

	this.lastSql = lastSql

	rows, err := query(lastSql, this.conn)
	if err != nil {
		this.lastErr = fmt.Sprintf("查询失败: %v", err)
		return nil, err
	}
	var result DbResult
	result.MapResult = make(map[string] TdbResult)
	result.Length = uint32(len(rows))

	for _, val := range rows {
		result.MapResult[val[key]] = val
	}

	return &result, nil
}



// 执行
func (this *DBPointer)Exec(sqlstr string, args... interface{}) (int64, error) {
//UPDATE game_user_msg SET `status` = '1' WHERE guid = 794585 and accid = 53998;

	lastSql := sqlstr
	if args != nil && len(args) != 0 {
		lastSql = fmt.Sprintf(sqlstr, args...)
	}
	if lastSql == "" {
		return 0, errors.New("SQL语句为空")
	}

	this.lastSql = lastSql
	if this.conn == nil {
		this.lastErr = "错误: SQL连线指针为空"
		return 0, errors.New("错误: SQL连线指针为nil pointer")
	}

	res, err := this.conn.Exec(lastSql)
	if err != nil {
		this.lastErr = fmt.Sprintf("执行失败! ERR:%v", err)
		this.lastSql = fmt.Sprintf(sqlstr, args...)
		return 0, err
	}

	if strings.HasPrefix(strings.ToLower(sqlstr), "insert") {
		return res.LastInsertId()
	}

	return res.RowsAffected()
}



// 获取数据库名称
func (this *DBPointer) GetDBName() string {
	return this.dbname
}


/**
 * 获取指定数据库下的所有表名字
 * @param dbname string 获取的数据库名
 * @param contain string 表名包含字符串，为空则取全部表
 * return
 * @1 数据表数组
 * @2 数据库
 */
func (this *DBPointer)GetTables(contain string) ([]string, error) {

	dbname := this.GetDBName()

	querySql := ""
	if contain == "" {
		querySql = "SHOW TABLES"
	} else {
		querySql = "SHOW TABLES LIKE '%"+contain+"%'"
	}

	res, err := this.Query(querySql)
	if err != nil {
		return []string{}, err
	}

	if res.Length == 0 {
		return nil, nil
	}

	tables := make([]string, res.Length)
	for i:=0;i<int(res.Length);i++ {
		if contain != "" {
			tables[i] = res.ArrResult[i]["Tables_in_"+dbname+" (%"+contain+"%)"]
		}else{
			tables[i] = res.ArrResult[i]["Tables_in_"+dbname]
		}
	}

	return tables, nil
}


/**
  是否存在这个表格
 */
func (this *DBPointer) HasTable(tablename string) (bool) {

	var tables, _ = this.GetTables(tablename)
	for _, val := range tables {
		if strings.EqualFold(val, tablename) {
			return true
		}
	}
	return false
}
