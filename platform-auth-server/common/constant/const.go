package constant

const (
	SCOPE_GLOBAL = "GLOBAL"

	AuthoritySubsystem = "SUB_SYSTEM"

	SubSystemNameSysPlatform = "SYS_PLATFORM"

	AuthSourceUm    = "UM"
	AuthSourceLocal = "LOCAL"

	UrlPrefix             = "/auth"
	AuthorizationHeader   = "Authorization"
	RefreshTokenHeader    = "Authorization-Info"
	UriLogin              = "/v1/api/login"
	UriUsersRegister      = "/v1/users/register"
	UriListApplyByApplier = "/v1/roles/apply/byapplier"
	UriRoles              = "/v1/roles"
	UriHealthCheck        = "/v1/health-check"

	Operator = "operator"

	ClientTypeUser      = "USER"
	ClientTypeSubSystem = "SUB_SYSTEM"

	TypeAccessToken  = "accessToken"
	TypeRefreshToken = "refreshToken"

	BearerTokenPrefix    = "Bearer "
	DefaultJwtSigningKey = "Platform+Auth+Server+Secret"
)
