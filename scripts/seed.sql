-- HEY!
-- You're expected to already migrate models in GORM before running this script.

CREATE EXTENSION IF NOT EXISTS pg_trgm;

TRUNCATE TABLE roles RESTART IDENTITY CASCADE;
truncate table products restart identity cascade;
TRUNCATE TABLE categories RESTART IDENTITY CASCADE;
truncate table products_categories restart identity cascade;
truncate table users restart identity cascade;
truncate table bids restart identity cascade;
truncate table questions restart identity cascade;
truncate table refresh_tokens restart identity cascade;

CREATE INDEX IF NOT EXISTS idx_products_name_trgm ON products USING gin (name gin_trgm_ops);

--- Add a trigger to update search vector

CREATE OR REPLACE FUNCTION update_product_search_vector()
RETURNS trigger AS $$
BEGIN
  NEW.search_vector :=
    to_tsvector(
      'simple',
      coalesce(NEW.name, '') || ' ' || coalesce(NEW.description, '')
    );
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER product_search_vector_trigger
BEFORE INSERT OR UPDATE
ON products
FOR EACH ROW
EXECUTE FUNCTION update_product_search_vector();

-- Main User

INSERT INTO users (name, email, password, oauth_type, created_at, updated_at) VALUES ('Admin', 'admin@example.com', '$argon2id$v=19$m=65536,t=2,p=4$y7rT7dfnqb4NJOeqhVHEoA$UIaO/jNFwRH/Oz1qXcTOjzMpOSM+il1p835bbliZ6IM', 'none', now(), now());
INSERT INTO users (name, email, password, oauth_type, created_at, updated_at) VALUES ('Luna', 'luna@example.com', '$argon2id$v=19$m=65536,t=2,p=4$y7rT7dfnqb4NJOeqhVHEoA$UIaO/jNFwRH/Oz1qXcTOjzMpOSM+il1p835bbliZ6IM', 'none', now(), now());
INSERT INTO users (name, email, password, oauth_type, created_at, updated_at) VALUES ('MrGeraffe', 'mrgeraffe@example.com', '$argon2id$v=19$m=65536,t=2,p=4$y7rT7dfnqb4NJOeqhVHEoA$UIaO/jNFwRH/Oz1qXcTOjzMpOSM+il1p835bbliZ6IM', 'none', now(), now());
INSERT INTO users (name, email, password, oauth_type, created_at, updated_at) VALUES ('Omi', 'omi@example.com', '$argon2id$v=19$m=65536,t=2,p=4$y7rT7dfnqb4NJOeqhVHEoA$UIaO/jNFwRH/Oz1qXcTOjzMpOSM+il1p835bbliZ6IM', 'none', now(), now());

-- Roles.
-- Run this after you setup GORM to make sure you can register.

INSERT INTO roles (id, description, created_at, updated_at) VALUES
  ('user', 'Default role for all users', now(), now()),
  ('moderator', 'Moderation role', now(), now()),
  ('admin', 'Administrative role', now(), now());

INSERT INTO user_roles VALUES (1, 'user'), (1, 'moderator'), (1, 'admin');
INSERT INTO user_roles VALUES (2, 'user'), (2, 'moderator');
INSERT INTO user_roles VALUES (3, 'user');
INSERT INTO user_roles VALUES (4, 'user');

-- Make MrGeraffe a seller.

-- Seeding for categories.

-- 2. PARENT CATEGORIES
INSERT INTO categories (name, created_at, updated_at) VALUES 
('Electronics', now(), now()),
('Fashion', now(), now()),
('Home & Garden', now(), now()),
('Collectibles', now(), now()),
('Sports & Outdoors', now(), now());

