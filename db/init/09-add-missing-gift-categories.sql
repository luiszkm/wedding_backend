-- file: db/init/09-add-missing-gift-categories.sql
-- Adicionar categorias de presente que estão sendo usadas nos testes e documentação

-- =================================================================
-- PROBLEMA IDENTIFICADO
-- =================================================================
-- Os testes e a aplicação estão usando categorias como "COZINHA" e "SALA"
-- mas o enum nome_rotulo_enum só tem: MAIN, CASAMENTO, LUADEMEL, HISTORIA, FAMILIA, OUTROS
-- 
-- Erro: invalid input value for enum nome_rotulo_enum: "COZINHA"

-- =================================================================
-- SOLUÇÃO
-- =================================================================

-- Adicionar as categorias de presente que faltam no enum
ALTER TYPE nome_rotulo_enum ADD VALUE IF NOT EXISTS 'COZINHA';
ALTER TYPE nome_rotulo_enum ADD VALUE IF NOT EXISTS 'SALA';
ALTER TYPE nome_rotulo_enum ADD VALUE IF NOT EXISTS 'QUARTO';
ALTER TYPE nome_rotulo_enum ADD VALUE IF NOT EXISTS 'BANHEIRO';
ALTER TYPE nome_rotulo_enum ADD VALUE IF NOT EXISTS 'JARDIM';
ALTER TYPE nome_rotulo_enum ADD VALUE IF NOT EXISTS 'DECORACAO';
ALTER TYPE nome_rotulo_enum ADD VALUE IF NOT EXISTS 'ELETRONICOS';
ALTER TYPE nome_rotulo_enum ADD VALUE IF NOT EXISTS 'UTENSILIOS';

-- =================================================================
-- VERIFICAÇÃO
-- =================================================================
-- Para verificar as categorias disponíveis:
-- SELECT unnest(enum_range(NULL::nome_rotulo_enum));