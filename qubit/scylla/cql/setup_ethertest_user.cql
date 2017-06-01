// login in using root/arastra
// cqlsh -u root -p <...> -f setup2.cql <node>
// ethertest: This user can do pretty much anything except add/rm users,
//            and grant new perms to users.
//            To be used with ether intg tests which create own keyspaces.
DROP USER IF EXISTS ethertest;

CREATE USER IF NOT EXISTS ethertest WITH PASSWORD 'ethertest' NOSUPERUSER;
GRANT ALL PERMISSIONS ON ALL KEYSPACES TO ethertest;
REVOKE AUTHORIZE ON ALL KEYSPACES FROM ethertest;


