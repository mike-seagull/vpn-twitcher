package main

type NotificationPlugin interface {
	SendMessage(msg string)
	Init(options map[string]string)
}