-- 3. SUB-CATEGORIES (2-level hierarchy)
INSERT INTO categories (name, parent_id, created_at, updated_at) VALUES 
('Laptops', 1, now(), now()),
('Mobile Devices', 1, now(), now()),
('Men''s Wear', 2, now(), now()),
('Watches', 2, now(), now()),
('Small Appliances', 3, now(), now()),
('Kitchenware', 3, now(), now()),
('Trading Cards', 4, now(), now()),
('Vintage Toys', 4, now(), now()),
('Camping Gear', 5, now(), now());

-- 4. PRODUCTS (Meaningful descriptions and pricing)
-- Prices follow a logic: starting_bid is ~10-20% of bin_price.

-- Electronics: Laptops
INSERT INTO products (starting_bid, step_bid_type, step_bid_value, bin_price, allows_unrated_buyers, auto_extends_time, expired_at, seller_id, thumbnail_url, name, description, created_at, updated_at) VALUES 
(150.00, 'fixed', 10.00, 850.00, false, true, now() + interval '5 days', 2, 'https://placehold.co/400x400?text=Macbook', 'MacBook Air M1 - 8GB RAM, 256GB SSD, Space Gray', 'Excellent condition with original charger.', now(), now()),
(200.00, 'percentage', 5.0, 1200.00, false, true, now() + interval '3 days', 2, 'https://placehold.co/400x400?text=Dell+XPS', 'Dell XPS 13 9310 - Intel i7, 16GB RAM, 512GB SSD', 'Perfect for developers and professionals.', now(), now());

-- Electronics: Mobile
INSERT INTO products (starting_bid, step_bid_type, step_bid_value, bin_price, allows_unrated_buyers, auto_extends_time, expired_at, seller_id, thumbnail_url, name, description, created_at, updated_at) VALUES 
(100.00, 'fixed', 5.00, 600.00, true, true, now() + interval '2 days', 3, 'https://placehold.co/400x400?text=iPhone+12', 'iPhone 12 128GB - Blue', 'Unlocked. No scratches on screen, battery health at 88%.', now(), now()),
(80.00, 'fixed', 5.00, 450.00, true, false, now() + interval '4 days', 1, 'https://placehold.co/400x400?text=Pixel+6', 'Google Pixel 6 - Stormy Black', 'Amazing camera quality with Google Tensor chip.', now(), now());

-- Fashion: Men's Wear & Watches
INSERT INTO products (starting_bid, step_bid_type, step_bid_value, bin_price, allows_unrated_buyers, auto_extends_time, expired_at, seller_id, thumbnail_url, name, description, created_at, updated_at) VALUES 
(45.00, 'fixed', 2.00, 150.00, true, false, now() + interval '6 days', 4, 'https://placehold.co/400x400?text=Leather+Jacket', 'Genuine Lambskin Leather Jacket - Size M', 'Classic biker style, worn only once.', now(), now()),
(120.00, 'percentage', 2.0, 550.00, false, true, now() + interval '24 hours', 4, 'https://placehold.co/400x400?text=Seiko+Watch', 'Seiko Prospex "Turtle" Automatic Diver', 'Water resistant 200m. Original box included.', now(), now());

-- Home & Garden: Appliances
INSERT INTO products (starting_bid, step_bid_type, step_bid_value, bin_price, allows_unrated_buyers, auto_extends_time, expired_at, seller_id, thumbnail_url, name, description, created_at, updated_at) VALUES 
(30.00, 'fixed', 5.00, 180.00, true, true, now() + interval '5 days', 1, 'https://placehold.co/400x400?text=Air+Fryer', 'Ninja Foodi Air Fryer - 4-Quart capacity', 'Great for healthy cooking with little to no oil.', now(), now()),
(25.00, 'fixed', 2.00, 120.00, true, false, now() + interval '7 days', 3, 'https://placehold.co/400x400?text=Espresso+Maker', 'De''Longhi Espresso Machine', 'Manual milk frother included. Perfect for home baristas.', now(), now());

