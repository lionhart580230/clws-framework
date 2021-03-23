package clmysql

import (
	"fmt"
	"errors"
	"reflect"
	"strings"
)

/**
   数据库语句生成器
 */
type SqlBuider struct {
	tablename string
	dbname string
	whereStr string
	fieldStr string
	updateData map[string] string
	orders string
	limit string
	group string
	having string
	expire uint32
	finalSql string

	addColumns [] MySqlColumns	// 保存的字段
	removeColumns [] string 	// 待删除的字段
	addIndexs []string			// 待添加的索引
	removeIndexs []string		// 待删除的索引

	lastColumns []MySqlColumns
	lastIndexs []string
	primaryKeys string

	lastTable string
	lastTableId string

	unionalls []string
	unions []string

	duplicateKey []string

	dbType uint32
	dbPointer *DBPointer
	dbTx *ClTranslate
}

type MySqlColumns struct {
	name string					// 字段名称
	typename string				// 字段类型
	null bool					// 是否为空
	defaults string				// 默认值
	autoInc bool				// 是否自动递增
	comment string				// 备注
}



func NewBuilder() *SqlBuider {

	sqlbuild := SqlBuider{}
	return &sqlbuild
}


// 使用DBPointer进行构建器创建
func (this *DBPointer) NewBuilder() *SqlBuider {

	sqlbuild := SqlBuider{
		dbPointer: this,
		dbType: 1,
		dbname: this.Dbname,
	}
	return &sqlbuild
}

/**
   设置表格名称
   @param tablename string  要设置的表格名称
 */
func (this *SqlBuider) Table(tablename string) (*SqlBuider) {
	this.tablename = tablename
	this.whereStr = ""
	this.fieldStr = ""
	this.updateData = make(map[string] string)
	this.orders = ""
	this.limit = ""
	this.group = ""
	this.having = ""
	this.expire = 0
	this.finalSql = ""

	this.addColumns = make([] MySqlColumns, 0)	// 保存的字段
	this.removeColumns = make([] string, 0) 	// 待删除的字段
	this.addIndexs = make([]string, 0)			// 待添加的索引
	this.removeIndexs = make([]string, 0)		// 待删除的索引

	this.lastColumns = make([]MySqlColumns, 0)
	this.lastIndexs = make([]string, 0)
	this.primaryKeys = ""

	this.lastTable = ""
	this.lastTableId = ""

	this.unions = make([]string, 0)
	this.unionalls = make([]string, 0)

	this.duplicateKey = make([]string, 0)

	return this
}


/**
   设置WHERE条件
   @param wherestr string WHERE条件文本
 */
func (this *SqlBuider) Where (wherestr string, args... interface{}) (*SqlBuider){
	this.whereStr = fmt.Sprintf(wherestr, args...)
	return this
}


/**
	设置重复要更新的key列表
 */
func (this *SqlBuider) OnDuplicateKey(keys []string) (*SqlBuider) {
	this.duplicateKey = keys
	return this
}


/**
   设置要查询的Field
   @param fiedlStr string FIELD字段列表
 */
func (this *SqlBuider) Field(fieldStr string) (*SqlBuider) {
	this.fieldStr = fieldStr
	return this
}


/**
   设置分组
   @param group string 分组内容
 */
func (this *SqlBuider) Group(groupStr string) (*SqlBuider) {
	this.group = groupStr
	return this
}

/**
   设置排序方式
   @param orders string 排序内容
 */
func (this *SqlBuider)Order(orders string) (*SqlBuider) {
	this.orders = orders
	return this
}

/**
   设置Cache 缓存时间
   @param expire int32 缓存有效期
 */
func (this *SqlBuider) Cache(expire uint32) (*SqlBuider) {
	this.expire = expire
	return this
}


/**
   设置LIMIT限制
   @param min int32 设置limit的最小值
   @param count int32 设置limit的数量
 */
func (this *SqlBuider) Limit(min int32, count int32) (*SqlBuider) {
	this.limit = fmt.Sprintf(" LIMIT %v, %v", min, count)
	return this
}


/**
	设置DB名字
	@param dbname string 设置DB名字
 */
