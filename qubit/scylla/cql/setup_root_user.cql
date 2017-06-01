// login in using cassandra and create root/arastra
// cqlsh -u cassdandra -p cassandra -f setup1.cql <node>
CREATE USER IF NOT EXISTS root      WITH PASSWORD 'arastra'   SUPERUSER ;

