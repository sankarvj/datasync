//Package adapter wraps common behaviour of sync operations.
package adapter

import (
	"log"
	"time"
)

//References for forign keys
type Reference struct {
	tablename string
	colname   string
}

//shaper interface
type Shaper interface {
	MarkAsLocal()
	UpdateLocalId(id int64)
	PushReferences(table string, col string)
}

type Localmodel struct {
	Id         int64       //local id
	Key        int64       //server id
	Updated    int64       //last updated time
	Synced     bool        //synced or not
	References []Reference //forignkey references
}

func (obj *Localmodel) MarkAsLocal() {
	obj.Synced = false
	obj.Updated = currentTime()
}

func (obj *Localmodel) UpdateLocalId(id int64) {
	obj.Id = id

	log.Println(" references ", obj.References)
}

func (obj *Localmodel) PushReferences(ref_table string, ref_key string) {
	reference := Reference{ref_table, ref_key}
	if obj.References == nil {
		obj.References = make([]Reference, 0)
	}
	obj.References = append(obj.References, reference)
}

func currentTime() int64 {
	return milliSeconds(time.Now())
}

func milliSeconds(now time.Time) int64 {
	return now.UTC().Unix() * 1000
}