-- Collectibles: Trading Cards & Toys
INSERT INTO products (starting_bid, step_bid_type, step_bid_value, bin_price, allows_unrated_buyers, auto_extends_time, expired_at, seller_id, thumbnail_url, name, description, created_at, updated_at) VALUES 
(50.00, 'percentage', 10.0, 500.00, false, true, now() + interval '12 hours', 2, 'https://placehold.co/400x400?text=Charizard', 'Base Set Shadowless Charizard - PSA 6', 'A rare centerpiece for any Pokemon collection.', now(), now()),
(15.00, 'fixed', 1.00, 75.00, true, false, now() + interval '3 days', 4, 'https://placehold.co/400x400?text=Vintage+Car', '1960s Matchbox Diecast Car - Red Ferrari', 'Minor paint wear, but rolls perfectly.', now(), now());

-- Sports & Outdoors: Camping
INSERT INTO products (starting_bid, step_bid_type, step_bid_value, bin_price, allows_unrated_buyers, auto_extends_time, expired_at, seller_id, thumbnail_url, name, description, created_at, updated_at) VALUES 
(40.00, 'fixed', 5.00, 220.00, true, true, now() + interval '4 days', 1, 'https://placehold.co/400x400?text=Tent', 'Coleman 4-Person Instant Cabin Tent', 'Sets up in 60 seconds. Weatherproof technology.', now(), now()),
(20.00, 'fixed', 2.00, 95.00, true, false, now() + interval '5 days', 2, 'https://placehold.co/400x400?text=Sleeping+Bag', 'Mummy Style Sleeping Bag - Rated for 0°F', 'Lightweight and compact for hiking.', now(), now());

-- Additional Mixed Products to reach 20
INSERT INTO products (starting_bid, step_bid_type, step_bid_value, bin_price, allows_unrated_buyers, auto_extends_time, expired_at, seller_id, thumbnail_url, name, description, created_at, updated_at) VALUES 
(299.00, 'percentage', 5.0, 1500.00, false, true, now() + interval '2 days', 2, 'https://placehold.co/400x400?text=Camera', 'Sony A7 III Camera Body', 'Full-frame mirrorless camera. Low shutter count, pristine condition.', now(), now()),
(25.00, 'fixed', 5.00, 110.00, true, true, now() + interval '6 days', 3, 'https://placehold.co/400x400?text=Cast+Iron', 'Cast Iron Skillet Set', 'Pre-seasoned 3-piece set (8, 10, 12 inch). Durable and heat-retentive.', now(), now()),
(10.00, 'fixed', 1.00, 45.00, true, false, now() + interval '2 days', 4, 'https://placehold.co/400x400?text=Keyboard', 'Mechanical Gaming Keyboard', 'RGB Backlit with Blue Switches. Tactile and clicky for typing enthusiasts.', now(), now()),
(60.00, 'percentage', 2.0, 300.00, false, true, now() + interval '4 days', 1, 'https://placehold.co/400x400?text=Wayfarer', 'Ray-Ban Wayfarer Classic', 'Polarized lenses with black frame. Includes original leather case.', now(), now()),
(35.00, 'fixed', 5.00, 130.00, true, false, now() + interval '7 days', 3, 'https://placehold.co/400x400?text=Grill', 'Portable Propane Grill', 'Tabletop design for tailgating and camping. 12,000 BTU burner.', now(), now()),
(5.00, 'fixed', 0.50, 25.00, true, true, now() + interval '3 days', 4, 'https://placehold.co/400x400?text=Water', 'Stainless Steel Water Bottle', 'Double-wall vacuum insulated. Keeps drinks cold for 24 hours.', now(), now()),
(12.00, 'fixed', 1.00, 60.00, true, false, now() + interval '5 days', 1, 'https://placehold.co/400x400?text=Bulb', 'Smart LED Light Bulb Pack', '4-pack of RGB bulbs compatible with Alexa and Google Home.', now(), now()),
(85.00, 'percentage', 5.0, 400.00, false, true, now() + interval '10 days', 2, 'https://placehold.co/400x400?text=Desk', 'Electric Standing Desk', 'Height adjustable frame only. Heavy-duty dual motor system.', now(), now());

