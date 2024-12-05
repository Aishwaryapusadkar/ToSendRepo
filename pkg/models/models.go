package models

type MessageRequestHeader struct {
	BodyLen      uint32
	TemplateId   uint16
	NetworkMsgId [8]byte
	Filler       [2]byte
}

type MessageResponseHeader struct {
	BodyLen    uint32
	TemplateId uint16
	Filler     [2]byte
}

type RequestHeader struct {
	MsgSeqNum uint32
	SenderId  uint32
}

type ResponseHeader struct {
	RequestTime [8]uint8
	SendingTime [8]uint8
	MsgSeqNum   uint32
	Filler      [4]byte
}

// ============================================Connection Models
type ConnectionGatewayRequest struct {
	MHeader    MessageRequestHeader
	RHeader    RequestHeader
	SessionId  uint32
	AppVersion [30]byte
	Password   [32]byte
	Filler     [6]byte
}

type ConnectionGatewayResponseBody struct {
	RHeader      ResponseHeader
	GatewayId    [4]uint8
	GatwaysubId  uint32
	SgatewayId   [4]uint8
	SgatwaysubId uint32
	SesionsMode  uint8
	TresMode     uint8
	Filler       [6]byte
}

type ConnectionGatewayResponseObject struct {
	Header MessageResponseHeader
	Body   ConnectionGatewayResponseBody
}

// ============================================Session Models
type SessionLoginRequest struct {
	MHeader               MessageRequestHeader
	RHeader               RequestHeader
	HeartbeatInterval     uint32
	SessionId             uint32
	AppVersion            [30]byte
	Password              [32]byte
	AppUsageOrders        [1]byte
	AppUsageQuotes        [1]byte
	OrderRoutingIndicator [1]byte
	FixEngineName         [30]byte
	FixEngineVersion      [30]byte
	FixEngineVendor       [30]byte
	AppSystemName         [30]byte
	AppSystemVersion      [30]byte
	AppSystemVendor       [30]byte
	Filler                [3]byte
}
type SessionLoginResponseBody struct {
	RHeader                   ResponseHeader
	ThorttleTimeInterval      [8]byte
	LastLoginTime             [8]byte
	LastloginIp               uint32
	ThrottleNoMsgs            uint32
	ThrottleDisconnectLimit   uint32
	HeartbeatInterval         uint32
	SessionInstanceId         uint32
	TresSesMode               uint8
	NoOfPartitions            uint8
	DaysLeftForPasswordExpiry uint8
	GraceLoginLeft            uint8
	AppVersion                [30]byte
	Filler                    [2]byte
}
type SessionLoginResponseObject struct {
	Header MessageResponseHeader
	Body   SessionLoginResponseBody
}

// ============================================Sunbscription Models
type SubscribeRequest struct {
	MHeader           MessageRequestHeader
	RHeader           RequestHeader
	SubscriptionScope uint32
	RefAppId          uint8
	Filler            [3]byte
}
type SubscribeResponseBody struct {
	RHeader  ResponseHeader
	AppSubId uint8
	Filler   [4]byte
}

type SubscribeResponseObject struct {
	Header MessageResponseHeader
	Body   SubscribeResponseBody
}

// ============================================Trade Models
type RBCHeaderSt struct {
	SendingTime      [8]byte
	AppSeqNum        [8]byte
	AppSubId         uint32
	PartionId        int16
	AppResendFlag    uint8
	AppId            uint8
	LastFragmentFlag uint8
	Filler           [7]byte
}

