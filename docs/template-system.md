# Sistema de Templates Híbrido

Esta documentação descreve o sistema de templates híbrido implementado conforme **ADR-002**, que permite personalização flexível das páginas públicas de eventos através de templates padrão e personalizados.

## Visão Geral

O sistema oferece dois modos de renderização:
- **Templates Padrão**: Layouts pré-definidos selecionáveis pelo usuário
- **Templates Bespoke**: Layouts únicos criados pelo desenvolvedor para clientes específicos

## Arquitetura

### Estrutura de Diretórios

```
templates/
├── standard/          # Templates padrão disponíveis para todos
│   ├── template_moderno.html
│   ├── template_classico.html
│   └── template_elegante.html
├── bespoke/           # Templates personalizados para clientes específicos
│   ├── cliente_premium.html
│   └── empresa_xyz.html
└── partials/          # Componentes reutilizáveis
    ├── header.html
    ├── navigation.html
    └── footer.html
```

### Componentes Principais

#### Template Engine
- **Implementação**: `html/template` nativo do Go
- **Features**: Auto-escape XSS, cache em memória, sistema de partials
- **Localização**: `internal/platform/template/engine.go`

#### Domain Models
- **EventPageData**: Estrutura completa dos dados para renderização
- **TemplateMetadata**: Metadados e configurações dos templates
- **PaletaCores**: Sistema de cores customizáveis
- **Localização**: `internal/pagetemplate/domain/`

## API Endpoints

### Públicos

#### Renderizar Página Pública
```http
GET /v1/eventos/{urlSlug}/pagina
```

Renderiza a página pública de um evento usando seu template configurado.

**Parâmetros:**
- `urlSlug`: URL amigável do evento

**Resposta:**
- `Content-Type: text/html`
- HTML completo da página do evento

**Exemplo:**
```bash
curl https://api.exemplo.com/v1/eventos/casamento-joao-maria/pagina
```

#### Listar Templates Disponíveis
```http
GET /v1/templates/disponiveis
```

Retorna lista de todos os templates padrão disponíveis.

**Resposta:**
```json
{
  "templates": [
    {
      "id": "template_moderno",
      "nome": "Moderno",
      "descricao": "Template moderno e minimalista",
      "tipo": "STANDARD",
      "paleta_default": {
        "primary": "#2563eb",
        "secondary": "#f1f5f9",
        "accent": "#10b981",
        "background": "#ffffff",
        "text": "#1f2937"
      },
      "suporta_gifts": true,
      "suporta_gallery": true,
      "suporta_messages": true,
      "suporta_rsvp": true
    }
  ],
  "total": 3
}
```

### Protegidos (Requerem Autenticação)

#### Atualizar Template do Evento
```http
PUT /v1/eventos/{eventId}/template
Authorization: Bearer {jwt-token}
Content-Type: application/json
```

Atualiza a configuração de template de um evento.

**Request Body:**
```json
{
  "is_bespoke": false,
  "standard_template_id": "template_moderno",
  "paleta_cores": {
    "primary": "#2563eb",
    "secondary": "#f1f5f9",
    "accent": "#10b981",
    "background": "#ffffff",
    "text": "#1f2937"
  }
}
```

Para template personalizado:
```json
{
  "is_bespoke": true,
  "bespoke_file_name": "cliente_premium.html",
  "paleta_cores": {
    "primary": "#8b5a3c",
    "secondary": "#f5f5dc",
    "accent": "#d4af37",
    "background": "#fdfdf8",
    "text": "#2c1810"
  }
}
```

**Resposta de Sucesso:**
```json
{
  "message": "Template atualizado com sucesso"
}
```

#### Obter Metadados de Template
```http
GET /v1/templates/{templateId}
```

Retorna informações detalhadas sobre um template específico.

**Exemplo:**
```bash
curl -H "Authorization: Bearer {token}" \
  https://api.exemplo.com/v1/templates/template_moderno
```

## Banco de Dados

### Alterações na Tabela `eventos`

