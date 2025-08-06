-- file: db/init/05-add-comunicados.sql
-- Migration para adicionar tabela de comunicados (ADR-003)

-- =================================================================
-- TABELA DE COMUNICADOS
-- =================================================================

CREATE TABLE comunicados (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_evento UUID NOT NULL,
    titulo VARCHAR(255) NOT NULL,
    mensagem TEXT NOT NULL,
    data_publicacao TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now()),

    -- Garante a relação com o evento ao qual o comunicado pertence
    CONSTRAINT fk_evento_comunicado
        FOREIGN KEY(id_evento) 
        REFERENCES eventos(id)
        ON DELETE CASCADE
);

-- Índices para melhorar performance
CREATE INDEX idx_comunicados_id_evento ON comunicados(id_evento);
CREATE INDEX idx_comunicados_data_publicacao ON comunicados(data_publicacao DESC);