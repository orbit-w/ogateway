package pb

// 所有协议ID常量
const (
	// Core 包协议ID
	PID_Core_Fail                   uint32 = 0xd03670ba // Core.Fail
	PID_Core_Notify_BeAttacked      uint32 = 0x8fee7235 // Core.Notify_BeAttacked
	PID_Core_OK                     uint32 = 0x0ece9291 // Core.OK
	PID_Core_Request_HeartBeat      uint32 = 0x95eee555 // Core.Request_HeartBeat
	PID_Core_Request_SearchBook     uint32 = 0xd3ecf693 // Core.Request_SearchBook
	PID_Core_Request_SearchBook_Rsp uint32 = 0xf1d19d0a // Core.Request_SearchBook_Rsp

	// Season 包协议ID
	PID_Season_Request_SeasonInfo uint32 = 0xd9714656 // Season.Request_SeasonInfo

)

// AllMessageNameToID 全局消息名称到ID的映射
var AllMessageNameToID = map[string]uint32{
	"Core-Fail":                   PID_Core_Fail,
	"Core-Notify_BeAttacked":      PID_Core_Notify_BeAttacked,
	"Core-OK":                     PID_Core_OK,
	"Core-Request_HeartBeat":      PID_Core_Request_HeartBeat,
	"Core-Request_SearchBook":     PID_Core_Request_SearchBook,
	"Core-Request_SearchBook_Rsp": PID_Core_Request_SearchBook_Rsp,
	"Season-Request_SeasonInfo":   PID_Season_Request_SeasonInfo,
}

// AllIDToMessageName 全局ID到消息名称的映射
var AllIDToMessageName = map[uint32]string{
	PID_Core_Fail:                   "Core-Fail",
	PID_Core_Notify_BeAttacked:      "Core-Notify_BeAttacked",
	PID_Core_OK:                     "Core-OK",
	PID_Core_Request_HeartBeat:      "Core-Request_HeartBeat",
	PID_Core_Request_SearchBook:     "Core-Request_SearchBook",
	PID_Core_Request_SearchBook_Rsp: "Core-Request_SearchBook_Rsp",
	PID_Season_Request_SeasonInfo:   "Season-Request_SeasonInfo",
}

// MessagePackageMap 消息名称到包名的映射
var MessagePackageMap = map[string]string{
	"Fail":                   "Core",
	"Notify_BeAttacked":      "Core",
	"OK":                     "Core",
	"Request_HeartBeat":      "Core",
	"Request_SearchBook":     "Core",
	"Request_SearchBook_Rsp": "Core",
	"Request_SeasonInfo":     "Season",
}

// GetProtocolID 获取指定消息名称的协议ID
func GetProtocolID(messageName string) (uint32, bool) {
	id, ok := AllMessageNameToID[messageName]
	return id, ok
}

// GetMessageName 获取指定协议ID的消息名称
func GetMessageName(pid uint32) (string, bool) {
	name, ok := AllIDToMessageName[pid]
	return name, ok
}
