package clmysql

/**
 *	高性能数据库封装类
 *
 * 
 */
import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)


var clconnections ClmysqlConnections
var cacheMgr *DbCacheMgr
var curdbname string


func init() {

	// 初始化连接池
	clconnections.master = make(map[string] *clmysqlStatus)
	clconnections.slaver = make(map[string] *clmysqlStatus)
	cacheMgr = NewCacheMgr()
	DBPointerPool = make(map[string] *DBPointer)

	go clconnections.StartToCheck()
}


// 新增主库
func AddMaster(key string, dbHost string, dbUser string, dbPass string, dbName string, limit uint16) error {
	return clconnections.AddMaster(key, dbHost, dbUser, dbPass, dbName, limit)
}

// 新增从库
func AddSlaver(key string, dbHost string, dbUser string, dbPass string, dbName string, readOnly bool, limit uint16) error {
	return clconnections.AddSlaver(key, dbHost, dbUser, dbPass, dbName, readOnly, limit)
}


// 内部查询用
func query(sqlstr string, curdb *sql.DB) ([]map[string] string , error) {
	row, err := curdb.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	columns, _ := row.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	// 将values的内存地址保存在scanArgs中
	for i := range values {
		scanArgs[i] = &values[i]
	}
	result := make([] map[string]string, 0)
	for row.Next() {
		records := make(map[string] string)
		row.Scan(scanArgs...)		// 获取扫描后的数组
		for i, col := range values {
			if col == nil {
				continue
			}

			records[columns[i]] = string(col.([] byte))

		}
		result = append(result, records)
	}
	return result, nil
}


// 内部查询用
func queryTx(sqlstr string, tx *sql.Tx) ([]map[string] string, error) {
	row, err := tx.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	columns, _ := row.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	// 将values的内存地址保存在scanArgs中
	for i := range values {
		scanArgs[i] = &values[i]
	}
	result := make([] map[string]string, 0)
	for row.Next() {
		records := make(map[string] string)
		row.Scan(scanArgs...)		// 获取扫描后的数组
		for i, col := range values {
			if col == nil {
				continue
			}

			records[columns[i]] = string(col.([] byte))

		}
		result = append(result, records)
	}
	return result, nil
}


/**
 * 数据库查询
 * 支持自定义缓存，指定缓存时间, 不指定或者设置为0 为不缓存
 * 可设置主键，如果设置了主键，将以map类型返回
 * 如果没有设置主键，将以切片类型返回
 * 
 * @param {[type]} sqlstr string 需要查询的数据库语句
 * @param {[type]} cache  uint32 缓存存在时间,单位: 秒
 * @param {[type]} key    string 主键名称
 *
 * return
 * @1  结果集
 * @2  结果条数
 * @3  如果有错误发生,这里是错误内容,否则为nil
 */
