-- TODO: Consider adding is_headquarters to table
-- TODO: Consider CHECKs

CREATE TABLE IF NOT EXISTS bank (
	swift_code VARCHAR(11) PRIMARY KEY,
	hq_swift_code VARCHAR(11) REFERENCES bank(swift_code) ON DELETE SET NULL,
	is_headquarter BOOL NOT NULL,
	bank_name TEXT NOT NULL,
	address TEXT NOT NULL,
	country_iso2_code VARCHAR(2) NOT NULL,
	country_name TEXT NOT NULL
);

CREATE INDEX idx_bank_hq_swift_code ON bank(hq_swift_code);
CREATE INDEX idx_bank_country_iso2_code ON bank(country_iso2_code);
