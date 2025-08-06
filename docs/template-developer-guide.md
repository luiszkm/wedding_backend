# Guia do Desenvolvedor - Sistema de Templates

Este guia técnico detalha como desenvolvedores podem trabalhar com o sistema de templates híbrido, criar novos templates e integrar funcionalidades personalizadas.

## Ambiente de Desenvolvimento

### Configuração Inicial

```bash
# 1. Aplicar migração do banco
psql -d wedding_db -f db/init/04-add-page-templates.sql

# 2. Criar estrutura de diretórios
mkdir -p templates/{standard,bespoke,partials}

# 3. Copiar templates base
cp templates/partials/* templates/partials/
cp templates/standard/template_moderno.html templates/standard/

# 4. Configurar hot reload (desenvolvimento)
export TEMPLATE_HOT_RELOAD=true
```

### Variáveis de Ambiente

```bash
# .env para desenvolvimento
TEMPLATE_HOT_RELOAD=true          # Recarrega templates automaticamente
TEMPLATE_DEBUG=true               # Logs detalhados de renderização  
TEMPLATE_CACHE_DISABLED=false     # Desabilita cache para debug
TEMPLATE_VALIDATE_STRICT=true     # Validação rigorosa
```

## Anatomia de um Template

### Estrutura Base

```html
<!DOCTYPE html>
<html lang="pt-BR">
{{template "partials/header.html" .}}

<body>
    <!-- SEÇÃO 1: Header do Evento -->
    <header class="event-header">
        <div class="container">
            <h1>{{.Event.Nome}}</h1>
            {{if .Event.Data}}
            <time datetime="{{.Event.Data.Format "2006-01-02"}}">
                {{.Event.Data.Format "02 de January de 2006"}}
            </time>
            {{end}}
        </div>
    </header>

    <!-- SEÇÃO 2: Navegação -->
    {{template "partials/navigation.html" .}}

    <!-- SEÇÃO 3: Conteúdo Principal -->
    <main>
        <!-- Hero Section -->
        <section id="inicio" class="hero">
            <!-- conteúdo hero -->
        </section>

        <!-- Seções Condicionais -->
        {{if .ShowRSVP}}{{template "sections/rsvp.html" .}}{{end}}
        {{if .ShowGifts}}{{template "sections/gifts.html" .}}{{end}}
        {{if .ShowMessages}}{{template "sections/messages.html" .}}{{end}}
        {{if .ShowGallery}}{{template "sections/gallery.html" .}}{{end}}
    </main>

    <!-- SEÇÃO 4: Footer -->
    {{template "partials/footer.html" .}}

    <!-- SEÇÃO 5: Scripts -->
    {{template "partials/scripts.html" .}}
</body>
</html>
```

### Dados Disponíveis

#### Objeto Event
```go
// .Event (eventDomain.Evento)
.Event.ID()        // uuid.UUID - ID do evento
.Event.Nome()      // string - Nome do evento
.Event.Data()      // time.Time - Data do evento
.Event.Tipo()      // TipoEvento - CASAMENTO, ANIVERSARIO, etc.
.Event.UrlSlug()   // string - URL amigável
.Event.IDTemplate() // string - Template padrão atual
.Event.IDTemplateArquivo() // *string - Template bespoke (pode ser nil)
.Event.PaletaCores() // PaletaCores - Cores customizadas
```

#### Coleções de Dados
```go
// .GuestGroups []GrupoDeConvidados  
{{range .GuestGroups}}
    .ID()              // uuid.UUID
    .ChaveDeAcesso()   // string
    .Convidados()      // []Convidado
{{end}}

// .Gifts []Presente
{{range .Gifts}}
    .ID()          // uuid.UUID
    .Nome()        // string
    .Descricao()   // string
    .FotoURL()     // string
    .EhFavorito()  // bool
    .Status()      // StatusPresente (DISPONIVEL/SELECIONADO)
{{end}}

// .Messages []Recado
{{range .Messages}}
    .ID()           // uuid.UUID
    .NomeDoAutor()  // string
    .Texto()        // string
    .Status()       // StatusRecado (APROVADO apenas na versão pública)
    .EhFavorito()   // bool
{{end}}

// .Photos []Foto
{{range .Photos}}
    .ID()           // uuid.UUID
    .UrlPublica()   // string
    .EhFavorito()   // bool
    .Rotulos()      // []Rotulo
{{end}}
```

