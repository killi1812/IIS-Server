package apiq

type InstagramUsername struct {
	Status   bool   `json:"status" xml:"status"`
	Username string `json:"username"`
	UserId   string `json:"user_id"`
	Attempts string `json:"attempts"`
}

//Example
/*
{
    "status": true,
    "username": "javan",
    "user_id": "18527",
    "attempts": "1"
}
*/
