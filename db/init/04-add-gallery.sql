-- =================================================================
-- CONTEXTO: GALERIA DE FOTOS (V3)
-- =================================================================

-- Novo tipo ENUM para os rótulos de fotos 
CREATE TYPE nome_rotulo_enum AS ENUM ('MAIN', 'CASAMENTO', 'LUADEMEL', 'HISTORIA', 'FAMILIA', 'OUTROS');

-- Nova tabela para os metadados das fotos 
CREATE TABLE fotos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_casamento UUID NOT NULL,
    storage_key TEXT NOT NULL,
    url_publica TEXT NOT NULL,
    eh_favorito BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now())
);

-- Nova tabela de junção para a relação muitos-para-muitos entre fotos e rótulos 
CREATE TABLE fotos_rotulos (
    id_foto UUID NOT NULL REFERENCES fotos(id) ON DELETE CASCADE,
    nome_rotulo nome_rotulo_enum NOT NULL,
    PRIMARY KEY (id_foto, nome_rotulo)
);