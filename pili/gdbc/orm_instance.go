package gdbc

import (
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/cihub/seelog"
	"github.com/daiguadaidai/pilipala/pili/config"
	"strings"
	"fmt"
)

var ormInstance *OrmInstance

type OrmInstance struct {
	DB *gorm.DB
	sync.Once
}

// 单例模式获取原生数据库链接
func GetOrmInstance() *OrmInstance {
	if ormInstance.DB == nil {
		// 实例化元数据库实例
		ormInstance.Once.Do(func() {
			// 获取元数据配置信息
			dbConfig := config.GetDBConfig()

			// 链接数据库
			var err error
			ormInstance.DB, err = gorm.Open("mysql", dbConfig.GetDataSource())
			if err != nil {
				seelog.Errorf("打开ORM数据库实例错误. %v", err)
			}

			ormInstance.DB.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)
			ormInstance.DB.DB().SetMaxIdleConns(dbConfig.MaxIdelConns)
			ormInstance.DB.DB().Ping()
		})
	}

	return ormInstance
}

func (this *OrmInstance) BatchInsert(db *gorm.DB, objArr []interface{}) error {
	// If there is no data, nothing to do.
	if len(objArr) == 0 {
		return nil
	}

	mainObj := objArr[0]
	mainScope := db.NewScope(mainObj)
	mainFields := mainScope.Fields()
	quoted := make([]string, 0, len(mainFields))
	for i := range mainFields {
		// If primary key has blank value (0 for int, "" for string, nil for interface ...), skip it.
		// If field is ignore field, skip it.
		if (mainFields[i].IsPrimaryKey && mainFields[i].IsBlank) || (mainFields[i].IsIgnored) {
			continue
		}
		quoted = append(quoted, mainScope.Quote(mainFields[i].DBName))
	}

	placeholdersArr := make([]string, 0, len(objArr))

	for _, obj := range objArr {
		scope := db.NewScope(obj)
		fields := scope.Fields()
		placeholders := make([]string, 0, len(fields))
		for i := range fields {
			if (fields[i].IsPrimaryKey && fields[i].IsBlank) || (fields[i].IsIgnored) {
				continue
			}
			placeholders = append(placeholders, scope.AddToVars(fields[i].Field.Interface()))
		}
		placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
		placeholdersArr = append(placeholdersArr, placeholdersStr)
		// add real variables for the replacement of placeholders' '?' letter later.
		mainScope.SQLVars = append(mainScope.SQLVars, scope.SQLVars...)
	}

	mainScope.Raw(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		mainScope.QuotedTableName(),
		strings.Join(quoted, ", "),
		strings.Join(placeholdersArr, ", "),
	))

	if _, err := mainScope.SQLDB().Exec(mainScope.SQL, mainScope.SQLVars...); err != nil {
		return err
	}
	return nil
}

func init() {
	// 初始化OrmInstance 实例
	ormInstance = new(OrmInstance)
}
