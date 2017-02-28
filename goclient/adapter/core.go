//Package adapter wraps common behaviour of sync operations.
package adapter

import (
	"database/sql"
	"log"
)

//References for forign keys
type Reference struct {
	tablename string
	colname   string
	ownname   string
	localkey  string
	serverkey string
}

//shaper interface
type Shaper interface {
	MarkAsLocal()
	UpdateLocalId(id int64)
	PushReferences(table string, col string, own_col string)
	MarkAsServer(db *sql.DB) error
	GetReferences() *[]Reference
}

type BaseModel struct {
	Id         int64       //local id
	Updated    int64       //last updated time
	Synced     bool        //synced or not
	References []Reference //forignkey references
}

//Shaper implementations
func (basemodel *BaseModel) MarkAsLocal() {
	basemodel.Synced = false
	basemodel.Updated = currentTime()
}

func (basemodel *BaseModel) UpdateLocalId(id int64) {
	basemodel.Id = id
}

func (basemodel *BaseModel) PushReferences(ref_table string, ref_col string, own_col string) {
	if basemodel.References == nil {
		basemodel.References = make([]Reference, 0)
	}
	basemodel.References = append(basemodel.References, Reference{ref_table, ref_col, own_col, "", ""})
}

func (basemodel *BaseModel) MarkAsServer(db *sql.DB) error {
	err := coreReference(db, basemodel)
	if err == nil {
		log.Println("references ", *basemodel.GetReferences())
	}
	return err
}

func (basemodel *BaseModel) GetReferences() *[]Reference {
	return &basemodel.References
}