```sql
-- Campos adicionados para suporte a templates
ALTER TABLE eventos ADD COLUMN id_template VARCHAR(100) DEFAULT 'template_moderno';
ALTER TABLE eventos ADD COLUMN id_template_arquivo VARCHAR(100) DEFAULT NULL;
ALTER TABLE eventos ADD COLUMN paleta_cores JSONB DEFAULT '{"primary": "#2563eb", "secondary": "#f1f5f9", "accent": "#10b981", "background": "#ffffff", "text": "#1f2937"}';
```

### Lógica de Precedência

1. Se `id_template_arquivo` estiver preenchido → usa template bespoke
2. Caso contrário → usa template padrão definido em `id_template`
3. Fallback → `template_moderno` se nenhum estiver definido

## Templates Padrão

### Template Moderno (`template_moderno`)

**Características:**
- Design minimalista e clean
- Layout responsivo com CSS Grid/Flexbox
- Paleta de cores azul/cinza
- Suporte completo a todas as funcionalidades

**Paleta Default:**
```json
{
  "primary": "#2563eb",
  "secondary": "#f1f5f9", 
  "accent": "#10b981",
  "background": "#ffffff",
  "text": "#1f2937"
}
```

### Template Clássico (`template_classico`)

**Características:**
- Design tradicional e elegante
- Tipografia clássica com serifa
- Paleta de cores marrom/bege
- Ideal para eventos formais

**Paleta Default:**
```json
{
  "primary": "#8b5a3c",
  "secondary": "#f5f5dc",
  "accent": "#d4af37", 
  "background": "#fdfdf8",
  "text": "#2c1810"
}
```

### Template Elegante (`template_elegante`)

**Características:**
- Design sofisticado e luxuoso
- Elementos escuros com acentos vibrantes
- Paleta dark com vermelho de destaque
- Para eventos premium

**Paleta Default:**
```json
{
  "primary": "#1a1a2e",
  "secondary": "#16213e",
  "accent": "#e94560",
  "background": "#0f0f23", 
  "text": "#ffffff"
}
```

## Sistema de Partials

### Header (`partials/header.html`)
- Meta tags e SEO
- Definição de CSS variables para cores
- Estilos base responsivos

### Navigation (`partials/navigation.html`)  
- Menu principal com scroll suave
- Links condicionais baseados no conteúdo disponível
- Navegação mobile-friendly

### Footer (`partials/footer.html`)
- Informações do evento e contato
- Copyright e branding
- Links sociais (se configurados)

## Personalização de Cores

### Variáveis CSS Disponíveis

```css
:root {
  --primary-color: /* Cor principal */
  --secondary-color: /* Cor secundária */  
  --accent-color: /* Cor de destaque */
  --background-color: /* Cor de fundo */
  --text-color: /* Cor do texto */
}
```

### Como Usar nos Templates

```html
<style>
  .minha-classe {
    background-color: var(--primary-color);
    color: var(--text-color);
    border: 2px solid var(--accent-color);
  }
</style>
```

## Criando Templates Bespoke

### Estrutura Básica

```html
<!DOCTYPE html>
<html lang="pt-BR">
{{template "partials/header.html" .}}

<body>
    <!-- Conteúdo personalizado -->
    <main>
        <h1>{{.Event.Nome}}</h1>
        {{if .Event.Data}}<p>{{.Event.Data.Format "02/01/2006"}}</p>{{end}}
        
        <!-- Seções condicionais -->
        {{if .ShowGifts}}
        <section class="gifts">
            {{range .Gifts}}
            <div class="gift-item">
                <h3>{{.Nome}}</h3>
                <p>{{.Descricao}}</p>
            </div>
            {{end}}
        </section>
        {{end}}
    </main>

    {{template "partials/footer.html" .}}
</body>
</html>
```

### Dados Disponíveis no Template

