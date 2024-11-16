DO
$body$
BEGIN
   IF NOT EXISTS (
      SELECT *
      FROM   pg_catalog.pg_user
      WHERE  usename = 'device_management') THEN

      CREATE USER "device_management" WITH PASSWORD 'device-management-pw';
   END IF;
END
$body$;

GRANT "devices-management-manager" TO "device_management";