func (this *SqlBuider) DB(dbname string) (*SqlBuider) {
	this.dbname = dbname
	return this
}

/**
   查询语句
 */
func (this *SqlBuider) Select() (*DbResult, error) {

	if this.tablename == "" {
		return nil, errors.New("EMPTY TABLE NAME")
	}

	if this.whereStr == "" {
		this.whereStr = "1"
	}

	if this.fieldStr == "" {
		this.fieldStr = "*"
	}

	var extraSql = ""
	if this.group != "" {
		extraSql += " GROUP BY "+this.group
	}


	var FinallySql = ""
	if len(this.unionalls) > 0 {
		// 使用一般联表查询
		for _, sub := range this.unionalls {
			if FinallySql == "" {
				FinallySql = fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v UNION ALL ", this.fieldStr, this.tablename, this.whereStr, extraSql)
			}else{
				FinallySql += " UNION ALL "
			}
			FinallySql += fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v", this.fieldStr, sub, this.whereStr, extraSql)
		}
	} else if len(this.unions) > 0 {
		// 使用强制联表查询
		for _, sub := range this.unions {
			if FinallySql == "" {
				FinallySql = fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v UNION ", this.fieldStr, this.tablename, this.whereStr, extraSql)
			}else{
				FinallySql += " UNION "
			}
			FinallySql += fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v", this.fieldStr, sub, this.whereStr, extraSql)
		}
	}

	if this.orders != "" {
		extraSql += " ORDER BY "+this.orders
	}

	if this.limit != "" {
		extraSql += this.limit
	}

	if FinallySql == "" {
		FinallySql = fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v", this.fieldStr, this.tablename, this.whereStr, extraSql)
	} else {
		FinallySql = fmt.Sprintf("SELECT %v FROM ( %v ) temp %v", this.fieldStr, FinallySql, extraSql)
	}
	this.finalSql = FinallySql

	var resp *DbResult = nil
	var err error


	switch this.dbType {
	case 0:		// 正常
		resp, err = Query(FinallySql, this.expire)
	case 1:		// Picker
		resp, err = this.dbPointer.Query(FinallySql)
	case 2:		// 事务
		resp, err = this.dbTx.QueryTx(FinallySql)
	}

	return resp, err
}

/**
    查询语句
	获取指定索引处的数据
	@param idx int 索引id
 */
func (this *SqlBuider) Find(idx int) (TdbResult, error) {

	if this.tablename == "" {
		return nil, errors.New("EMPTY TABLE NAME")
	}

	if this.whereStr == "" {
		this.whereStr = "1"
	}

	if this.fieldStr == "" {
		this.fieldStr = "*"
	}

	var extraSql = ""
	if this.group != "" {
		extraSql += "GROUP BY "+this.group
	}

	var FinallySql = ""
	if len(this.unionalls) > 0 {
		// 使用一般联表查询
		for _, sub := range this.unionalls {
			if FinallySql == "" {
				FinallySql = fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v UNION ALL ", this.fieldStr, this.tablename, this.whereStr, extraSql)
			}else{
				FinallySql += " UNION ALL "
			}
			FinallySql += fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v ", this.fieldStr, sub, this.whereStr, extraSql)
		}
	} else if len(this.unions) > 0 {
		// 使用强制联表查询
		for _, sub := range this.unions {
			if FinallySql == "" {
				FinallySql = fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v UNION ", this.fieldStr, this.tablename, this.whereStr, extraSql)
			}else{
				FinallySql += " UNION "
			}
			FinallySql += fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v ", this.fieldStr, sub, this.whereStr, extraSql)
		}
	}

	if this.orders != "" {
		extraSql += " ORDER BY "+this.orders
	}

	if this.limit != "" {
		extraSql += this.limit
	}

	if FinallySql == "" {
		FinallySql = fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v", this.fieldStr, this.tablename, this.whereStr, extraSql)
	} else {
		FinallySql = fmt.Sprintf("SELECT %v FROM ( %v ) temp %v", this.fieldStr, FinallySql, extraSql)
	}
	this.finalSql = FinallySql

	var resp *DbResult = nil
	var err error

	switch this.dbType {
	case 0:		// 正常
		resp, err = Query(FinallySql, this.expire)
	case 1:		// Picker
		resp, err = this.dbPointer.Query(FinallySql)
	case 2:		// 事务
		resp, err = this.dbTx.QueryTx(FinallySql)
	}

	if err != nil {
		return nil, err
	}

	if resp == nil || resp.Length == 0{
		return nil, nil
	}

	return resp.ArrResult[0], nil
}

