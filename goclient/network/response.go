package network

type Response struct {
	Id      int16 //Basically 0 for success and any other value is error
	Msg     string
	Outcome []interface{}
}

/****
* IMPORTANT - Each API in this project should use this method to send json back to the client
* this will make the response so common across the client end
*
* VERY IMPORTANT - Add new objects to the end, because adding in front will break
* the excisting flow
******/
func MakeResponse(id int16, msg string, results ...interface{}) Response {
	// for key, value := range results {
	// 	fmt.Println("key :::::::::::::::::::::::::::::: ", key)
	// 	fmt.Println("value :::::::::::::::::::::::::::::: ", value)
	// }
	response := Response{
		Id:      id,
		Msg:     msg,
		Outcome: results,
	}
	return response
}

const (
	ResponseSuccess        = 1
	ResponseError          = 1000
	ResponsePrivilegeError = 5000
	ResponseAuthError      = 6000
	ResponseNetworkError   = 7000
	ResponseInternalError  = 8000
)
