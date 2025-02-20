CREATE TABLE wiki_categories (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    parent_id int REFERENCES wiki_categories (id) DEFAULT null,
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE wiki_pages (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    last_updated timestamp without time zone NOT NULL DEFAULT now(),
    category_id int REFERENCES wiki_categories (id)
);

CREATE TABLE wiki_content (
    page_id int REFERENCES wiki_pages (id),
    language varchar(255) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    last_updated timestamp without time zone NOT NULL DEFAULT now(),
    content text NOT NULL,
    search tsvector NOT NULL GENERATED ALWAYS AS (
        setweight(to_tsvector('simple', coalesce(content, '')), 'A') :: tsvector
    ) STORED,
    PRIMARY KEY (page_id, language)
);

CREATE TABLE wiki_outlinks (
    page_id int REFERENCES wiki_pages (id),
    target_id int REFERENCES wiki_pages (id),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (page_id, target_id)
);

CREATE INDEX wiki_pages_name_idx ON wiki_pages (name);
CREATE INDEX wiki_content_page_id_language_idx ON wiki_content (page_id, language);
CREATE INDEX wiki_content_search_idx ON wiki_content USING gin (search);
CREATE INDEX wiki_outlinks_page_id ON wiki_outlinks (page_id);
CREATE INDEX wiki_outlinks_target_id ON wiki_outlinks (target_id);