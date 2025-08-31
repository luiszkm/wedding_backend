-- file: db/init/08-fix-presentes-status-column.sql
-- Hotfix migration para corrigir a coluna status removida inadvertidamente na migration 06

-- =================================================================
-- PROBLEMA IDENTIFICADO
-- =================================================================
-- A migration 06-add-fractional-gifts.sql executou:
-- DROP TYPE IF EXISTS status_presente CASCADE;
-- 
-- O CASCADE removeu a coluna status da tabela presentes, mas a migration
-- não incluiu o comando para re-adicionar a coluna após recriar o enum.

-- =================================================================
-- SOLUÇÃO
-- =================================================================

-- Re-adicionar a coluna status à tabela presentes
ALTER TABLE presentes 
ADD COLUMN IF NOT EXISTS status status_presente NOT NULL DEFAULT 'DISPONIVEL';

-- Atualizar registros existentes que possam ter valores nulos ou inconsistentes
-- (caso existam registros criados antes desta correção)
UPDATE presentes 
SET status = 'DISPONIVEL' 
WHERE status IS NULL;

-- Recriar índice de performance que pode ter sido perdido
CREATE INDEX IF NOT EXISTS idx_presentes_status ON presentes(status);

-- =================================================================
-- VERIFICAÇÃO
-- =================================================================
-- Para verificar se a correção funcionou, execute:
-- SELECT column_name, data_type, is_nullable, column_default 
-- FROM information_schema.columns 
-- WHERE table_name = 'presentes' AND column_name = 'status';