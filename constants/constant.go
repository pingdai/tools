package constants

// ENV
const (
	// test
	ENV_TEST = "test"
	// pre
	ENV_PRE = "pre"

	// 单独控制向外输出config信息
	ENV_PRINT_CONFIG = "print_config"

	// 服务名
	EnvVarKeyProjectName = "PROJECT_NAME"
	// 服务环境变量
	EnvVarKeyEnvFlag = "ENV_FLAG"
	// gin框架环境，线上版本定义为 release
	I_EnvVarKeyGinMode = "GIN_MODE"
)

// LOG
const (
	// text
	LOG_FORMAT_TEXT = "text"
	// json
	LOG_FORMAT_JSON = "json"

	LOG_LEVEL_DEBUG = "debug"
	LOG_LEVEL_INFO  = "info"
	LOG_LEVEL_WARN  = "warn"
	LOG_LEVEL_ERROR = "error"
)

// GIN
const (
	// log_id
	HEADER_LOD_ID = "Log-ID"
	// remote_service
	HEADER_REMOTE_SERVICE = "Remote-Service"
)

// MQ类型
type MQType int

const (
	// 生产者
	MQ_TYPE_PRODUCER MQType = 1
	// 消费者
	MQ_TYPE_CONSUMER MQType = 2
)
