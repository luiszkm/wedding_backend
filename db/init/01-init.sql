-- file: db/init/01-init.sql

-- =================================================================
-- CONTEXTO: GESTÃO DE CONVIDADOS
-- =================================================================

-- Definição de Tipos Enum para consistência
CREATE TYPE status_rsvp AS ENUM ('PENDENTE', 'CONFIRMADO', 'RECUSADO');

-- Tabela para o Agregado GrupoDeConvidados
CREATE TABLE grupos_de_convidados (
    id UUID PRIMARY KEY,
    id_casamento UUID NOT NULL,
    chave_de_acesso VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now()),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now()),
    UNIQUE(id_casamento, chave_de_acesso)
);

-- Tabela para a Entidade Convidado
CREATE TABLE convidados (
    id UUID PRIMARY KEY,
    id_grupo UUID NOT NULL REFERENCES grupos_de_convidados(id) ON DELETE CASCADE,
    nome VARCHAR(255) NOT NULL,
    status_rsvp status_rsvp NOT NULL DEFAULT 'PENDENTE'
);

-- =================================================================
-- CONTEXTO: LISTA DE PRESENTES
-- =================================================================

-- Definição de Tipos Enum para consistência
CREATE TYPE tipo_detalhe_presente AS ENUM ('PRODUTO_EXTERNO', 'PIX');
CREATE TYPE status_presente AS ENUM ('DISPONIVEL', 'SELECIONADO');

-- Tabela para a Entidade Presente
CREATE TABLE presentes (
    id UUID PRIMARY KEY,
    id_casamento UUID NOT NULL,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    foto_url TEXT,
    eh_favorito BOOLEAN NOT NULL DEFAULT FALSE,
    status status_presente NOT NULL DEFAULT 'DISPONIVEL',
    detalhes_tipo tipo_detalhe_presente NOT NULL,
    detalhes_link_loja TEXT,
    detalhes_chave_pix VARCHAR(255),
    id_selecao UUID,
    CONSTRAINT chk_detalhes CHECK (
        (detalhes_tipo = 'PRODUTO_EXTERNO' AND detalhes_link_loja IS NOT NULL) OR
        (detalhes_tipo = 'PIX' AND detalhes_chave_pix IS NOT NULL)
    )
);

-- Tabela para o Agregado SelecaoDePresentes
CREATE TABLE selecoes_de_presentes (
    id UUID PRIMARY KEY,
    id_casamento UUID NOT NULL,
    id_grupo_de_convidados UUID NOT NULL REFERENCES grupos_de_convidados(id),
    data_da_selecao TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now())
);

-- Adiciona a referência que faltava na tabela de presentes
ALTER TABLE presentes ADD CONSTRAINT fk_selecao FOREIGN KEY (id_selecao) REFERENCES selecoes_de_presentes(id);