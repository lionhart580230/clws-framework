package clmysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)


type ClTranslate struct {
	tx *sql.Tx
	DBName string
}

// 查询事务
func (this *ClTranslate) QueryTx(sqlstr string, args... interface{}) (*DbResult, error) {

	if this.tx == nil {
		return nil, errors.New("错误: 事务指针为空")
	}

	if args != nil && len(args) > 0 {
		sqlstr = fmt.Sprintf(sqlstr, args...)
	}

	rows, err := queryTx(sqlstr, this.tx)
	if err != nil {
		return nil, err
	}
	var result DbResult
	result.ArrResult = make([] TdbResult, 0)
	result.Length = uint32(len(rows))

	for _, val := range rows {
		result.ArrResult = append(result.ArrResult, val)
	}

	return &result, nil
}



// 执行事务
func (this *ClTranslate)ExecTx(sqlstr string, args... interface{}) (int64, error) {
	//ss := strings.ToLower(sqlstr)
	//if strings.Contains(ss, "delete") {
	//	str_arr := []string{}
	//	for i:=5;i>=0;i--{
	//		_, file, line, ok := runtime.Caller(i)
	//		if ok {
	//			str_arr = append(str_arr,fmt.Sprintf("%v|%v",file,line))
	//		}
	//	}
	//	cllog.Log.LogWarning("","%v 删除语句 %v",strings.Join(str_arr,"->"),sqlstr)
	//}
	if this.tx == nil {
		return 0, errors.New("错误: 事务指针为 nil pointer")
	}

	if args != nil && len(args) != 0 {
		sqlstr = fmt.Sprintf(sqlstr, args...)
	}

	res, err := this.tx.Exec(sqlstr)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("%v, SQL:%v", err, sqlstr))
	}

	if strings.HasPrefix(strings.ToLower(sqlstr), "insert") {
		return res.LastInsertId()
	}

	return res.RowsAffected()
}



// 提交事务
func (this *ClTranslate) Commit() error {
	return this.tx.Commit()
}


// 回滚事务
func (this *ClTranslate) Rollback() error {
	return this.tx.Rollback()
}


// 使用DBPointer进行构建器创建
func (this *ClTranslate) NewBuilder() *SqlBuider {

	sqlbuild := SqlBuider{
		dbTx: this,
		dbType: 2,
		dbname: this.DBName,
	}
	return &sqlbuild
}


/**
  是否存在这个表格
 */
func (this *ClTranslate) HasTable(tablename string) (bool) {

	var tables, _ = this.GetTables(tablename)
	for _, val := range tables {
		if strings.EqualFold(val, tablename) {
			return true
		}
	}
	return false
}


/**
 * 获取指定数据库下的所有表名字
 * @param dbname string 获取的数据库名
 * @param contain string 表名包含字符串，为空则取全部表
 * return
 * @1 数据表数组
 * @2 数据库
 */
func (this *ClTranslate)GetTables(contain string) ([]string, error) {

	querySql := ""
	if contain == "" {
		querySql = "SHOW TABLES"
	} else {
		querySql = "SHOW TABLES LIKE '%"+contain+"%'"
	}

	res, err := this.QueryTx(querySql)
	if err != nil {
		return []string{}, err
	}

	if res.Length == 0 {
		return nil, nil
	}

	tables := make([]string, res.Length)
	for i:=0;i<int(res.Length);i++ {
		if contain != "" {
			tables[i] = res.ArrResult[i]["Tables_in_"+this.DBName+" (%"+contain+"%)"]
		}else{
			tables[i] = res.ArrResult[i]["Tables_in_"+this.DBName]
		}
	}

	return tables, nil
}
