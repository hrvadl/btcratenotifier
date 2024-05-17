CREATE TABLE subscribers (
  id int PRIMARY KEY,
  email varchar(50) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE mailing_list (
  id int PRIMARY KEY,
  last_sent timestamp NOT NULL
);

ALTER TABLE subscribers
ADD CONSTRAINT UC_subscriber_email UNIQUE (email);

ALTER TABLE mailing_list
ADD CONSTRAINT UC_mailing_list_last_sent UNIQUE (last_sent);
