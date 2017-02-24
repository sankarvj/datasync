#Sync adapter

This project has three parts mobile client, sync adapter and a web server 

Goals : 

 1)  The sync adapter should act as a pluggable adapter to an already operating wired system. At some point of time if we add/remove the adapter, the old code should work with minimum change. 
 
                live setup    : mobile client ----> controller ----> api ----> network ----> web server
                offlibe setup : mobile client ----> controller ----> adapter ----> api ----> network ----> web server
                    
 How to implement the sync adapter with the existing system ?
 
 1) Existing models in the old system should implement the basemodel from the syncadapter
 
 2) unplug the connection between controller and api in the old setup and connect the adapter in betwen controller and api
 
                            mobile client ----> controller ----> adapter ----> api ----> network ----> web server
 
 
