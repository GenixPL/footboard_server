package main

type Client struct {
	// TODO: add socket stuff here

	// If it's an empty string then the object should be treated as nil.
	id string

	secondsLeft uint16
}
