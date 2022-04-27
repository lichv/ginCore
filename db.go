package ginCore

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"time"
)

type DBConfig struct {
	Type    string `json:"type"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Name    string `json:"name"`
	User    string `json:"user"`
	Passwd  string `json:"passwd"`
	Charset string `json:"charset"`
	Prefix  string `json:"prefix"`
}

type DBService struct {
	Config *DBConfig
	Client *gorm.DB
}

func GetDBByConfig(c *DBConfig) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	namingStrategy := schema.NamingStrategy{
		TablePrefix:   c.Prefix,                          // 表名前缀，`User`表为`t_users`
		SingularTable: true,                              // 使用单数表名，启用该选项后，`User` 表将是`user`
		NameReplacer:  strings.NewReplacer("CID", "Cid"), // 在转为数据库名称之前，使用NameReplacer更改结构/字段名称。
	}
	if c.Type == "mysql" {
		return GetMysqlDB(c.Host, c.Port, c.User, c.Passwd, c.Name, "utf8", newLogger, namingStrategy)
	} else if c.Type == "postgres" {
		return GetPostgresDB(c.Host, c.Port, c.User, c.Passwd, c.Name, newLogger, namingStrategy)
	} else if c.Type == "sqlserver" {
		return GetSqlserverDB(c.Host, c.Port, c.User, c.Passwd, c.Name, newLogger, namingStrategy)
	} else if c.Type == "sqlite" {
		return GetSqliteDB(c.Name, newLogger, namingStrategy)
	}
	return nil, nil
}

func (s *DBService) Count(tableName string, query map[string]interface{}) (int64, error) {
	var count int64
	err := s.Client.Table(tableName).Scopes(WhereBuild(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func WhereBuild(where map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var whereSQL string
		var vals []interface{}
		for k, v := range where {
			ks := strings.Split(k, " ")
			if len(ks) > 2 {
				return db
			}

			if whereSQL != "" {
				whereSQL += " AND "
			}

			fmt.Println(strings.Join(ks, ","))
			switch len(ks) {
			case 1:
				switch v := v.(type) {
				case NullType:
					fmt.Println()
					if v == IsNotNull {
						whereSQL += fmt.Sprint("`", k, "` IS NOT NULL")
					} else {
						whereSQL += fmt.Sprint("`", k, "` IS NULL")
					}
				default:
					whereSQL += fmt.Sprint("`", k, "`=?")
					vals = append(vals, v)
				}
			case 2:
				k = ks[0]
				switch ks[1] {
				case "=":
					whereSQL += fmt.Sprint("`", k, "`= ?")
					vals = append(vals, v)
				case ">":
					whereSQL += fmt.Sprint("`", k, "`>?")
					vals = append(vals, v)
				case ">=":
					whereSQL += fmt.Sprint("`", k, "`>=?")
					vals = append(vals, v)
				case "<":
					whereSQL += fmt.Sprint("`", k, "`<?")
					vals = append(vals, v)
				case "<=":
					whereSQL += fmt.Sprint("`", k, "`<=?")
					vals = append(vals, v)
				case "!=":
					whereSQL += fmt.Sprint("`", k, "`!=?")
					vals = append(vals, v)
				case "<>":
					whereSQL += fmt.Sprint("`", k, "`!=?")
					vals = append(vals, v)
				case "in":
					whereSQL += fmt.Sprint("`", k, "` in (?)")
					vals = append(vals, v)
				case "like":
					whereSQL += fmt.Sprint("`", k, "` like ?")
					vals = append(vals, v)
				}
			}

		}
		return db.Where(whereSQL, vals...)
	}
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
func NothingDone() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}
func NameHandle(keyword string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name like ?", "%"+keyword+"%")
	}
}
func LabelHandle(keyword string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("label like ?", "%"+keyword+"%")
	}
}
func NoteHandle(keyword string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("note like ?", "%"+keyword+"%")
	}
}

func KeyWordHandle(keyword string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name like ?", "%"+keyword+"%").Or("label like ?", "%"+keyword+"%").Or("note like ?", "%"+keyword+"%")
	}
}
func GetMysqlDB(host string, port int, username string, passwd string, dbname string, charset string, logger logger.Interface, namingStrategy schema.NamingStrategy) (*gorm.DB, error) {
	var datetimePrecision = 2
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", username, passwd, host, port, dbname, charset)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DefaultDatetimePrecision:  &datetimePrecision,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger:         logger,
		NamingStrategy: namingStrategy,
	})
	if err != nil {
		log.Panicln("Init mysql failed", err.Error())
	}

	return db, nil
}

func GetPostgresDB(host string, port int, username, passwd, dbname string, logger logger.Interface, namingStrategy schema.NamingStrategy) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", host, username, passwd, dbname, port)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:         logger,
		NamingStrategy: namingStrategy,
	})
	if err != nil {
		log.Panicln("Init mysql failed", err.Error())
	}

	return db, nil
}

func GetSqlserverDB(host string, port int, username, passwd, dbname string, logger logger.Interface, namingStrategy schema.NamingStrategy) (*gorm.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%S:%d?database=%s", username, passwd, host, port, dbname)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger:         logger,
		NamingStrategy: namingStrategy,
	})
	if err != nil {
		log.Panicln("Init mysql failed", err.Error())
	}

	return db, nil

}

func GetSqliteDB(filepath string, logger logger.Interface, namingStrategy schema.NamingStrategy) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s", filepath)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger:         logger,
		NamingStrategy: namingStrategy,
	})
	if err != nil {
		log.Panicln("Init mysql failed", err.Error())
	}
	return db, nil
}
