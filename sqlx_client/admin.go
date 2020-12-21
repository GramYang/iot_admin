package sqlx_client

func GetAdminByName(name string) (*Admin,error){
	var admin Admin
	err:=db.Get(&admin,"select * from `admin` where `user_name`=?",name)
	if err!=nil{
		return nil,err
	}
	return &admin,err
}