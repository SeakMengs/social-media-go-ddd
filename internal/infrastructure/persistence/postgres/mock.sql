-- Posts
INSERT INTO posts (id, user_id, content)
VALUES
(uuid_generate_v4(), '2453e890-515d-4403-b8ff-3f2dd9be32d0', 'Hello world! This is my first post.'),
(uuid_generate_v4(), '26c251da-0613-4671-b414-e6c1a86e2e60', 'Enjoying a sunny day at the park.'),
(uuid_generate_v4(), 'efe61ee3-784d-4718-b47e-3f57d3dc675c', 'Just finished reading a great book.');

-- Likes
INSERT INTO likes (id, user_id, post_id)
VALUES
(uuid_generate_v4(), '2453e890-515d-4403-b8ff-3f2dd9be32d0', (SELECT id FROM posts WHERE user_id='26c251da-0613-4671-b414-e6c1a86e2e60' LIMIT 1)),
(uuid_generate_v4(), '26c251da-0613-4671-b414-e6c1a86e2e60', (SELECT id FROM posts WHERE user_id='efe61ee3-784d-4718-b47e-3f57d3dc675c' LIMIT 1)),
(uuid_generate_v4(), 'efe61ee3-784d-4718-b47e-3f57d3dc675c', (SELECT id FROM posts WHERE user_id='2453e890-515d-4403-b8ff-3f2dd9be32d0' LIMIT 1));

-- Favorites
INSERT INTO favorites (id, user_id, post_id)
VALUES
(uuid_generate_v4(), '2453e890-515d-4403-b8ff-3f2dd9be32d0', (SELECT id FROM posts WHERE user_id='efe61ee3-784d-4718-b47e-3f57d3dc675c' LIMIT 1)),
(uuid_generate_v4(), '26c251da-0613-4671-b414-e6c1a86e2e60', (SELECT id FROM posts WHERE user_id='2453e890-515d-4403-b8ff-3f2dd9be32d0' LIMIT 1)),
(uuid_generate_v4(), 'efe61ee3-784d-4718-b47e-3f57d3dc675c', (SELECT id FROM posts WHERE user_id='26c251da-0613-4671-b414-e6c1a86e2e60' LIMIT 1));

-- Reposts
INSERT INTO reposts (id, user_id, post_id, comment)
VALUES
(uuid_generate_v4(), '2453e890-515d-4403-b8ff-3f2dd9be32d0', (SELECT id FROM posts WHERE user_id='efe61ee3-784d-4718-b47e-3f57d3dc675c' LIMIT 1), 'Sharing this interesting post!'),
(uuid_generate_v4(), '26c251da-0613-4671-b414-e6c1a86e2e60', (SELECT id FROM posts WHERE user_id='2453e890-515d-4403-b8ff-3f2dd9be32d0' LIMIT 1), 'Check this out!'),
(uuid_generate_v4(), 'efe61ee3-784d-4718-b47e-3f57d3dc675c', (SELECT id FROM posts WHERE user_id='26c251da-0613-4671-b414-e6c1a86e2e60' LIMIT 1), 'Loved this!');