DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS chat_members CASCADE;
DROP TABLE IF EXISTS chats CASCADE;
DROP TABLE IF EXISTS users_posts_likes CASCADE;
DROP TABLE IF EXISTS users_posts CASCADE;
DROP FUNCTION IF EXISTS fullname(varchar, varchar) CASCADE;
DROP TABLE IF EXISTS friends_invites CASCADE;
DROP TABLE IF EXISTS users_info CASCADE;
DROP TABLE IF EXISTS tokens CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP EXTENSION IF EXISTS pg_trgm;