type TradeBody struct {
	RHeader           RBCHeaderSt
	SecurityId        int64
	RelatedSecurityId int64
	Price             int64
	LastPrice         int64
	SideLastPrice     int64
	//ClearingTradePrice         int64   //m
	//Yeild                      int64   //m
	//UnderlyingDirtyPrice       int64    //m
	TransactionTime  uint64
	OrderId          uint64
	SenderLocationId uint64
	CIOrdId          uint64
	//Activitytime               uint64//
	Filler1             uint64 //
	Filler2             uint32
	MsgTag              int32
	TradeId             uint32
	OrigTradeId         uint32
	BusinessUnitId      uint32 //**
	SessionId           uint32
	OwnerUserId         uint32 //**
	PartyIdClearingUnit uint32
	CumQty              int32
	LeavesQty           int32
	MarketSegmentId     int32
	RelatedSymbol       int32
	LastQty             int32
	SideLastQty         int32
	//ClearingTradeQty           int32
	SideTradeId           uint32
	MatchDate             uint32
	TradeMatch            uint32
	StrategyLinkId        uint32
	TotNumTradeReports    int32
	Filler4               uint16
	MultiLegReportingType uint8
	TradeReportType       uint8
	TrasnferReason        uint8
	//RollOverFlag               uint8
	PartyIdBeneficiery         [9]byte
	PartyIdTakeupTradingfirm   [5]byte
	PartyIdOrderOrignatingFirm [7]byte
	AccountType                uint8
	AggresorSide               uint8
	MatchType                  uint8
	MatchSubType               uint8
	Side                       uint8
	AggresorIndicator          uint8
	TradingCapacity            uint8
	Account                    [2]byte
	PositionEffect             [1]byte
	//CustOrderHandlingInst      [1]byte
	//AlgoId                     [16]byte
	//ClientCode                 [12]byte
	//CPCCode                    [12]byte
	FreeText1                 [12]byte
	FreeText2                 [12]byte
	FreeText3                 [12]byte
	OrderCategory             [1]byte
	OrderType                 uint8
	RelatedproductComplex     uint8
	OrderSide                 uint8
	PartyClearingOrganisation [4]byte
	PartyExecutingFirm        [5]byte
	PartyExecutingTrader      [6]byte
	PartyClearingFirm         [5]byte
	Filler5                   [7]byte
}
type TradeObject struct {
	Header MessageResponseHeader
	Body   TradeBody
}

// ============================================User Login Models
type UserLoginRequest struct {
	MHeader  MessageRequestHeader
	RHeader  RequestHeader
	Username uint32
	Password [32]byte
	Filler   [4]byte
}
type UserLoginResponseBody struct {
	RHeader                   ResponseHeader
	LastLoginTime             uint64
	DaysLeftForPasswordExpiry uint8
	GraceLoginLeft            uint8
	Filler                    [6]byte
}
type UserLoginResponseObject struct {
	Header MessageResponseHeader
	Body   UserLoginResponseBody
}

// ============================================Retransmission Models

type ReTransmissionRequest struct {
	MHeader      MessageRequestHeader
	RHeader      RequestHeader
	AppBegnSeqno uint64
	AppEndSeqno  uint64
	Scope        uint32
	PartionId    uint16
	RefAppId     uint8
	Filler       [1]byte
}
type ReTransmissionResponseBody struct {
	RHeader          ResponseHeader
	AppEndSeqno      uint64
	RefAppSeqNo      uint64
	AppTotalMsgCount uint16
	Filler           [6]byte
}
type ReTransmissionResponseObject struct {
	Header MessageResponseHeader
	Body   ReTransmissionResponseBody
}

type PartitionInfo struct {
	Id       uint16
	BegSeqNo uint64
	EndSeqno uint64
}
type RetranmistTradeObject struct {
	UserId    uint32
	Partitons []PartitionInfo
}

// ============================================User Password change Models
// type UserPasswordChangeRequest struct {
// 	MHeader     MessageRequestHeader
// 	RHeader     RequestHeader
// 	Username    uint32
// 	OldPassword [32]byte
// 	NewPassword [32]byte
// 	Filler      [4]byte
// }

// type UserPasswordChangeResponse struct {
// 	MHeader MessageResponseHeader
// 	RHeader ResponseHeader
// 	Filler  [4]byte
// }

// ============================================Reject Models
type RejectBody struct {
	RequestTime      [8]uint8
	RequestOut       [8]uint8
	TrdRegTSTTimeIn  [8]uint8
	TrdRegTSTTimeOut [8]uint8
	ResponseIn       [8]uint8
	SendingTime      [8]uint8
	MsgSeqNum        uint32
	LastFragment     uint8
	Filler1          [3]byte
	Reason           uint32
	TextLen          uint16
	SessionStatus    uint8
	Filler2          [1]byte
}
type RejectObject struct {
	Header MessageResponseHeader
	Body   RejectBody
}