/**
    查询语句
	获取指定索引处的数据
	@param idx int 索引id
 */
func (this *SqlBuider) Count() (int32, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	if this.whereStr == "" {
		this.whereStr = "1"
	}

	if this.fieldStr == "" {
		this.fieldStr = "*"
	}

	var extraSql = ""
	if this.group != "" {
		extraSql += "GROUP BY "+this.group
	}

	var FinallySql = ""
	if len(this.unionalls) > 0 {
		for _, val := range this.unionalls {
			if FinallySql == "" {
				FinallySql = fmt.Sprintf("SELECT COUNT(*) as t_count FROM %v WHERE %v UNION ALL ", this.tablename, this.whereStr)
			} else {
				FinallySql += " UNION ALL "
			}
			FinallySql += fmt.Sprintf("SELECT COUNT(*) as t_count FROM %v WHERE %v ", val, this.whereStr)
		}
	} else if len(this.unions) > 0 {
		for _, val := range this.unions {
			if FinallySql == "" {
				FinallySql = fmt.Sprintf("SELECT COUNT(*) as t_count FROM %v WHERE %v UNION ALL ", this.tablename, this.whereStr)
			} else {
				FinallySql += " UNION ALL "
			}
			FinallySql += fmt.Sprintf("SELECT COUNT(*) as t_count FROM %v WHERE %v ", val, this.whereStr)
		}
	}

	if FinallySql == "" {
		FinallySql = fmt.Sprintf("SELECT COUNT(*) as t_count FROM %v WHERE %v %v", this.tablename, this.whereStr, extraSql)
	} else {
		FinallySql = fmt.Sprintf("SELECT SUM(t_count) as t_count FROM (%v) temp", FinallySql)
	}
	this.finalSql = FinallySql

	var resp *DbResult = nil
	var err error
	if this.dbPointer != nil {
		resp, err = this.dbPointer.Query(FinallySql)
	} else {
		resp, err = Query(FinallySql, this.expire)
	}

	if err != nil {
		return 0, err
	}
	if resp == nil || resp.Length == 0 {
		return 0, err
	}

	return resp.ArrResult[0].GetInt32("t_count", 0), nil
}



/**
    查询语句
	获取指定索引处的数据
	@param idx int 索引id
 */
func (this *SqlBuider) Max(_field string) (uint64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	if this.whereStr == "" {
		this.whereStr = "1"
	}

	if this.fieldStr == "" {
		this.fieldStr = "*"
	}

	var extraSql = ""
	if this.group != "" {
		extraSql += "GROUP BY "+this.group
	}

	var FinallySql = strings.Builder{}
	if len(this.unionalls) > 0 {
		for _, val := range this.unionalls {
			if FinallySql.Len() > 0 {
				FinallySql.WriteString(" UNION ALL ")
			}
			FinallySql.WriteString(fmt.Sprintf("SELECT Max(`%v`) as max_id FROM %v WHERE %v ", _field, val, this.whereStr))
		}
	} else if len(this.unions) > 0 {
		for _, val := range this.unions {
			if FinallySql.Len() > 0 {
				FinallySql.WriteString( " UNION ALL " )
			}
			FinallySql.WriteString( fmt.Sprintf("SELECT Max(`%v`) as max_id FROM %v WHERE %v ", _field, val, this.whereStr) )
		}
	}

	if FinallySql.Len() == 0 {
		FinallySql.Reset()
		FinallySql.WriteString(fmt.Sprintf("SELECT Max(`%v`) as max_id FROM %v WHERE %v %v", _field, this.tablename, this.whereStr, extraSql))
	} else {
		_finnal := FinallySql.String()
		FinallySql.Reset()
		FinallySql.WriteString(fmt.Sprintf("SELECT Max(`%v`) as max_id FROM (%v) temp", _field, _finnal))
	}
	this.finalSql = FinallySql.String()

	var resp *DbResult = nil
	var err error
	if this.dbPointer != nil {
		resp, err = this.dbPointer.Query(FinallySql.String())
	} else {
		resp, err = Query(FinallySql.String(), this.expire)
	}

	if err != nil {
		return 0, err
	}
	if resp == nil || resp.Length == 0 {
		return 0, err
	}

	return resp.ArrResult[0].GetUint64("max_id", 0), nil
}


