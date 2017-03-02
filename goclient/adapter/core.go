//Package adapter wraps common behaviour of sync operations.
package adapter

//shaper interface
type Shaper interface {
	MarkAsLocal()
	UpdateLocalId(id int64)
}

type BaseModel struct {
	Id      int64 //local id
	Updated int64 //last updated time
	Synced  bool  //synced or not
}

//Shaper implementations
func (basemodel *BaseModel) MarkAsLocal() {
	basemodel.Synced = false
	basemodel.Updated = currentTime()
}

func (basemodel *BaseModel) UpdateLocalId(id int64) {
	basemodel.Id = id
}
