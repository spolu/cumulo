package schemas

import "github.com/spolu/cumulo/lib/db"

const (
	devicesSQL = `
CREATE TABLE IF NOT EXISTS devices(
  token VARCHAR(256) NOT NULL,
  created TIMESTAMP NOT NULL,

  user VARCHAR(32) NOT NULL,      -- user token
  status VARCHAR(32) NOT NULL,    -- status (active, canceled)

  agent VARCHAR(512) NOT NULL,
  ip_address VARCHAR(32) NOT NULL,

  secret VARCHAR(256) NOT NULL,   -- secret used to connect

  PRIMARY KEY(token),
  FOREIGN KEY(user) REFERENCES users(token)
);
`
)

func init() {
	db.RegisterSchema(
		"api",
		"devices",
		devicesSQL,
	)
}
