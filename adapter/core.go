//Package adapter wraps common behaviour of sync operations.
package adapter

type BaseModel struct {
	Id      int64 //local id
	Key     int64 //server id
	Updated int64 //last updated time
	Synced  bool  //synced or not
}

//Cooker interface
type Cooker interface {
	MarkAsLocal()
	MarkAsPureLocal()
	UpdateLocalId(id int64)
}

//Cooker implementations
func (basemodel *BaseModel) MarkAsLocal() {
	if basemodel.Id == 0 { //storing ticket originally created at client
		basemodel.MarkAsPureLocal()
	} else { //storing ticket originally created at server
		basemodel.Synced = true
	}

}

func (basemodel *BaseModel) MarkAsPureLocal() {
	basemodel.Synced = false
	basemodel.Updated = currentTime()
}

func (basemodel *BaseModel) UpdateLocalId(id int64) {
	basemodel.Id = id
}

//Passer interface
type Passer interface {
	GetServerId() int64
	GetUpdatedAt() int64
	GetLSId() int64
}

//Passer implementation
func (basemodel BaseModel) GetServerId() int64 {
	return basemodel.Key
}

func (basemodel BaseModel) GetUpdatedAt() int64 {
	return basemodel.Updated
}

func (basemodel BaseModel) GetLSId() int64 {
	return basemodel.Id
}