#### Configurações e Flags
```go
// Flags de exibição (bool)
.ShowGifts      // true se há presentes para exibir
.ShowGallery    // true se há fotos para exibir  
.ShowMessages   // true se há recados aprovados
.ShowRSVP       // true se há grupos de convidados

// Paleta de cores (map[string]string)
.PaletaCores.primary      // Cor principal
.PaletaCores.secondary    // Cor secundária
.PaletaCores.accent       // Cor de destaque
.PaletaCores.background   // Cor de fundo
.PaletaCores.text         // Cor do texto

// Contato (opcional)
{{if .Contact}}
.Contact.Nome      // string
.Contact.Email     // string
.Contact.Telefone  // string
{{end}}

// Dados customizados
.CustomData["key"]  // interface{} - dados extras
```

## Funções Template Disponíveis

### Funções de Formatação

```html
<!-- Formatação de texto -->
{{.Event.Nome | upper}}           <!-- JOÃO & MARIA -->
{{.Event.Nome | lower}}           <!-- joão & maria -->
{{.Event.Nome | title}}           <!-- João & Maria -->

<!-- Truncar texto -->
{{.Message.Texto | truncate 100}} <!-- Corta em 100 caracteres -->

<!-- Formatação de data (personalizada) -->
{{.Event.Data | formatDate "02/01/2006"}}     <!-- 15/06/2024 -->
{{.Event.Data | formatDate "January 2, 2006"}} <!-- June 15, 2024 -->
```

### Funções de Utilidade

```html
<!-- Criar mapa/dict -->
{{$styles := dict "color" .PaletaCores.primary "background" .PaletaCores.secondary}}
<div style="color: {{$styles.color}}; background: {{$styles.background}}">

<!-- Operações condicionais -->
{{if eq .Event.Tipo "CASAMENTO"}}
    <i class="icon-rings"></i>
{{else if eq .Event.Tipo "ANIVERSARIO"}}
    <i class="icon-birthday"></i>
{{end}}

<!-- Loops com index -->
{{range $index, $gift := .Gifts}}
    <div class="gift-item gift-{{$index}}">
        <h3>{{$gift.Nome}}</h3>
    </div>
{{end}}
```

## Criando Templates Personalizados

### Processo de Desenvolvimento

#### 1. Análise de Requisitos
```markdown
- Qual o público-alvo?
- Que funcionalidades são necessárias?
- Há requisitos de design específicos?
- Precisa de responsividade mobile?
- Há restrições de performance?
```

#### 2. Setup do Template
```bash
# Criar arquivo base
touch templates/bespoke/cliente_premium.html

# Copiar estrutura base
cp templates/standard/template_moderno.html templates/bespoke/cliente_premium.html

# Editar conforme necessário
code templates/bespoke/cliente_premium.html
```

#### 3. Desenvolvimento Iterativo

```bash
# Terminal 1: Rodar servidor com hot reload
export TEMPLATE_HOT_RELOAD=true
go run cmd/api/main.go

# Terminal 2: Watch de mudanças
watch "curl -s http://localhost:3000/v1/eventos/teste-evento/pagina > /tmp/output.html"

# Terminal 3: Abrir no browser
open /tmp/output.html
```

### Exemplo Completo: Template Premium

