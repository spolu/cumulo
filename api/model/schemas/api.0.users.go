package schemas

import "github.com/spolu/cumulo/lib/db"

const (
	usersSQL = `
CREATE TABLE IF NOT EXISTS users(
  token VARCHAR(256) NOT NULL,
  created TIMESTAMP NOT NULL,

  phone VARCHAR(32) NOT NULL,     -- phone number
  pubkey VARCHAR(128),            -- lightning network pubkey

  PRIMARY KEY(token),
  CONSTRAINT users_phone_u UNIQUE (phone)
);
`
)

func init() {
	db.RegisterSchema(
		"api",
		"users",
		usersSQL,
	)
}
