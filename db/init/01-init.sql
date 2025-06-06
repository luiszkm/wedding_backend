-- file: db/init/01-init.sql

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

-- Você pode adicionar aqui as outras tabelas do projeto (presentes, selecoes_de_presentes)
-- para que todo o esquema seja criado de uma só vez.