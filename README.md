
#Sync adapter

    Performs datasync between client and server by some simple design and it helps the developer by extracting the monotonous code needed for sync in a single place.
    User has to write their own implementation to store local data and retrive server data. In between that they can call their methods via this adapter to achive the sync.

        oldsetup : mobile_client ----> controller ----> api ----> network ----> web_server

        newsetup : mobile_client ----> controller ----> adapter ----> api ----> network ----> web_server

    Adapter will run any one of the following sync logics based up on the situation  

    Specific Sync : Run a sync adapter when data changes on the device.
        (This option is straightforward to implement. All api requests will come under this)
    General Sync : Allows automate data transfer based on a variety of criteria, including network changes, elapsed time, or time of day.
        (Whatever fails in the "Specific Sync" would be aggregated and it will be synced in one shot)
    Remote Sync : Run the sync adapter in response to a message from a server, indicating that server-based data has changed.
        (Sync via push notification)
    Erase Sync : Sync deleted rows from the local to remote and vice versa 
        (It uses some peculiar logics from the rest of the above) - It has some edge cases :(
    
#Goals:
    
    1) Automated network checking & Automated execution

    2) Remote Sync tend to send the push notification to all devices at the same time. This situation can cause multiple instance of your sync adapter to run at the same time, causing server and network overload. Avoid this situation by starting the "Remote Sync" specific to a device.
        
    3) Android sync adapter is only in charge of running your code, it has no idea how your data should be synced with the server. But this sync adapter is somewhat intelligent since it knows some basic context of your model and by using that it should take care of sync operations of its own.

    4) If someone break rules, handle it gracefully and intimate the developer with meaningful log messages.
    
                
#Problems

    1) Forign key reference in the server table should be modified with its corresponding local_id. Upon the api request, the local instance should bring back the remote_id to form the api params.
        
    2) Sync attachements have problems such as deleting the local copy,rename and client validation has to be done. 

    3) Deciding the priority of action (Local/Server) during "General Sync". Right now it will give the priority to local

    4) Too many rules, reduce it


#5Rules to follow

    1) Server table should contain id & updated column

    2) Client tablename be the plural of its model name & colname should be as same as model field name (Could remove this check if a utility added to create tables from the model)

    3) Client model must embed "basemodel" struct (see detail in #HowToImplementSection1)

    4) Client model must add tags to the column if that column has any reference (see detail in #HowToImplementSection4)

                    
#How to implement the sync adapter with the existing system ?
 
    1) Existing models in the client system should inherit the basemodel from the syncadapter
 
                            type BaseModel struct {
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
                            	adapter.BaseModel //Embed BaseModel
                            }
                            
    2) So that all the methods declared under BaseModel is promoted to be accessed via other models inherit it
 
                            ticket.MarkAsLocal()
                            ticket.UpdateLocalId(id int64)
                            etc...
                            
    3) Invidual models can override promoted methods in case if any modification needed
                            
                            func (obj *BaseModel) UpdateLocalId(id int64) {
                            	obj.Id = id
                            }
                            
                            func (obj *Ticket) UpdateLocalId(id int64) {
                            	obj.Id = id * 10
                            }

    4) If a column in a table references a id from the other table than that column must be tagged like the below structure 
    
                            type Note struct {
                                Ticketid int64 `rt:"trips" rk:"id"`
                                Name     string
                                Desc     string
                                created  time.Time
                                adapter.BaseModel
                            }
                            
        Note : ticketid column of notes table references id column of Ticket table. 
    


#Others

    A) Specific Sync - (POST,PUT) (LOCAL --> SERVER)
    B) General Sync - (POST,PUT) (LOCAL --> SERVER)
    C) Remote Sync - (POST,PUT) (SERVER --> LOCAL)
    D) Specific Sync - (GET) 
    E) Erase Sync - (DELETE) (LOCAL --> SERVER) && (SERVER --> LOCAL)

    So basically : Whatever missed in A will be handled by B
                   Whatever missed in C will be handled by D

