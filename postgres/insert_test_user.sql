INSERT INTO cq_server_users
(username, password, created_at)
VALUES
('admin', 'admin', NOW());

-- ONLY RUN THIS IN DEV ENVIRONMENT
ALTER USER admin WITH SUPERUSER;