/**
   事务查询语句
 */
func (this *SqlBuider) SelectTx(tx *DbTransform) (*DbResult, error) {

	var extraSql = ""
	if this.group != "" {
		extraSql += "GROUP BY "+this.group
	}
	if this.limit != "" {
		extraSql += this.limit
	}

	sqlStr := fmt.Sprintf("SELECT %v FROM %v WHERE %v %v", this.fieldStr, this.tablename, this.whereStr, extraSql)
	this.finalSql = sqlStr
	var resp *DbResult = nil
	var err error

	switch this.dbType {
	case 0:		// 正常
		resp, err = Query(sqlStr, 0)
	case 1:		// Picker
		resp, err = this.dbPointer.Query(sqlStr)
	case 2:		// 事务
		resp, err = this.dbTx.QueryTx(sqlStr)
	}

	return resp, err
}


/**
   更新语句.
   @param data map[string] string 需要更新字段列表
   @return 修改成功个数, 错误
 */
func (this *SqlBuider) Save(data map[string] interface{}) (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	fieldstr := ""
	for key, val := range data {
		if fieldstr != "" {
			fieldstr += ","
		}
		fieldstr += fmt.Sprintf("`%v` = '%v'", key, val)
	}

	if fieldstr == "" {
		return 0, errors.New("EMPTY UPDATE COLUMN LIST")
	}


	sqlStr := fmt.Sprintf("UPDATE %v SET %v WHERE %v", this.tablename, fieldstr, this.whereStr)
	this.finalSql = sqlStr

	var resp int64
	var err error

	switch this.dbType {
	case 0:		// 正常
		resp, err = Exec(sqlStr)
	case 1:		// Picker
		resp, err = this.dbPointer.Exec(sqlStr)
	case 2:		// 事务
		resp, err = this.dbTx.ExecTx(sqlStr)
	}
	return resp, err
}

/**
   删除语句.
   @return 删除个数, 错误
 */
func (this *SqlBuider) Del() (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	sqlStr := fmt.Sprintf("DELETE FROM %v WHERE %v", this.tablename, this.whereStr)
	this.finalSql = sqlStr

	var resp int64
	var err error

	switch this.dbType {
	case 0:		// 正常
		resp, err = Exec(sqlStr)
	case 1:		// Picker
		resp, err = this.dbPointer.Exec(sqlStr)
	case 2:		// 事务
		resp, err = this.dbTx.ExecTx(sqlStr)
	}

	return resp, err
}


/**
   事务更新语句.
   @param data map[string] string 需要更新字段列表
   @return 修改成功个数, 错误
 */
func (this *SqlBuider) SaveTx(data map[string] interface{}) (int64, error){

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	fieldstr := ""
	for key, val := range data {
		if fieldstr != "" {
			fieldstr += ","
		}
		fieldstr += fmt.Sprintf("%v = '%v'", key, val)
	}

	if fieldstr == "" {
		return 0, errors.New("EMPTY UPDATE COLUMN LIST")
	}


	sqlStr := fmt.Sprintf("UPDATE %v SET %v WHERE %v", this.tablename, fieldstr, this.whereStr)
	this.finalSql = sqlStr

	var resp int64
	var err error

	switch this.dbType {
	case 0:		// 正常
		resp, err = Exec(sqlStr)
	case 1:		// Picker
		resp, err = this.dbPointer.Exec(sqlStr)
	case 2:		// 事务
		resp, err = this.dbTx.ExecTx(sqlStr)
	}
	return resp, err
}


/**
   初始化整个表.
   @return 返回是否发生错误
 */
