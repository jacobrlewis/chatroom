package shared

import (
	"fmt"
	"net/url"
)

func GetInitEndpoint() string {
	return "/welcome"
}

func GetInitUrl(host string) string {
	return fmt.Sprintf("http://%s%s", host, GetInitEndpoint())
}

func GetRoomWelcomeEndpoint(room string) string {
	return fmt.Sprintf("/room/%s/welcome", room)
}

func GetRoomWelcomeUrl(host string, room string) string {
	return fmt.Sprintf("http://%s%s", host, GetRoomWelcomeEndpoint(room))
}

func GetRoomSendChatEndpoint(room string) string {
	return fmt.Sprintf("/room/%s", room)
}

func GetRoomWsEndpoint(room string) string {
	return fmt.Sprintf("/room/%s/ws", room)
}

func GetRoomWsUrl(host string, room string) url.URL {
	return url.URL{Scheme: "ws", Host: host, Path: GetRoomWsEndpoint(room)}
}
