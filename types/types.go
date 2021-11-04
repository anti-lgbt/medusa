package types

const (
	AuthzInvalidSession = "authz.invalid_session"

	JWTDecodeAndVerify     = "jwt.decode_and_verify"
	ServerInvalidBody      = "server.invalid_body"
	ServerInvalidQuery     = "server.invalid_query"
	RecordNotFound         = "record.not_found"
	AuthZInvalidPermission = "authz.invalid_permission"
)

type UserRole string

const (
	UserRoleMember   UserRole = "member"
	UserRoleMusician UserRole = "musician"
	UserRoleSinger   UserRole = "singer"
	UserRoleAdmin    UserRole = "admin"
)

type UserState string

const (
	UserStatePending UserState = "pending"
	UserStateActive  UserState = "active"
	UserStateDeleted UserState = "deleted"
	UserStateBanned  UserState = "banned"
)

type GroupRole string

const (
	GroupRoleMember GroupRole = "member"
	GroupRoleAdmin  GroupRole = "admin"
)

type GroupState string

const (
	GroupStateActive  GroupState = "active"
	GroupStateDeleted GroupState = "deleted"
)

type LabelState string

const (
	LabelStatePending  LabelState = "pending"
	LabelStateVerified LabelState = "verified"
	LabelStateDeleted  LabelState = "deleted"
)

type MusicState string

const (
	MusicStatePending  MusicState = "pending"
	MusicStateActive   MusicState = "active"
	MusicStateRejected MusicState = "rejected"
	MusicStateDeleted  MusicState = "deleted"
)

type Error struct {
	Error string `json:"error"`
}

type FileType string

const (
	FileTypeImage FileType = "image"
	FileTypeAudio FileType = "audio"
)

type EngineMailerPayload struct {
	Key      string                 `json:"key"`
	Language string                 `json:"language"`
	To       string                 `json:"to"`
	Record   map[string]interface{} `json:"record"`
}

type MailerConfig struct {
	Events []MailerConfigEvent `json:"events"`
}

type MailerConfigEvent struct {
	Name      string                                `yaml:"name"`
	Key       string                                `yaml:"key"`
	Templates map[string]MailerConfigEventTemplates `yaml:"templates"`
}

type MailerConfigEventTemplates struct {
	Subject      string `yaml:"subject"`
	TemplatePath string `yaml:"template_path"`
}

type Ordering string

var (
	OrderingAsc  Ordering = "asc"
	OrderingDesc Ordering = "desc"
)
