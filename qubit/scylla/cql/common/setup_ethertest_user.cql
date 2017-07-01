// login in using root/arastra
// cqlsh -u root -p <...> -f setup2.cql <node>
// ethertest: This user can do pretty much anything to "ethertest" keyspace,
//            except add/rm users and grant new perms to users.
//            To be used with ether intg tests which create own keyspaces.
CONSISTENCY ALL;
DROP USER IF EXISTS ethertest;

CREATE USER IF NOT EXISTS ethertest WITH PASSWORD 'ethertest' NOSUPERUSER;
GRANT ALL PERMISSIONS ON keyspace ethertest TO ethertest;
REVOKE AUTHORIZE ON keyspace ethertest FROM ethertest;