-- 5. MAPPING PRODUCTS TO CATEGORIES
INSERT INTO products_categories (product_id, category_id) VALUES 
(1, 6), (2, 6),   -- Laptops
(3, 7), (4, 7),   -- Mobile
(5, 8),           -- Men's Wear
(6, 9),           -- Watches
(7, 10), (8, 11), -- Home Appliances/Kitchen
(9, 12), (10, 13),-- Collectibles
(11, 14), (12, 14),-- Camping
(13, 1),          -- Electronics (Parent)
(14, 11),         -- Kitchenware
(15, 6),          -- Electronics
(16, 9),          -- Watches/Fashion
(17, 14),         -- Camping
(18, 14),         -- Camping
(19, 10),         -- Home/Appliances
(20, 10);         -- Home/Appliances

-- Update TS Vector.
update products set search_vector = to_tsvector('simple', name || ' ' || description);

-- Seeding PRODUCT IMAGES
-- Product 1: MacBook Air
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(1, 'http://dummyimage.com/600x400/2c3e50/ffffff.png&text=Macbook+Side+View', 'Side profile of the thin MacBook Air.'),
(1, 'http://dummyimage.com/600x400/2c3e50/ffffff.png&text=Macbook+Keyboard', 'Close up of the backlit Magic Keyboard.');

-- Product 2: Dell XPS
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(2, 'http://dummyimage.com/600x400/34495e/ffffff.png&text=XPS+Screen', 'InfinityEdge display of the Dell XPS.'),
(2, 'http://dummyimage.com/600x400/34495e/ffffff.png&text=XPS+Ports', 'View of the USB-C and Thunderbolt ports.');

-- Product 3: iPhone 12
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(3, 'http://dummyimage.com/600x400/003366/ffffff.png&text=iPhone+Back', 'The glass back of the Blue iPhone 12.'),
(3, 'http://dummyimage.com/600x400/003366/ffffff.png&text=iPhone+Box', 'Original packaging and included cables.');

-- Product 4: Pixel 6
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(4, 'http://dummyimage.com/600x400/222222/ffffff.png&text=Pixel+Camera+Bar', 'Close up of the iconic Pixel camera bar.'),
(4, 'http://dummyimage.com/600x400/222222/ffffff.png&text=Pixel+UI', 'Android interface running on Pixel 6.');

-- Product 5: Leather Jacket
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(5, 'http://dummyimage.com/600x400/3d2b1f/ffffff.png&text=Leather+Texture', 'Detailed shot of the lambskin leather grain.'),
(5, 'http://dummyimage.com/600x400/3d2b1f/ffffff.png&text=Jacket+Inside', 'Silk lining and internal pockets.');

-- Product 6: Seiko Watch
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(6, 'http://dummyimage.com/600x400/1a1a1a/ffffff.png&text=Seiko+Lume', 'The watch face glowing in low light.'),
(6, 'http://dummyimage.com/600x400/1a1a1a/ffffff.png&text=Seiko+Clasp', 'Stainless steel bracelet and diver clasp.');

-- Product 7: Air Fryer
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(7, 'http://dummyimage.com/600x400/555555/ffffff.png&text=Fryer+Basket', 'The non-stick basket inside the air fryer.'),
(7, 'http://dummyimage.com/600x400/555555/ffffff.png&text=Cooked+Fries', 'Demo shot of crispy fries made in the fryer.');

-- Product 8: Espresso Maker
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(8, 'http://dummyimage.com/600x400/8b4513/ffffff.png&text=Espresso+Shot', 'Fresh espresso pouring into a cup.'),
(8, 'http://dummyimage.com/600x400/8b4513/ffffff.png&text=Steam+Wand', 'Detail of the manual milk frothing wand.');