// Example struct for BXT trade data
type Bxttrade struct {
	//RHeader                        RequestHeader
	BxtSecurityId                  int64
	BxtRelatedSecurityId           int64
	BxtPrice                       float64
	BxtLastPrice                   float64
	BxtSideLastPrice               float64
	BxtTransactionTime             int64
	BxtOrderId                     int64
	BxtSenderLocationId            int64
	BxtCIOrdId                     int64
	BxtMsgTag                      int32
	BxtTradeId                     int64
	BxtOrigTradeId                 int64
	BxtBusinessUnitId              int64
	BxtSessionId                   int64
	BxtOwnerUserId                 int64
	BxtPartyIdClearingUnit         int64
	BxtCumQty                      int32
	BxtLeavesQty                   int32
	BxtMarketSegmentId             int32
	BxtRelatedSymbol               int32
	BxtLastQty                     int32
	BxtSideLastQty                 int32
	BxtSideTradeId                 int64
	BxtMatchDate                   int32
	BxtTradeMatch                  int32
	BxtStrategyLinkId              int32
	BxtTotNumTradeReports          int32
	BxtMultiLegReportingType       int32
	BxtTradeReportType             int32
	BxtTransferReason              int32
	BxtPartyIdBeneficiary          string
	BxtPartyIdTakeupTradingFirm    string
	BxtPartyIdOrderOriginatingFirm string
	BxtAccountType                 int32
	BxtAggressorSide               int32
	BxtMatchType                   int32
	BxtMatchSubType                int32
	BxtSide                        int32
	BxtAggressorIndicator          int32
	BxtTradingCapacity             int32
	BxtAccount                     string
	BxtPositionEffect              string
	BxtFreeText1                   string
	BxtFreeText2                   string
	BxtFreeText3                   string
	BxtOrderCategory               string
	BxtOrderType                   int32
	BxtRelatedProductComplex       int32
	BxtOrderSide                   int32
	BxtPartyClearingOrganisation   string
	BxtPartyExecutingFirm          string
	BxtPartyExecutingTrader        string
	BxtPartyClearingFirm           string
	BxtFiller5                     string
}

type TradeHeader struct {
	SecurityId                 int64
	RelatedSecurityId          int64
	Price                      int64
	LastPrice                  int64
	SideLastPrice              int64
	TransactionTime            uint64
	OrderId                    uint64
	SenderLocationId           uint64
	CIOrdId                    uint64
	Filler1                    uint64 //
	Filler2                    uint32
	MsgTag                     int32
	TradeId                    uint32
	OrigTradeId                uint32
	BusinessUnitId             uint32 //**
	SessionId                  uint32
	OwnerUserId                uint32 //**
	PartyIdClearingUnit        uint32
	CumQty                     int32
	LeavesQty                  int32
	MarketSegmentId            int32
	RelatedSymbol              int32
	LastQty                    int32
	SideLastQty                int32
	SideTradeId                uint32
	MatchDate                  uint32
	TradeMatch                 uint32
	StrategyLinkId             uint32
	TotNumTradeReports         int32
	Filler4                    uint16
	MultiLegReportingType      uint8
	TradeReportType            uint8
	TrasnferReason             uint8
	PartyIdBeneficiery         [9]byte
	PartyIdTakeupTradingfirm   [5]byte
	PartyIdOrderOrignatingFirm [7]byte
	AccountType                uint8
	AggresorSide               uint8
	MatchType                  uint8
	MatchSubType               uint8
	Side                       uint8
	AggresorIndicator          uint8
	TradingCapacity            uint8
	Account                    [2]byte
	PositionEffect             [1]byte
	FreeText1                  [12]byte
	FreeText2                  [12]byte
	FreeText3                  [12]byte
	OrderCategory              [1]byte
	OrderType                  uint8
	RelatedproductComplex      uint8
	OrderSide                  uint8
	PartyClearingOrganisation  [4]byte
	PartyExecutingFirm         [5]byte
	PartyExecutingTrader       [6]byte
	PartyClearingFirm          [5]byte
	Filler5                    [7]byte
}