func (this *SqlBuider) Truncate() error {

	if this.tablename == "" {
		return errors.New("EMPTY TABLE NAME")
	}

	var err error

	if this.dbPointer != nil {
		_, err = this.dbPointer.Exec("TRUNCATE TABLE %v", this.tablename)
	} else {
		_, err = Exec("TRUNCATE TABLE %v", this.tablename)
	}

	return err
}


/**
   添加语句
   @param data map[string] string 需要添加的字段列表
   @return 最后一条添加的id
 */
func (this *SqlBuider) Add(data map[string] interface{}) (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	// 拼接字段区和值字段区
	fieldstr := strings.Builder{}
	valuestr := strings.Builder{}
	for key, val := range data {
		if valuestr.Len() > 0 {
			valuestr.WriteString( "," )
			fieldstr.WriteString( "," )
		}

		fieldstr.WriteString( fmt.Sprintf("`%v`", key) )
		if key == "guid" || key == "uid" || key == "id" {
			valuestr.WriteString( fmt.Sprintf("%v", val) )
		} else {
			valuestr.WriteString( fmt.Sprintf("'%v'", val) )
		}
	}

	if fieldstr.Len() == 0 {
		return 0, errors.New("EMPTY UPDATE COLUMN LIST")
	}

	// 拼接重复区
	onDuplicateStr := strings.Builder{}
	if this.duplicateKey != nil && len(this.duplicateKey) > 0 {
		onDuplicateStr.WriteString(" ON DUPLICATE KEY UPDATE ")

		for i, val := range this.duplicateKey {
			if i > 0 {
				onDuplicateStr.WriteString(",")
			}
			onDuplicateStr.WriteString(fmt.Sprintf("`%[1]v` = VALUES(`%[1]v`)", val))
		}
	}

	sqlStr := fmt.Sprintf("INSERT INTO %v (%v) VALUES(%v) %v", this.tablename, fieldstr.String(), valuestr.String(), onDuplicateStr.String())
	this.finalSql = sqlStr
	var resp int64
	var err error

	switch this.dbType {
	case 0:		// 正常
		resp, err = Exec(sqlStr)
	case 1:		// Picker
		resp, err = this.dbPointer.Exec(sqlStr)
	case 2:		// 事务
		resp, err = this.dbTx.ExecTx(sqlStr)
	}

	return resp, err
}


/**
   添加语句
   @param data map[string] string 需要添加的字段列表
   @return 最后一条添加的id
 */
func (this *SqlBuider) Replace(data map[string] interface{}) (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	fieldstr := strings.Builder{}
	valuestr := strings.Builder{}
	for key, val := range data {
		if valuestr.Len() > 0 {
			valuestr.WriteString( "," )
			fieldstr.WriteString( "," )
		}
		fieldstr.WriteString( fmt.Sprintf("`%v`", key) )
		if key == "guid" || key == "uid" || key == "id" {
			valuestr.WriteString( fmt.Sprintf("%v", val) )
		} else {
			valuestr.WriteString( fmt.Sprintf("'%v'", val) )
		}
	}

	if fieldstr.Len() == 0 {
		return 0, errors.New("EMPTY UPDATE COLUMN LIST")
	}

	sqlStr := fmt.Sprintf("REPLACE INTO %v (%v) VALUES(%v)", this.tablename, fieldstr.String(), valuestr.String())
	this.finalSql = sqlStr

	var resp int64
	var err error

	switch this.dbType {
	case 0:		// 正常
		resp, err = Exec(sqlStr)
	case 1:		// Picker
		resp, err = this.dbPointer.Exec(sqlStr)
	case 2:		// 事务
		resp, err = this.dbTx.ExecTx(sqlStr)
	}

	return resp, err
}

/**
   事务添加语句
   @param data map[string] string 需要添加的字段列表
   @return 最后一条添加的id
 */
