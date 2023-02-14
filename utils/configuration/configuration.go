package configuration

const (
	PNAME_SC_HOME_DIR_NAME                = "sip.communicator.SC_HOME_DIR_NAME"
	PNAME_SC_HOME_DIR_LOCATION            = "sip.communicator.SC_HOME_DIR_LOCATION"
	PNAME_SC_CACHE_DIR_LOCATION           = "sip.communicator.SC_CACHE_DIR_LOCATION"
	PNAME_SC_LOG_DIR_LOCATION             = "sip.communicator.SC_LOG_DIR_LOCATION"
	PNAME_CONFIGURATION_FILE_IS_READ_ONLY = "sip.communicator.CONFIGURATION_FILE_IS_READ_ONLY"
	PNAME_CONFIGURATION_FILE_NAME         = "sip.communicator.CONFIGURATION_FILE_NAME"
)

type Configuration interface {
	setProperty(propertyName string, property interface{})
	setSystemProperty(propertyName string, property interface{})
	setProperties(properties map[string]interface{})
	getProperty(propertyName string) interface{}
	removePropertyByPrefix(prefix string)
	getAllPropertyNames() []string
	getPropertyNamesByPrefix(prefix string) []string
	getPropertyNamesByPrefixExactMatch(prefix string) []string
	getPropertyNamesBySuffix(suffix string) []string
	getPropertyNamesBySuffixExactMatch(suffix string) []string
	getString(propertyName string, defaultValue string) string
	getBool(propertyName string, defaultVal bool) bool
	getInt32(propertyName string, defaultVal int32) int32
	getInt64(propertyName string, defaultVal int64) int64
	getFloat64(propertyName string, defaultVal float64) float64
}
