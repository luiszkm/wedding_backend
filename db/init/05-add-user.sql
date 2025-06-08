CREATE TYPE status_assinatura AS ENUM ('PENDENTE', 'ATIVA', 'EXPIRADA', 'CANCELADA');

CREATE TABLE usuarios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    telefone VARCHAR(20),
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now())
);

CREATE TABLE planos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome VARCHAR(100) NOT NULL UNIQUE,
    preco_em_centavos INTEGER NOT NULL,
    numero_maximo_eventos INTEGER NOT NULL,
    duracao_em_dias INTEGER NOT NULL
);

CREATE TABLE assinaturas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_usuario UUID NOT NULL REFERENCES usuarios(id),
    id_plano UUID NOT NULL REFERENCES planos(id),
    data_inicio TIMESTAMP WITH TIME ZONE,
    data_fim TIMESTAMP WITH TIME ZONE,
    status status_assinatura NOT NULL DEFAULT 'PENDENTE'
);

CREATE TABLE casamentos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_usuario UUID NOT NULL REFERENCES usuarios(id),
    nome_evento VARCHAR(255) NOT NULL,
    data_evento DATE,
    url_slug VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now())
);

