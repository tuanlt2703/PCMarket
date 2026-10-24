REVOKE ALL PRIVILEGES ON DATABASE :"db_name" FROM PUBLIC;

\if :has_migrator
SELECT format('CREATE ROLE %I LOGIN NOINHERIT PASSWORD %L', :'migrator_user', :'migrator_password')
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = :'migrator_user')
\gexec

ALTER ROLE :"migrator_user" LOGIN NOINHERIT PASSWORD :'migrator_password';
GRANT CONNECT ON DATABASE :"db_name" TO :"migrator_user";
GRANT :"db_user" TO :"migrator_user";
\endif