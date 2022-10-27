package dbs

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	//数据库驱动
	_ "github.com/go-sql-driver/mysql"
)

//Db db结构
type Db struct {
	MysqlDb *sql.DB
	tx      *sql.Tx
}

//数据库连接配置
// const (
// 	UserName = "root"
// 	PassWord = "root"
// 	HOST     = "localhost"
// 	PORT     = "3306"
// 	DATABASE = "wb"
// 	CHARSET  = "utf8"
// )

//NewDb db对象
func NewDb(UserName, PassWord, HOST, PORT, DATABASE, CHARSET string) *Db {
	db := &Db{}
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", UserName, PassWord, HOST, PORT, DATABASE, CHARSET)

	var MysqlErr error
	// 打开连接失败
	db.MysqlDb, MysqlErr = sql.Open("mysql", dbDSN)
	//defer MysqlDb.Close();
	if MysqlErr != nil {
		log.Println("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + MysqlErr.Error())
	}

	// 最大连接数
	db.MysqlDb.SetMaxOpenConns(50)
	// 闲置连接数
	db.MysqlDb.SetMaxIdleConns(1)
	// 最大连接周期
	db.MysqlDb.SetConnMaxLifetime(100 * time.Second)

	if MysqlErr = db.MysqlDb.Ping(); nil != MysqlErr {
		panic("数据库链接失败: " + MysqlErr.Error())
	}

	return db
}

//初始话，可以利用来创建默认的数据库对象，以及直接使用默认数据库对象及其相关函数，不用通过指定数据库对象来调用，跟方面
//func init()  {
//	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", UserName, PassWord, HOST, PORT, DATABASE, CHARSET)
//
//	// 打开连接失败
//	MysqlDb, MysqlErr = sql.Open("mysql", dbDSN)
//	//defer MysqlDb.Close();
//	if MysqlErr != nil {
//		log.Println("dbDSN: " + dbDSN)
//		panic("数据源配置不正确: " + MysqlErr.Error())
//	}
//
//	// 最大连接数
//	MysqlDb.SetMaxOpenConns(100)
//	// 闲置连接数
//	MysqlDb.SetMaxIdleConns(20)
//	// 最大连接周期
//	MysqlDb.SetConnMaxLifetime(100*time.Second)
//
//	if MysqlErr = MysqlDb.Ping(); nil != MysqlErr {
//		panic("数据库链接失败: " + MysqlErr.Error())
//	}
//}

//SetArgs 为FindOne设置参数组
func (d *Db) SetArgs(args ...interface{}) []interface{} {
	s := []interface{}{}
	return append(s, args...)
}

//FindOne 查询单个数据
//SQL，SQL参数切片请使用SetArgs()
//每个字段值一一保存到变量：&a,&b...【注意必须是变量的指针，且个数能字段对应上】
//
//if err == sql.ErrNoRows 找不到记录
//if err == nil 找到了且没有出错
func (d *Db) FindOne(sql string, args []interface{}, res ...interface{}) error {
	//预处理
	stmt, err := d.MysqlDb.Prepare(sql)

	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRow(args...)

	//row := MysqlDb.QueryRow(sql,args...) //普通查询

	if err := row.Scan(res...); err != nil {
		//fmt.Printf("scan failed, err:%v", err)
		return err
	}

	return nil
}

//FindAll 查询所有数据：数据映射的切片
//SQL，SQL参数列表，请使用者自己解析出相应的结果集中的数据类型
//TODO 后期加上自动类型转换功能
/*
	type tier struct {
		Tid  int
		Name string
		Pid  int
		D    int
		T    int
		G    int
	}
	res, err := db.FindAll("SELECT * FROM tier")
	if err != nil {
		fmt.Println(err.Error())
	}
*通过转换和断言取出正确类型的值
	for _, re := range res {
		fmt.Println(re)
		t := &tier{}
		t.Tid = int(re["tid"].(int64))
		t.Name = string(re["name"].([]byte))
		t.Pid = int(re["pid"].(int64))
		t.D = int(re["d"].(int64))
		t.T = int(re["t"].(int64))
		t.G = int(re["g"].(int64))
		fmt.Println(t)
	}
*/
func (d *Db) FindAll(sql string, args ...interface{}) ([]map[string]interface{}, error) {
	//预处理
	stmt, err := d.MysqlDb.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	//rows, _:= MysqlDb.Query(sql,args...) //普通查询

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	cLen := len(columns)
	all := []map[string]interface{}{}
	resPs := make([]interface{}, cLen) //存放这条记录的地址
	// 遍历
	for rows.Next() {
		//转结果
		resSlice := make([]interface{}, cLen) //存放一条记录
		for i := 0; i < cLen; i++ {
			resPs[i] = &resSlice[i] //取地址
		}
		err := rows.Scan(resPs...) //扫描到相应的地址上
		if err != nil {
			return nil, err
		}

		//组织结果
		var resArr = make(map[string]interface{}, cLen)
		for k, v := range columns {
			resArr[v] = resSlice[k]
		}
		all = append(all, resArr)
	}

	return all, nil
}

//Insert 插入数据：插入的主键ID 影响行数
func (d *Db) Insert(sql string, args ...interface{}) (int64, int64, error) {
	//预处理
	stmt, err := d.MysqlDb.Prepare(sql)
	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(args...)
	//ret,err := d.MysqlDb.Exec(sql,args...) //普通执行

	if err != nil {
		return 0, 0, err
	}

	//插入数据的主键id
	lastInsertID, err := ret.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	//影响行数
	rowsAffected, err := ret.RowsAffected()
	if err != nil {
		return lastInsertID, 0, err
	}

	return lastInsertID, rowsAffected, nil
}

//Update 更新数据：影响的行数
func (d *Db) Update(sql string, args ...interface{}) (int64, error) {
	//预处理
	stmt, err := d.MysqlDb.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	//ret,_ := MysqlDb.Exec(sql,args...)

	updNums, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}

	return updNums, nil
}

//Del 删除数据：影响的行数
func (d *Db) Del(sql string, args ...interface{}) (int64, error) {
	//预处理
	stmt, err := d.MysqlDb.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	//ret,_ := MysqlDb.Exec(sql,args...)

	updNums, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}

	return updNums, nil
}

//Close 关闭数据库
func (d *Db) Close() {
	d.MysqlDb.Close()
}

//Begin 启动事务
func (d *Db) Begin() error {
	var ok error
	d.tx, ok = d.MysqlDb.Begin()

	return ok
}

//Commit 提交事务
func (d *Db) Commit() error {
	return d.tx.Commit()
}

//Rollback 回滚事务
func (d *Db) Rollback() error {
	return d.tx.Rollback()
}

//-------------------------------事务处理
//tx, _ := MysqlDb.Begin()
//tx.Commit()
//tx.Rollback()