```html
<!DOCTYPE html>
<html lang="pt-BR">
{{template "partials/header.html" .}}

<head>
    <!-- Custom CSS para template premium -->
    <style>
        /* Paleta de cores luxuosa */
        :root {
            --primary-color: {{.PaletaCores.primary | default "#1a1a2e"}};
            --secondary-color: {{.PaletaCores.secondary | default "#16213e"}};
            --accent-color: {{.PaletaCores.accent | default "#e94560"}};
            --background-color: {{.PaletaCores.background | default "#0f0f23"}};
            --text-color: {{.PaletaCores.text | default "#ffffff"}};
            --gold-accent: #d4af37;
        }

        /* Animações premium */
        .fade-in-up {
            animation: fadeInUp 0.8s ease-out;
        }

        @keyframes fadeInUp {
            from {
                opacity: 0;
                transform: translateY(30px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        /* Header premium */
        .premium-header {
            background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            position: relative;
            overflow: hidden;
        }

        .premium-header::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: url('data:image/svg+xml,<svg>...</svg>');
            opacity: 0.1;
            pointer-events: none;
        }

        .premium-title {
            font-family: 'Playfair Display', serif;
            font-size: clamp(2rem, 8vw, 6rem);
            color: var(--text-color);
            text-align: center;
            margin: 0;
            position: relative;
        }

        .premium-subtitle {
            font-size: 1.5rem;
            color: var(--accent-color);
            text-align: center;
            margin-top: 1rem;
            letter-spacing: 2px;
            text-transform: uppercase;
        }

        /* Seções premium */
        .premium-section {
            padding: 5rem 0;
            position: relative;
        }

        .premium-section:nth-child(even) {
            background: rgba(255, 255, 255, 0.02);
        }

        /* Cards premium para presentes */
        .premium-gift-card {
            background: linear-gradient(145deg, 
                rgba(255, 255, 255, 0.1), 
                rgba(255, 255, 255, 0.05)
            );
            backdrop-filter: blur(20px);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 20px;
            padding: 2rem;
            text-align: center;
            transition: all 0.3s ease;
            position: relative;
            overflow: hidden;
        }

        .premium-gift-card:hover {
            transform: translateY(-10px);
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
        }

        .premium-gift-card::before {
            content: '';
            position: absolute;
            top: 0;
            left: -100%;
            width: 100%;
            height: 100%;
            background: linear-gradient(90deg, 
                transparent, 
                rgba(255, 255, 255, 0.1), 
                transparent
            );
            transition: left 0.5s ease;
        }

        .premium-gift-card:hover::before {
            left: 100%;
        }
    </style>

    <!-- Fontes customizadas -->
    <link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;700&display=swap" rel="stylesheet">
</head>

<body>
    <!-- Header Premium -->
    <header class="premium-header">
        <div class="container">
            <div class="fade-in-up">
                <h1 class="premium-title">{{.Event.Nome}}</h1>
                {{if .Event.Data}}
                <p class="premium-subtitle">
                    {{.Event.Data.Format "02 . 01 . 2006"}}
                </p>
                {{end}}
            </div>
        </div>
    </header>

    <!-- Navegação Premium -->
    {{template "partials/navigation.html" .}}

    <main>
        <!-- Hero Premium -->
        <section id="inicio" class="premium-section">
            <div class="container fade-in-up">
                <div class="hero-content">
                    <h2>Uma Celebração Inesquecível</h2>
                    <p>Junte-se a nós neste momento único e especial de nossas vidas.</p>
                    
                    {{if .Event.Data}}
                    <div class="countdown-timer" data-date="{{.Event.Data.Format "2006-01-02T15:04:05Z07:00"}}">
                        <!-- Contador regressivo JavaScript -->
                    </div>
                    {{end}}
                </div>
            </div>
        </section>

        <!-- Lista de Presentes Premium -->
        {{if .ShowGifts}}
        <section id="presentes" class="premium-section">
            <div class="container">
                <h2 class="section-title">Lista de Presentes</h2>
                <div class="gifts-grid">
                    {{range $index, $gift := .Gifts}}
                    <div class="premium-gift-card fade-in-up" style="animation-delay: {{mul $index 0.1}}s">
                        {{if $gift.FotoURL}}
                        <div class="gift-image">
                            <img src="{{$gift.FotoURL}}" alt="{{$gift.Nome}}" loading="lazy">
                        </div>
                        {{end}}
                        
                        <h3>{{$gift.Nome}}</h3>
                        {{if $gift.Descricao}}
                        <p>{{$gift.Descricao | truncate 100}}</p>
                        {{end}}
                        
                        {{if eq $gift.Status "DISPONIVEL"}}
                        <button class="premium-btn" data-gift-id="{{$gift.ID}}">
                            Selecionar Presente
                        </button>
                        {{else}}
                        <span class="gift-selected">
                            <i class="icon-check"></i> Já Selecionado
                        </span>
                        {{end}}
                        
                        {{if $gift.EhFavorito}}
                        <div class="favorite-badge">
                            <i class="icon-heart"></i>
                        </div>
                        {{end}}
                    </div>
                    {{end}}
                </div>
            </div>
        </section>
        {{end}}

        <!-- Recados Premium -->
        {{if .ShowMessages}}
        <section id="recados" class="premium-section">
            <div class="container">
                <h2 class="section-title">Mensagens dos Nossos Queridos</h2>
                <div class="messages-carousel">
                    {{range .Messages}}
                    {{if eq .Status "APROVADO"}}
                    <div class="premium-message-card {{if .EhFavorito}}featured{{end}}">
                        <blockquote>{{.Texto}}</blockquote>
                        <cite>— {{.NomeDoAutor}}</cite>
                        
                        {{if .EhFavorito}}
                        <div class="featured-badge">
                            <i class="icon-star"></i>
                        </div>
                        {{end}}
                    </div>
                    {{end}}
                    {{end}}
                </div>
            </div>
        </section>
        {{end}}

        <!-- Galeria Premium -->
        {{if .ShowGallery}}
        <section id="fotos" class="premium-section">
            <div class="container">
                <h2 class="section-title">Momentos Especiais</h2>
                <div class="premium-gallery">
                    {{range $index, $photo := .Photos}}
                    <div class="gallery-item" style="animation-delay: {{mul $index 0.05}}s">
                        <img src="{{$photo.UrlPublica}}" 
                             alt="Foto do evento" 
                             loading="lazy"
                             onclick="openLightbox({{$index}})">
                        
                        {{if $photo.EhFavorito}}
                        <div class="photo-favorite">
                            <i class="icon-heart-filled"></i>
                        </div>
                        {{end}}
                    </div>
                    {{end}}
                </div>
            </div>
        </section>
        {{end}}

        <!-- RSVP Premium -->
        {{if .ShowRSVP}}
        <section id="confirmacao" class="premium-section">
            <div class="container">
                <h2 class="section-title">Confirme Sua Presença</h2>
                <div class="premium-rsvp-form">
                    <form id="rsvpForm" class="rsvp-form">
                        <div class="form-group">
                            <input type="text" 
                                   id="accessKey" 
                                   placeholder="Sua chave de acesso"
                                   class="premium-input"
                                   required>
                        </div>
                        
                        <button type="submit" class="premium-btn premium-btn-lg">
                            Confirmar Presença
                            <i class="icon-arrow-right"></i>
                        </button>
                    </form>
                    
                    <div id="rsvpResult" class="rsvp-result" style="display: none;">
                        <!-- Resultado será preenchido via JavaScript -->
                    </div>
                </div>
            </div>
        </section>
        {{end}}
    </main>

    <!-- Footer Premium -->
    {{template "partials/footer.html" .}}

    <!-- Scripts Personalizados -->
    <script>
        // Animações de entrada
        const observerOptions = {
            threshold: 0.1,
            rootMargin: '0px 0px -50px 0px'
        };

        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.style.opacity = '1';
                    entry.target.style.transform = 'translateY(0)';
                }
            });
        }, observerOptions);

        // Aplicar observer a elementos com fade-in-up
        document.querySelectorAll('.fade-in-up').forEach(el => {
            el.style.opacity = '0';
            el.style.transform = 'translateY(30px)';
            el.style.transition = 'all 0.8s ease-out';
            observer.observe(el);
        });

        // Contador regressivo
        const countdownEl = document.querySelector('.countdown-timer');
        if (countdownEl) {
            const targetDate = new Date(countdownEl.dataset.date).getTime();
            
            const updateCountdown = () => {
                const now = new Date().getTime();
                const distance = targetDate - now;
                
                const days = Math.floor(distance / (1000 * 60 * 60 * 24));
                const hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
                const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
                
                countdownEl.innerHTML = `
                    <div class="countdown-item">
                        <span class="countdown-number">${days}</span>
                        <span class="countdown-label">Dias</span>
                    </div>
                    <div class="countdown-item">
                        <span class="countdown-number">${hours}</span>
                        <span class="countdown-label">Horas</span>
                    </div>
                    <div class="countdown-item">
                        <span class="countdown-number">${minutes}</span>
                        <span class="countdown-label">Minutos</span>
                    </div>
                `;
                
                if (distance < 0) {
                    countdownEl.innerHTML = '<h3>O grande dia chegou!</h3>';
                }
            };
            
            updateCountdown();
            setInterval(updateCountdown, 60000); // Atualizar a cada minuto
        }

        // Lightbox para galeria
        function openLightbox(index) {
            // Implementar lightbox premium
            console.log('Abrir lightbox para foto', index);
        }

        // RSVP Form
        document.getElementById('rsvpForm')?.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const accessKey = document.getElementById('accessKey').value;
            const resultDiv = document.getElementById('rsvpResult');
            
            try {
                const response = await fetch('/v1/rsvps', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        chave_acesso: accessKey,
                        confirmacoes: [] // Será preenchido após validar chave
                    })
                });
                
                if (response.ok) {
                    resultDiv.innerHTML = '<p class="success">Chave válida! Confirme sua presença.</p>';
                    resultDiv.style.display = 'block';
                } else {
                    resultDiv.innerHTML = '<p class="error">Chave inválida. Verifique e tente novamente.</p>';
                    resultDiv.style.display = 'block';
                }
            } catch (error) {
                resultDiv.innerHTML = '<p class="error">Erro ao processar. Tente novamente.</p>';
                resultDiv.style.display = 'block';
            }
        });
    </script>
</body>
</html>
```

