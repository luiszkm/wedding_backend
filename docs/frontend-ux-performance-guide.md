# Guia de Frontend - UX e Performance para Aplicação de Casamento

## Visão Geral

Este documento define as diretrizes para o desenvolvimento frontend da plataforma de casamento, priorizando experiência do usuário excepcional e performance otimizada. O objetivo é criar uma aplicação rápida, intuitiva e que emocione os usuários durante um dos momentos mais importantes de suas vidas.

## Princípios de Design e UX

### 1. Experiência Centrada no Usuário

#### Personas Principais
- **Noivos**: Precisam de controle total e facilidade de configuração
- **Convidados**: Buscam informações rápidas e claras sobre o evento
- **Família**: Querem participar e contribuir de forma simples

#### Jornada do Usuário

**Noivos (Administradores)**
1. Cadastro e configuração inicial do evento
2. Personalização da página (template, cores, fotos)
3. Gestão de convidados e RSVPs
4. Configuração de lista de presentes
5. Criação do roteiro do evento
6. Publicação de comunicados
7. Monitoramento de confirmações

**Convidados (Visualização Pública)**
1. Acesso via link único do evento
2. Visualização das informações principais
3. Confirmação de presença (RSVP)
4. Seleção de presentes
5. Envio de mensagens aos noivos
6. Consulta do roteiro do evento
7. Upload de fotos na galeria

### 2. Princípios de Interface

#### Simplicidade e Clareza
- Interface limpa com foco no conteúdo essencial
- Máximo 3 ações principais por tela
- Navegação intuitiva com no máximo 3 níveis
- Textos claros e diretos

#### Emocionalidade
- Design que transmite alegria e celebração
- Uso de cores suaves e elegantes
- Tipografia que comunica sofisticação
- Micro-interações que surpreendem positivamente

#### Acessibilidade
- Contraste mínimo 4.5:1 para textos
- Navegação completa por teclado
- Textos alternativos para imagens
- Suporte a leitores de tela

## Arquitetura Frontend Recomendada

### Stack Tecnológico

**Framework Base**
```
Next.js 14+ (App Router)
- Renderização híbrida (SSR/SSG/CSR)
- Otimizações automáticas de performance
- SEO otimizado
- Image optimization nativo
```

**Gerenciamento de Estado**
```
Zustand + TanStack Query
- Estado global leve (Zustand)
- Cache inteligente para dados remotos (TanStack Query)
- Sincronização automática
- Otimistic updates
```

**Estilização**
```
Tailwind CSS + CSS Modules
- Utility-first para desenvolvimento rápido
- CSS Modules para componentes complexos
- Design system consistente
- Bundle size otimizado
```

**Performance**
```
- Bundle analyzer para monitoramento
- Tree shaking automático
- Code splitting por rotas
- Service Workers para cache
```

### Estrutura de Pastas

```
src/
├── app/                    # App Router (Next.js 14)
│   ├── (auth)/            # Grupo de rotas autenticadas
│   ├── (public)/          # Grupo de rotas públicas
│   ├── evento/[slug]/     # Páginas dinâmicas do evento
│   └── layout.tsx
├── components/
│   ├── ui/                # Componentes base do design system
│   ├── forms/             # Componentes de formulários
│   ├── layout/            # Layout components
│   └── features/          # Componentes por feature
├── hooks/                 # Custom hooks
├── lib/                   # Utilitários e configurações
├── stores/               # Estado global (Zustand)
├── types/                # TypeScript types
└── styles/               # Estilos globais
```

## Performance - Métricas Alvo

### Core Web Vitals
- **LCP (Largest Contentful Paint)**: < 1.2s
- **FID (First Input Delay)**: < 100ms
- **CLS (Cumulative Layout Shift)**: < 0.1
- **TTFB (Time to First Byte)**: < 600ms

### Métricas Adicionais
- **First Contentful Paint**: < 1.0s
- **Time to Interactive**: < 2.5s
- **Bundle Size**: < 250KB inicial
- **Lighthouse Score**: 90+ em todas as categorias

## Estratégias de Performance

### 1. Renderização e Carregamento

#### Server-Side Rendering (SSR)
```typescript
// Para páginas públicas do evento
export default async function EventPage({ params }: { params: { slug: string } }) {
  const event = await getEventBySlug(params.slug);
  
  return (
    <EventLayout event={event}>
      <EventContent />
    </EventLayout>
  );
}
```

#### Static Site Generation (SSG)
```typescript
// Para páginas que mudam pouco
export async function generateStaticParams() {
  const events = await getPublishedEvents();
  
  return events.map((event) => ({
    slug: event.slug,
  }));
}
```

