-- file: db/init/01-init.sql - VERSÃO CORRIGIDA

-- =================================================================
-- TIPOS ENUM GERAIS
-- =================================================================
CREATE TYPE status_rsvp AS ENUM ('PENDENTE', 'CONFIRMADO', 'RECUSADO');
CREATE TYPE tipo_detalhe_presente AS ENUM ('PRODUTO_EXTERNO', 'PIX');
CREATE TYPE status_presente AS ENUM ('DISPONIVEL', 'SELECIONADO');
CREATE TYPE status_recado AS ENUM ('PENDENTE', 'APROVADO', 'REJEITADO');
CREATE TYPE nome_rotulo_enum AS ENUM ('MAIN', 'CASAMENTO', 'LUADEMEL', 'HISTORIA', 'FAMILIA', 'OUTROS');
CREATE TYPE status_assinatura AS ENUM ('PENDENTE', 'ATIVA', 'EXPIRADA', 'CANCELADA');
CREATE TYPE tipo_evento AS ENUM ('CASAMENTO', 'ANIVERSARIO', 'CHA_DE_BEBE', 'OUTRO');


-- =================================================================
-- TABELAS DE PLATAFORMA (V4)
-- =================================================================

CREATE TABLE usuarios (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    telefone VARCHAR(20),
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now())
);

CREATE TABLE planos (
    id UUID PRIMARY KEY,
    nome VARCHAR(100) NOT NULL UNIQUE,
    preco_em_centavos INTEGER NOT NULL,
    numero_maximo_eventos INTEGER NOT NULL,
    duracao_em_dias INTEGER NOT NULL,
    id_stripe_price VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE assinaturas (
    id UUID PRIMARY KEY,
    id_usuario UUID NOT NULL REFERENCES usuarios(id),
    id_plano UUID NOT NULL REFERENCES planos(id),
    data_inicio TIMESTAMP WITH TIME ZONE,
    data_fim TIMESTAMP WITH TIME ZONE,
    status status_assinatura NOT NULL DEFAULT 'PENDENTE'
);

CREATE TABLE eventos (
    id UUID PRIMARY KEY,
    id_usuario UUID NOT NULL REFERENCES usuarios(id),
    nome VARCHAR(255) NOT NULL,
    data DATE,
    tipo tipo_evento NOT NULL,
    url_slug VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now())
);


-- =================================================================
-- TABELAS DE EVENTO (V1, V2, V3)
-- =================================================================

CREATE TABLE convidados_grupos (
    id UUID PRIMARY KEY,
    id_evento UUID NOT NULL REFERENCES eventos(id) ON DELETE CASCADE,
    chave_de_acesso VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now()),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now()),
    UNIQUE(id_evento, chave_de_acesso)
);

CREATE TABLE convidados (
    id UUID PRIMARY KEY,
    id_grupo UUID NOT NULL REFERENCES convidados_grupos(id) ON DELETE CASCADE,
    nome VARCHAR(255) NOT NULL,
    status_rsvp status_rsvp NOT NULL DEFAULT 'PENDENTE'
);

CREATE TABLE presentes_selecoes (
    id UUID PRIMARY KEY,
    id_evento UUID NOT NULL REFERENCES eventos(id) ON DELETE CASCADE,
    id_grupo_de_convidados UUID NOT NULL REFERENCES convidados_grupos(id),
    data_da_selecao TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now())
);

CREATE TABLE presentes (
    id UUID PRIMARY KEY,
    id_evento UUID NOT NULL REFERENCES eventos(id) ON DELETE CASCADE,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    foto_url TEXT,
    eh_favorito BOOLEAN NOT NULL DEFAULT FALSE,
    status status_presente NOT NULL DEFAULT 'DISPONIVEL',
    categoria nome_rotulo_enum,
    detalhes_tipo tipo_detalhe_presente NOT NULL,
    detalhes_link_loja TEXT,
    detalhes_chave_pix VARCHAR(255),
    id_selecao UUID REFERENCES presentes_selecoes(id),
    CONSTRAINT chk_detalhes CHECK (
        (detalhes_tipo = 'PRODUTO_EXTERNO' AND detalhes_link_loja IS NOT NULL) OR
        (detalhes_tipo = 'PIX' AND detalhes_chave_pix IS NOT NULL)
    )
);

CREATE TABLE recados (
    id UUID PRIMARY KEY,
    id_evento UUID NOT NULL REFERENCES eventos(id) ON DELETE CASCADE,
    id_grupo_de_convidados UUID NOT NULL REFERENCES convidados_grupos(id),
    nome_do_autor VARCHAR(255) NOT NULL,
    texto TEXT NOT NULL,
    status status_recado NOT NULL DEFAULT 'PENDENTE',
    eh_favorito BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now())
);

CREATE TABLE fotos (
    id UUID PRIMARY KEY,
    id_evento UUID NOT NULL REFERENCES eventos(id) ON DELETE CASCADE,
    storage_key TEXT NOT NULL,
    url_publica TEXT NOT NULL,
    eh_favorito BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now())
);

CREATE TABLE fotos_rotulos (
    id_foto UUID NOT NULL REFERENCES fotos(id) ON DELETE CASCADE,
    nome_rotulo nome_rotulo_enum NOT NULL,
    PRIMARY KEY (id_foto, nome_rotulo)
);