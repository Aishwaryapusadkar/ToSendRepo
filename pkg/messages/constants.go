package messages

const (
	//pkt size
	ResponseHeaderSize      = 8
	ConnectionGatewayPktLen = uint32(96)
	SessionLoginPktLen      = uint32(280)
	UserLoginPktLen         = uint32(64)
	RetransmissionPktLen    = uint32(48)
	SubscriptionPktLen      = uint32(32)

	//mType
	CONNECTION_GATEWAY = "connection-gateway"
	SESSION_LOGIN      = "session-login"
	USER_LOGIN         = "user-login"
	RETRANSMISSION     = "retransmission"
	SUBSCRIPTION       = "subscription"
	TRADE              = "TRADE"

	//REQUEST TEMPLATE_ID
	RejectedTemplateId                 = uint16(10010)
	ConnetionGatewayRequestTemplateId  = uint16(10020)
	ConnetionGatewayResponseTemplateId = uint16(10021)
	SessionLoginRequestTemplateId      = uint16(10000)
	SessionLoginResponseTemplateId     = uint16(10001)
	UserLoginRequestTemplateId         = uint16(10018)
	UserLoginResponseTemplateId        = uint16(10019)
	RetransmissionRequestTemplateId    = uint16(10008)
	RetransmissionResponseTemplateId   = uint16(10009)
	SubscriptionRequestTemplateId      = uint16(10025)
	SubscriptionResponseTemplateId     = uint16(10005)
	HeartbeatTemplateId                = uint16(10023)
	TradeTemplateId                    = uint16(10500)
	SessionLogoutTemplateId            = uint16(10012)
)

var (
	AccountTypeMapping = map[uint]string{20: "OWN", 30: "CLIENT", 0: "UNKNOWN"}
	OrderTypeMapping   = map[uint8]string{2: "LIMIT", 3: "STOP MARKET", 4: "STOP LIMIT", 5: "MARKET", 6: "BLOCK DEAL"}
	BuySellMapping     = map[uint]string{1: "B", 2: "S", 3: "RECALL", 4: "EARLY RETURN"}
	AOPOFlagMapping    = map[uint8]string{1: "0", 2: "1"}
	EmptyCPCCode       = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	CSVHeaders         = []string{ // please add new field in last to maintained the sequence
		"Member Id",
		"Trader ID",
		"Session ID",
		"Client ID",
		"Client Type",
		"CP Code",
		"CP confirmation Code",
		"Order ID",
		"Trade ID",
		"Scrip Code",
		"Rate",
		"Qty",
		"Buy/Sell",
		"AO/PO Flag",
		"Location ID",
		"Order Time Stamp",
		"Time",
		"Date",
		"Order Type",
		"Trade Modification Time",
	}
)