-- Product 9: Charizard Card
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(9, 'http://dummyimage.com/600x400/ff8c00/ffffff.png&text=Card+Back', 'The back of the Charizard card showing condition.'),
(9, 'http://dummyimage.com/600x400/ff8c00/ffffff.png&text=PSA+Label', 'Close up of the PSA 6 authentication label.');

-- Product 10: Vintage Car
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(10, 'http://dummyimage.com/600x400/cc0000/ffffff.png&text=Car+Underside', 'Matchbox branding on the base of the toy.'),
(10, 'http://dummyimage.com/600x400/cc0000/ffffff.png&text=Car+Front', 'The front grille and headlight detail.');

-- Product 11: Camping Tent
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(11, 'http://dummyimage.com/600x400/228b22/ffffff.png&text=Tent+Interior', 'Internal space showing gear pockets.'),
(11, 'http://dummyimage.com/600x400/228b22/ffffff.png&text=Tent+Bag', 'The tent packed down in its carry bag.');

-- Product 12: Sleeping Bag
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(12, 'http://dummyimage.com/600x400/000080/ffffff.png&text=Sleeping+Bag+Zipper', 'Heavy duty zipper and draft tube detail.'),
(12, 'http://dummyimage.com/600x400/000080/ffffff.png&text=Compressed+Bag', 'The sleeping bag inside its compression sack.');

-- Product 13: Sony Camera
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(13, 'http://dummyimage.com/600x400/000000/ffffff.png&text=Camera+Sensor', 'The pristine full-frame sensor view.'),
(13, 'http://dummyimage.com/600x400/000000/ffffff.png&text=Camera+LCD', 'The tilt-able LCD screen at the back.');

-- Product 14: Skillet Set
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(14, 'http://dummyimage.com/600x400/333333/ffffff.png&text=Stacked+Skillets', 'The three skillets nested together.'),
(14, 'http://dummyimage.com/600x400/333333/ffffff.png&text=Skillet+Handle', 'Close up of the cast iron handles.');

-- Product 15: Gaming Keyboard
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(15, 'http://dummyimage.com/600x400/7b68ee/ffffff.png&text=Keyboard+RGB', 'The RGB lighting in a dark room.'),
(15, 'http://dummyimage.com/600x400/7b68ee/ffffff.png&text=Keycaps', 'Close up of the double-shot PBT keycaps.');

-- Product 16: Ray-Bans
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(16, 'http://dummyimage.com/600x400/2f4f4f/ffffff.png&text=Glasses+Case', 'The original brown leather Ray-Ban case.'),
(16, 'http://dummyimage.com/600x400/2f4f4f/ffffff.png&text=Glasses+Side', 'The temple logo and hinge detail.');

-- Product 17: Propane Grill
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(17, 'http://dummyimage.com/600x400/d2691e/ffffff.png&text=Grill+Open', 'The internal cooking grate and burner.'),
(17, 'http://dummyimage.com/600x400/d2691e/ffffff.png&text=Grill+Knobs', 'The temperature control dial and igniter.');

-- Product 18: Water Bottle
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(18, 'http://dummyimage.com/600x400/4682b4/ffffff.png&text=Bottle+Cap', 'The leak-proof screw cap and handle.'),
(18, 'http://dummyimage.com/600x400/4682b4/ffffff.png&text=Bottle+Bottom', 'The non-slip silicone base.');

-- Product 19: LED Bulbs
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(19, 'http://dummyimage.com/600x400/ffff00/333333.png&text=Bulb+Colors', 'Demonstration of different RGB colors.'),
(19, 'http://dummyimage.com/600x400/ffff00/333333.png&text=App+Control', 'Smartphone app controlling the light bulbs.');

-- Product 20: Standing Desk
INSERT INTO product_images (product_id, url, alt_text) VALUES 
(20, 'http://dummyimage.com/600x400/ffffff/000000.png&text=Desk+Controller', 'Digital display showing height settings.'),
(20, 'http://dummyimage.com/600x400/ffffff/000000.png&text=Desk+Motor', 'View of the dual motor assembly.');

