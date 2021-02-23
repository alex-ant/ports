-- +migrate Up

CREATE TABLE ports (
  id varchar(5) PRIMARY KEY,
  name varchar(48) NOT NULL,
  city varchar(48) NOT NULL,
  country varchar(48) NOT NULL,
  alias text[],
  regions text[],
  lat float,
  lng float,
  province varchar(48) NOT NULL,
  timezone varchar(32) NOT NULL,
  unlocs text[] NOT NULL,
  code varchar(5) NOT NULL
);

-- +migrate Down

DROP TABLE ports;
