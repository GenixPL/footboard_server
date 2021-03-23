package game

import u "footboard_server/utils"

// Removes client at given index.
//
// Will also clear player fields if needed.
func (game *Game) RemoveClientAtIndex(index int) {
	u.LogV("RemoveClientAtIndex", "removing client at index: "+string(index))

	client := game.Clients[index]
	client.Connection.Close()
	game.Clients = append(game.Clients[:index], game.Clients[index+1:]...)

	// Clear players field if needed.
	if game.Player1.Id == client.Id {
		game.Player1 = nil
	} else if game.Player2.Id == client.Id {
		game.Player2 = nil
	}
}

// Removes client with provided id.
//
// Will also clear player fields if needed.
func (game *Game) RemoveClientWithId(id string) {
	u.LogV("RemoveClientWithId", "removing client with id: "+string(id))

	// Clear players field if needed.
	if game.Player1.Id == id {
		game.Player1 = nil
	} else if game.Player2.Id == id {
		game.Player2 = nil
	}

	// Get index of the client that we want to remove.
	index := -1
	for i, client := range game.Clients {
		if client.Id == id {
			index = i
			break
		}
	}

	// No such client.
	if index == -1 {
		u.LogV("RemoveClientWithId", "client with id: "+id+" doesn't exist")
		return
	}

	game.RemoveClientAtIndex(index)
}
