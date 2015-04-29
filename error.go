package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "fmt"

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

// errors which does not have error codes.
var (
	ColumnCreateError = Error{Message: "column create error"}
	CtxOpenError      = Error{Message: "ctx open error"}
	DBCreateError     = Error{Message: "db create error"}
	TableCreateError  = Error{Message: "table create error"}
	ObjColumnError    = Error{Message: "column not found error"}
)

// errors which have error codes.
var (
	EndOfDataError                       = Error{Code: END_OF_DATA, Message: "end of data"}
	UnknownError                         = Error{Code: UNKNOWN_ERROR, Message: "unknown error"}
	OperationNotPermittedError           = Error{Code: OPERATION_NOT_PERMITTED, Message: "operation not permitted"}
	NoSuchFileOrDirectoryError           = Error{Code: NO_SUCH_FILE_OR_DIRECTORY, Message: "no such file or directory"}
	NoSuchProcessError                   = Error{Code: NO_SUCH_PROCESS, Message: "no such process"}
	InterruptedFunctionCallError         = Error{Code: INTERRUPTED_FUNCTION_CALL, Message: "interrupted function call"}
	InputOutputError                     = Error{Code: INPUT_OUTPUT_ERROR, Message: "input output error"}
	NoSuchDeviceOrAddressError           = Error{Code: NO_SUCH_DEVICE_OR_ADDRESS, Message: "no such device or address"}
	ArgListTooLongError                  = Error{Code: ARG_LIST_TOO_LONG, Message: "arg list too long"}
	ExecFormatError                      = Error{Code: EXEC_FORMAT_ERROR, Message: "exec format error"}
	BadFileDescriptorError               = Error{Code: BAD_FILE_DESCRIPTOR, Message: "bad file descriptor"}
	NoChildProcessesError                = Error{Code: NO_CHILD_PROCESSES, Message: "no child processes"}
	ResourceTemporarilyUnavailableError  = Error{Code: RESOURCE_TEMPORARILY_UNAVAILABLE, Message: "resource temporarily unavailable"}
	NotEnoughSpaceError                  = Error{Code: NOT_ENOUGH_SPACE, Message: "not enough space"}
	PermissionDeniedError                = Error{Code: PERMISSION_DENIED, Message: "permission denied"}
	BadAddressError                      = Error{Code: BAD_ADDRESS, Message: "bad address"}
	ResourceBusyError                    = Error{Code: RESOURCE_BUSY, Message: "resource busy"}
	FileExistsError                      = Error{Code: FILE_EXISTS, Message: "file exists"}
	ImproperLinkError                    = Error{Code: IMPROPER_LINK, Message: "improper link"}
	NoSuchDeviceError                    = Error{Code: NO_SUCH_DEVICE, Message: "no such device"}
	NotADirectoryError                   = Error{Code: NOT_A_DIRECTORY, Message: "not a directory"}
	IsADirectoryError                    = Error{Code: IS_A_DIRECTORY, Message: "is a directory"}
	InvalidArgumentError                 = Error{Code: INVALID_ARGUMENT, Message: "invalid argument"}
	TooManyOpenFilesInSystemError        = Error{Code: TOO_MANY_OPEN_FILES_IN_SYSTEM, Message: "too many open files in system"}
	TooManyOpenFilesError                = Error{Code: TOO_MANY_OPEN_FILES, Message: "too many open files"}
	InappropriateIOControlOperationError = Error{Code: INAPPROPRIATE_I_O_CONTROL_OPERATION, Message: "inappropriate i o control operation"}
	FileTooLargeError                    = Error{Code: FILE_TOO_LARGE, Message: "file too large"}
	NoSpaceLeftOnDeviceError             = Error{Code: NO_SPACE_LEFT_ON_DEVICE, Message: "no space left on device"}
	InvalidSeekError                     = Error{Code: INVALID_SEEK, Message: "invalid seek"}
	ReadOnlyFileSystemError              = Error{Code: READ_ONLY_FILE_SYSTEM, Message: "read only file system"}
	TooManyLinksError                    = Error{Code: TOO_MANY_LINKS, Message: "too many links"}
	BrokenPipeError                      = Error{Code: BROKEN_PIPE, Message: "broken pipe"}
	DomainError                          = Error{Code: DOMAIN_ERROR, Message: "domain error"}
	ResultTooLargeError                  = Error{Code: RESULT_TOO_LARGE, Message: "result too large"}
	ResourceDeadlockAvoidedError         = Error{Code: RESOURCE_DEADLOCK_AVOIDED, Message: "resource deadlock avoided"}
	NoMemoryAvailableError               = Error{Code: NO_MEMORY_AVAILABLE, Message: "no memory available"}
	FilenameTooLongError                 = Error{Code: FILENAME_TOO_LONG, Message: "filename too long"}
	NoLocksAvailableError                = Error{Code: NO_LOCKS_AVAILABLE, Message: "no locks available"}
	FunctionNotImplementedError          = Error{Code: FUNCTION_NOT_IMPLEMENTED, Message: "function not implemented"}
	DirectoryNotEmptyError               = Error{Code: DIRECTORY_NOT_EMPTY, Message: "directory not empty"}
	IllegalByteSequenceError             = Error{Code: ILLEGAL_BYTE_SEQUENCE, Message: "illegal byte sequence"}
	SocketNotInitializedError            = Error{Code: SOCKET_NOT_INITIALIZED, Message: "socket not initialized"}
	OperationWouldBlockError             = Error{Code: OPERATION_WOULD_BLOCK, Message: "operation would block"}
	AddressIsNotAvailableError           = Error{Code: ADDRESS_IS_NOT_AVAILABLE, Message: "address is not available"}
	NetworkIsDownError                   = Error{Code: NETWORK_IS_DOWN, Message: "network is down"}
	NoBufferError                        = Error{Code: NO_BUFFER, Message: "no buffer"}
	SocketIsAlreadyConnectedError        = Error{Code: SOCKET_IS_ALREADY_CONNECTED, Message: "socket is already connected"}
	SocketIsNotConnectedError            = Error{Code: SOCKET_IS_NOT_CONNECTED, Message: "socket is not connected"}
	SocketIsAlreadyShutdownedError       = Error{Code: SOCKET_IS_ALREADY_SHUTDOWNED, Message: "socket is already shutdowned"}
	OperationTimeoutError                = Error{Code: OPERATION_TIMEOUT, Message: "operation timeout"}
	ConnectionRefusedError               = Error{Code: CONNECTION_REFUSED, Message: "connection refused"}
	RangeError                           = Error{Code: RANGE_ERROR, Message: "range error"}
	TokenizerError                       = Error{Code: TOKENIZER_ERROR, Message: "tokenizer error"}
	FileCorruptError                     = Error{Code: FILE_CORRUPT, Message: "file corrupt"}
	InvalidFormatError                   = Error{Code: INVALID_FORMAT, Message: "invalid format"}
	ObjectCorruptError                   = Error{Code: OBJECT_CORRUPT, Message: "object corrupt"}
	TooManySymbolicLinksError            = Error{Code: TOO_MANY_SYMBOLIC_LINKS, Message: "too many symbolic links"}
	NotSocketError                       = Error{Code: NOT_SOCKET, Message: "not socket"}
	OperationNotSupportedError           = Error{Code: OPERATION_NOT_SUPPORTED, Message: "operation not supported"}
	AddressIsInUseError                  = Error{Code: ADDRESS_IS_IN_USE, Message: "address is in use"}
	ZlibError                            = Error{Code: ZLIB_ERROR, Message: "zlib error"}
	Lz4Error                             = Error{Code: LZ4_ERROR, Message: "lz4 error"}
	StackOverFlowError                   = Error{Code: STACK_OVER_FLOW, Message: "stack over flow"}
	SyntaxError                          = Error{Code: SYNTAX_ERROR, Message: "syntax error"}
	RetryMaxError                        = Error{Code: RETRY_MAX, Message: "retry max"}
	IncompatibleFileFormatError          = Error{Code: INCOMPATIBLE_FILE_FORMAT, Message: "incompatible file format"}
	UpdateNotAllowedError                = Error{Code: UPDATE_NOT_ALLOWED, Message: "update not allowed"}
	TooSmallOffsetError                  = Error{Code: TOO_SMALL_OFFSET, Message: "too small offset"}
	TooLargeOffsetError                  = Error{Code: TOO_LARGE_OFFSET, Message: "too large offset"}
	TooSmallLimitError                   = Error{Code: TOO_SMALL_LIMIT, Message: "too small limit"}
	CasError                             = Error{Code: CAS_ERROR, Message: "cas error"}
	UnsupportedCommandVersionError       = Error{Code: UNSUPPORTED_COMMAND_VERSION, Message: "unsupported command version"}
	NormalizerError                      = Error{Code: NORMALIZER_ERROR, Message: "normalizer error"}
	TokenFilterError                     = Error{Code: TOKEN_FILTER_ERROR, Message: "token filter error"}
	CommandError                         = Error{Code: COMMAND_ERROR, Message: "command error"}
	PluginError                          = Error{Code: PLUGIN_ERROR, Message: "plugin error"}
	ScorerError                          = Error{Code: SCORER_ERROR, Message: "scorer error"}
)

