# 💳 Documentação de Integração Frontend - Módulo de Cobrança

## Visão Geral

O módulo de Cobrança integra com Stripe para gestão de assinaturas e pagamentos. Inclui planos de assinatura, processamento de pagamentos e gerenciamento de assinaturas ativas.

## Endpoints da API

### Base URL
```
http://localhost:3000/v1
```

### 1. 📋 **Listar Planos**

**Endpoint:** `GET /planos`

**Resposta:**
```json
{
  "planos": [
    {
      "id": "plan_basic",
      "nome": "Básico",
      "descricao": "Plano básico para casamentos pequenos",
      "preco": 99.99,
      "moeda": "BRL",
      "intervalo": "month",
      "features": [
        "100 convidados",
        "Lista de presentes",
        "Galeria básica",
        "1 template"
      ],
      "stripe_price_id": "price_1234567890"
    }
  ]
}
```

### 2. 💳 **Criar Assinatura**

**Endpoint:** `POST /assinaturas`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer {jwt_token}
```

**Body:**
```json
{
  "id_plano": "plan_basic",
  "payment_method_id": "pm_1234567890abcdef",
  "dados_cobranca": {
    "nome": "João Silva",
    "email": "joao@exemplo.com",
    "endereco": {
      "linha1": "Rua das Flores, 123",
      "cidade": "São Paulo",
      "estado": "SP",
      "cep": "01234-567",
      "pais": "BR"
    }
  }
}
```

**Resposta de Sucesso (201):**
```json
{
  "id": "sub_1234567890",
  "status": "active",
  "client_secret": "pi_1234_secret_abcd",
  "proxima_cobranca": "2024-02-15T00:00:00Z"
}
```

### 3. 📄 **Obter Assinatura Ativa**

**Endpoint:** `GET /usuarios/{idUsuario}/assinatura`

**Resposta:**
```json
{
  "id": "sub_1234567890",
  "plano": {
    "nome": "Básico",
    "preco": 99.99
  },
  "status": "active",
  "data_inicio": "2024-01-15T00:00:00Z",
  "proxima_cobranca": "2024-02-15T00:00:00Z",
  "cancelar_no_fim_periodo": false
}
```

### 4. ❌ **Cancelar Assinatura**

**Endpoint:** `DELETE /assinaturas/{idAssinatura}`

## 🎨 Componentes React

### Hook para Billing

```javascript
// hooks/useBilling.js
import { useState, useEffect } from 'react';

