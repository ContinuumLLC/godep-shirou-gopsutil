package model

const (

	//Generic error codes to wrap errors generated at downstream layers.

	//ErrAssetService to handle error at service layer
	ErrAssetService = "ErrAssetService"
	//ErrAssetDal to handle error at dal layer
	ErrAssetDal = "ErrAssetDal"
	//ErrAssetMsgListener to handle error at message listener layer
	ErrAssetMsgListener = "ErrAssetMsgListener"
	//ErrAssetInstallDate to handle invalid or blank date format
	ErrAssetInstallDate = "ErrAssetInstallDate"
	//ErrNotImplemented error for missing implementation
	ErrNotImplemented = "ErrNotImplemented"
)