func errorFromRc(code C.grn_rc) error {
	switch code {
	case SUCCESS:
		panic("must not call errorFromRc for SUCCESS")
	case END_OF_DATA:
		return EndOfDataError
	case UNKNOWN_ERROR:
		return UnknownError
	case OPERATION_NOT_PERMITTED:
		return OperationNotPermittedError
	case NO_SUCH_FILE_OR_DIRECTORY:
		return NoSuchFileOrDirectoryError
	case NO_SUCH_PROCESS:
		return NoSuchProcessError
	case INTERRUPTED_FUNCTION_CALL:
		return InterruptedFunctionCallError
	case INPUT_OUTPUT_ERROR:
		return InputOutputError
	case NO_SUCH_DEVICE_OR_ADDRESS:
		return NoSuchDeviceOrAddressError
	case ARG_LIST_TOO_LONG:
		return ArgListTooLongError
	case EXEC_FORMAT_ERROR:
		return ExecFormatError
	case BAD_FILE_DESCRIPTOR:
		return BadFileDescriptorError
	case NO_CHILD_PROCESSES:
		return NoChildProcessesError
	case RESOURCE_TEMPORARILY_UNAVAILABLE:
		return ResourceTemporarilyUnavailableError
	case NOT_ENOUGH_SPACE:
		return NotEnoughSpaceError
	case PERMISSION_DENIED:
		return PermissionDeniedError
	case BAD_ADDRESS:
		return BadAddressError
	case RESOURCE_BUSY:
		return ResourceBusyError
	case FILE_EXISTS:
		return FileExistsError
	case IMPROPER_LINK:
		return ImproperLinkError
	case NO_SUCH_DEVICE:
		return NoSuchDeviceError
	case NOT_A_DIRECTORY:
		return NotADirectoryError
	case IS_A_DIRECTORY:
		return IsADirectoryError
	case INVALID_ARGUMENT:
		return InvalidArgumentError
	case TOO_MANY_OPEN_FILES_IN_SYSTEM:
		return TooManyOpenFilesInSystemError
	case TOO_MANY_OPEN_FILES:
		return TooManyOpenFilesError
	case INAPPROPRIATE_I_O_CONTROL_OPERATION:
		return InappropriateIOControlOperationError
	case FILE_TOO_LARGE:
		return FileTooLargeError
	case NO_SPACE_LEFT_ON_DEVICE:
		return NoSpaceLeftOnDeviceError
	case INVALID_SEEK:
		return InvalidSeekError
	case READ_ONLY_FILE_SYSTEM:
		return ReadOnlyFileSystemError
	case TOO_MANY_LINKS:
		return TooManyLinksError
	case BROKEN_PIPE:
		return BrokenPipeError
	case DOMAIN_ERROR:
		return DomainError
	case RESULT_TOO_LARGE:
		return ResultTooLargeError
	case RESOURCE_DEADLOCK_AVOIDED:
		return ResourceDeadlockAvoidedError
	case NO_MEMORY_AVAILABLE:
		return NoMemoryAvailableError
	case FILENAME_TOO_LONG:
		return FilenameTooLongError
	case NO_LOCKS_AVAILABLE:
		return NoLocksAvailableError
	case FUNCTION_NOT_IMPLEMENTED:
		return FunctionNotImplementedError
	case DIRECTORY_NOT_EMPTY:
		return DirectoryNotEmptyError
	case ILLEGAL_BYTE_SEQUENCE:
		return IllegalByteSequenceError
	case SOCKET_NOT_INITIALIZED:
		return SocketNotInitializedError
	case OPERATION_WOULD_BLOCK:
		return OperationWouldBlockError
	case ADDRESS_IS_NOT_AVAILABLE:
		return AddressIsNotAvailableError
	case NETWORK_IS_DOWN:
		return NetworkIsDownError
	case NO_BUFFER:
		return NoBufferError
	case SOCKET_IS_ALREADY_CONNECTED:
		return SocketIsAlreadyConnectedError
	case SOCKET_IS_NOT_CONNECTED:
		return SocketIsNotConnectedError
	case SOCKET_IS_ALREADY_SHUTDOWNED:
		return SocketIsAlreadyShutdownedError
	case OPERATION_TIMEOUT:
		return OperationTimeoutError
	case CONNECTION_REFUSED:
		return ConnectionRefusedError
	case RANGE_ERROR:
		return RangeError
	case TOKENIZER_ERROR:
		return TokenizerError
	case FILE_CORRUPT:
		return FileCorruptError
	case INVALID_FORMAT:
		return InvalidFormatError
	case OBJECT_CORRUPT:
		return ObjectCorruptError
	case TOO_MANY_SYMBOLIC_LINKS:
		return TooManySymbolicLinksError
	case NOT_SOCKET:
		return NotSocketError
	case OPERATION_NOT_SUPPORTED:
		return OperationNotSupportedError
	case ADDRESS_IS_IN_USE:
		return AddressIsInUseError
	case ZLIB_ERROR:
		return ZlibError
	case LZ4_ERROR:
		return Lz4Error
	case STACK_OVER_FLOW:
		return StackOverFlowError
	case SYNTAX_ERROR:
		return SyntaxError
	case RETRY_MAX:
		return RetryMaxError
	case INCOMPATIBLE_FILE_FORMAT:
		return IncompatibleFileFormatError
	case UPDATE_NOT_ALLOWED:
		return UpdateNotAllowedError
	case TOO_SMALL_OFFSET:
		return TooSmallOffsetError
	case TOO_LARGE_OFFSET:
		return TooLargeOffsetError
	case TOO_SMALL_LIMIT:
		return TooSmallLimitError
	case CAS_ERROR:
		return CasError
	case UNSUPPORTED_COMMAND_VERSION:
		return UnsupportedCommandVersionError
	case NORMALIZER_ERROR:
		return NormalizerError
	case TOKEN_FILTER_ERROR:
		return TokenFilterError
	case COMMAND_ERROR:
		return CommandError
	case PLUGIN_ERROR:
		return PluginError
	case SCORER_ERROR:
		return ScorerError
	default:
		panic(fmt.Sprintf("unknown grn_rc: %d", code))
	}
}
