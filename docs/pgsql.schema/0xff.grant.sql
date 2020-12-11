--
--
-- GRANT Permissions
GRANT   SELECT
   ON   ALL TABLES IN SCHEMA view, public
   TO   tomsvc;
GRANT   SELECT,
        INSERT,
        UPDATE,
        DELETE
   ON   ALL TABLES IN SCHEMA asset, bulk, filter, ix, meta, yp
   TO   tomsvc;
GRANT   USAGE,
        SELECT
   ON   ALL SEQUENCES IN SCHEMA asset, bulk, filter, ix, meta, yp
   TO   tomsvc;