#### Client-Side Rendering (CSR)
```typescript
// Para áreas administrativas
'use client'
import { useQuery } from '@tanstack/react-query';

export default function AdminDashboard() {
  const { data, isLoading } = useQuery({
    queryKey: ['dashboard'],
    queryFn: fetchDashboardData
  });
  
  if (isLoading) return <DashboardSkeleton />;
  return <Dashboard data={data} />;
}
```

### 2. Otimização de Assets

#### Imagens
```typescript
import Image from 'next/image';

// Otimização automática com Next.js
<Image
  src="/hero-photo.jpg"
  alt="Foto dos noivos"
  width={800}
  height={600}
  priority // Para imagens above the fold
  placeholder="blur"
  blurDataURL="data:image/jpeg;base64,..."
/>

// Lazy loading para galerias
<Image
  src={photo.url}
  alt={photo.description}
  width={300}
  height={300}
  loading="lazy"
  sizes="(max-width: 768px) 100vw, 50vw"
/>
```

#### Fonts
```typescript
// app/layout.tsx
import { Inter, Playfair_Display } from 'next/font/google';

const inter = Inter({
  subsets: ['latin'],
  display: 'swap',
  variable: '--font-inter'
});

const playfair = Playfair_Display({
  subsets: ['latin'],
  display: 'swap',
  variable: '--font-playfair'
});
```

#### Bundle Splitting
```typescript
// Lazy loading de componentes pesados
import dynamic from 'next/dynamic';

const PhotoGallery = dynamic(
  () => import('./PhotoGallery'),
  {
    loading: () => <GallerySkeleton />,
    ssr: false
  }
);

const MapComponent = dynamic(
  () => import('./MapComponent'),
  { ssr: false }
);
```

### 3. Cache e Estado

#### API Cache com TanStack Query
```typescript
// hooks/useEvent.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

export const useEvent = (eventId: string) => {
  return useQuery({
    queryKey: ['event', eventId],
    queryFn: () => fetchEvent(eventId),
    staleTime: 5 * 60 * 1000, // 5 minutos
    cacheTime: 10 * 60 * 1000, // 10 minutos
  });
};

export const useUpdateEvent = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: updateEvent,
    onSuccess: (data) => {
      queryClient.setQueryData(['event', data.id], data);
      queryClient.invalidateQueries(['events']);
    },
  });
};
```

#### Estado Global Otimizado
```typescript
// stores/eventStore.ts
import { create } from 'zustand';
import { subscribeWithSelector } from 'zustand/middleware';

interface EventState {
  currentEvent: Event | null;
  isEditing: boolean;
  setEvent: (event: Event) => void;
  toggleEdit: () => void;
}

export const useEventStore = create<EventState>()(
  subscribeWithSelector((set) => ({
    currentEvent: null,
    isEditing: false,
    setEvent: (event) => set({ currentEvent: event }),
    toggleEdit: () => set((state) => ({ isEditing: !state.isEditing })),
  }))
);
```

## Design System e Componentes

### 1. Tokens de Design

```typescript
// lib/tokens.ts
export const tokens = {
  colors: {
    primary: {
      50: '#fdf2f8',
      500: '#ec4899',
      900: '#831843',
    },
    neutral: {
      50: '#fafafa',
      500: '#737373',
      900: '#171717',
    }
  },
  spacing: {
    xs: '0.5rem',
    sm: '1rem',
    md: '1.5rem',
    lg: '2rem',
    xl: '3rem',
  },
  typography: {
    h1: {
      fontSize: '2.25rem',
      lineHeight: '2.5rem',
      fontWeight: '700',
    }
  }
} as const;
```

### 2. Componentes Base

#### Button Component
```typescript
// components/ui/Button.tsx
import { cva, type VariantProps } from 'class-variance-authority';

const buttonVariants = cva(
  'inline-flex items-center justify-center rounded-md font-medium transition-colors',
  {
    variants: {
      variant: {
        default: 'bg-primary-500 text-white hover:bg-primary-600',
        outline: 'border border-primary-500 text-primary-500 hover:bg-primary-50',
        ghost: 'text-primary-500 hover:bg-primary-50',
      },
      size: {
        sm: 'h-9 px-3 text-sm',
        md: 'h-10 px-4',
        lg: 'h-11 px-8 text-lg',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'md',
    },
  }
);

interface ButtonProps extends VariantProps<typeof buttonVariants> {
  children: React.ReactNode;
  loading?: boolean;
}

export const Button = ({ variant, size, loading, children, ...props }: ButtonProps) => {
  return (
    <button
      className={buttonVariants({ variant, size })}
      disabled={loading}
      {...props}
    >
      {loading && <Spinner className="mr-2 h-4 w-4" />}
      {children}
    </button>
  );
};
```

