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