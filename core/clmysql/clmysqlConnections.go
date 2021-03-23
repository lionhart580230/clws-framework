package clmysql

import (
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"fmt"
	"log"
)

/**
   数据库连接池管理，主从控制。故障侦测等
 */

const (
	MYSQL_STATUS_FAILED = 0
	MYSQL_STATUS_NORMAL = 1
	MYSQL_STATUS_BUSY = 2
)

// 单个链接的状态
type clmysqlStruct struct {
	db *sql.DB // 数据库指针
	lastCheck time.Time			// 上次检测时间
	status int8
	dbname string				// 当前绑定的数据库名称
	mkey string
	skey int
}

// 单个数据库状态
type clmysqlStatus struct {
	dbpool []*clmysqlStruct		// 连接池
	onlyread bool				// 是否只读
	maxLimit uint16				// 最大链接数
	Host string					// 数据库地址
	User string					// 数据库用户名
	Pass string					// 数据库密码
	DBName string				// 数据库名称
}


// 数据库连接池
type ClmysqlConnections struct {
	master map[string] *clmysqlStatus
	slaver map[string] *clmysqlStatus
	masterLock sync.RWMutex
	slaverLock sync.RWMutex
}

/**
	DB 事务
 */
type DbTransform struct {
	tx *sql.Tx
	st *clmysqlStruct
}

// 开始检查每个链接的状态
func (clConn *ClmysqlConnections)StartToCheck () {
	for {
		if len(clConn.master) > 0 {
			clConn.masterLock.Lock()
			for key, val := range clConn.master {
				for mkey, conn := range val.dbpool {
					if time.Since(conn.lastCheck).Minutes() >= 3 {
						err := conn.db.Ping()
						if err != nil {
							clConn.master[key].dbpool[mkey].status = MYSQL_STATUS_FAILED
							clConn.master[key].dbpool[mkey].lastCheck = time.Now()
							continue
						}

						clConn.master[key].dbpool[mkey].status = MYSQL_STATUS_NORMAL
						clConn.master[key].dbpool[mkey].lastCheck = time.Now()
					}
				}

			}
			clConn.masterLock.Unlock()
		}

		if len(clConn.slaver) > 0 {
			clConn.slaverLock.Lock()
			for key, val := range clConn.slaver {
				for mkey, conn := range val.dbpool {
					err := conn.db.Ping()
					if err != nil {
						clConn.slaver[key].dbpool[mkey].status = MYSQL_STATUS_FAILED
						clConn.slaver[key].dbpool[mkey].lastCheck = time.Now()
						continue
					}

					clConn.slaver[key].dbpool[mkey].status = MYSQL_STATUS_NORMAL
					clConn.slaver[key].dbpool[mkey].lastCheck = time.Now()
				}

			}
			clConn.slaverLock.Unlock()
		}

		DBPointerPoolLocker.Lock()
		for k, v := range DBPointerPool {
			if v.lastupdate < uint32(time.Now().Unix()) - 1800 {
				DBPointerPool[k].conn.Close()
				delete(DBPointerPool, k)
			}
		}
		DBPointerPoolLocker.Unlock()

		<-time.After(3 * time.Minute)
	}

}


// 添加新的主库, 如果存在多个主库，则每次操作都将自动选择一个优质的主库链接进行数据库操作
// @param key string 链接键
// @param host string 数据库的IP:PORT
// @param user string 数据库的帐号
// @param pass string 数据库的密码
// @param dbname string 数据库的名称
func (clConn *ClmysqlConnections)AddMaster(key string, host string, user string, pass string, dbname string, limit uint16) error {

	clConn.masterLock.Lock()
	clConn.master[key] = &clmysqlStatus{
		onlyread: false,
		Host: host,
		User: user,
		Pass: pass,
		DBName: dbname,
		maxLimit: limit,
	}
	clConn.master[key].OpenNew(key)
	clConn.masterLock.Unlock()
	return nil
}


// 开启一个新的数据库连线
func (Conn *clmysqlStatus) OpenNew(mkey string) *clmysqlStruct {

	dbp, err := sql.Open("mysql", Conn.User+":"+Conn.Pass+"@tcp("+Conn.Host+")/"+Conn.DBName+"?charset=utf8")
	if err != nil {
		log.Printf(">> ERR: %v\n", err)
		return nil
	}

	err = dbp.Ping()
	if err != nil {
		log.Printf(">> ERR: %v\n", err)
		return nil
	}

	dbp.SetMaxOpenConns(1)
	dbp.SetConnMaxLifetime(0)

	object := clmysqlStruct{
		db: dbp,
		status: MYSQL_STATUS_NORMAL,
		lastCheck: time.Now(),
		mkey: mkey,
		skey: len(Conn.dbpool),
	}

	Conn.dbpool = append(Conn.dbpool, &object)
	return &object
}