#### Form Components
```typescript
// components/ui/Form.tsx
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';

interface FormFieldProps {
  label: string;
  error?: string;
  required?: boolean;
  children: React.ReactNode;
}

export const FormField = ({ label, error, required, children }: FormFieldProps) => {
  return (
    <div className="space-y-1">
      <label className="block text-sm font-medium text-gray-700">
        {label}
        {required && <span className="text-red-500 ml-1">*</span>}
      </label>
      {children}
      {error && <p className="text-sm text-red-600">{error}</p>}
    </div>
  );
};

export const Input = ({ error, ...props }: InputProps) => {
  return (
    <input
      className={cn(
        'block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500',
        error && 'border-red-300 focus:border-red-500 focus:ring-red-500'
      )}
      {...props}
    />
  );
};
```

## Fluxos de UX Otimizados

### 1. Onboarding dos Noivos

#### Wizard de Configuração
```typescript
// components/onboarding/SetupWizard.tsx
const steps = [
  { id: 'basic', title: 'Informações Básicas', component: BasicInfoStep },
  { id: 'design', title: 'Design da Página', component: DesignStep },
  { id: 'features', title: 'Funcionalidades', component: FeaturesStep },
];

export const SetupWizard = () => {
  const [currentStep, setCurrentStep] = useState(0);
  const [formData, setFormData] = useState({});

  const handleNext = (stepData: any) => {
    setFormData(prev => ({ ...prev, ...stepData }));
    setCurrentStep(prev => prev + 1);
  };

  return (
    <div className="max-w-2xl mx-auto">
      <ProgressBar current={currentStep} total={steps.length} />
      
      <AnimatePresence mode="wait">
        <motion.div
          key={currentStep}
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          exit={{ opacity: 0, x: -20 }}
        >
          <steps[currentStep].component
            onNext={handleNext}
            data={formData}
          />
        </motion.div>
      </AnimatePresence>
    </div>
  );
};
```

### 2. RSVP para Convidados

#### Fluxo Simplificado
```typescript
// components/rsvp/RSVPFlow.tsx
export const RSVPFlow = ({ accessKey }: { accessKey: string }) => {
  const { data: guestGroup, isLoading } = useGuestGroup(accessKey);
  const [confirmations, setConfirmations] = useState<Record<string, 'CONFIRMADO' | 'RECUSADO'>>({});
  
  if (isLoading) return <RSVPSkeleton />;
  if (!guestGroup) return <InvalidKeyMessage />;

  return (
    <div className="space-y-6">
      <WelcomeMessage groupName={guestGroup.chaveDeAcesso} />
      
      <div className="space-y-4">
        {guestGroup.convidados.map((guest) => (
          <GuestConfirmation
            key={guest.id}
            guest={guest}
            value={confirmations[guest.id]}
            onChange={(status) => 
              setConfirmations(prev => ({
                ...prev,
                [guest.id]: status
              }))
            }
          />
        ))}
      </div>
      
      <ConfirmButton
        onConfirm={() => submitRSVP(accessKey, confirmations)}
        disabled={Object.keys(confirmations).length === 0}
      />
    </div>
  );
};
```

### 3. Seleção de Presentes

