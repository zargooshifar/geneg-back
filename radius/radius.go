package radius

import (
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"log"
	"msgv2-back/database"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/models"
)

func Setup() {
	handler := func(w radius.ResponseWriter, r *radius.Request) {
		username := rfc2865.UserName_GetString(r.Packet)
		password := rfc2865.UserPassword_GetString(r.Packet)

		var code radius.Code

		user := new(models.User)

		if count := database.DB.Where(&models.User{Username: username}).First(&user).RowsAffected; count == 0 {
			code = radius.CodeAccessReject
		} else
		//check password
		if !utils.VerifyPassword(user.Password, password) {
			code = radius.CodeAccessReject
		} else {
			code = radius.CodeAccessAccept
		}

		log.Printf("Writing %v to %v", code, r.RemoteAddr)
		w.Write(r.Response(code))
	}

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(`secret`)),
	}
	log.Printf("Starting Radius server on :1812")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
