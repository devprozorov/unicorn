# Скрытие чатов и откликов - Спецификация

## Обзор

Система позволяет пользователям и компаниям скрывать отклики и чаты. Скрытые записи автоматически удаляются через месяц после того, как обе стороны их скрыли.

## Изменения в моделях

### Application Model

Добавлены новые поля:

```go
type Application struct {
    // ... существующие поля ...
    
    // Hidden state: скрыто ли для пользователя или компании
    Hidden   Hidden    `bson:"hidden" json:"-"`
    HiddenAt time.Time `bson:"hiddenAt,omitempty" json:"-"`
}

type Hidden struct {
    User    bool `bson:"user" json:"user"`       // скрыто пользователем
    Company bool `bson:"company" json:"company"` // скрыто компанией
}
```

## API Endpoints

### Скрыть отклик/чат

**Endpoint**: `POST /api/applications/:id/hide`

**Описание**: Скрывает отклик для текущего пользователя (user или company). Чат остается доступным для другой стороны.

**Authentication**: Required ✓

**Path Parameters**:
- `id` (string, required): applicationId отклика

**Поведение**:
- Для пользователей: устанавливает `hidden.user = true`
- Для компаний: устанавливает `hidden.company = true`
- Обновляет поле `hiddenAt` текущим временем

**Response (200 OK)**:
```json
{
  "ok": true
}
```

**Response (500 Error)**:
```json
{
  "ok": false,
  "error": "server_error"
}
```

**Status Codes**:
- `200` - Отклик успешно скрыт
- `401` - Не авторизован
- `500` - Внутренняя ошибка сервера

---

## Логика работы

### 1. Скрытие откликов

- Пользователь или компания могут скрыть отклик независимо друг от друга
- Скрытый отклик больше не отображается в списках:
  - `GET /api/applications/inbox` (для компаний)
  - `GET /api/applications/my` (для пользователей)
  - `GET /api/chats/my`
  - `GET /api/user/chats`

### 2. Автоматическое раскрытие

При отправке нового сообщения в чат:
- Автоматически снимается скрытие для обеих сторон
- `hidden.user` и `hidden.company` устанавливаются в `false`
- Чат снова становится видимым для обеих сторон

### 3. Автоматическое удаление

Система запускает фоновую задачу очистки каждые 24 часа:
- Удаляются записи, где **обе** стороны скрыли отклик
- С момента скрытия прошло более 30 дней
- Условие: `hidden.user = true` AND `hidden.company = true` AND `hiddenAt < now() - 30 days`

## Примеры использования

### JavaScript/Fetch

```javascript
// Скрыть отклик
const hideApplication = async (applicationId, token) => {
  const response = await fetch(
    `http://localhost:8080/api/applications/${applicationId}/hide`,
    {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      }
    }
  );
  
  const data = await response.json();
  if (!data.ok) {
    throw new Error(data.error);
  }
};

// Пример использования
hideApplication('01ARZ3NDEKTSV4RRFFQ69G5FAV', userToken);
```

### Nuxt Composable

```typescript
export const useApplications = () => {
  const config = useRuntimeConfig();
  const baseUrl = config.public.apiUrl;
  const { token } = useAuth();

  const hideApplication = async (applicationId: string) => {
    return await $fetch(
      `${baseUrl}/api/applications/${applicationId}/hide`,
      {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      }
    );
  };

  return {
    hideApplication
  };
};
```

## Миграция данных

Для существующих записей в базе данных поля `hidden` автоматически получат значение по умолчанию:
```json
{
  "hidden": {
    "user": false,
    "company": false
  }
}
```

Никаких дополнительных миграций не требуется - MongoDB автоматически обрабатывает отсутствующие поля.

## Индексы

Рекомендуется добавить следующие индексы для оптимизации:

```javascript
// MongoDB shell
db.applications.createIndex({ "userId": 1, "hidden.user": 1, "createdAt": -1 });
db.applications.createIndex({ "companyId": 1, "hidden.company": 1, "createdAt": -1 });
db.applications.createIndex({ "hidden.user": 1, "hidden.company": 1, "hiddenAt": 1 });
```

## Логирование

Система логирует:
- Количество удаленных записей при очистке
- Ошибки при очистке старых записей
- Ошибки при снятии скрытия (не прерывают основной процесс)

## Сценарии использования

### Сценарий 1: Пользователь скрывает неактуальный отклик
1. Пользователь отправил отклик на вакансию
2. Компания не ответила или отклонила
3. Пользователь скрывает отклик через `POST /api/applications/:id/hide`
4. Отклик исчезает из списка пользователя, но остается у компании

### Сценарий 2: Компания архивирует старые отклики
1. Компания закрыла вакансию
2. Компания скрывает все отклики по этой вакансии
3. Отклики исчезают из inbox компании
4. Если пользователи тоже скроют - через месяц записи удалятся

### Сценарий 3: Возобновление диалога
1. Обе стороны скрыли отклик
2. Компания решила вернуться к кандидату
3. Компания находит отклик через прямую ссылку или другой способ
4. Компания отправляет сообщение
5. Скрытие автоматически снимается для обеих сторон
6. Чат снова появляется в списках

## Технические детали

### Фоновая задача очистки

```go
// Запускается при старте сервера
cleaner := cleanup.NewCleaner(apps)
go cleaner.Start(context.Background(), 24*time.Hour)
```

### Условие удаления

```go
filter := bson.M{
    "$and": []bson.M{
        {"hidden.user": true},
        {"hidden.company": true},
        {"hiddenAt": bson.M{"$lt": oneMonthAgo}},
    },
}
```

### Фильтрация в запросах

Все методы списков автоматически фильтруют скрытые записи:

```go
// Для пользователей
filter := bson.M{"userId": userID, "hidden.user": bson.M{"$ne": true}}

// Для компаний
filter := bson.M{"companyId": companyID, "hidden.company": bson.M{"$ne": true}}
```
