# URL'ы для настройки Robokassa

## Frontend URL'ы (где находятся страницы)

### 1. Success Page
**Файл**: `/unicorn/app/pages/subscribe/success.vue`
**URL для пользователя**: `https://yourdomain.com/subscribe/success`
**Описание**: Страница успешной оплаты. Пользователь видит это, когда платеж прошел успешно.

### 2. Fail Page  
**Файл**: `/unicorn/app/pages/subscribe/fail.vue`
**URL для пользователя**: `https://yourdomain.com/subscribe/fail`
**Описание**: Страница отмены/ошибки оплаты. Пользователь видит это, если отменил платеж или произошла ошибка.

---

## Backend URL'ы (API endpoints)

### 3. Result URL (обязательный callback от Robokassa)
**Метод**: `POST`
**Backend endpoint**: `/api/subscription/robokassa/result`
**Полный URL**: `https://yourdomain.com/api/subscription/robokassa/result`
**Описание**: Endpoint для обработки callback от Robokassa. Robokassa отправляет подтверждение оплаты сюда.
**Параметры**: OutSum, InvId, SignatureValue, Shp_userId
**Ответ**: OK{InvId} или код ошибки

---

## Как настроить в личном кабинете Robokassa

1. Зайдите в личный кабинет Robokassa
2. Найдите раздел "Настройки магазина" или "Shop settings"
3. Укажите следующие URL'ы:

```
Result URL (обязательный, POST):
https://yourdomain.com/api/subscription/robokassa/result

Success URL (GET):
https://yourdomain.com/subscribe/success

Fail URL (GET):
https://yourdomain.com/subscribe/fail
```

---

## Переменные окружения

В `.env` файл добавьте:

```env
# Robokassa Configuration
ROBOKASSA_MERCHANT_LOGIN=your_merchant_login
ROBOKASSA_PASSWORD1=your_password_1  # Для проверки callback'ов
ROBOKASSA_PASSWORD2=your_password_2  # Для результата
ROBOKASSA_TEST_MODE=true             # false в production
SUBSCRIPTION_PRICE=990.00
SUBSCRIPTION_DURATION_DAYS=30
```

---

## Frontend Configuration

В файле `nuxt.config.ts` используется:

```typescript
runtimeConfig: {
  public: {
    apiBase: process.env.NUXT_PUBLIC_API_BASE || 'https://yourdomain.com/api'
  }
}
```

Убедитесь, что `apiBase` указывает на ваш backend API.

---

## Тестирование локально

Для локальной разработки используйте:

```env
NUXT_PUBLIC_API_BASE=http://localhost:8080/api
```

Robokassa callback'и нельзя получить локально (нет доступа в интернет), но вы можете:
1. Вручную протестировать endpoint POST запросом
2. Использовать тестовый режим Robokassa (ROBOKASSA_TEST_MODE=true)
3. Использовать ngrok для проброса локального сервера в интернет

---

## Процесс платежа

```
1. Пользователь кликает "Подписаться" → SubscriptionBanner.vue
2. useSubscription.goToPayment() создает платеж → POST /api/subscription/create-payment
3. Backend возвращает paymentUrl
4. Пользователь переходит на Robokassa (window.location.href = paymentUrl)
5. Пользователь оплачивает на странице Robokassa
6. Robokassa отправляет callback → POST /api/subscription/robokassa/result (backend)
7. Backend проверяет подпись, активирует подписку
8. Robokassa редиректит пользователя:
   - При успехе → /subscribe/success
   - При ошибке → /subscribe/fail
```

---

## Файлы компонентов

| Файл | Описание |
|------|---------|
| `unicorn/app/composables/useSubscription.ts` | API composable для работы с подписками |
| `unicorn/app/components/SubscriptionBanner.vue` | Баннер с предложением подписки |
| `unicorn/app/pages/subscribe/success.vue` | Страница успеха (Success URL) |
| `unicorn/app/pages/subscribe/fail.vue` | Страница ошибки (Fail URL) |
| `unicorn/app/components/PremiumCard.vue` | Компонент для отображения премиум контента |

---

## Примеры полных URL'ов для production

### Если домен: `unicornstar.online`

```
Success URL:    https://unicornstar.online/subscribe/success
Fail URL:       https://unicornstar.online/subscribe/fail
Result URL:     https://unicornstar.online/api/subscription/robokassa/result
```

### Если домен: `api.example.com`

```
Success URL:    https://api.example.com/subscribe/success
Fail URL:       https://api.example.com/subscribe/fail
Result URL:     https://api.example.com/api/subscription/robokassa/result
```

---

## Важно!

- ✅ **Result URL** - самый критичный, должен быть настроен первым
- ✅ **Success URL** и **Fail URL** - для пользовательского опыта (редирект)
- ✅ Используйте HTTPS для production
- ✅ Проверьте правильность пути `/api/subscription/robokassa/result`
- ✅ Тестируйте на тестовом сервере Robokassa перед production