func QueryByKey(sqlstr string, cache uint32, key string) (*DbResult, error) {

	cacheKey := string(md5.New().Sum([]byte(fmt.Sprintf("%v|%v", sqlstr, key))))
	// 检查缓存
	if cache > 0 {
		cacheMgr.DbLock.Lock()
		defer cacheMgr.DbLock.Unlock()

		resp := cacheMgr.GetCache(cacheKey)
		if resp != nil {
			return resp, nil
		}
	}

	dbObject := clconnections.SelectSlaver("")
	if dbObject == nil {
		dbObject = clconnections.SelectMaster("")
	}

	if dbObject == nil {
		return nil, errors.New("Connections Is OverFlow Limit!!\n")
	}

	rows, err := query(sqlstr, dbObject.db)
	dbObject.Close()
	if err != nil {
		return nil, err
	}


	var result DbResult
	result.MapResult = make(map[string] TdbResult)
	result.Length = uint32(len(rows))

	for _, val := range rows {
		result.MapResult[val[key]] = val
	}

	// 更新缓存
	if cache > 0 {
		cacheMgr.SetCache(cacheKey, &result, cache)
	}

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
func Query(sqlstr string, cache uint32) (*DbResult, error){

	cacheKey := string(md5.New().Sum([]byte(fmt.Sprintf("%v", sqlstr))))
	// 检查缓存
	if cache > 0 {
		resp := cacheMgr.GetCache(cacheKey)
		if resp != nil {
			return resp, nil
		}
	}

	dbObject := clconnections.SelectSlaver("")
	if dbObject == nil {
		dbObject = clconnections.SelectMaster("")
	}

	if dbObject == nil {
		return nil, errors.New("Connections Is OverFlow Limit!!\n")
	}

	//if dbObject.dbname != curdbname {
//		dbObject.dbname = curdbname
//		dbObject.db.Exec("USE "+curdbname)
//	}

	rows, err := query(sqlstr, dbObject.db)
	dbObject.Close()
	if err != nil {
		return  nil, err
	}
	var result DbResult
	result.ArrResult = make([] TdbResult, 0)
	result.Length = uint32(len(rows))

	for _, val := range rows {
		result.ArrResult = append(result.ArrResult, val)
	}

	// 更新缓存
	if cache > 0 {
		cacheMgr.SetCache(cacheKey, &result, cache)
	}

	return &result, nil
}

/**
 * 事务查询
 * 当查询出错的时候将会自动回滚
 * @param {[type]} sqlstr string [查询语句]
 * @param {[type]} tx     *Tx    [事务指针]
 */
func QueryTx (sqlstr string, dtf *DbTransform) (*DbResult, error) {

	row, err := dtf.tx.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	columns, _ := row.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	// 将values的内存地址保存在scanArgs中
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var result DbResult
	result.ArrResult = make([]TdbResult, 0)
	result.Length = uint32(0)
	for row.Next() {

		records := make(map[string] string, len(columns))
		row.Scan(scanArgs...)		// 获取扫描后的数组
		for i, col := range values {
			if col == nil {
				continue
			}
			records[columns[i]] = string(col.([] byte))
		}
		// 把所有结果装载到一个结构里面
		result.ArrResult = append(result.ArrResult, records)
		result.Length++
	}

	if len(result.ArrResult) == 0 {
		return nil, nil
	}

	return &result, nil
}


/**
 * 事务执行
 * 当执行错误的时候将会自动回滚
 * @param {[type]} sqlstr string [description]
 * @param {[type]} tx     *Tx    [description]
 */
func ExecTx(sqlStr string, dtf *DbTransform) (int64, error) {
	//sqlStr = strings.TrimSpace(sqlStr)
	//ss := strings.ToLower(sqlStr)
	//if strings.Contains(ss, "delete") {
	//	str_arr := []string{}
	//	for i:=3;i>0;i--{
	//		_, file, line, ok := runtime.Caller(i)
	//		if ok {
	//			str_arr = append(str_arr,fmt.Sprintf("%v|%v",file,line))
	//		}
	//	}
	//	cllog.Log.LogWarning("","%v 删除语句 %v",strings.Join(str_arr,"->"),sqlStr)
	//}

	res, err := dtf.tx.Exec(sqlStr)
	if err != nil {
		return 0, err
	}

	if strings.HasPrefix(strings.ToLower(sqlStr), "insert") {
		return res.LastInsertId()
	}
	return res.RowsAffected()
}

// 开启事务
func StartTrans() (*DbTransform, error) {
	return clconnections.StartTrans()
}

// 提交事务
func Commit(tx *DbTransform) bool {
	err := tx.Commit()
	if err != nil {
		return false
	}
	return true
}

// 回滚事务
func RollBack(tx *DbTransform) {
	tx.st.Close()
	tx.tx.Rollback()
}

/**
 * 执行语句
 */
func Exec(sqlStr string, args... interface{}) ( int64, error) {
	sqlStr = strings.TrimSpace(fmt.Sprintf(sqlStr, args...))
	//ss := strings.ToLower(sqlStr)
	//if strings.Contains(ss, "delete") {
	//	str_arr := []string{}
	//	for i:=3;i>0;i--{
	//		_, file, line, ok := runtime.Caller(i)
	//		if ok {
	//			str_arr = append(str_arr,fmt.Sprintf("%v|%v",file,line))
	//		}
	//	}
	//	cllog.Log.LogWarning("","%v 删除语句 %v",strings.Join(str_arr,"->"),sqlStr)
	//}
	dbObject := clconnections.SelectMaster("")

	if dbObject == nil {
		return 0, errors.New("Connections Is OverFlow Limit!!\n")
	}

//	if dbObject.dbname != curdbname {//
//		dbObject.dbname = curdbname
//		dbObject.db.Exec("USE "+curdbname)
//	}


	res, err := dbObject.db.Exec(sqlStr)
	dbObject.Close()

	if err != nil {
		return 0, err
	}

	if strings.HasPrefix(strings.ToLower(sqlStr), "insert") {
		return res.LastInsertId()
	}
	return res.RowsAffected()
}


/**
   更新, 返回更新行数
 */
func Save(sqlStr string) (int64, error) {
	dbObject := clconnections.SelectMaster("")
	if dbObject == nil {
		return 0, errors.New("Connections Is OverFlow Limit!!\n")
	}
	res, err := dbObject.db.Exec(sqlStr)
	dbObject.Close()

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

/**
 * 添加新的数据
 * @param {[type]} tablename string                   [要添加数据的表名]
 * @param {[type]} data      map[string]interface{}   [要添加数据的数据结构]
 */
func AddData(tablename string, data map[string]interface{}) (int64, error) {
	var sqlStr = "INSERT INTO %v (%v) VALUES(%v)"
	var sqlField = ""
	var sqlValues = ""

	for _k, _v := range data {
		if sqlField != "" {
			sqlField += ","
		}

		if sqlValues != "" {
			sqlValues += ","
		}

		sqlField += fmt.Sprintf("`%v`", _k)
		sqlValues += fmt.Sprintf("'%v'", _v)
	}

	strSql := fmt.Sprintf(sqlStr, tablename, sqlField, sqlValues)
	return Exec(strSql)
}

/**
 * 事务添加新的数据
 * 当发生错误的时候将会自动回滚
 * @param {[type]} tablename string                   [要添加数据的表名]
 * @param {[type]} data      map[string]interface{}   [要添加数据的数据结构]
 */
func AddDataTx(tablename string, data map[string]interface{}, tx *DbTransform) (int64, error) {
	var sqlStr = "INSERT INTO %v (%v) VALUES(%v)"
	var sqlField = ""
	var sqlValues = ""

	for _k, _v := range data {
		if sqlField != "" {
			sqlField += ","
		}

		if sqlValues != "" {
			sqlValues += ","
		}

		sqlField += fmt.Sprintf("`%v`", _k)
		sqlValues += fmt.Sprintf("'%v'", _v)
	}

	strSql := fmt.Sprintf(sqlStr, tablename, sqlField, sqlValues)
	return ExecTx(strSql, tx)
}


/**
 * 获取指定数据库下的所有表名字
 * @param dbname string 获取的数据库名
 * @param contain string 表名包含字符串，为空则取全部表
 * return
 * @1 数据表数组
 * @2 数据库
 */
func GetTables(dbname string, contain string) ([]string, error) {

	if dbname != "" {
		ToggleDBName(dbname)
	}

	querySql := ""
	if contain == "" {
		querySql = "SHOW TABLES"
	} else {
		querySql = "SHOW TABLES LIKE '%"+contain+"%'"
	}

	res, err := Query(querySql, 0)
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
func HasTable(dbname string, tablename string) (bool) {
	if dbname == "" {
		dbname = curdbname
	}
	var tables, _ = GetTables(dbname, tablename)
	for _, val := range tables {
		if strings.EqualFold(val, tablename) {
			return true
		}
	}
	return false
}


// 切换数据库
func ToggleDBName(name string) {
	curdbname = name
	Exec("USE "+curdbname)
}

// 获取当前使用的数据库名称
func GetDBName() string {
	return curdbname
}

/**
 * 获取此连接下的所有可见数据库名字
 * return
 * @1 数据库列表数组
 * @2 数据库
 */
func GetDatabases() ([]string, error) {
	res, err := Query("show databases", 0)
	if err != nil {
		return nil, err
	}

	if res.Length == 0 {
		return nil, nil
	}

	databases := make([]string, res.Length)
	for i:=0; i < int(res.Length); i++ {
		databases[i] = res.ArrResult[i]["Database"]
	}
	return databases, nil
}

/**
 * 获取此连接下的所有可见数据库名字
 * return
 * @1 数据库列表数组
 * @2 数据库
 * param : prex  数据库前缀
 */
func GetDatabasesByPrex(prex string) (map[int]string, error) {
	res, err := Query("show databases", 0)
	if err != nil {
		return nil, err
	}

	if res.Length == 0 {
		return nil, nil
	}

	databases := make(map[int]string)
	k := 0
	for i:=0; i < int(res.Length); i++ {
		if !strings.HasPrefix(res.ArrResult[i]["Database"],prex) {
			continue
		}
		databases[i] = res.ArrResult[i]["Database"]
		k++
	}
	return databases, nil
}
