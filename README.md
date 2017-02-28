#Sync adapter

    A cross platform utility encapsulates the code for data sync between the client and the remote server.

    #Datasync between servers is taken care of : General sync, Specific sync, Remote sync and Erase sync.
    
    Specific Sync : Run a sync adapter when data changes on the device.
        (This option is straightforward to implement. All api requests will come under this.)
    General Sync : Allows you to automate data transfer based on a variety of criteria, including network changes, elapsed time, or time of day.
        (Whatever fails in the "Specific Sync" would be aggregated and it will be synced in one shot. )
    Remote Sync : Run the sync adapter in response to a message from a server, indicating that server-based data has changed.
        (Sync via push notification)
    Erase Sync : Sync deleted rows from the local to remote and vice versa 
        (It uses some peculiar logics from the rest of the above)
        
    In otherways they are categorized by CRUD
    A) Specific Sync - (POST,PUT) (LOCAL --> SERVER)
    B) General Sync - (POST,PUT) (LOCAL --> SERVER)
    C) Remote Sync - (POST,PUT) (SERVER --> LOCAL)
    D) Specific Sync - (GET)  (SERVER --> LOCAL)
    E) Erase Sync - (DELETE) (LOCAL --> SERVER) && (SERVER --> LOCAL)

    So basically : Whatever missed in A will be handled by B
                   Whatever missed in C will be handled by D
    
#Goals: (Few points taken from the android sync adater)
    
    1) Plug-in architecture. The sync adapter should act as a adapter to an already operating system. At some point if I remove the adapter, the old code should work. I mean it should hit the API directly and bring back the results.
 
                live setup    : mobile client ----> controller ----> api ----> network ----> web server
                offline setup : mobile client ----> controller ----> adapter ----> api ----> network ----> web server

    2) Automated network checking & Automated execution

    3) Remote Sync tend to send the push notification to all devices at the same time. This situation can cause multiple instance of your sync adapter to run at the same time, causing server and network overload. Avoid this situation by starting the "Remote Sync" specific to a device.
    
    4) All local_id to server_id and server_id to local_id conversions should happen in one place inside the adapter.
    
    5) In the above pipe, before the "adapter" the object should maintain local scope and after that the object should maintain server scope.

    6) Android sync adapter is only in charge of running your code, it has no idea how your data should be synced with the server. But this sync adapter is somewhat intelligent since it knows some basic context of your model and by using that it should take care sync operations of its own.
    
                
#Problems

    1) Forign key reference in the server table should be modified with its corresponding local_id. Upon the api request, the local instance should bring back remote id to form the api params.
        
    2) Calculate device specific time to run "Remote Sync" at different intervals across devices to reduce server load.

    3) Sync attachements have problems such as deleting the local copy,rename and client validation has to be done. 

    4) Deciding the priority of action (Local/Server) during "General Sync". Right now it will give the priority to local

#Rules of thumb

    1) Server table should contain id & updated column

    2) Client model must implement "localmodel" struct (see detail in #HowToImplementSection1)

    3) Client model must add tags to the column if any that column reference other table columns (see detail in #HowToImplementSection4)

                    
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

    4) If a column in a table references a id from the other table than that column must be tagged like the below structure 
    
                            type Note struct {
                                Ticketid int64 `rt:"trips" rk:"id"`
                                Name     string
                                Desc     string
                                created  time.Time
                                adapter.Localmodel
                            }
                            
        Note : ticketid column of notes table references id column of Ticket table. 
    
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