#### Interface Intuitiva
```typescript
// components/gifts/GiftSelection.tsx
export const GiftSelection = ({ eventId }: { eventId: string }) => {
  const { data: gifts, isLoading } = useGifts(eventId);
  const [selectedGifts, setSelectedGifts] = useState<string[]>([]);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {gifts?.map((gift) => (
        <GiftCard
          key={gift.id}
          gift={gift}
          selected={selectedGifts.includes(gift.id)}
          onSelect={(giftId) => {
            setSelectedGifts(prev => 
              prev.includes(giftId)
                ? prev.filter(id => id !== giftId)
                : [...prev, giftId]
            );
          }}
        />
      ))}
      
      <FixedCartSummary
        selectedGifts={selectedGifts}
        onCheckout={() => proceedToCheckout(selectedGifts)}
      />
    </div>
  );
};

const GiftCard = ({ gift, selected, onSelect }: GiftCardProps) => {
  return (
    <motion.div
      whileHover={{ y: -4 }}
      className={cn(
        'relative overflow-hidden rounded-lg border-2 transition-all',
        selected ? 'border-primary-500 shadow-lg' : 'border-gray-200'
      )}
    >
      <Image
        src={gift.imagemUrl}
        alt={gift.nome}
        width={300}
        height={200}
        className="w-full h-48 object-cover"
      />
      
      <div className="p-4">
        <h3 className="font-semibold">{gift.nome}</h3>
        <p className="text-sm text-gray-600 mt-1">{gift.descricao}</p>
        <p className="text-lg font-bold text-primary-500 mt-2">
          R$ {gift.preco.toFixed(2)}
        </p>
        
        <Button
          onClick={() => onSelect(gift.id)}
          variant={selected ? 'default' : 'outline'}
          className="w-full mt-3"
        >
          {selected ? 'Selecionado' : 'Selecionar'}
        </Button>
      </div>
      
      {selected && (
        <div className="absolute top-2 right-2 bg-primary-500 text-white rounded-full p-1">
          <CheckIcon className="w-4 h-4" />
        </div>
      )}
    </motion.div>
  );
};
```

## Estados de Carregamento e Feedback

### 1. Skeleton Screens

```typescript
// components/ui/Skeletons.tsx
export const EventPageSkeleton = () => {
  return (
    <div className="animate-pulse space-y-6">
      {/* Header */}
      <div className="h-64 bg-gray-200 rounded-lg" />
      
      {/* Content */}
      <div className="space-y-4">
        <div className="h-8 bg-gray-200 rounded w-3/4" />
        <div className="h-4 bg-gray-200 rounded w-1/2" />
        <div className="space-y-2">
          <div className="h-4 bg-gray-200 rounded" />
          <div className="h-4 bg-gray-200 rounded w-5/6" />
        </div>
      </div>
    </div>
  );
};

export const GiftCardSkeleton = () => {
  return (
    <div className="animate-pulse space-y-4">
      <div className="h-48 bg-gray-200 rounded-lg" />
      <div className="space-y-2">
        <div className="h-4 bg-gray-200 rounded w-3/4" />
        <div className="h-4 bg-gray-200 rounded w-1/2" />
        <div className="h-6 bg-gray-200 rounded w-1/4" />
      </div>
    </div>
  );
};
```

### 2. Toast Notifications

```typescript
// hooks/useToast.ts
import { toast } from 'sonner';

export const useToast = () => {
  return {
    success: (message: string) => toast.success(message),
    error: (message: string) => toast.error(message),
    info: (message: string) => toast.info(message),
    loading: (message: string) => toast.loading(message),
  };
};

// Uso nos componentes
const { success, error } = useToast();

const handleSave = async () => {
  try {
    await saveEvent(eventData);
    success('Evento salvo com sucesso!');
  } catch (err) {
    error('Erro ao salvar evento');
  }
};
```

## Otimizações Mobile

### 1. Touch-Friendly Interface

```css
/* Botões maiores para touch */
.btn-touch {
  min-height: 44px;
  min-width: 44px;
  padding: 12px 16px;
}

/* Espaçamento adequado entre elementos tocáveis */
.touch-target {
  margin: 8px 0;
}

/* Feedback visual para toques */
.btn:active {
  transform: scale(0.98);
  transition: transform 0.1s;
}
```

### 2. Responsive Design

```typescript
// hooks/useBreakpoint.ts
import { useState, useEffect } from 'react';

export const useBreakpoint = () => {
  const [breakpoint, setBreakpoint] = useState('md');

  useEffect(() => {
    const updateBreakpoint = () => {
      const width = window.innerWidth;
      if (width < 640) setBreakpoint('sm');
      else if (width < 768) setBreakpoint('md');
      else if (width < 1024) setBreakpoint('lg');
      else setBreakpoint('xl');
    };

    updateBreakpoint();
    window.addEventListener('resize', updateBreakpoint);
    return () => window.removeEventListener('resize', updateBreakpoint);
  }, []);

  return breakpoint;
};

// Uso responsivo
const EventCard = () => {
  const breakpoint = useBreakpoint();
  const isMobile = breakpoint === 'sm';

  return (
    <div className={cn(
      'rounded-lg overflow-hidden',
      isMobile ? 'p-4' : 'p-6'
    )}>
      {/* Conteúdo adaptado */}
    </div>
  );
};
```

## SEO e Acessibilidade

### 1. Meta Tags Dinâmicas

