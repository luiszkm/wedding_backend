-- file: db/init/10-create-missing-cotas-table.sql
-- Hotfix para criar a tabela cotas_de_presentes que estava faltando

-- =================================================================
-- PROBLEMA IDENTIFICADO
-- =================================================================
-- A migration 06-add-fractional-gifts.sql foi parcialmente aplicada:
-- - O enum tipo_presente existe
-- - A coluna tipo foi adicionada à tabela presentes
-- - Mas a tabela cotas_de_presentes não foi criada
-- 
-- Erro: relation "cotas_de_presentes" does not exist

-- =================================================================
-- SOLUÇÃO
-- =================================================================

-- Criar a tabela cotas_de_presentes que estava faltando
CREATE TABLE IF NOT EXISTS cotas_de_presentes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_presente UUID NOT NULL,
    numero_cota INTEGER NOT NULL,
    valor_cota NUMERIC(10,2) NOT NULL,
    status status_presente NOT NULL DEFAULT 'DISPONIVEL',
    id_selecao UUID,
    
    -- Foreign keys
    CONSTRAINT fk_cota_presente 
        FOREIGN KEY (id_presente) 
        REFERENCES presentes(id) 
        ON DELETE CASCADE,
        
    CONSTRAINT fk_cota_selecao 
        FOREIGN KEY (id_selecao) 
        REFERENCES presentes_selecoes(id) 
        ON DELETE SET NULL,
    
    -- Constraints de negócio
    CONSTRAINT uk_presente_numero_cota 
        UNIQUE (id_presente, numero_cota),
        
    CONSTRAINT chk_numero_cota_positivo 
        CHECK (numero_cota > 0),
        
    CONSTRAINT chk_valor_cota_positivo 
        CHECK (valor_cota > 0)
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_cotas_id_presente ON cotas_de_presentes(id_presente);
CREATE INDEX IF NOT EXISTS idx_cotas_status ON cotas_de_presentes(status);
CREATE INDEX IF NOT EXISTS idx_cotas_id_selecao ON cotas_de_presentes(id_selecao);

-- =================================================================
-- VERIFICAÇÃO
-- =================================================================
-- Para verificar se a tabela foi criada:
-- \dt cotas_de_presentes