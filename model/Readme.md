## 数据库表目录

+ 下面放各服务目录，每个服务的表放在自己的服务目录下面
+ 统一使用GROM V2,
```golang
import (
    "gorm.io/gorm"
)
``` 

## GORM文档
> https://gorm.io/zh_CN/docs/index.html



## MySQL数据库设计规范
+ 每个表中必须要有主键
+ 用尽量少的存储空间来存数一个字段的数据，能使用tinyint就不用int/bigint，能使用varchar(16)就不用varchar(255)
+ 索引建立在区分度高以及查询命中高的字段上，多字段根据业务建立复合索引(遵循最左前缀)
```golang
//唯一索引
type User struct {
    Name1 string `gorm:"uniqueIndex"`
    Name2 string `gorm:"uniqueIndex:idx_name"`
}

//复合索引
type User struct {
    Name    string  `gorm:"index:idx_name_age,priority:1"`   //优先级越小的建立复合索引时顺序在越前面       
    Age     uint8   `gorm:"index:idx_name_age,priority:2"`   //name与age构成了复合索引idx_name_age(name,age)
}

//全文索引
type User struct {
    Name string `gorm:"index:,class:FULLTEXT,option:WITH PARSER ngram"`
}
```
+ 表与表之间的相关联字段名称要求尽可能的相同
+ 表的各个字段最好设置默认值，不要为NULL，整型可以设置`gorm:"default:0"`，字符串可以设置`gorm:"default:''"`



## GORM操作
### 超时处理
`gorm`提供了`WithContext()`方法，因此对于长 Sql 查询，你可以传入一个带超时的 `context` 给 `db.WithContext` 来设置超时时间
```golang
//单会话
db.WithContext(ctx).Find(&users)

//持续会话
tx := db.WithContext(ctx)
tx.First(&user, 1)
tx.Model(&user).Update("role", "admin")

//超时处理
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

db.WithContext(ctx).Find(&users)
```

### 事务
数据库事务，保证数据的一致性，用法有两种:
> 注意：事务一旦开始，就应该使用事务方法返回/生成的DB，如下为tx
```golang
//用法1
db.Transaction(func(tx *gorm.DB) error {
    // 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
  if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
    // 返回任何错误都会回滚事务
    return err
  }

  if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
    return err
  }

  // 返回 nil 提交事务
  return nil
}

//用法2
// 开始事务，事务内都要使用tx
tx := db.Begin()

// 遇到错误时回滚事务
if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
    tx.Rollback()
    return 
}

//...

if err := httpcall(); err != nil {
    tx.Rollback()
    return 
}

// 否则，提交事务
tx.Commit()
```

### 表迁移&更新
#### AutoMigrate
可以自动更新表新增的字段/索引，但是不会去删除数据库表中已经存在的字段/索引，使用
```golang
db.AutoMigrate(&User{})
```


#### Migrator 
GORM 提供了 Migrator 接口，该接口为每个数据库提供了统一的 API 接口，可用来为您的数据库构建独立迁移，例如：
+ 判断表是否存在，创建/删除/重命名表
+ 判断表中的列是否存在，创建/删除/修改/重命名列
+ 判断表中的索引是否存在，创建/删除/重命名索引
```golang
type Migrator interface {
  // AutoMigrate
  AutoMigrate(dst ...interface{}) error

  // Database
  CurrentDatabase() string
  FullDataTypeOf(*schema.Field) clause.Expr

  // Table
  CreateTable(dst ...interface{}) error
  DropTable(dst ...interface{}) error
  HasTable(dst interface{}) bool
  RenameTable(oldName, newName interface{}) error
  GetTables() (tableList []string, err error)

  // Columns
  AddColumn(dst interface{}, field string) error
  DropColumn(dst interface{}, field string) error
  AlterColumn(dst interface{}, field string) error
  MigrateColumn(dst interface{}, field *schema.Field, columnType ColumnType) error
  HasColumn(dst interface{}, field string) bool
  RenameColumn(dst interface{}, oldName, field string) error
  ColumnTypes(dst interface{}) ([]ColumnType, error)

  // Constraints
  CreateConstraint(dst interface{}, name string) error
  DropConstraint(dst interface{}, name string) error
  HasConstraint(dst interface{}, name string) bool

  // Indexes
  CreateIndex(dst interface{}, name string) error
  DropIndex(dst interface{}, name string) error
  HasIndex(dst interface{}, name string) bool
  RenameIndex(dst interface{}, oldName, newName string) error
}
```

### Hint
+ 强制选择索引
```golang
import "gorm.io/hints"

db.Clauses(hints.UseIndex("idx_user_name")).Find(&User{})
// SELECT * FROM `users` USE INDEX (`idx_user_name`)
```