```typescript
// app/evento/[slug]/page.tsx
import type { Metadata } from 'next';

export async function generateMetadata({ params }: { params: { slug: string } }): Promise<Metadata> {
  const event = await getEventBySlug(params.slug);
  
  return {
    title: `${event.nomeNoivo} & ${event.nomeNoiva} - ${event.dataDoEvento}`,
    description: `Página do casamento de ${event.nomeNoivo} e ${event.nomeNoiva}. Confirme sua presença e veja todas as informações do evento.`,
    openGraph: {
      title: `Casamento ${event.nomeNoivo} & ${event.nomeNoiva}`,
      description: `${event.nomeNoivo} e ${event.nomeNoiva} se casam em ${event.dataDoEvento}`,
      images: [event.fotoCapaUrl],
    },
    robots: {
      index: true,
      follow: true,
    },
  };
}
```

### 2. Estrutura Semântica

```typescript
// components/EventPage.tsx
export const EventPage = ({ event }: { event: Event }) => {
  return (
    <main role="main" aria-labelledby="event-title">
      <header>
        <h1 id="event-title">
          Casamento {event.nomeNoivo} & {event.nomeNoiva}
        </h1>
        <time dateTime={event.dataDoEvento}>
          {formatDate(event.dataDoEvento)}
        </time>
      </header>
      
      <nav aria-label="Seções do evento" className="sticky top-0 z-10">
        <ul role="list">
          <li><a href="#rsvp">Confirmar Presença</a></li>
          <li><a href="#gifts">Lista de Presentes</a></li>
          <li><a href="#itinerary">Roteiro</a></li>
          <li><a href="#gallery">Galeria</a></li>
        </ul>
      </nav>
      
      <section id="rsvp" aria-labelledby="rsvp-title">
        <h2 id="rsvp-title">Confirmar Presença</h2>
        <RSVPForm />
      </section>
      
      {/* Outras seções */}
    </main>
  );
};
```

## Monitoramento e Analytics

### 1. Performance Monitoring

```typescript
// lib/performance.ts
export const trackPerformance = () => {
  if (typeof window !== 'undefined') {
    // Core Web Vitals
    import('web-vitals').then(({ getCLS, getFID, getFCP, getLCP, getTTFB }) => {
      getCLS(console.log);
      getFID(console.log);
      getFCP(console.log);
      getLCP(console.log);
      getTTFB(console.log);
    });
  }
};

// app/layout.tsx
export default function RootLayout({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    trackPerformance();
  }, []);

  return (
    <html lang="pt-BR">
      <body>{children}</body>
    </html>
  );
}
```

### 2. User Analytics

```typescript
// lib/analytics.ts
interface EventData {
  action: string;
  category: string;
  label?: string;
  value?: number;
}

export const trackEvent = ({ action, category, label, value }: EventData) => {
  if (typeof window !== 'undefined' && window.gtag) {
    window.gtag('event', action, {
      event_category: category,
      event_label: label,
      value: value,
    });
  }
};

// Uso nos componentes
const handleRSVPSubmit = () => {
  trackEvent({
    action: 'rsvp_submitted',
    category: 'engagement',
    label: eventId,
  });
};
```

## Checklist de Performance

### Desenvolvimento
- [ ] Usar Next.js Image component para todas as imagens
- [ ] Implementar lazy loading para componentes pesados
- [ ] Configurar bundle analyzer
- [ ] Usar dynamic imports para código não crítico
- [ ] Implementar Service Workers para cache

### Deploy
- [ ] Configurar CDN para assets estáticos
- [ ] Habilitar compressão Gzip/Brotli
- [ ] Configurar cache headers apropriados
- [ ] Implementar preload para recursos críticos
- [ ] Usar HTTP/2 Push quando apropriado

### Monitoramento
- [ ] Configurar Core Web Vitals tracking
- [ ] Implementar error boundary global
- [ ] Configurar logging de erros
- [ ] Monitorar bundle size ao longo do tempo
- [ ] Implementar alertas para degradação de performance

## Considerações de Negócio

### 1. Conversão
- Formulários simples e diretos
- CTAs claros e visíveis
- Redução de fricção no processo
- Feedback imediato nas ações

### 2. Engajamento
- Micro-interações deliciosas
- Carregamento progressivo
- Estados vazios informativos
- Gamificação sutil (progresso de configuração)

### 3. Retenção
- Performance consistente
- Funcionalidade offline básica
- Sincronização automática
- Backup de dados do usuário

Esta documentação serve como guia para criar uma experiência frontend excepcional, priorizando performance e usabilidade para garantir que os noivos e convidados tenham a melhor experiência possível durante o uso da plataforma.