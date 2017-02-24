#Sync adapter

    This project has three parts mobile client, sync adapter and a web server 

#Goals: 

    1) The sync adapter should act as a pluggable adapter to an already operating system. At some point of time if we add/remove the adapter, the old code should work with minimal change. 
 
                live setup    : mobile client ----> controller ----> api ----> network ----> web server
                offline setup : mobile client ----> controller ----> adapter ----> api ----> network ----> web server
                
#Rule of thumb

    1) Server table should contain : id & updated column
    2) Client table must implement localmodel struct
                    
#How to implement the sync adapter with the existing system ?
 
    1) Existing models in the client system should embed the localmodel from the syncadapter
 
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
                            


 
 
 
