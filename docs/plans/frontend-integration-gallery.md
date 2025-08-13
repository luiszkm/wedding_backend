# 📸 Documentação de Integração Frontend - Módulo de Galeria

## Visão Geral

O módulo de Galeria permite upload, gerenciamento e visualização de fotos do casamento com integração ao AWS S3/R2. Inclui sistema de favoritos, rótulos e interface pública para convidados.

## Endpoints da API

### Base URL
```
http://localhost:3000/v1
```

### 1. 📋 **Listar Fotos Públicas**

**Endpoint:** `GET /casamentos/{idCasamento}/fotos/publico`

**Resposta de Sucesso (200):**
```json
{
  "fotos": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "url": "https://storage.example.com/wedding/foto1.jpg",
      "descricao": "Cerimônia no altar",
      "favorito": true,
      "rotulos": ["cerimonia", "altar"],
      "data_upload": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### 2. 📤 **Upload de Foto (Autenticado)**

**Endpoint:** `POST /casamentos/{idCasamento}/fotos`

**Headers:**
```
Authorization: Bearer {jwt_token}
Content-Type: multipart/form-data
```

**Body (multipart/form-data):**
```
file: [arquivo da imagem]
descricao: "Cerimônia no altar"
```

### 3. ⭐ **Alternar Favorito (Autenticado)**

**Endpoint:** `POST /fotos/{idFoto}/favoritar`

### 4. 🏷️ **Gerenciar Rótulos (Autenticado)**

**Adicionar:** `POST /fotos/{idFoto}/rotulos`
**Remover:** `DELETE /fotos/{idFoto}/rotulos/{rotulo}`

## 🎨 Componente React

### Hook para Galeria

```javascript
// hooks/useGallery.js
import { useState } from 'react';

export const useGallery = (weddingId, token = null) => {
  const [photos, setPhotos] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchPublicPhotos = async () => {
    setLoading(true);
    try {
      const response = await fetch(`/v1/casamentos/${weddingId}/fotos/publico`);
      const data = await response.json();
      setPhotos(data.fotos);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const uploadPhoto = async (file, description) => {
    if (!token) throw new Error('Token necessário');
    
    const formData = new FormData();
    formData.append('file', file);
    formData.append('descricao', description);

    try {
      const response = await fetch(`/v1/casamentos/${weddingId}/fotos`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        },
        body: formData
      });
      
      if (!response.ok) throw new Error('Erro no upload');
      return await response.json();
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  const toggleFavorite = async (photoId) => {
    if (!token) return;
    
    try {
      await fetch(`/v1/fotos/${photoId}/favoritar`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      setPhotos(prev => prev.map(photo => 
        photo.id === photoId 
          ? { ...photo, favorito: !photo.favorito }
          : photo
      ));
    } catch (err) {
      setError(err.message);
    }
  };

  return {
    photos,
    loading,
    error,
    fetchPublicPhotos,
    uploadPhoto,
    toggleFavorite
  };
};
```

### Componente de Galeria Pública

```javascript
// components/PublicGallery.jsx
import React, { useEffect, useState } from 'react';
import { useGallery } from '../hooks/useGallery';

const PublicGallery = ({ weddingId }) => {
  const { photos, loading, error, fetchPublicPhotos } = useGallery(weddingId);
  const [selectedPhoto, setSelectedPhoto] = useState(null);
  const [filter, setFilter] = useState('all');

  useEffect(() => {
    fetchPublicPhotos();
  }, [weddingId]);

  const getUniqueLabels = () => {
    const labels = new Set();
    photos.forEach(photo => {
      photo.rotulos?.forEach(label => labels.add(label));
    });
    return Array.from(labels);
  };

  const filteredPhotos = photos.filter(photo => {
    if (filter === 'all') return true;
    if (filter === 'favorites') return photo.favorito;
    return photo.rotulos?.includes(filter);
  });

  if (loading) return <div className="loading">Carregando galeria...</div>;
  if (error) return <div className="error">Erro: {error}</div>;

  return (
    <div className="public-gallery">
      <h2>Galeria de Fotos</h2>
      
      {/* Filtros */}
      <div className="gallery-filters">
        <button 
          className={filter === 'all' ? 'active' : ''}
          onClick={() => setFilter('all')}
        >
          Todas
        </button>
        <button 
          className={filter === 'favorites' ? 'active' : ''}
          onClick={() => setFilter('favorites')}
        >
          ⭐ Favoritas
        </button>
        {getUniqueLabels().map(label => (
          <button
            key={label}
            className={filter === label ? 'active' : ''}
            onClick={() => setFilter(label)}
          >
            {label}
          </button>
        ))}
      </div>

      {/* Grid de Fotos */}
      <div className="photos-grid">
        {filteredPhotos.map((photo) => (
          <div 
            key={photo.id} 
            className="photo-item"
            onClick={() => setSelectedPhoto(photo)}
          >
            <img src={photo.url} alt={photo.descricao} />
            {photo.favorito && <div className="favorite-badge">⭐</div>}
            <div className="photo-overlay">
              <p>{photo.descricao}</p>
            </div>
          </div>
        ))}
      </div>

      {/* Modal para foto selecionada */}
      {selectedPhoto && (
        <PhotoModal 
          photo={selectedPhoto}
          onClose={() => setSelectedPhoto(null)}
        />
      )}
    </div>
  );
};

const PhotoModal = ({ photo, onClose }) => (
  <div className="modal-overlay" onClick={onClose}>
    <div className="photo-modal" onClick={(e) => e.stopPropagation()}>
      <button className="close-button" onClick={onClose}>✕</button>
      <img src={photo.url} alt={photo.descricao} />
      <div className="photo-info">
        <h3>{photo.descricao}</h3>
        {photo.rotulos && (
          <div className="labels">
            {photo.rotulos.map(label => (
              <span key={label} className="label">{label}</span>
            ))}
          </div>
        )}
      </div>
    </div>
  </div>
);

export default PublicGallery;
```

## 🎨 Estilos CSS

```css
.public-gallery {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.gallery-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 30px;
  flex-wrap: wrap;
}

.gallery-filters button {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s;
}

.gallery-filters button.active {
  background: #2196F3;
  color: white;
  border-color: #2196F3;
}

.photos-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 15px;
}

.photo-item {
  position: relative;
  aspect-ratio: 1;
  overflow: hidden;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.3s;
}

.photo-item:hover {
  transform: scale(1.05);
}

.photo-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.favorite-badge {
  position: absolute;
  top: 10px;
  right: 10px;
  background: rgba(255,255,255,0.9);
  border-radius: 50%;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.photo-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0,0,0,0.7));
  color: white;
  padding: 20px 15px 10px;
  transform: translateY(100%);
  transition: transform 0.3s;
}