func (this *SqlBuider) AddTx(data map[string] interface{}) (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	fieldstr := strings.Builder{}
	valuestr := strings.Builder{}
	for key, val := range data {
		if valuestr.Len() > 0 {
			valuestr.WriteString( "," )
			fieldstr.WriteString( "," )
		}
		fieldstr.WriteString( fmt.Sprintf("`%v`", key) )
		if key == "guid" || key == "uid" || key == "id" {
			valuestr.WriteString( fmt.Sprintf("%v", val) )
		} else {
			valuestr.WriteString( fmt.Sprintf("'%v'", val) )

		}
	}

	if fieldstr.Len() == 0 {
		return 0, errors.New("EMPTY UPDATE COLUMN LIST")
	}

	sqlStr := fmt.Sprintf("INSERT INTO %v (%v) VALUES(%v)", this.tablename, fieldstr.String(), valuestr.String())
	this.finalSql = sqlStr
	var resp int64
	var err error

	switch this.dbType {
	case 0:		// 正常
		resp, err = Exec(sqlStr)
	case 1:		// Picker
		resp, err = this.dbPointer.Exec(sqlStr)
	case 2:		// 事务
		resp, err = this.dbTx.ExecTx(sqlStr)
	}
	if err == nil {
		return resp, nil
	}
	return resp, errors.New(fmt.Sprintf("%v,SQL:%v", err, sqlStr))
}


/**
  添加字段
  @param col string 字段名称
  @param tpname string 字段类型
  @param isnull bool 是否为空
  @param comment string 备注
 */
func (this *SqlBuider) AddColumn(col string, tpname string, isnull bool, defval string, comment string ) (*SqlBuider) {

	this.addColumns = append(this.addColumns, MySqlColumns{
		name: col,
		typename: tpname,
		null: isnull,
		defaults: defval,
		comment: comment,
		autoInc: false,
	})
	return this
}


/**
  删除字段
  @param col 字段名称
 */
func (this *SqlBuider) RemoveColumn (col string) (*SqlBuider) {
	this.removeColumns = append(this.removeColumns, col)
	return this
}


/**
  设置为主键
  col : 需要设置为主键的字段名称
  auto : 是否自动递增
*/
func (this *SqlBuider) SetId(col string, autoInc bool) (*SqlBuider) {
	for key, val := range this.addColumns {
		if val.name == col {
			if autoInc {
				this.addColumns[key].autoInc = true
			}
			break
		}else{
			if autoInc {
				this.addColumns[key].autoInc = false
			}
		}
	}

	this.primaryKeys = col
	return this
}


/**
  添加索引
  cols: 要添加的索引，用逗号隔开
  unique: 是否不重复
 */
func (this *SqlBuider) AddIndex(cols string, unique bool) (*SqlBuider) {
	indexStr := ""
	if unique {
		indexStr = "UNIQUE KEY"
	}else{
		indexStr = "KEY"
	}
	indexStr += " ("+cols+")"
	this.addIndexs = append(this.addIndexs, indexStr)
	return this
}


/**
  移除索引
*/
func (this *SqlBuider) RemoveIndex(cols string) (*SqlBuider) {
	this.removeIndexs = append(this.removeIndexs, cols)
	return this
}


/**
  创建表格
  overWrite: 是否覆盖, 如果为true则会先删除原先的表
*/
func (this *SqlBuider) CreateTable(overWrite bool) bool {
	//生成整个表格的SQL
	if this.tablename == "" {
		return false
	}

	if len(this.addColumns) == 0 {
		return false
	}

	sqlStr := strings.Builder{}
	sqlStr.WriteString( "CREATE TABLE IF NOT EXISTS "+this.tablename + "(" )
	for key, val := range this.addColumns {

		sqlStr.WriteString( "`" + val.name + "` " + val.typename )
		if val.autoInc == true {
			sqlStr.WriteString( " AUTO_INCREMENT" )
		} else {
			if val.null == false && val.typename != "text" {
				sqlStr.WriteString( " NOT NULL DEFAULT '" + val.defaults+"'" )
			}
		}
		sqlStr.WriteString( " COMMENT '" + val.comment + "'" )

		if key < len(this.addColumns)-1 {
			sqlStr.WriteString( "," )
		}
	}
	if this.primaryKeys != "" {
		sqlStr.WriteString( ", PRIMARY KEY (" + this.primaryKeys + ")" )
	}

	if len(this.addIndexs) > 0 {
		sqlStr.WriteString( "," + strings.Join(this.addIndexs, ",") )
	}
	sqlStr.WriteString( ") ENGINE=INNODB DEFAULT CHARSET=UTF8" )

	if overWrite {
		Exec("DROP TABLE IF EXISTS " + this.tablename)
	}

	this.finalSql = sqlStr.String()

	if this.dbPointer != nil {
		_, err := this.dbPointer.Exec( sqlStr.String() )
		if err != nil {
			fmt.Printf(">> 添加Exec :sql失败! %v\n", err)
			return false
		}
	} else if this.dbTx != nil {
		_, err := this.dbTx.ExecTx( sqlStr.String() )
		if err != nil {
			fmt.Printf(">> 添加ExecTx sql失败! %v\n", err)
			return false
		}
	}else {
		_, err := Exec( sqlStr.String() )

		if err != nil {
			fmt.Printf(">> 添加默认 sql失败! %v\n", err)
			return false
		}
	}
	return true
}