-- Seed bids

-- Product 1: MacBook Air (Start: 150, BIN: 850)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(160.00, false, 1, 3, now() - interval '4 days'),
(170.00, false, 1, 4, now() - interval '3 days'),
(185.00, true, 1, 2, now() - interval '2 days'),
(210.00, false, 1, 3, now() - interval '1 day'),
(250.00, false, 1, 4, now() - interval '10 hours');

-- Product 2: Dell XPS (Start: 200, BIN: 1200)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(220.00, false, 2, 1, now() - interval '2 days'),
(250.00, true, 2, 3, now() - interval '1 day'),
(300.00, false, 2, 4, now() - interval '5 hours');

-- Product 3: iPhone 12 (Start: 100, BIN: 600)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(110.00, false, 3, 1, now() - interval '1 day'),
(125.00, false, 3, 2, now() - interval '12 hours'),
(135.00, true, 3, 4, now() - interval '2 hours');

-- Product 4: Pixel 6 (Start: 80, BIN: 450)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(85.00, false, 4, 3, now() - interval '3 days'),
(95.00, false, 4, 2, now() - interval '2 days');

-- Product 5: Leather Jacket (Start: 45, BIN: 150)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(50.00, false, 5, 1, now() - interval '5 days'),
(55.00, false, 5, 2, now() - interval '4 days'),
(65.00, true, 5, 3, now() - interval '2 days');

-- Product 6: Seiko Watch (Start: 120, BIN: 550)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(130.00, false, 6, 1, now() - interval '18 hours'),
(150.00, false, 6, 2, now() - interval '10 hours'),
(180.00, false, 6, 3, now() - interval '2 hours');

-- Product 7: Air Fryer (Start: 30, BIN: 180)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(35.00, false, 7, 2, now() - interval '4 days'),
(45.00, false, 7, 3, now() - interval '2 days'),
(50.00, true, 7, 4, now() - interval '1 day');

-- Product 8: Espresso Maker (Start: 25, BIN: 120)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(30.00, false, 8, 1, now() - interval '6 days'),
(35.00, false, 8, 2, now() - interval '4 days'),
(40.00, false, 8, 4, now() - interval '1 day');

-- Product 9: Charizard Card (Start: 50, BIN: 500)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(100.00, false, 9, 1, now() - interval '10 hours'),
(200.00, false, 9, 3, now() - interval '8 hours'),
(350.00, true, 9, 4, now() - interval '4 hours'),
(360.00, false, 9, 1, now() - interval '2 hours');

-- Product 10: Vintage Car (Start: 15, BIN: 75)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(16.00, false, 10, 2, now() - interval '2 days'),
(20.00, false, 10, 1, now() - interval '1 day');

-- Product 11: Camping Tent (Start: 40, BIN: 220)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(45.00, false, 11, 2, now() - interval '3 days'),
(55.00, false, 11, 3, now() - interval '2 days'),
(70.00, true, 11, 4, now() - interval '1 day');

-- Product 12: Sleeping Bag (Start: 20, BIN: 95)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(22.00, false, 12, 1, now() - interval '4 days'),
(26.00, false, 12, 4, now() - interval '2 days');

-- Product 13: Sony Camera (Start: 299, BIN: 1500)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(320.00, false, 13, 3, now() - interval '1 day'),
(350.00, false, 13, 4, now() - interval '12 hours'),
(400.00, true, 13, 1, now() - interval '5 hours');

-- Product 14: Skillet Set (Start: 25, BIN: 110)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(30.00, false, 14, 2, now() - interval '5 days'),
(40.00, false, 14, 4, now() - interval '3 days');

-- Product 15: Gaming Keyboard (Start: 10, BIN: 45)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(12.00, false, 15, 1, now() - interval '1 day'),
(15.00, false, 15, 2, now() - interval '18 hours'),
(18.00, true, 15, 3, now() - interval '6 hours');

