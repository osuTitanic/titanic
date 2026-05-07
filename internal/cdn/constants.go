package cdn

const (
	DefaultStreamingBaseUrl = "https://cdn.titanic.sh"
	DefaultRedirectBaseUrl  = "https://s3.titanic.sh"
	DefaultUserAgent        = "osuTitanic/titanic"
)

type Permission string

const (
	PermissionList   Permission = "list"
	PermissionUpload Permission = "upload"
	PermissionDelete Permission = "delete"
)