## Debugging e Troubleshooting

### Logs de Template

```bash
# Ativar logs detalhados
export TEMPLATE_DEBUG=true

# Ver logs em tempo real
tail -f /var/log/wedding_backend.log | grep "template"

# Logs específicos
grep "template.*render" /var/log/wedding_backend.log
grep "template.*error" /var/log/wedding_backend.log
```

### Validação de Templates

```bash
# Validar sintaxe de um template específico
go run tools/validate-template.go templates/bespoke/cliente_premium.html

# Validar todos os templates
go run tools/validate-all-templates.go

# Testar renderização com dados mock
go run tools/test-render.go templates/bespoke/cliente_premium.html
```

### Comandos de Debug

```bash
# Limpar cache de templates
curl -X DELETE "http://localhost:3000/v1/debug/template-cache"

# Ver templates em cache
curl "http://localhost:3000/v1/debug/template-cache"

# Forçar reload de template específico  
curl -X POST "http://localhost:3000/v1/debug/template-reload/template_moderno"

# Dump de dados para template
curl "http://localhost:3000/v1/debug/template-data/evento-teste" | jq
```

## Testes Automatizados

### Estrutura de Testes

```go
// internal/pagetemplate/domain/template_integration_test.go
func TestTemplateRendering(t *testing.T) {
    engine := template.NewGoTemplateEngine("../../../templates")
    
    tests := []struct {
        name       string
        templateID string
        data       *domain.EventPageData
        wantError  bool
        contains   []string
    }{
        {
            name:       "template moderno renderiza corretamente",
            templateID: "template_moderno", 
            data:       createMockEventData(),
            wantError:  false,
            contains:   []string{"<!DOCTYPE html>", "Celebre Conosco"},
        },
        // mais casos de teste...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            html, err := engine.RenderEventPage(tt.templateID, tt.data)
            
            if tt.wantError {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            htmlStr := string(html)
            
            for _, expected := range tt.contains {
                assert.Contains(t, htmlStr, expected)
            }
        })
    }
}
```

