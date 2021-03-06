CREATE TABLE price_settings
(
  key        INTEGER                                         NOT NULL
  CONSTRAINT settings_pkey
  PRIMARY KEY,
  value      NUMERIC(10, 2)                                  NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()          NOT NULL
);

INSERT INTO price_settings (key, value) VALUES(1 /* BTC_USD_RATE */, 0) ON CONFLICT DO NOTHING; 
INSERT INTO price_settings (key, value) VALUES(2 /* BUY_CSGOKEY_PRICE */, 0) ON CONFLICT DO NOTHING; 
INSERT INTO price_settings (key, value) VALUES(3 /* SELL_CSGOKEY_PRICE */, 0) ON CONFLICT DO NOTHING; 