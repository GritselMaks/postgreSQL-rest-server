CREATE TABLE IF NOT EXISTS articles(
   id serial PRIMARY KEY,
   title VARCHAR(200) NOT NULL,
   full_text VARCHAR(1000) NOT NULL,
   price integer NOT NULL,
   url_foto text[3] NOT NULL,
   data_at date NOT NULL
);