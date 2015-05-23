package goroonga

// builtin types
const (
	DB_VOID = iota
	DB_DB
	DB_OBJECT
	DB_BOOL
	DB_INT8
	DB_UINT8
	DB_INT16
	DB_UINT16
	DB_INT32
	DB_UINT32
	DB_INT64
	DB_UINT64
	DB_FLOAT
	DB_TIME
	DB_SHORT_TEXT
	DB_TEXT
	DB_LONG_TEXT
	DB_TOKYO_GEO_POINT
	DB_WGS84_GEO_POINT
)

// obj flags
const (
	OBJ_TABLE_TYPE_MASK = 0x07
	OBJ_TABLE_HASH_KEY  = 0x00
	OBJ_TABLE_PAT_KEY   = 0x01
	OBJ_TABLE_DAT_KEY   = 0x02
	OBJ_TABLE_NO_KEY    = 0x03

	OBJ_KEY_MASK      = 0x07 << 3
	OBJ_KEY_UINT      = 0x00 << 3
	OBJ_KEY_INT       = 0x01 << 3
	OBJ_KEY_FLOAT     = 0x02 << 3
	OBJ_KEY_GEO_POINT = 0x03 << 3

	OBJ_KEY_WITH_SIS  = 0x01 << 6
	OBJ_KEY_NORMALIZE = 0x01 << 7

	OBJ_COLUMN_TYPE_MASK = 0x07
	OBJ_COLUMN_SCALAR    = 0x00
	OBJ_COLUMN_VECTOR    = 0x01
	OBJ_COLUMN_INDEX     = 0x02

	OBJ_COMPRESS_MASK = 0x07 << 4
	OBJ_COMPRESS_NONE = 0x00 << 4
	OBJ_COMPRESS_ZLIB = 0x01 << 4
	OBJ_COMPRESS_LZ4  = 0x02 << 4

	OBJ_WITH_SECTION  = 0x01 << 7
	OBJ_WITH_WEIGHT   = 0x01 << 8
	OBJ_WITH_POSITION = 0x01 << 9
	OBJ_RING_BUFFER   = 0x01 << 10

	OBJ_UNIT_MASK              = 0x0f << 8
	OBJ_UNIT_DOCUMENT_NONE     = 0x00 << 8
	OBJ_UNIT_DOCUMENT_SECTION  = 0x01 << 8
	OBJ_UNIT_DOCUMENT_POSITION = 0x02 << 8
	OBJ_UNIT_SECTION_NONE      = 0x03 << 8
	OBJ_UNIT_SECTION_POSITION  = 0x04 << 8
	OBJ_UNIT_POSITION_NONE     = 0x05 << 8
	OBJ_UNIT_USERDEF_DOCUMENT  = 0x06 << 8
	OBJ_UNIT_USERDEF_SECTION   = 0x07 << 8
	OBJ_UNIT_USERDEF_POSITION  = 0x08 << 8

	OBJ_NO_SUBREC   = 0x00 << 13
	OBJ_WITH_SUBREC = 0x01 << 13

	OBJ_KEY_VAR_SIZE = 0x01 << 14

	OBJ_TEMPORARY  = 0x00 << 15
	OBJ_PERSISTENT = 0x01 << 15
)

// error codes
const (
	SUCCESS                             = 0
	END_OF_DATA                         = 1
	UNKNOWN_ERROR                       = -1
	OPERATION_NOT_PERMITTED             = -2
	NO_SUCH_FILE_OR_DIRECTORY           = -3
	NO_SUCH_PROCESS                     = -4
	INTERRUPTED_FUNCTION_CALL           = -5
	INPUT_OUTPUT_ERROR                  = -6
	NO_SUCH_DEVICE_OR_ADDRESS           = -7
	ARG_LIST_TOO_LONG                   = -8
	EXEC_FORMAT_ERROR                   = -9
	BAD_FILE_DESCRIPTOR                 = -10
	NO_CHILD_PROCESSES                  = -11
	RESOURCE_TEMPORARILY_UNAVAILABLE    = -12
	NOT_ENOUGH_SPACE                    = -13
	PERMISSION_DENIED                   = -14
	BAD_ADDRESS                         = -15
	RESOURCE_BUSY                       = -16
	FILE_EXISTS                         = -17
	IMPROPER_LINK                       = -18
	NO_SUCH_DEVICE                      = -19
	NOT_A_DIRECTORY                     = -20
	IS_A_DIRECTORY                      = -21
	INVALID_ARGUMENT                    = -22
	TOO_MANY_OPEN_FILES_IN_SYSTEM       = -23
	TOO_MANY_OPEN_FILES                 = -24
	INAPPROPRIATE_I_O_CONTROL_OPERATION = -25
	FILE_TOO_LARGE                      = -26
	NO_SPACE_LEFT_ON_DEVICE             = -27
	INVALID_SEEK                        = -28
	READ_ONLY_FILE_SYSTEM               = -29
	TOO_MANY_LINKS                      = -30
	BROKEN_PIPE                         = -31
	DOMAIN_ERROR                        = -32
	RESULT_TOO_LARGE                    = -33
	RESOURCE_DEADLOCK_AVOIDED           = -34
	NO_MEMORY_AVAILABLE                 = -35
	FILENAME_TOO_LONG                   = -36
	NO_LOCKS_AVAILABLE                  = -37
	FUNCTION_NOT_IMPLEMENTED            = -38
	DIRECTORY_NOT_EMPTY                 = -39
	ILLEGAL_BYTE_SEQUENCE               = -40
	SOCKET_NOT_INITIALIZED              = -41
	OPERATION_WOULD_BLOCK               = -42
	ADDRESS_IS_NOT_AVAILABLE            = -43
	NETWORK_IS_DOWN                     = -44
	NO_BUFFER                           = -45
	SOCKET_IS_ALREADY_CONNECTED         = -46
	SOCKET_IS_NOT_CONNECTED             = -47
	SOCKET_IS_ALREADY_SHUTDOWNED        = -48
	OPERATION_TIMEOUT                   = -49
	CONNECTION_REFUSED                  = -50
	RANGE_ERROR                         = -51
	TOKENIZER_ERROR                     = -52
	FILE_CORRUPT                        = -53
	INVALID_FORMAT                      = -54
	OBJECT_CORRUPT                      = -55
	TOO_MANY_SYMBOLIC_LINKS             = -56
	NOT_SOCKET                          = -57
	OPERATION_NOT_SUPPORTED             = -58
	ADDRESS_IS_IN_USE                   = -59
	ZLIB_ERROR                          = -60
	LZ4_ERROR                           = -61
	STACK_OVER_FLOW                     = -62
	SYNTAX_ERROR                        = -63
	RETRY_MAX                           = -64
	INCOMPATIBLE_FILE_FORMAT            = -65
	UPDATE_NOT_ALLOWED                  = -66
	TOO_SMALL_OFFSET                    = -67
	TOO_LARGE_OFFSET                    = -68
	TOO_SMALL_LIMIT                     = -69
	CAS_ERROR                           = -70
	UNSUPPORTED_COMMAND_VERSION         = -71
	NORMALIZER_ERROR                    = -72
	TOKEN_FILTER_ERROR                  = -73
	COMMAND_ERROR                       = -74
	PLUGIN_ERROR                        = -75
	SCORER_ERROR                        = -76
)
