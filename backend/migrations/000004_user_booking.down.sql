-- 000004_user_booking.down.sql
-- 回滚 Sprint 3 新增表：按依赖顺序删除用户、预订相关表。

DROP TABLE IF EXISTS cabin_holds;
DROP TABLE IF EXISTS booking_passengers;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS passengers;
DROP TABLE IF EXISTS users;
