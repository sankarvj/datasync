//Package adapter wraps common behaviour of sync operations.
package adapter

type BaseModel struct {
	Id      int64 //local id
	Key     int64 //local id
	Updated int64 //last updated time
	Synced  bool  //synced or not
}

//Cooker interface
type Cooker interface {
	MarkAsLocal()
	UpdateLocalId(id int64)
	GetLocalId() int64
}

//Cooker implementations
func (basemodel *BaseModel) MarkAsLocal() {
	if basemodel.Id == 0 { //storing ticket originally created at client
		basemodel.Synced = false
		basemodel.Updated = currentTime()
	} else { //storing ticket originally created at server
		basemodel.Synced = true
	}

}

func (basemodel *BaseModel) UpdateLocalId(id int64) {
	basemodel.Id = id
}

func (basemodel *BaseModel) GetLocalId() int64 {
	return basemodel.Id
}

//Passer interface
type Passer interface {
	GetServerId() int64
	GetUpdated() int64
	GetId() int64
}

//Passer implementation
func (basemodel BaseModel) GetServerId() int64 {
	return basemodel.Key
}

func (basemodel BaseModel) GetUpdated() int64 {
	return basemodel.Updated
}

func (basemodel BaseModel) GetId() int64 {
	return basemodel.Id
}