.photo-item:hover .photo-overlay {
  transform: translateY(0);
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.photo-modal {
  position: relative;
  max-width: 90vw;
  max-height: 90vh;
  background: white;
  border-radius: 8px;
  overflow: hidden;
}

.photo-modal img {
  width: 100%;
  height: auto;
  display: block;
}

.close-button {
  position: absolute;
  top: 15px;
  right: 15px;
  background: rgba(0,0,0,0.5);
  color: white;
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  cursor: pointer;
  font-size: 18px;
}

.photo-info {
  padding: 20px;
}

.labels {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-top: 10px;
}

.label {
  background: #e0e0e0;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
}

@media (max-width: 768px) {
  .photos-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  }
  
  .photo-modal {
    max-width: 95vw;
    max-height: 95vh;
  }
}
```

## ⚠️ Tratamento de Erros

### Códigos de Status HTTP

| Status | Descrição | Quando Ocorre |
|--------|-----------|---------------|
| 200 | Sucesso | Fotos carregadas com sucesso |
| 201 | Criado | Foto enviada com sucesso |
| 400 | Bad Request | Arquivo inválido ou muito grande |
| 401 | Unauthorized | Token JWT inválido |
| 413 | Payload Too Large | Arquivo excede tamanho máximo |
| 415 | Unsupported Media Type | Tipo de arquivo não suportado |

## 📱 Considerações para UX

### 1. **Performance**
- Lazy loading para imagens
- Thumbnails para grid
- Compressão automática

### 2. **Usabilidade**
- Sistema de filtros intuitivo
- Modal para visualização ampliada
- Suporte a gestos em mobile

### 3. **Acessibilidade**
- Alt text descritivo
- Navegação por teclado
- Contraste adequado

## 🔐 Segurança

### 1. **Upload Seguro**
- Validação de tipo de arquivo
- Limite de tamanho
- Sanitização de nomes

### 2. **Armazenamento**
- URLs assinadas para acesso
- Controle de acesso via JWT
- Backup automático

---

Para testes, consulte `tests/http/gallery.http`.