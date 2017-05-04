
#Sync adapter

    Performs datasync between client and server with some minimal design changes in the database. In the flow of data fetching it extracts the invariant code needed for sync logic in to a common place. 
    
    User has to write their own implementation to store local data and retrive server data. In between that they can call their methods via this adapter to achive the sync.
1GeMxb9cgLpHkXJOmDvy90A3BQn-p774IRR40ng65Ki8
        oldsetup : mobile_client ----> controller ----> api ----> network ----> web_server
        newsetup : mobile_client ----> controller ----> adapter ----> api ----> network ----> web_server

    Adapter will run any one of the following sync techniques based up on the situation  

    Basic  : Run a sync adapter right after the dataset changes in the client.
        (This option is straightforward to implement. All api requests will come under this)
    Periodic  : Allows automate data transfer based on a variety of criteria, including network changes, elapsed time, or time of day.
        (Whatever fails in the "Specific Sync" would be aggregated and it will be synced in one shot)
    Impulse  : Run the sync adapter in response to a message from a server, indicating that server-based data has changed.
        (Sync via push notification, Synchronize and store data before it is needed)
    Eject  : Sync deleted rows from the local to remote and vice versa 
        (It uses some peculiar logics from the rest of the above) - It has some edge cases :(
    Arrest : Reduce periodic sync occurances
    
#Goals:
    
    1) Automated network checking & Automated execution

    2) Remote Sync tend to send the push notification to all devices at the same time. This situation can cause multiple instance of your sync adapter to run at the same time, causing server and network overload. Avoid this situation by starting the "Remote Sync" specific to a device.
        
    3) Android sync adapter is only in charge of running your code, it has no idea how your data should be synced with the server. But this sync adapter is somewhat intelligent since it knows some basic context of your model and by using that it should take care of sync operations of its own.

    4) If someone break rules, handle it gracefully and intimate the developer with meaningful log messages.

#5Rules to follow

    1) Server table should contain id & updated column - updated column must use UTC

    2) Client tablename be the plural of its model name & colname should be as same as model field name (Could remove this check if a utility added to create tables from the model)

    3) Client model must embed "basemodel" struct (see detail in #HowToImplementSection1)

    4) Client model must add tags to the column if that column has any reference (see detail in #HowToImplementSection4)

#Problems
        
    1) Deciding the priority of action (Local/Server) during "Periodic Sync". Right now it will give the priority to local

    2) Too many rules, reduce it.

    3) Data conflits. What if a agent added a note in offline but when it tries to sync that ticket is resolved already.

                    
#How to implement the sync adapter with the existing system ?
 
    1) Existing models in the client system should inherit the basemodel from the syncadapter
 
                            type BaseModel struct {
                            	Id      int64    //local id
                            	Key     int64    //server id
                            	Updated int64    //last updated time
                            	Synced  bool     //synced or not
                            }
                            
                            type Ticket struct {
                            	Subject   string
                            	Desc      string
                            	requester string
                            	agent     string
                            	created   time.Time
                            	core.BaseModel //Embed BaseModel
                            }
                            
    2) So that all the methods declared under BaseModel is promoted to be accessed via other models inherit it
 
                            ticket.PrepareLocal()
                            ticket.SetLocalId(id int64)
                            etc...
                            
    3) Invidual models can override promoted methods in case if any modification needed
                            
                            func (obj *BaseModel) SetLocalId(id int64) {
                            	obj.Id = id
                            }
                            
                            func (obj *Ticket) SetLocalId(id int64) {
                            	obj.Id = id * 10
                            }

    4) If a column in a table references a id from the other table than that column must be tagged like the below structure 
    
                            type Note struct {
                                Ticketid int64 `rt:"tickets" rk:"id"`
                                Name     string
                                Desc     string
                                created  time.Time
                                core.BaseModel
                            }
                            
        Note : ticketid column of notes table references id column of Ticket table. 
    


#Basically

    #Create
    Create a model --> Update the core keys --> Save the model to local db --> Update the Id --> Send to Client --> Heat the model --> Call API --> Update the core keys

    #Edit 
    Create a model --> Update the core keys --> Save the model to local db --> Send to Client --> Heat the model --> Call API --> Update the core keys

    #Read
    Parse to a model --> Cool the model --> Check what to do --> Update/Create/Nothing 
    
#Retry Logic
    A) Basic Sync - (POST,PUT) (LOCAL --> SERVER)
    B) Periodic Sync - (POST,PUT) (LOCAL --> SERVER)
    C) Remote Sync - (POST,PUT) (SERVER --> LOCAL)
    D) Basic Sync - (GET) 
    E) Erase Sync - (DELETE) (LOCAL --> SERVER) && (SERVER --> LOCAL)

    So basically : Whatever missed in A will be handled by B
                   Whatever missed in C will be handled by D


#EdgeCases 

    1) What if the reference of one table points to other table and local db don't have that table ? (Say ticket table has assignee id as a column and in the server side assignee id is the userid of users table)

    2) If an API updates two model in the server. (Create Ticket API create a ticket and adds the current user to assignee table implicitly)
       Is that format is good REST API format ? Or is it not ? 

    3) What if it really updates the server end and fails to update the changes in the localdb ?

    4) What if user changes his time in the local device ?

        Solution 1) Server will update the local time on each response.

    5) Conflict of two users needs to be shown in UI. How ? What is the purpose ?

    6) Tags vs Interface (Tags is too brutal.Use Interface in each class to convert server to db item. But I think tags are best)

    7) Periodic needs to include GET first and then only it should post,put


#Useful Links

https://docs.microsoft.com/en-us/azure/app-service-mobile/app-service-mobile-offline-data-sync
https://docs.google.com/document/d/1GeMxb9cgLpHkXJOmDvy90A3BQn-p774IRR40ng65Ki8/pub
https://github.com/JohnGoodstadt/mobileSyncIOS
https://culturedcode.com/things/blog/2010/12/state-of-sync-part-1/