```go
type EventPageData struct {
    Event        *eventDomain.Evento          // Dados do evento
    GuestGroups  []*guestDomain.GrupoDeConvidados // Grupos de convidados
    Gifts        []*giftDomain.Presente       // Lista de presentes
    Messages     []*mbDomain.Recado           // Recados aprovados
    Photos       []*galleryDomain.Foto        // Fotos da galeria
    PaletaCores  eventDomain.PaletaCores      // Cores personalizadas
    
    // Flags de controle
    ShowGifts    bool
    ShowGallery  bool  
    ShowMessages bool
    ShowRSVP     bool
    
    // Dados extras
    Contact      *ContactInfo
    CustomData   map[string]interface{}
}
```

### Processo de Deploy

1. **Criar template**: Desenvolver arquivo `.html` no padrão
2. **Validar**: Testar renderização e responsividade  
3. **Deploy**: Colocar arquivo em `templates/bespoke/`
4. **Configurar**: Atualizar evento via API para usar o novo template
5. **Testar**: Verificar página pública do evento

## Segurança

### Proteções Implementadas

1. **Auto-escape XSS**: `html/template` escapa automaticamente todo conteúdo injetado
2. **Validação de Input**: Todos os dados são validados antes da renderização
3. **Controle de Acesso**: Apenas desenvolvedores podem criar templates bespoke
4. **Isolamento**: Templates não têm acesso a funções perigosas do sistema

### Boas Práticas

- ✅ **Sempre** usar `{{.Campo}}` para output seguro
- ✅ **Never** usar `{{.Campo | safeHTML}}` sem validação prévia  
- ✅ **Validar** todos os dados customizados antes de injetar
- ✅ **Testar** templates com dados maliciosos

## Performance

### Otimizações Implementadas

1. **Cache de Templates**: Templates compilados ficam em memória
2. **Lazy Loading**: Templates são carregados sob demanda
3. **Compressão**: Respostas HTML são comprimidas via gzip
4. **CSS Inline**: Estilos críticos inline para render rápido

### Cache Headers

```http
# Páginas renderizadas
Cache-Control: public, max-age=300  # 5 minutos

# Lista de templates  
Cache-Control: public, max-age=3600 # 1 hora
```

## Monitoramento e Logs

### Métricas Importantes

- Tempo de renderização por template
- Taxa de cache hit/miss
- Erros de renderização
- Templates mais utilizados

### Logs Estruturados

```json
{
  "level": "info",
  "msg": "template rendered successfully", 
  "event_slug": "casamento-joao-maria",
  "template_id": "template_moderno",
  "render_time_ms": 45,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Troubleshooting

### Problemas Comuns

#### Template não renderiza
```bash
# Verificar se arquivo existe
ls -la templates/standard/template_moderno.html

# Verificar logs do aplicação
docker logs wedding_backend | grep "template"
```

#### Dados não aparecem na página
- Verificar se flags `Show*` estão corretas
- Confirmar se dados existem no banco
- Validar estrutura do template

#### Cores não aplicam
- Verificar sintaxe das variáveis CSS
- Confirmar se paleta está no formato correto
- Testar com paleta padrão

### Comandos Úteis

```bash
# Validar templates
go run cmd/template-validator/main.go

# Testar renderização
curl -v "http://localhost:3000/v1/eventos/meu-evento/pagina"

# Ver cache de templates
curl "http://localhost:3000/v1/debug/template-cache"
```

## Roadmap

### Funcionalidades Futuras

- [ ] Editor visual de templates no admin
- [ ] Sistema de versionamento de templates
- [ ] Prévia em tempo real
- [ ] Temas sazonais automáticos
- [ ] A/B testing de templates
- [ ] Métricas de conversão por template

### Melhorias Técnicas

- [ ] Hot reload de templates em desenvolvimento  
- [ ] Compressão avançada de assets
- [ ] CDN para templates estáticos
- [ ] Renderização server-side com React/Vue
- [ ] Progressive Web App (PWA)

---

**Última atualização**: Janeiro 2024  
**Versão**: 1.0  
**Status**: ✅ Implementado e em produção