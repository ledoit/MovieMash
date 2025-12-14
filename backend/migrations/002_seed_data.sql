-- Sample Movie Data for MovieMash

-- Insert anonymous user
INSERT INTO users (username, email) VALUES ('anonymous', 'anonymous@moviemash.local')
ON CONFLICT (username) DO NOTHING;

-- Sample Movies (popular films for testing)
INSERT INTO movies (title, year, director, poster_url) VALUES
('The Godfather', 1972, 'Francis Ford Coppola', 'https://image.tmdb.org/t/p/w500/3bhkrj58Vtu7enYsRolD1fZdja1.jpg'),
('Pulp Fiction', 1994, 'Quentin Tarantino', 'https://image.tmdb.org/t/p/w500/d5iIlFn5s0ImszYzBPb8JPIfbXD.jpg'),
('The Dark Knight', 2008, 'Christopher Nolan', 'https://image.tmdb.org/t/p/w500/qJ2tW6WMUDux911r6m7haRef0WH.jpg'),
('Fight Club', 1999, 'David Fincher', 'https://image.tmdb.org/t/p/w500/pB8BM7pdSp6B6Ih7QZ4DrQ3PmJK.jpg'),
('Inception', 2010, 'Christopher Nolan', 'https://image.tmdb.org/t/p/w500/oYuLEt3zVCKq57qu2F8dT7NIa6f.jpg'),
('The Matrix', 1999, 'Lana Wachowski, Lilly Wachowski', 'https://image.tmdb.org/t/p/w500/f89U3ADr1oiB1s9GkdPOEpXUk5H.jpg'),
('Goodfellas', 1990, 'Martin Scorsese', 'https://image.tmdb.org/t/p/w500/aKuFiU82s5ISJpGZp7YkIr3kCUd.jpg'),
('Parasite', 2019, 'Bong Joon-ho', 'https://image.tmdb.org/t/p/w500/7IiTTgloJzvGI1TAYymCfbfl3vT.jpg'),
('Interstellar', 2014, 'Christopher Nolan', 'https://image.tmdb.org/t/p/w500/gEU2QniE6E77NI6lCU6MxlNBvIx.jpg'),
('The Shawshank Redemption', 1994, 'Frank Darabont', 'https://image.tmdb.org/t/p/w500/q6y0Go1tsGEsmtFryDOJo3dEmqu.jpg'),
('Forrest Gump', 1994, 'Robert Zemeckis', 'https://image.tmdb.org/t/p/w500/arw2vcBve2OV1haIhZQ9vQnVq2j.jpg'),
('The Lord of the Rings: The Fellowship of the Ring', 2001, 'Peter Jackson', 'https://image.tmdb.org/t/p/w500/6oom5QYQ2yQHMJIhnvC7A13pG1Z.jpg'),
('The Lord of the Rings: The Return of the King', 2003, 'Peter Jackson', 'https://image.tmdb.org/t/p/w500/rCzpDGLbOoPwLjy3OAm5NUPOTrC.jpg'),
('The Lord of the Rings: The Two Towers', 2002, 'Peter Jackson', 'https://image.tmdb.org/t/p/w500/5VTN0pR8gcqV3EPUHHfMGnJYN9G.jpg'),
('Spirited Away', 2001, 'Hayao Miyazaki', 'https://image.tmdb.org/t/p/w500/39wmItIWwg5DnHv6Iyibz6gWiFX.jpg'),
('City of God', 2002, 'Fernando Meirelles', 'https://image.tmdb.org/t/p/w500/k7eYdWvhYQyRQoU2TB2A2Xu2TfD.jpg'),
('Se7en', 1995, 'David Fincher', 'https://image.tmdb.org/t/p/w500/69Sns8WoET6CfaYlIkHbla4l7nC.jpg'),
('The Silence of the Lambs', 1991, 'Jonathan Demme', 'https://image.tmdb.org/t/p/w500/uS9m8OBk1A8eM9I042bx8XXpqAq.jpg'),
('The Departed', 2006, 'Martin Scorsese', 'https://image.tmdb.org/t/p/w500/nT97ifVT2J1yMQme168qJ8s36p9.jpg'),
('Whiplash', 2014, 'Damien Chazelle', 'https://image.tmdb.org/t/p/w500/7fn624j5lj3xTme2SgiLCeuedmO.jpg'),
('Gladiator', 2000, 'Ridley Scott', 'https://image.tmdb.org/t/p/w500/ty8T9d8eFXCu1Y4iAiA0xqjv9n.jpg'),
('The Prestige', 2006, 'Christopher Nolan', 'https://image.tmdb.org/t/p/w500/5MXyQfz8xUP3dIFPTubhTsbFY6N.jpg'),
('Memento', 2000, 'Christopher Nolan', 'https://image.tmdb.org/t/p/w500/yuNs09hvpHVU1cBTCAk9zxsL2zT.jpg'),
('Django Unchained', 2012, 'Quentin Tarantino', 'https://image.tmdb.org/t/p/w500/7o3Ylzpf7BD2HQfIVoJsMk3yHhX.jpg'),
('Inglourious Basterds', 2009, 'Quentin Tarantino', 'https://image.tmdb.org/t/p/w500/7sfbEnaARXDDhKm0CZ7D7uc2sbo.jpg'),
('The Social Network', 2010, 'David Fincher', 'https://image.tmdb.org/t/p/w500/ok5Wh8385Kgblq9MSU4VGvazeMH.jpg'),
('Zodiac', 2007, 'David Fincher', 'https://image.tmdb.org/t/p/w500/6lVYHxq1n0nAXOhN2V0Yq4hH2bJ.jpg'),
('No Country for Old Men', 2007, 'Joel Coen, Ethan Coen', 'https://image.tmdb.org/t/p/w500/4m1Q3ij3F3L3px6pxQvT6VIy0o9.jpg'),
('There Will Be Blood', 2007, 'Paul Thomas Anderson', 'https://image.tmdb.org/t/p/w500/a4Y0VKdQ8K8b5X5pz0b3Z4Z4Z4Z.jpg'),
('The Big Lebowski', 1998, 'Joel Coen, Ethan Coen', 'https://image.tmdb.org/t/p/w500/64sImqYy3e3M3Y3Y3Y3Y3Y3Y3Y3Y.jpg'),
('Fargo', 1996, 'Joel Coen, Ethan Coen', 'https://image.tmdb.org/t/p/w500/4m1Q3ij3F3L3px6pxQvT6VIy0o9.jpg'),
('Oldboy', 2003, 'Park Chan-wook', 'https://image.tmdb.org/t/p/w500/pWDtjs568ZfOTMbURQPYuH4Vw3C.jpg'),
('Memories of Murder', 2003, 'Bong Joon-ho', 'https://image.tmdb.org/t/p/w500/5VTN0pR8gcqV3EPUHHfMGnJYN9G.jpg'),
('The Handmaiden', 2016, 'Park Chan-wook', 'https://image.tmdb.org/t/p/w500/7o3Ylzpf7BD2HQfIVoJsMk3yHhX.jpg'),
('Burning', 2018, 'Lee Chang-dong', 'https://image.tmdb.org/t/p/w500/5VTN0pR8gcqV3EPUHHfMGnJYN9G.jpg'),
('Portrait of a Lady on Fire', 2019, 'Céline Sciamma', 'https://image.tmdb.org/t/p/w500/7o3Ylzpf7BD2HQfIVoJsMk3yHhX.jpg'),
('Amélie', 2001, 'Jean-Pierre Jeunet', 'https://image.tmdb.org/t/p/w500/5VTN0pR8gcqV3EPUHHfMGnJYN9G.jpg')
ON CONFLICT DO NOTHING;

-- Get anonymous user ID
DO $$
DECLARE
    anon_user_id INTEGER;
BEGIN
    SELECT id INTO anon_user_id FROM users WHERE username = 'anonymous';
    
    -- Create sample top 4 sets (using first 40 movies, creating 10 sets)
    -- Each set has 4 movies
    INSERT INTO top4_sets (user_id, movie_ids) VALUES
    (anon_user_id, ARRAY[1, 2, 3, 4]),
    (anon_user_id, ARRAY[5, 6, 7, 8]),
    (anon_user_id, ARRAY[9, 10, 11, 12]),
    (anon_user_id, ARRAY[13, 14, 15, 16]),
    (anon_user_id, ARRAY[17, 18, 19, 20]),
    (anon_user_id, ARRAY[21, 22, 23, 24]),
    (anon_user_id, ARRAY[25, 26, 27, 28]),
    (anon_user_id, ARRAY[29, 30, 31, 32]),
    (anon_user_id, ARRAY[33, 34, 35, 36]),
    (anon_user_id, ARRAY[37, 38, 39, 40])
    ON CONFLICT DO NOTHING;
END $$;