// 添加新的从库, 从库可以指定是否只读，只读的从库将不会对除了select之外的任何请求做响应
// @param key string 从库的key值
// @param host string 从库的IP:PORT
// @param user string 从库的帐号
// @param pass string 从库的密码
// @param dbname string 从库的数据库名称

func (clConn *ClmysqlConnections)AddSlaver(key string, host string, user string, pass string, dbname string, readOnly bool, limit uint16) error {

	clConn.slaverLock.Lock()
	clConn.slaver[key] = &clmysqlStatus{
		onlyread: readOnly,
		Host: host,
		User: user,
		Pass: pass,
		DBName: dbname,
		maxLimit: limit,
	}
	clConn.slaver[key].OpenNew(key)
	clConn.slaverLock.Unlock()
	return nil
}

// 选择一个从数据库链接
func (clConn *ClmysqlConnections)SelectSlaver(key string) (*clmysqlStruct) {

	// 没有从数据库
	if len(clConn.slaver) == 0 {
		return nil
	}

	// 没有指定使用哪一个从库
	if key == "" {
		// 先扫描一遍，查看是否有空闲的链接可以使用
		clConn.slaverLock.Lock()
		defer clConn.slaverLock.Unlock()

		for skey, val := range clConn.slaver {
			for vkey, conn := range val.dbpool {
				// 找到空闲链接
				if conn.status == MYSQL_STATUS_NORMAL || (conn.status == MYSQL_STATUS_BUSY && time.Since(conn.lastCheck).Minutes() >= 1 ) {
					clConn.slaver[skey].dbpool[vkey].status = MYSQL_STATUS_BUSY
					clConn.slaver[skey].dbpool[vkey].lastCheck = time.Now()
					return clConn.slaver[skey].dbpool[vkey]
				}
			}
		}

		// 尝试创建链接
		for mkey, val := range clConn.slaver {

			// 存在连接池未满的
			if len(val.dbpool) < int(val.maxLimit) {
				return val.OpenNew(mkey)
			}
		}
	} else {
		// 有传入key
		// 先扫描一遍，查看是否有空闲的链接可以使用
		clConn.slaverLock.Lock()
		defer clConn.slaverLock.Unlock()

		if _, find := clConn.slaver[key]; !find {
			return nil
		}

		for skey, conn := range clConn.slaver[key].dbpool{
			if conn.status == MYSQL_STATUS_NORMAL || (conn.status == MYSQL_STATUS_BUSY && time.Since(conn.lastCheck).Minutes() >= 5) {
				clConn.slaver[key].dbpool[skey].status = MYSQL_STATUS_BUSY
				clConn.slaver[key].dbpool[skey].lastCheck = time.Now()
				return clConn.slaver[key].dbpool[skey]
			}
		}
		// 连接池未满, 添加一条新链接
		if len(clConn.slaver[key].dbpool) < int(clConn.slaver[key].maxLimit) {
			return clConn.slaver[key].OpenNew(key)
		}
	}

	return nil
}

// 选择一个主数据库链接
func (clConn *ClmysqlConnections) SelectMaster(key string) *clmysqlStruct {
	// 没有从数据库
	if len(clConn.master) == 0 {
		fmt.Printf(">> NO MASTER!!\n")
		return nil
	}

	// 没有指定使用哪一个从库
	if key == "" {
		// 先扫描一遍，查看是否有空闲的链接可以使用
		for {
			for mkey, val := range clConn.master {
				for vkey, conn := range val.dbpool {
					// 找到空闲链接
					if conn.status == MYSQL_STATUS_NORMAL || (conn.status == MYSQL_STATUS_BUSY && time.Since(conn.lastCheck).Minutes() >= 1 ) {
						clConn.masterLock.Lock()
						defer clConn.masterLock.Unlock()

						clConn.master[mkey].dbpool[vkey].status = MYSQL_STATUS_BUSY
						clConn.master[mkey].dbpool[vkey].lastCheck = time.Now()
						return clConn.master[mkey].dbpool[vkey]
					}
				}
			}

			// 尝试创建链接
			for mkey, val := range clConn.master {
				// 存在连接池未满的
				if len(val.dbpool) < int(val.maxLimit) {
					return val.OpenNew(mkey)
				}
			}
		}
	}
	return nil
}

// 操作完毕，释放资源
func (cl *clmysqlStruct) Close() {
	cl.status = MYSQL_STATUS_NORMAL
}


// 创建事务
func (clConn *ClmysqlConnections)StartTrans() (*DbTransform, error) {
	var mconn = clConn.SelectMaster("")
	tx, err := mconn.db.Begin()
	if err != nil {
		return nil, err
	}
	var clt = DbTransform{
		tx : tx,
		st : mconn,
	}
	return &clt, nil
}

// 提交事务
func (ct *DbTransform) Commit() error {
	ct.st.Close()
	return ct.tx.Commit()
}