### Testes de Performance

```go
func BenchmarkTemplateRendering(b *testing.B) {
    engine := template.NewGoTemplateEngine("../../../templates")
    data := createMockEventData()
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _, err := engine.RenderEventPage("template_moderno", data)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Testes E2E

```bash
#!/bin/bash
# scripts/test-templates-e2e.sh

set -e

echo "🧪 Testando templates E2E..."

# 1. Iniciar servidor de teste
go run cmd/api/main.go &
SERVER_PID=$!
sleep 5

# 2. Testar renderização de cada template
TEMPLATES=("template_moderno" "template_classico" "template_elegante")

for template in "${TEMPLATES[@]}"; do
    echo "Testando $template..."
    
    # Atualizar evento para usar o template
    curl -X PUT "http://localhost:3000/v1/eventos/test-event/template" \
         -H "Authorization: Bearer $TEST_JWT" \
         -H "Content-Type: application/json" \
         -d "{\"is_bespoke\": false, \"standard_template_id\": \"$template\"}"
    
    # Renderizar página
    response=$(curl -s "http://localhost:3000/v1/eventos/test-event/pagina")
    
    # Validações básicas
    if [[ $response == *"<!DOCTYPE html>"* ]]; then
        echo "✅ $template renderizou HTML válido"
    else
        echo "❌ $template não renderizou HTML válido"
        exit 1
    fi
    
    if [[ $response == *"var(--primary-color)"* ]]; then
        echo "✅ $template aplicou cores corretamente"
    else
        echo "❌ $template não aplicou cores"
        exit 1
    fi
