package storage

import (
	"database/sql"
	"fmt"
	"log"
	"ProxyPool/models"
	_ "github.com/mattn/go-sqlite3"
)

var GlobalSqliteDB, err = sql.Open("sqlite3", string(Config.Sqlite.Addr))

// Storage struct is used for storeing persistent data of alerts
type SqliteStorage struct {
	database string
	table    string
}

// NewStorage creates and returns new Storage instance
func NewSqliteStorage() *SqliteStorage {
	return &SqliteStorage{database: Config.Sqlite.DB, table: Config.Sqlite.Table}
}

// Create insert new item
func (s *SqliteStorage) Create(item *models.IP) error {
	log.Println("error:" , err)
	stmt, err := GlobalSqliteDB.Prepare("INSERT INTO ip('ip','type') values(?,?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	result, err := stmt.Exec(item.Data, item.Type)
	if err != nil {
		fmt.Printf("add error: %v", err)
		return err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted id is ", lastID)
	return nil
}

// GetOne Finds and returns one data from storage
func (s *SqliteStorage) GetOne(value string) (*models.IP, error) {
	//ses := s.GetDBSession()
	//defer ses.Close()
	//t := models.NewIP()
	//err := ses.DB(s.database).C(s.table).Find(bson.M{"data": value}).One(t)
	//if err != nil {
	//	return nil, err
	//}
	//return t, nil
	//rows, err := GlobalSqliteDB.Query("SELECT * FROM ip WHERE id >= ((SELECT MAX(id) FROM ip)-(SELECT MIN(id) FROM ip)) * RANDOM() + (SELECT MIN(id) FROM ip)  LIMIT 1")
	rows, err := GlobalSqliteDB.Query("SELECT * FROM ip ORDER BY RANDOM() LIMIT 1")
	if err != nil {
		fmt.Println("error get:", err.Error())
		return nil, err
	}
	defer rows.Close()
	rows.Next()
	newIp := models.NewIP()
	err = rows.Scan(&newIp.IDInt, &newIp.Data, &newIp.Type)
	if err != nil {
		fmt.Println("error scan:", err)
		return nil, err
	}
	fmt.Println(newIp.IDInt, newIp.Data, newIp.Type)
	return newIp, nil
}

// Count all collections
func (s *SqliteStorage) Count() int {
	//ses := s.GetDBSession()
	//defer ses.Close()
	//num, err := ses.DB(s.database).C(s.table).Count()
	//if err != nil {
	//	num = 0
	//}
	//return num
	rows, err := GlobalSqliteDB.Query("SELECT COUNT(*) as count FROM  ip")
	if err != nil {
		fmt.Println("error get:", err.Error())
		return 0
	}
	defer rows.Close()
	return checkCount(rows)
}

// Delete .
func (s *SqliteStorage) Delete(ip *models.IP) error {
	//ses := s.GetDBSession()
	//defer ses.Close()
	//err := ses.DB(s.database).C(s.table).RemoveId(ip.ID)
	//if err != nil {
	//	return err
	//}
	//return nil
	stmt, err := GlobalSqliteDB.Prepare("delete from ip where id=?")
	checkErr(err)

	_, err = stmt.Exec(ip.IDInt)
	checkErr(err)
	defer stmt.Close()
	return err
}

// Update .
func (s *SqliteStorage) Update(ip *models.IP) error {
	//ses := s.GetDBSession()
	//defer ses.Close()
	//err := ses.DB(s.database).C(s.table).Update(bson.M{"_id": ip.ID}, ip)
	//if err != nil {
	//	return err
	//}
	//return nil
	//更新数据
	stmt, err := GlobalSqliteDB.Prepare("update ip set 'ip'=?,'type'=? where id=?")
	checkErr(err)

	res, err := stmt.Exec(ip.Data, ip.Type, ip.IDInt)
	checkErr(err)
	defer stmt.Close()

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)
	return err
}

// GetAll .
func (s *SqliteStorage) GetAll() ([]*models.IP, error) {
	//ses := s.GetDBSession()
	//defer ses.Close()
	//var ips []*models.IP
	//err := ses.DB(s.database).C(s.table).Find(nil).All(&ips)
	//if err != nil {
	//	return nil, err
	//}
	//return ips, nil
	rows, err := GlobalSqliteDB.Query("SELECT * FROM ip")
	checkErr(err)
	defer rows.Close()
	ips := make([]*models.IP, 0)
	for rows.Next() {
		newIp := models.NewIP()
		err = rows.Scan(&newIp.IDInt, &newIp.Data, &newIp.Type)
		if err != nil {
			fmt.Println("error scan:", err)
			return nil, err
		}
		ips = append(ips, newIp)
	}
	return ips, nil
}

//// FindAll .
//func (s *SqliteStorage) FindAll(value string) ([]*models.IP, error) {
//	ses := s.GetDBSession()
//	defer ses.Close()
//	var ips []*models.IP
//	err := ses.DB(s.database).C(s.table).Find(bson.M{"type": bson.M{"$regex": value, "$options": "$i"}}).All(&ips)
//	if err != nil {
//		return nil, err
//	}
//	return ips, nil
//
//}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err:= rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}