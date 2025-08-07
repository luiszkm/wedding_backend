-- Migration: Add itinerary (roteiro) functionality
-- ADR-005: Criação do Contexto de Roteiro para o Módulo de Itinerário do Evento

CREATE TABLE itens_roteiro (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_evento UUID NOT NULL,
    horario TIMESTAMP WITH TIME ZONE NOT NULL,
    titulo_atividade VARCHAR(255) NOT NULL,
    descricao_atividade TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Garante a relação com o evento ao qual o roteiro pertence
    CONSTRAINT fk_evento_roteiro
        FOREIGN KEY(id_evento) 
        REFERENCES eventos(id)
        ON DELETE CASCADE
);

-- Index para otimizar consultas por evento
CREATE INDEX idx_itens_roteiro_id_evento ON itens_roteiro(id_evento);

-- Index para otimizar ordenação por horário
CREATE INDEX idx_itens_roteiro_horario ON itens_roteiro(horario);

-- Index composto para consultas eficientes de itens de um evento ordenados por horário
CREATE INDEX idx_itens_roteiro_evento_horario ON itens_roteiro(id_evento, horario);