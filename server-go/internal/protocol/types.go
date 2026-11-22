package protocol

// Types de paquets du protocole applicatif JSON length-prefixed (Ndue).
type MessageType string

const (
	TypeRegister         MessageType = "REGISTER"
	TypeLogin            MessageType = "LOGIN"
	TypeSendMessage      MessageType = "SEND_MESSAGE"
	TypeGetMessages      MessageType = "GET_MESSAGES"
	TypeMessageDelivered MessageType = "MESSAGE_DELIVERED"
	TypeMessageRead      MessageType = "MESSAGE_READ"
	TypeCreateGroup      MessageType = "CREATE_GROUP"
	TypeJoinGroup        MessageType = "JOIN_GROUP"
	TypeGetAnnouncements MessageType = "GET_ANNOUNCEMENTS"
	TypeNotification     MessageType = "NOTIFICATION"
	TypePresence         MessageType = "PRESENCE"
	TypeCallSignal       MessageType = "CALL_SIGNAL"
	TypePing             MessageType = "PING"
	TypePong             MessageType = "PONG"
	TypeLogout           MessageType = "LOGOUT"
	TypeError            MessageType = "ERROR"
)
