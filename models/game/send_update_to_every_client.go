package game

import (
	u "footboard_server/models/utils"
)

func (game *Game) SendUpdateToEveryClient() {
	gameJsonString, err := game.ToJsonString()
	if err != nil {
		for _, client := range game.Clients {
			err := client.Connection.WriteMessage(1, u.JsonedErr(errorCouldntParseGameToJson))
			if err != nil {
				u.LogE("InformEveryClient", "sending msg failed")
				continue
			}
		}
		return
	}

	for _, client := range game.Clients {
		u.LogV("InformEveryClient", "informing client: "+client.Id)

		err := client.Connection.WriteMessage(1, u.JsonedMsg(gameJsonString))
		if err != nil {
			u.LogE("InformEveryClient", "sending msg failed")
			continue
		}
	}

}
