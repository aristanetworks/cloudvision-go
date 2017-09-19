CONSISTENCY ALL;
// login in using root/arastra
// cqlsh -u root -p <...> -f setup2.cql <node>
// DROP USER IF EXISTS cassandra;
// DROP USER IF EXISTS ether;
// DROP USER IF EXISTS arastra;

// ether: This user can modify data keyspace "ether" but not change schema.
//        To be used with Quantumfs and other qubit tools.
CREATE USER IF NOT EXISTS ether     WITH PASSWORD 'ether'     NOSUPERUSER;
REVOKE ALL PERMISSIONS ON ALL KEYSPACES FROM ether;
GRANT MODIFY ON KEYSPACE qubit0 to ether;
GRANT SELECT ON KEYSPACE qubit0 to ether;

// arastra: This user can only read data.
CREATE USER IF NOT EXISTS arastra   WITH PASSWORD 'arastra'   NOSUPERUSER;
REVOKE ALL PERMISSIONS ON ALL KEYSPACES FROM arastra;
GRANT SELECT ON ALL KEYSPACES TO arastra;


