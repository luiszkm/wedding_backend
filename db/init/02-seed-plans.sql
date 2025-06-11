-- file: db/init/02-seed-plans.sql - VERSÃO CORRIGIDA

-- Limpa a tabela planos e, em cascata, qualquer tabela que dependa dela (como assinaturas)
TRUNCATE TABLE planos RESTART IDENTITY CASCADE;

-- ATENÇÃO: Substitua os valores 'price_...' pelos IDs de Preço reais do seu painel da Stripe.
-- Adicionamos o valor para a coluna 'numero_maximo_eventos' em cada linha.
INSERT INTO planos (id, nome, preco_em_centavos, numero_maximo_eventos, duracao_em_dias, id_stripe_price) VALUES
('a1a1a1a1-1111-1111-1111-111111111111', 'Mensal', 9990, 1, 30, 'price_1RYrGiF0zWP39VVFEhj9TuNd'),
('b2b2b2b2-2222-2222-2222-222222222222', 'Trimestral', 27990, 3, 90, 'price_1RYrGDF0zWP39VVF5VnvKUNa'),
('c3c3c3c3-3333-3333-3333-333333333333', 'Semestral', 53990, 5, 180, 'price_1RYqXGF0zWP39VVFJkJEzZoW');