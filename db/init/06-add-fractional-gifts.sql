-- file: db/init/06-add-fractional-gifts.sql
-- Migration para suportar presentes fracionados (ADR-004)

-- =================================================================
-- ATUALIZAÇÃO DE TIPOS ENUM
-- =================================================================

-- Adicionar novo tipo de presente
CREATE TYPE tipo_presente AS ENUM ('INTEGRAL', 'FRACIONADO');

-- Atualizar status do presente para incluir PARCIALMENTE_SELECIONADO
DROP TYPE IF EXISTS status_presente CASCADE;
CREATE TYPE status_presente AS ENUM ('DISPONIVEL', 'SELECIONADO', 'PARCIALMENTE_SELECIONADO');

-- =================================================================
-- ALTERAÇÕES NA TABELA PRESENTES
-- =================================================================

-- Adicionar campos para suportar presentes fracionados
ALTER TABLE presentes 
ADD COLUMN tipo tipo_presente NOT NULL DEFAULT 'INTEGRAL',
ADD COLUMN valor_total_presente NUMERIC(10,2);

-- Recriar constraint de status (foi removida quando dropamos o enum)
ALTER TABLE presentes 
ALTER COLUMN status SET DEFAULT 'DISPONIVEL';

-- =================================================================
-- NOVA TABELA DE COTAS
-- =================================================================

CREATE TABLE cotas_de_presentes (
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
CREATE INDEX idx_cotas_id_presente ON cotas_de_presentes(id_presente);
CREATE INDEX idx_cotas_status ON cotas_de_presentes(status);
CREATE INDEX idx_cotas_id_selecao ON cotas_de_presentes(id_selecao);

-- =================================================================
-- ATUALIZAÇÃO DE DADOS EXISTENTES
-- =================================================================

-- Todos os presentes existentes são do tipo INTEGRAL
UPDATE presentes SET tipo = 'INTEGRAL' WHERE tipo IS NULL;