export const useBilling = (token) => {
  const [plans, setPlans] = useState([]);
  const [subscription, setSubscription] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchPlans = async () => {
    try {
      const response = await fetch('/v1/planos');
      const data = await response.json();
      setPlans(data.planos);
    } catch (err) {
      setError(err.message);
    }
  };

  const createSubscription = async (subscriptionData) => {
    if (!token) throw new Error('Token necessário');
    
    setLoading(true);
    try {
      const response = await fetch('/v1/assinaturas', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(subscriptionData)
      });
      
      if (!response.ok) throw new Error('Erro ao criar assinatura');
      return await response.json();
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const fetchUserSubscription = async (userId) => {
    if (!token) return;
    
    try {
      const response = await fetch(`/v1/usuarios/${userId}/assinatura`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (response.ok) {
        const data = await response.json();
        setSubscription(data);
      }
    } catch (err) {
      // Usuário pode não ter assinatura ativa
      setSubscription(null);
    }
  };

  useEffect(() => {
    fetchPlans();
  }, []);

  return {
    plans,
    subscription,
    loading,
    error,
    createSubscription,
    fetchUserSubscription
  };
};
```

### Componente de Seleção de Planos

```javascript
// components/PlanSelector.jsx
import React, { useEffect, useState } from 'react';
import { loadStripe } from '@stripe/stripe-js';
import { Elements, CardElement, useStripe, useElements } from '@stripe/react-stripe-js';
import { useBilling } from '../hooks/useBilling';
import { useAuth } from '../hooks/useAuth';

const stripePromise = loadStripe(process.env.REACT_APP_STRIPE_PUBLISHABLE_KEY);

const PlanSelector = () => {
  const { token, user } = useAuth();
  const { plans, subscription, loading, fetchUserSubscription } = useBilling(token);
  const [selectedPlan, setSelectedPlan] = useState(null);
  const [showPayment, setShowPayment] = useState(false);

  useEffect(() => {
    if (user?.id) {
      fetchUserSubscription(user.id);
    }
  }, [user]);

  const formatPrice = (price) => {
    return new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'BRL'
    }).format(price);
  };

  if (subscription) {
    return (
      <div className="current-subscription">
        <h2>Sua Assinatura Ativa</h2>
        <div className="subscription-card">
          <h3>{subscription.plano.nome}</h3>
          <p className="price">{formatPrice(subscription.plano.preco)}/mês</p>
          <p>Status: <span className="status active">{subscription.status}</span></p>
          <p>Próxima cobrança: {new Date(subscription.proxima_cobranca).toLocaleDateString('pt-BR')}</p>
          
          <button className="cancel-button">
            Cancelar Assinatura
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="plan-selector">
      <h2>Escolha seu Plano</h2>
      
      <div className="plans-grid">
        {plans.map((plan) => (
          <div 
            key={plan.id}
            className={`plan-card ${selectedPlan?.id === plan.id ? 'selected' : ''}`}
            onClick={() => setSelectedPlan(plan)}
          >
            <h3>{plan.nome}</h3>
            <div className="price">
              {formatPrice(plan.preco)}
              <span className="period">/{plan.intervalo === 'month' ? 'mês' : 'ano'}</span>
            </div>
            <p className="description">{plan.descricao}</p>
            
            <ul className="features">
              {plan.features.map((feature, index) => (
                <li key={index}>✓ {feature}</li>
              ))}
            </ul>
            
            <button 
              className="select-button"
              onClick={(e) => {
                e.stopPropagation();
                setSelectedPlan(plan);
                setShowPayment(true);
              }}
            >
              Selecionar Plano
            </button>
          </div>
        ))}
      </div>
      
      {showPayment && selectedPlan && (
        <Elements stripe={stripePromise}>
          <PaymentForm 
            plan={selectedPlan}
            onSuccess={() => {
              setShowPayment(false);
              fetchUserSubscription(user.id);
            }}
            onCancel={() => setShowPayment(false)}
          />
        </Elements>
      )}
    </div>
  );
};

const PaymentForm = ({ plan, onSuccess, onCancel }) => {
  const stripe = useStripe();
  const elements = useElements();
  const { createSubscription, loading } = useBilling();
  const [formData, setFormData] = useState({
    nome: '',
    email: '',
    endereco: {
      linha1: '',
      cidade: '',
      estado: '',
      cep: ''
    }
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!stripe || !elements) return;
    
    const cardElement = elements.getElement(CardElement);
    
    try {
      // Criar Payment Method
      const { error, paymentMethod } = await stripe.createPaymentMethod({
        type: 'card',
        card: cardElement,
        billing_details: {
          name: formData.nome,
          email: formData.email,
          address: {
            line1: formData.endereco.linha1,
            city: formData.endereco.cidade,
            state: formData.endereco.estado,
            postal_code: formData.endereco.cep,
            country: 'BR'
          }
        }
      });
      
      if (error) {
        console.error('Erro no pagamento:', error);
        return;
      }
      
      // Criar assinatura
      const result = await createSubscription({
        id_plano: plan.id,
        payment_method_id: paymentMethod.id,
        dados_cobranca: formData
      });
      
      // Se precisar de confirmação 3D Secure
      if (result.client_secret) {
        const { error: confirmError } = await stripe.confirmCardPayment(result.client_secret);
        if (confirmError) {
          console.error('Erro na confirmação:', confirmError);
          return;
        }
      }
      
      onSuccess();
    } catch (error) {
      console.error('Erro ao processar pagamento:', error);
    }
  };

  return (
    <div className="payment-modal">
      <div className="payment-content">
        <h3>Finalizar Assinatura - {plan.nome}</h3>
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Nome Completo:</label>
            <input
              type="text"
              value={formData.nome}
              onChange={(e) => setFormData({...formData, nome: e.target.value})}
              required
            />
          </div>
          
          <div className="form-group">
            <label>Email:</label>
            <input
              type="email"
              value={formData.email}
              onChange={(e) => setFormData({...formData, email: e.target.value})}
              required
            />
          </div>
          
          <div className="form-group">
            <label>Cartão de Crédito:</label>
            <div className="card-element">
              <CardElement
                options={{
                  style: {
                    base: {
                      fontSize: '16px',
                      color: '#424770',
                      '::placeholder': {
                        color: '#aab7c4',
                      },
                    },
                  },
                }}
              />
            </div>
          </div>
          
          <div className="form-actions">
            <button type="button" onClick={onCancel}>
              Cancelar
            </button>
            <button type="submit" disabled={!stripe || loading}>
              {loading ? 'Processando...' : `Assinar por ${new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(plan.preco)}/mês`}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default PlanSelector;
```

## 🎨 Estilos CSS

```css
.plan-selector {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.plans-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.plan-card {
  border: 2px solid #eee;
  border-radius: 8px;
  padding: 30px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  background: white;
}

.plan-card:hover,
.plan-card.selected {
  border-color: #2196F3;
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.2);
}

.plan-card h3 {
  margin: 0 0 15px 0;
  color: #333;
  font-size: 24px;
}

.price {
  font-size: 36px;
  font-weight: bold;
  color: #2196F3;
  margin-bottom: 10px;
}

.period {
  font-size: 16px;
  color: #666;
}

.description {
  color: #666;
  margin-bottom: 20px;
}

.features {
  list-style: none;
  padding: 0;
  margin: 20px 0;
  text-align: left;
}

.features li {
  padding: 5px 0;
  color: #555;
}

.select-button {
  width: 100%;
  background: #2196F3;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.payment-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.payment-content {
  background: white;
  padding: 30px;
  border-radius: 8px;
  max-width: 500px;
  width: 90%;
}

.card-element {
  border: 1px solid #ddd;
  padding: 12px;
  border-radius: 4px;
  background: white;
}

.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 20px;
}

.form-actions button {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.current-subscription {
  max-width: 500px;
  margin: 0 auto;
  padding: 20px;
}

.subscription-card {
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 20px;
  text-align: center;
}

.status.active {
  color: #28a745;
  font-weight: bold;
}

.cancel-button {
  background: #dc3545;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 15px;
}
```

## 📱 Considerações para UX

### 1. **Segurança**
- Integração oficial com Stripe
- PCI compliance automático
- 3D Secure quando necessário

### 2. **Usabilidade**
- Comparação clara entre planos
- Processo de checkout simplificado
- Feedback visual durante pagamento

### 3. **Gestão**
- Status da assinatura em tempo real
- Facilidade para cancelamento
- Histórico de pagamentos

---

Para testes, consulte `tests/http/billing.http` e configure as variáveis do Stripe.