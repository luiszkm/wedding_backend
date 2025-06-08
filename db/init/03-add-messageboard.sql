-- =================================================================
-- CONTEXTO: MURAL DE RECADOS (V2)
-- =================================================================

-- Novo tipo ENUM para o status do recado
CREATE TYPE status_recado AS ENUM ('PENDENTE', 'APROVADO', 'REJEITADO');

-- Nova tabela para o Agregado Recado
CREATE TABLE recados (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_casamento UUID NOT NULL,
    id_grupo_de_convidados UUID NOT NULL REFERENCES grupos_de_convidados(id),
    nome_do_autor VARCHAR(255) NOT NULL,
    texto TEXT NOT NULL,
    status status_recado NOT NULL DEFAULT 'PENDENTE',
    eh_favorito BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now())
);