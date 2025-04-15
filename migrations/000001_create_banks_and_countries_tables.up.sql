-- TODO: Add indexes

CREATE TABLE IF NOT EXISTS country (
	iso2_code varchar(2) PRIMARY KEY,
	country_name text NOT NULL,
	time_zone text NOT NULL
);

CREATE TABLE IF NOT EXISTS bank (
	swift_code varchar(11) PRIMARY KEY,
	hq_swift_code varchar(11) REFERENCES bank (swift_code),
	country_iso2_code varchar(2) REFERENCES country (iso2_code) NOT NULL,
	bank_name text NOT NULL,
	address text NOT NULL
);
