package shared

import (
	"fmt"
)

func GetInitEndpoint() string {
	return "/welcome"
}

func GetRoomWelcomeEndpoint(id string) string {
	return fmt.Sprintf("/room/%s/welcome", id)
}

func GetRoomSendChatEndpoint(id string) string {
	return fmt.Sprintf("/room/%s", id)
}