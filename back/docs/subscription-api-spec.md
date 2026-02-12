# Subscription API - Документация

## Обзор
Система подписки через Robokassa, которая снимает ограничения на создание вакансий и резюме.

### Преимущества подписки:
- ✅ **16 вакансий/резюме** вместо 2
- ✅ **Premium статус** - контент отображается выше обычного
- ✅ **Цветовая маркировка** (#FFD700 - золотой цвет)
- ✅ **Приоритет в поиске** для работодателей

---

## Эндпоинты

### 1. Получить статус подписки
**GET** `/api/subscription/status`

**Authentication**: Required ✓

**Response (200 OK)**:
```json
{
  "ok": true,
  "active": true,
  "endDate": "2026-02-09T10:30:00Z",
  "daysLeft": 28
}
```

**Response (без активной подписки)**:
```json
{
  "ok": true,
  "active": false
}
```

---

### 2. Создать ссылку на оплату
**POST** `/api/subscription/create-payment`

**Authentication**: Required ✓

**Response (200 OK)**:
```json
{
  "ok": true,
  "paymentUrl": "https://auth.robokassa.ru/Merchant/Index.aspx?MerchantLogin=...",
  "invId": 1704794400,
  "amount": "990.00"
}
```

**Response (503 Service Unavailable)** - если Robokassa не настроена:
```json
{
  "ok": false,
  "error": "payment_disabled"
}
```

---

### 3. Robokassa Result Callback
**POST** `/api/subscription/robokassa/result`

**Description**: Этот endpoint вызывается Robokassa для подтверждения оплаты.

**Form Data**:
- `OutSum` - сумма платежа
- `InvId` - ID инвойса
- `SignatureValue` - подпись для проверки
- `Shp_userId` - ID пользователя

**Response**: `OK{InvId}` или код ошибки

---

### 4. Success URL
**GET** `/api/subscription/robokassa/success`

**Description**: Страница успешной оплаты (перенаправление от Robokassa).

**Query Parameters**:
- `OutSum` - сумма
- `InvId` - ID инвойса
- `SignatureValue` - подпись
- `Shp_userId` - ID пользователя

**Response**:
```json
{
  "ok": true,
  "message": "Оплата успешно завершена!",
  "invId": 1704794400
}
```

---

### 5. Fail URL
**GET** `/api/subscription/robokassa/fail`

**Description**: Страница отмены/ошибки оплаты.

**Response**:
```json
{
  "ok": false,
  "message": "Оплата отменена или произошла ошибка",
  "invId": "1704794400"
}
```

---

## Изменения в моделях

### Vacancy & Resume
Добавлены новые поля:
```json
{
  "isPremium": true,
  "colorCode": "#FFD700"
}
```

### Лимиты
- **Без подписки**: 2 вакансии/резюме
- **С подпиской**: 16 вакансий/резюме

### Сортировка
Премиум вакансии/резюме отображаются **выше** обычных в списках.

---

## Настройка в .env

```env
# Robokassa Settings
ROBOKASSA_MERCHANT_LOGIN=your_merchant_login
ROBOKASSA_PASSWORD1=your_password_1
ROBOKASSA_PASSWORD2=your_password_2
ROBOKASSA_TEST_MODE=true
SUBSCRIPTION_PRICE=990.00
SUBSCRIPTION_DURATION_DAYS=30
```

---

## Интеграция с фронтом (Nuxt)

### Composable для подписки

```typescript
// composables/useSubscription.ts
export const useSubscription = () => {
  const config = useRuntimeConfig();
  const baseUrl = config.public.apiUrl;
  const token = useState('accessToken');

  const getStatus = async () => {
    const data = await $fetch(`${baseUrl}/api/subscription/status`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token.value}`
      }
    });
    return data;
  };

  const createPayment = async () => {
    const data = await $fetch(`${baseUrl}/api/subscription/create-payment`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token.value}`
      }
    });
    
    // Перенаправляем на страницу оплаты
    if (data.ok && data.paymentUrl) {
      window.location.href = data.paymentUrl;
    }
    
    return data;
  };

  return {
    getStatus,
    createPayment
  };
};
```

### Пример использования на странице

```vue
<template>
  <div>
    <div v-if="!status?.active" class="subscription-banner">
      <h2>Получите Premium подписку!</h2>
      <ul>
        <li>✅ До 16 вакансий/резюме</li>
        <li>✅ Приоритетное отображение</li>
        <li>✅ Золотая маркировка</li>
      </ul>
      <button @click="subscribe">Подписаться за 990₽</button>
    </div>
    
    <div v-else class="active-subscription">
      <p>✅ Premium активна до {{ formatDate(status.endDate) }}</p>
      <p>Осталось дней: {{ status.daysLeft }}</p>
    </div>
  </div>
</template>

<script setup>
const { getStatus, createPayment } = useSubscription();
const status = ref(null);

onMounted(async () => {
  status.value = await getStatus();
});

const subscribe = async () => {
  await createPayment();
};
</script>
```

### Отображение премиум вакансий

```vue
<template>
  <div 
    class="vacancy-card"
    :class="{ 'premium': vacancy.isPremium }"
    :style="{ borderLeftColor: vacancy.colorCode }"
  >
    <div v-if="vacancy.isPremium" class="premium-badge">
      ⭐ PREMIUM
    </div>
    <h3>{{ vacancy.title }}</h3>
    <p>{{ vacancy.description }}</p>
  </div>
</template>

<style scoped>
.vacancy-card.premium {
  border-left: 4px solid var(--color);
  background: linear-gradient(to right, #fff9e6, white);
}

.premium-badge {
  color: #FFD700;
  font-weight: bold;
  font-size: 12px;
}
</style>
```

---

## Robokassa URLs для настройки

В личном кабинете Robokassa укажите:

**Result URL** (обязательный):
```
https://yourdomain.com/api/subscription/robokassa/result
```

**Success URL**:
```
https://yourdomain.com/api/subscription/robokassa/success
```

**Fail URL**:
```
https://yourdomain.com/api/subscription/robokassa/fail
```

**Метод**: `POST` для Result URL, `GET` для Success/Fail

---

## Безопасность

1. ✅ Проверка подписи от Robokassa (MD5)
2. ✅ Валидация InvID перед обработкой
3. ✅ Проверка существования подписки в БД
4. ✅ Атомарное обновление пользователя и контента
5. ✅ Логирование всех операций

---

## База данных

### Коллекция `subscriptions`

```javascript
{
  "_id": ObjectId("..."),
  "subscriptionId": "01ARZ3NDEKTSV4RRFFQ69G5FAV",
  "userId": "01ARZ3NDEKTSV4RRFFQ69G5FAW",
  "amount": 990.00,
  "currency": "RUB",
  "status": "paid", // pending/paid/cancelled
  "invId": 1704794400,
  "outSum": "990.00",
  "startDate": ISODate("2026-01-09T10:30:00Z"),
  "endDate": ISODate("2026-02-09T10:30:00Z"),
  "createdAt": ISODate("2026-01-09T10:25:00Z"),
  "updatedAt": ISODate("2026-01-09T10:30:00Z")
}
```

### Обновление User при активации подписки

```javascript
{
  "subscription": {
    "active": true,
    "until": ISODate("2026-02-09T10:30:00Z")
  }
}
```

### Обновление Vacancy/Resume

```javascript
{
  "isPremium": true,
  "colorCode": "#FFD700"
}
```

---

## Тестирование

### Тестовый режим Robokassa

При `ROBOKASSA_TEST_MODE=true` все платежи проходят в тестовом режиме.

**Тестовые данные**:
- Любая карта: успешный платеж
- Специальные карты Robokassa для проверки ошибок

### Ручное тестирование

1. Создайте пользователя
2. Проверьте лимит (должно быть 2)
3. Создайте платеж через API
4. Перейдите по ссылке оплаты
5. Оплатите (тестовая карта)
6. Проверьте статус подписки
7. Убедитесь, что лимит = 16
8. Создайте вакансию/резюме - должен быть isPremium: true

---

## Возможные ошибки

| Код | Ошибка | Причина |
|-----|--------|---------|
| 401 | unauthorized | Нет токена доступа |
| 403 | limit_reached | Достигнут лимит вакансий/резюме |
| 503 | payment_disabled | Robokassa не настроена |
| 400 | invalid_signature | Неверная подпись от Robokassa |
| 404 | subscription not found | Подписка не найдена по InvID |

---

## Мониторинг

Все callback'и от Robokassa логируются:
- ✅ Успешные платежи
- ❌ Ошибки проверки подписи
- ⚠️ Не найденные подписки
- 🔄 Обновления статуса пользователя
