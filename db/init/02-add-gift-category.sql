-- file: db/init/02-add-gift-category.sql

-- Cria um novo tipo para as categorias de presente
CREATE TYPE categoria_presente AS ENUM (
    'SALA',
    'COZINHA',
    'BANHEIRO',
    'QUARTO',
    'PRODUTO_EXTERNO', -- Mantido para retrocompatibilidade ou uso duplo
    'PIX',             -- Mantido para retrocompatibilidade ou uso duplo
    'OUTROS'
);

-- Adiciona a nova coluna 'categoria' à tabela de presentes
ALTER TABLE presentes ADD COLUMN categoria categoria_presente;