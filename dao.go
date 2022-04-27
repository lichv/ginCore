package ginCore

type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

type Dao struct {
	Containter *Container
}

type DaoInterface interface {
	GetContainter() *Container
	SetContainter(ct *Container) *Dao
	GetList(query map[string]interface{}, orderby string, limit int, offset int) ([]*BaseModel, error)
	GetPage(query map[string]interface{}, orderby string, page int, size int) ([]*BaseModel, int, int, error)
	GetOne(query map[string]interface{}, orderby string) (*BaseModel, int, int, error)
	Find(query map[string]interface{}, orderby string) (*BaseModel, error)
	FindWithField(column string, value interface{}) (*BaseModel, error)
	Insert(m *BaseModel) (int64, error)
	Modify(id int, m *BaseModel) (int64, error)
	ModifyField(id uint, field string, value interface{}) (int64, error)
	Update(query map[string]interface{}, m *BaseModel) (int64, error)
	Delete(id int) (int64, error)
	Clear(query map[string]interface{}) (int64, error)
}

func (d *Dao) GetContainter() *Container {
	return d.Containter
}

func (d *Dao) SetDB(containter *Container) *Dao {
	d.Containter = containter
	return d
}
