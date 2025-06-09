-- file: db/init/02-seed-plans.sql

-- Insere alguns planos de exemplo na tabela 'planos'
-- Os preços estão em centavos (ex: 9990 = R$ 99,90)
INSERT INTO planos (id, nome, preco_em_centavos, numero_maximo_eventos, duracao_em_dias) VALUES
('a1a1a1a1-1111-1111-1111-111111111111', 'Básico', 9990, 1, 365),
('b2b2b2b2-2222-2222-2222-222222222222', 'Premium', 19990, 3, 365),
('c3c3c3c3-3333-3333-3333-333333333333', 'Deluxe', 29990, 5, 730);