package utils

func JsonedErr(err string) []byte {
	msg := "{"

	msg += `"error": "` + err + `", `
	msg += `"game": null`

	msg += "}"

	return []byte(msg)
}

func JsonedMsg(game string) []byte {
	msg := "{"

	msg += `"error": null,`
	msg += `"game": ` + game

	msg += "}"

	return []byte(msg)
}

func JsonedMsgWithUid(game string, uid string) []byte {
	msg := "{"

	msg += `"error": null,`
	msg += `"game": ` + game + ", "
	msg += `"uid": "` + uid + `"`

	msg += "}"

	return []byte(msg)
}

func JsonedErrWithUid(err string, uid string) []byte {
	msg := "{"

	msg += `"error": "` + err + `", `
	msg += `"game": null, `
	msg += `"uid": "` + uid + `"`

	msg += "}"

	return []byte(msg)
}