// 修改表格结构
func (this *SqlBuider) SaveTable() error {

	//生成整个表格的SQL
	if this.tablename == "" {
		return errors.New("TABLE NAME IS EMPTY")
	}

	SqlStr := strings.Builder{}
	SqlStr.WriteString( "ALTER TABLE "+this.tablename )
	AnyThing := false

	if len(this.addColumns) > 0 {
		// 添加字段
		AnyThing = true
		for key, val := range this.addColumns {
			SqlStr.WriteString(  "ADD COLUMN "+val.name+" "+val.typename )
			if !val.null {
				SqlStr.WriteString( " NOT NULL" )
			}

			SqlStr.WriteString( " DEFAULT '" + val.defaults + "' COMMENT '" + val.comment + "'" )
			if key < len(this.addColumns) - 1 {
				SqlStr.WriteString( "," )
			}
		}
	}

	if len(this.removeColumns) > 0 {
		if AnyThing {
			SqlStr.WriteString( "," )
		}
		AnyThing = true
		for key, val := range this.removeColumns {
			SqlStr.WriteString( "DROP COLUMN " + val )
			if key < len(this.removeColumns) - 1 {
				SqlStr.WriteString( "," )
			}
		}
	}

	if len(this.addIndexs) > 0 {
		if AnyThing {
			SqlStr.WriteString( "," )
		}
		AnyThing = true
		for key, val := range this.addIndexs {
			SqlStr.WriteString( "ADD INDEX " + val )
			if key < len(this.addIndexs) - 1 {
				SqlStr.WriteString( "," )
			}
		}
	}

	if len(this.removeIndexs) > 0 {
		if AnyThing {
			SqlStr.WriteString( "," )
		}
		for key, val := range this.removeIndexs {
			SqlStr.WriteString( "DROP INDEX " + val )
			if key < len(this.removeIndexs) - 1 {
				SqlStr.WriteString( "," )
			}
		}
	}

	this.finalSql = SqlStr.String()

	if this.dbPointer != nil {
		_, err := this.dbPointer.Exec( SqlStr.String() )
		if err != nil {
			return err
		}
	} else if this.dbTx != nil {
		_, err := this.dbTx.ExecTx( SqlStr.String() )
		if err != nil {
			return err
		}
	}else {
		_, err := Exec( SqlStr.String() )
		if err != nil {
			return err
		}
	}
	return nil
}


// 联表查询
// @param tablename string 联表名称
func (this *SqlBuider) UnionAll(tablename string) *SqlBuider {
	for _, val := range this.unionalls {
		if val == tablename {
			return this
		}
	}
	this.unionalls = append(this.unionalls, tablename)
	return this
}

// 强制联表查询
// @param tablename string 联表名称
func (this *SqlBuider) Union(tablename string) *SqlBuider {
	for _, val := range this.unions {
		if val == tablename {
			return this
		}
	}
	this.unions = append(this.unionalls, tablename)
	return this
}

/*
	获取sql语句
 */
func (this *SqlBuider)GetLastSql() string {
	return this.finalSql
}


