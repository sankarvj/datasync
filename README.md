#Sync adapter

    This project helps maintain & sync localdata with the remote server
    It stores data locally and update the database when something changes in the remote server.

    We categorize sync logic into : general sync, specific sync, remote sync and trash sync

    Specific Sync : sync a specific row from the local db to the remote server 
        (All api requests will come under this)
    General Sync : sync data from the local db to the remote server on the event of network toggle / Time based / Manual user req
        (Whatever fails in the "Specific Sync" will be aggregated and sync in one shot)
    Remote Sync : sync on request from the server 
        (Sync via push notification)
    Trash Sync : sync deleted rows from the local to remote and vice versa 
        (It uses some peculiar logics from the rest of the above)
        
    In otherways they are categorized by CRUD
    A) Specific Sync - (POST,PUT) (LOCAL ---> SERVER)
    B) General Sync - (POST,PUT) (LOCAL --> SERVER)
    C) Remote Sync - (POST,PUT) (SERVER --> LOCAL)
    D) Specific Sync - (GET)  (SERVER --> LOCAL)
    E) Trash Sync - (DELETE) (LOCAL ---> SERVER) && (SERVER --> LOCAL)

    So basically : Whatever missed in A will be handled by B
                   Whatever missed in C will be handled by D
    

#Goals: 

    1) The sync adapter should act as a pluggable adapter to an already operating system. At some point of time if we add/remove the adapter, the old code should work with minimal change. 
 
                live setup    : mobile client ----> controller ----> api ----> network ----> web server
                offline setup : mobile client ----> controller ----> adapter ----> api ----> network ----> web server
    
    2) All local_id to server_id and server_id to local_id conversions should happen in one place inside the adapter.
    
    3) In the above pipe, before the "adapter" the object should maintain local scope and after that the object should maintain server scope.
                
#Problems

    1)  Forign key reference in the server table should be modified with its corresponding local_id.
        Also upon the api request, the local instance should bring back remote id to form the api params.
        
    2)  Code repeatation during each request should be handled.

#Rule of thumb

    1) Server table should contain : id & updated column

    2) Client table must implement "localmodel" struct

    3) If a column in a table references a id from the other table than that column must be tagged like the below structure 
    
                            type Note struct {
                            	Ticketid int64 `rt:"trips" rk:"id"`
                            	Name     string
                            	Desc     string
                            	created  time.Time
                            	adapter.Localmodel
                            }
                            
        Note : ticketid column of notes table references id column of Ticket table.
    
                    
#How to implement the sync adapter with the existing system ?
 
    1) Existing models in the client system should inherit the localmodel from the syncadapter
 
                            type Localmodel struct {
                            	Id      int64    //local id
                            	Key     int64    //server id
                            	Updated int64    //last updated time
                            	Synced  bool     //synced or not
                            	Baseids []Baseid //forignkey ids
                            }
                            
                            type Ticket struct {
                            	Subject   string
                            	Desc      string
                            	requester string
                            	agent     string
                            	created   time.Time
                            	adapter.Localmodel //Embed localmodel
                            }
                            
    2) So that all the methods declared under Localmodel is promoted to be accessed via other models inherit it
 
                            ticket.MarkAsLocal()
                            ticket.UpdateLocalId(id int64)
                            etc...
                            
    3) Invidual models can override promoted methods in case if any modification needed
                            
                            func (obj *Localmodel) UpdateLocalId(id int64) {
                            	obj.Id = id
                            }
                            
                            func (obj *Ticket) UpdateLocalId(id int64) {
                            	obj.Id = id * 10
                            }
    
#How it works 

#Create

    1) create the object
    2) store it in the db
    3) send the response back to client
    4) hit the api with the created local object
    5) onsuccess - update serverid and forign ids, set synced as true and call client with response updated 
    6) onerror   - if network error do nothing, if server fails delete the row and update the client

#Update

#Get

#Delete

