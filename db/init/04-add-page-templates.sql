-- file: db/init/04-add-page-templates.sql
-- Migração para adicionar suporte a templates híbridos conforme ADR-002

-- Adiciona campos de template à tabela eventos
ALTER TABLE eventos ADD COLUMN id_template VARCHAR(100) DEFAULT 'template_moderno';
ALTER TABLE eventos ADD COLUMN id_template_arquivo VARCHAR(100) DEFAULT NULL;
ALTER TABLE eventos ADD COLUMN paleta_cores JSONB DEFAULT '{"primary": "#2563eb", "secondary": "#f1f5f9", "accent": "#10b981", "background": "#ffffff", "text": "#1f2937"}';

-- Comentários para documentação
COMMENT ON COLUMN eventos.id_template IS 'ID do template padrão (template_moderno, template_classico, template_elegante)';
COMMENT ON COLUMN eventos.id_template_arquivo IS 'Nome do arquivo de template personalizado (bespoke). Quando preenchido, tem precedência sobre id_template';
COMMENT ON COLUMN eventos.paleta_cores IS 'Paleta de cores em formato JSON para personalização visual do template';

-- Índice para otimizar buscas por template
CREATE INDEX idx_eventos_id_template_arquivo ON eventos(id_template_arquivo) WHERE id_template_arquivo IS NOT NULL;