// 获取查找
func (this *SqlBuider) FindAll(_resp interface{}) error {

	_value := reflect.ValueOf(_resp)
	_valueE := _value.Elem()
	_valueE = _valueE.Slice(0, _valueE.Cap())

	_element := _valueE.Type().Elem()
	fieldList := GetAllField(reflect.New(_element).Interface())

	var extraSql = ""
	if this.group != "" {
		extraSql += "GROUP BY "+this.group
	}

	if this.orders != "" {
		extraSql += " ORDER BY "+this.orders
	}

	if this.limit != "" {
		extraSql += this.limit
	}

	where_str := ""
	if this.whereStr != "" {
		where_str = "WHERE " + this.whereStr
	}


	sqlStr := fmt.Sprintf("SELECT `%v` FROM `%v`.`%v` %v %v", strings.Join(fieldList, "`,`"), this.dbname, this.tablename, where_str, extraSql)

	var resp *DbResult = nil
	var err error

	switch this.dbType {
	case 0:		// 正常
		resp, err = Query(sqlStr, 0)
	case 1:		// Picker
		resp, err = this.dbPointer.Query(sqlStr)
	case 2:		// 事务
		resp, err = this.dbTx.QueryTx(sqlStr)
	}

	if err != nil {
		return err
	}

	if resp != nil && resp.Length > 0 {
		i := 0
		for idx, row := range resp.ArrResult {

			// 需要添加
			if _valueE.Len() == idx {
				elemp := reflect.New(_element)
				Unmarsha(row, elemp.Interface())
				_valueE = reflect.Append(_valueE, elemp.Elem())
			}

			i++
		}

		_value.Elem().Set(_valueE.Slice(0, i))
	}

	return nil
}


// 获取查找
func (this *SqlBuider) FindOne(_resp interface{}) error {

	fieldList := GetAllField(_resp)

	var extraSql = ""
	if this.group != "" {
		extraSql += "GROUP BY "+this.group
	}

	if this.limit != "" {
		extraSql += this.limit
	}

	where_str := ""
	if this.whereStr != "" {
		where_str = "WHERE " + this.whereStr
	}


	sqlStr := fmt.Sprintf("SELECT `%v` FROM `%v`.`%v` %v %v", strings.Join(fieldList, "`,`"), this.dbname, this.tablename, where_str, extraSql)

	var resp *DbResult = nil
	var err error

	switch this.dbType {
	case 0:		// 正常
		resp, err = Query(sqlStr, 0)
	case 1:		// Picker
		resp, err = this.dbPointer.Query(sqlStr)
	case 2:		// 事务
		resp, err = this.dbTx.QueryTx(sqlStr)
	}

	if err != nil {
		return err
	}

	if resp == nil || resp.Length == 0 {
		return errors.New("not found")
	}

	Unmarsha(resp.ArrResult[0], _resp)
	return nil
}




// 获取查找
func (this *SqlBuider) AddObj(_resp interface{}, _include_primary bool) (int64, error) {

	fieldList, valuesList := GetInsertSql(_resp, _include_primary)

	// 拼接重复区
	onDuplicateStr := strings.Builder{}
	if this.duplicateKey != nil && len(this.duplicateKey) > 0 {
		onDuplicateStr.WriteString(" ON DUPLICATE KEY UPDATE ")

		for i, val := range this.duplicateKey {
			if i > 0 {
				onDuplicateStr.WriteString(",")
			}
			onDuplicateStr.WriteString(fmt.Sprintf("`%[1]v` = VALUES(`%[1]v`)", val))
		}
	}

	sqlStr := fmt.Sprintf("INSERT INTO `%v`.`%v` (`%v`) VALUES('%v') %v", this.dbname, this.tablename, strings.Join(fieldList, "`,`"), strings.Join(valuesList, "','"), onDuplicateStr.String())

	var resp int64
	var err error

	switch this.dbType {
	case 0:		// 正常
		resp, err = Exec(sqlStr, 0)
	case 1:		// Picker
		resp, err = this.dbPointer.Exec(sqlStr)
	case 2:		// 事务
		resp, err = this.dbTx.ExecTx(sqlStr)
	}

	if err != nil {
		return 0, errors.New(fmt.Sprintf("%v SQL(%v)", err, sqlStr))
	}

	return resp, nil
}

