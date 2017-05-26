DROP KEYSPACE IF EXISTS sanity_test;

CREATE KEYSPACE sanity_test WITH REPLICATION =
{'class': 'NetworkTopologyStrategy', 'DC1': 3};

USE sanity_test;

CREATE TABLE sanity_test (
  id text,
  test_value text,
  PRIMARY KEY (id)
);

INSERT INTO sanity_test (id, test_value) VALUES ('1', 'one');
INSERT INTO sanity_test (id, test_value) VALUES ('2', 'two');
INSERT INTO sanity_test (id, test_value) VALUES ('3', 'three');