done

# 3. Limpar
kill $SERVER_PID

echo "✅ Todos os templates passaram nos testes E2E!"
```

## Best Practices

### Performance

1. **Minimize HTML**: Remova espaços desnecessários em produção
2. **CSS Crítico**: Inclua CSS essencial inline, carregar resto assíncrono
3. **Lazy Loading**: Implemente lazy loading para imagens
4. **Compressão**: Use gzip/brotli para reduzir tamanho
5. **Cache**: Aproveite o cache do browser com headers adequados

### Acessibilidade

1. **Semantic HTML**: Use tags semânticas (`<main>`, `<article>`, `<section>`)
2. **Alt Text**: Sempre forneça alt text para imagens
3. **Contraste**: Garanta contraste adequado entre texto e fundo
4. **Keyboard Navigation**: Todos os elementos devem ser navegáveis via teclado
5. **ARIA Labels**: Use ARIA labels para elementos interativos

### SEO

1. **Meta Tags**: Inclua title, description e Open Graph tags
2. **Structured Data**: Adicione JSON-LD para eventos
3. **URLs Amigáveis**: Use slugs descritivos
4. **Sitemap**: Gere sitemap automaticamente
5. **Performance**: Otimize Core Web Vitals

### Código Limpo

1. **Comentários**: Documente seções complexas do template
2. **Consistência**: Mantenha padrões de indentação e nomenclatura
3. **Modularização**: Use partials para código reutilizável
4. **Validação**: Sempre valide dados antes de usar
5. **Error Handling**: Implemente fallbacks para dados ausentes

---

**Última atualização**: Janeiro 2024  
**Versão**: 1.0  
**Maintainer**: Equipe de Desenvolvimento