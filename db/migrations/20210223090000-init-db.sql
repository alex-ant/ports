-- +migrate Up

CREATE TABLE ports (
  id varchar(5) PRIMARY KEY,
  name varchar(32) NOT NULL,
  city varchar(32) NOT NULL,
  country varchar(32) NOT NULL,
  alias text[],
  regions text[],
  lat float NOT NULL,
  lng float NOT NULL,
  province varchar(32) NOT NULL,
  timezone varchar(32) NOT NULL,
  unlocs text[] NOT NULL,
  code varchar(5) NOT NULL
);

-- +migrate Down

DROP TABLE ports;