-- Product 16: Ray-Bans (Start: 60, BIN: 300)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(70.00, false, 16, 2, now() - interval '3 days'),
(85.00, false, 16, 3, now() - interval '2 days'),
(100.00, true, 16, 4, now() - interval '1 day');

-- Product 17: Propane Grill (Start: 35, BIN: 130)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(40.00, false, 17, 1, now() - interval '6 days'),
(45.00, false, 17, 2, now() - interval '4 days');

-- Product 18: Water Bottle (Start: 5, BIN: 25)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(6.00, false, 18, 3, now() - interval '2 days'),
(7.50, false, 18, 4, now() - interval '1 day'),
(9.00, true, 18, 1, now() - interval '12 hours');

-- Product 19: LED Bulbs (Start: 12, BIN: 60)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(15.00, false, 19, 2, now() - interval '4 days'),
(20.00, false, 19, 3, now() - interval '2 days');

-- Product 20: Standing Desk (Start: 85, BIN: 400)
INSERT INTO bids (price, automated, product_id, user_id, created_at) VALUES 
(100.00, false, 20, 1, now() - interval '8 days'),
(125.00, false, 20, 3, now() - interval '5 days'),
(150.00, true, 20, 4, now() - interval '2 days');

update products p set bids_count = (select count(1) from bids where product_id = p.id);
update products p set current_highest_bid_id = (select id from bids where product_id = p.id order by price desc limit 1);

-- Seed questions

-- Product 1: MacBook Air
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('Is the battery original or has it been replaced?', 'It is the original battery, currently at 92% health.', 1, 4, now() - interval '2 days'),
('Does it come with the original box?', NULL, 1, 3, now() - interval '1 day');

-- Product 2: Dell XPS
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('Are there any dead pixels on the 4K screen?', 'No dead pixels at all, the display is pristine.', 2, 4, now() - interval '1 day');

-- Product 3: iPhone 12
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('Is FaceID working perfectly?', 'Yes, all sensors including FaceID work as expected.', 3, 2, now() - interval '12 hours'),
('Can you ship it via express mail?', NULL, 3, 4, now() - interval '5 hours');

-- Product 5: Leather Jacket
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('What is the pit-to-pit measurement in inches?', 'It is exactly 21 inches.', 5, 1, now() - interval '3 days');

-- Product 6: Seiko Watch
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('Does the bezel align perfectly at 12?', 'Yes, alignment is spot on for this specific unit.', 6, 3, now() - interval '10 hours');

-- Product 7: Air Fryer
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('Is the basket dishwasher safe?', NULL, 7, 2, now() - interval '1 day');

-- Product 8: Espresso Maker
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('Does it include the portafilter and tampt?', 'Yes, everything in the original retail box is included.', 8, 4, now() - interval '2 days');

-- Product 9: Charizard Card
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('Are there any visible scratches on the PSA case?', 'The case has some very light surface scuffs but no cracks.', 9, 1, now() - interval '5 hours'),
('Will this be shipped in a box or a bubble mailer?', NULL, 9, 3, now() - interval '2 hours');

-- Product 11: Camping Tent
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('How heavy is the tent when packed?', 'It weighs approximately 18 lbs.', 11, 4, now() - interval '2 days');

-- Product 13: Sony Camera
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('What is the exact shutter count?', 'Shutter count is currently 12,450.', 13, 1, now() - interval '1 day');

-- Product 15: Gaming Keyboard
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('Are the switches hot-swappable?', 'No, these are soldered Blue switches.', 15, 2, now() - interval '10 hours');

-- Product 16: Ray-Bans
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('Are these the 50mm or 54mm size?', NULL, 16, 3, now() - interval '1 day');

-- Product 20: Standing Desk
INSERT INTO questions (content, answer, product_id, user_id, created_at) VALUES 
('What is the maximum weight capacity?', 'It can support up to 250 lbs easily.', 20, 1, now() - interval '3 days');
