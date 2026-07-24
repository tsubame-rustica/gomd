CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    dir_name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    category_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    UNIQUE (category_id, file_name),
    UNIQUE (category_id, title)
);

-- テストデータ
INSERT INTO categories (name, dir_name) VALUES 
('Git', 'git'),
('Linux', 'linux');

INSERT INTO posts (title, file_name, category_id) VALUES 
('Gitについて', 'default', 1),
('ブランチの命名規則あれこれ', 'branch', 1),
('Linuxのコマンドまとめ', 'default', 2);