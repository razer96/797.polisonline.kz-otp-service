package database

type OTP struct {
	ID          string `gorm:"type:varchar(100); unique_index; PRIMARY_KEY" json:"id"`
	Secret      string `gorm:"type:varchar(16); not null" json:"secret"`
	PhoneNumber string `gorm:"type:varchar(12); not null" json:"phone_number"`
	Attempts    int    `gorm:"type:int; default 0" json:"attempts"`
	Status      string `gorm:"type:boolean; varchar(1) not null default '1'" json:"status"`
	SendAt      int64  `gorm:"not null" json:"send_at"`
	ConfirmedAt int64  `gorm:"default null" json:"confirmed_at"`
}

var OTPStructSchema = "CREATE TABLE IF NOT EXISTS otps ( id varchar(100), secret varchar(16) not null, phone_number varchar(12) not null, attempts int default 0, status varchar(1) not null default '1', send_at integer not null, confirmed_at integer not null, PRIMARY KEY(id));"
