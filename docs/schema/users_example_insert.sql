-- users 表联调初始化数据
-- 执行前请确认已创建 users 表（见 docs/schema/users.sql）

TRUNCATE TABLE `users`;

-- 默认明文密码: password
INSERT INTO `users` (`name`, `email`, `password`, `created_at`, `updated_at`) VALUES
  ('Alice', 'alice@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', NOW(), NOW()),
  ('Bob', 'bob@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', NOW(), NOW()),
  ('Carol', 'carol@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', NOW(), NOW()),
  ('David', 'david@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', NOW(), NOW());
