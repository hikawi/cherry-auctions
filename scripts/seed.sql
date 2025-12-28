TRUNCATE TABLE roles RESTART IDENTITY CASCADE;

INSERT INTO roles (id, description, created_at, updated_at) VALUES
  ('user', 'Default role for all users', now(), now()),
  ('moderator', 'Moderation role', now(),)
  ('admin', 'Administrative role', now(), now());

TRUNCATE TABLE categories RESTART IDENTITY CASCADE;

INSERT INTO categories (name, parent_id, created_at, updated_at) VALUES
  ('Electronics', NULL, now(), now()),
  ('Cards', NULL, now(), now()),
  ('Books', NULL, now(), now()),
  ('Cardcaptor Sakura', 2, now(), now()),
  ('Laptops', 1, now(), now()),
  ('Phones', 1, now(), now()),
  ('Fiction', 3, now(), now()),
  ('Non-Fiction', 3, now(), now()),
  ('Horror', 3, now(), now()),
  ('Fantasy', 3, now(), now()),
  ('Clow Cards', 4, now(), now()),
  ('Sakura Cards', 4, now(), now()),
  ('Clear Cards', 4, now(), now());

