package main

import (
	"os/user"
)

// Получить UID данного пользователя (Example: root - UID = 0)
var currentUserUid = func() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	/*
		sudoUid, exists := os.LookupEnv("SUDO_UID")
		if exists {
			return sudoUid
		}
	*/
	return u.Uid
}()

/*
// Получить LoginName данного пользователя (example: root)
var currentUserLogin = func() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	sudoRealLogin, exists := os.LookupEnv("SUDO_USER") // Example: 'locadmin'
	if exists {
		return sudoRealLogin
	}
	return u.Username
}()
*/
