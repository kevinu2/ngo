package Jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Blacklist struct {
	ID        uint       `json:"id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Jwt       string     `json:"jwt,omitempty"`
}

type Claims struct {
	ID                 int    `json:"id,omitempty"`
	OrgId              int    `json:"org_id,omitempty"`
	OrgName            string `json:"org_name,omitempty"`
	ClientGid          string `json:"client_gid,omitempty"`
	UserName           string `json:"user_name,omitempty"`
	AppRights          string `json:"app_rights,omitempty"`
	AuthorityId        string `json:"authority_id,omitempty"`
	jwt.StandardClaims `json:"jwt_._standard_claims"`
	SessionId          string `json:"session_id,omitempty"`
}
