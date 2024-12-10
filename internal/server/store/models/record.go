package models

import "time"

type RecordByUUID struct {
	Uuid    string
	Version uint32
}

type AttributeDB struct {
	Uuid     string `db:"uuid"`
	ItemUuid string `db:"item_uuid"`
	Name     string `db:"name"`
	Value    string `db:"value"`
}

type FileDB struct {
	Uuid      string    `db:"uuid"`
	ItemUuid  string    `db:"item_uuid"`
	Path      *string   `db:"path"`
	Hash      string    `db:"hash"`
	Size      int       `db:"size"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Meta      *string   `db:"meta"`
}

type RecordDB struct {
	Uuid        string     `db:"uuid"`
	UserUuid    string     `db:"user_uuid"`
	Name        string     `db:"name"`
	Username    *string    `db:"username"`
	Url         *string    `db:"url"`
	Password    *string    `db:"password"`
	ExpiredAt   *time.Time `db:"expired_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	Cardnum     *string    `db:"cardnum"`
	Description *string    `db:"description"`
	Version     uint32     `db:"version"`
}
