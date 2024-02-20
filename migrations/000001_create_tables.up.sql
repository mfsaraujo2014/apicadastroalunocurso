CREATE TABLE IF NOT EXISTS aluno (
    codigo SERIAL PRIMARY KEY,
    nome VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS curso (
    codigo SERIAL PRIMARY KEY,
    descricao VARCHAR(50) NOT NULL,
    ementa TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS curso_aluno (
    codigo SERIAL PRIMARY KEY,
    codigo_aluno INT NOT NULL,
    codigo_curso INT NOT NULL,
    FOREIGN KEY (codigo_aluno) REFERENCES aluno(codigo),
    FOREIGN KEY (codigo_curso) REFERENCES curso(codigo)
);

INSERT INTO aluno (nome) VALUES
    ('Matheus'),
    ('Felipe'),
    ('JOAO PAULO MANTOVANI'),
    ('Maria Clara Almeida');

INSERT INTO curso (descricao, ementa) VALUES
    ('Curso 1', 'Ementa do Curso 1'),
    ('Curso 2', 'Ementa do Curso 2'),
    ('Curso 3', 'Ementa do Curso 3');