# API Админ Панели - Полное описание

## Базовая информация
- **Base URL**: `/api/admin`
- **Authentication**: Bearer Token (в заголовке `Authorization: Bearer {accessToken}`)
- **Content-Type**: `application/json`
- **Response Format**: JSON

---

## Аутентификация

### 1. Вход в админ панель
**Endpoint**: `POST /login`

**Описание**: Получить access token для доступа к админ панели. Требует валидные учетные данные админа.

**Request Body**:
```json
{
  "login": "string",
  "password": "string"
}
```

**Parameters**:
- `login` (string, required): Логин администратора
- `password` (string, required): Пароль администратора

**Response (200 OK)**:
```json
{
  "ok": true,
  "accessToken": "string (JWT token, действителен 30 минут)"
}
```

**Response (401 Unauthorized)**:
```json
{
  "ok": false,
  "error": "invalid_credentials"
}
```

**Status Codes**:
- `200` - Успешный вход
- `401` - Неверные учетные данные
- `500` - Внутренняя ошибка сервера

**Headers для последующих запросов**:
```
Authorization: Bearer {accessToken}
```

---

## Управление пользователями

### 2. Получить список пользователей
**Endpoint**: `GET /users`

**Описание**: Получить список пользователей с фильтрацией по типу, поиском и другими параметрами. Возвращает не более 50 пользователей.

**Authentication**: Required ✓

**Query Parameters**:
- `type` (string, optional): Фильтр по типу пользователя
  - `user` - обычные пользователи
  - `company` - компании
  - Если не указан - возвращаются пользователи всех типов
- `search` (string, optional): Поиск по логину или displayName (регистронезависимый)
- `blocked` (string, optional): Фильтр по статусу блокировки
  - `true` - только заблокированные
  - `false` - только незаблокированные
- `deleted` (string, optional): Фильтр по статусу удаления
  - `true` - только удаленные
  - `false` - только неудаленные
- `premium` (string, optional): Фильтр по наличию активной подписки
  - `true` - только с активной подпиской
  - `false` - только без активной подписки

**Response (200 OK)**:
```json
{
  "ok": true,
  "items": [
    {
      "userId": "string (ULID)",
      "login": "string",
      "displayName": "string",
      "type": "user | company",
      "createdAt": "ISO 8601 timestamp (недоступна для клиента, только на бэке)"
    }
  ]
}
```

**Response Example**:
```json
{
  "ok": true,
  "items": [
    {
      "userId": "01ARZ3NDEKTSV4RRFFQ69G5FAV",
      "login": "john_doe",
      "displayName": "John Doe",
      "type": "user"
    },
    {
      "userId": "01ARZ3NDEKTSV4RRFFQ69G5FAW",
      "login": "tech_corp",
      "displayName": "Tech Corporation LLC",
      "type": "company"
    }
  ]
}
```

**Status Codes**:
- `200` - Успешно
- `401` - Не авторизован
- `500` - Внутренняя ошибка сервера

**Примеры запросов**:
- `GET /api/admin/users?type=user` - все обычные пользователи
- `GET /api/admin/users?search=john` - поиск по "john"
- `GET /api/admin/users?blocked=true` - только заблокированные
- `GET /api/admin/users?premium=true&type=user` - пользователи с активной подпиской

**Pagination**: Текущая реализация возвращает фиксированный лимит (50 пользователей). Для реализации пагинации на фронте можно добавить query параметры `limit` и `offset`.

---

### 3. Получить детальную информацию о пользователе
**Endpoint**: `GET /users/:userId`

**Описание**: Получить подробную информацию об отдельном пользователе, включая статус блокировки, подписки и MFA.

**Authentication**: Required ✓

**Path Parameters**:
- `userId` (string, required): ULID идентификатор пользователя

**Response (200 OK)**:
```json
{
  "ok": true,
  "userId": "string (ULID)",
  "login": "string",
  "displayName": "string",
  "type": "user | company",
  "status": {
    "deleted": false,
    "blocked": false
  },
  "subscription": {
    "active": false,
    "until": "ISO 8601 timestamp or zero value"
  },
  "mfa": {
    "totpEnabled": false
  },
  "createdAt": "ISO 8601 timestamp",
  "updatedAt": "ISO 8601 timestamp"
}
```

**Response (404 Not Found)**:
```json
{
  "ok": false,
  "error": "user_not_found"
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
- `200` - Пользователь найден
- `401` - Не авторизован
- `404` - Пользователь не найден
- `500` - Внутренняя ошибка сервера

---

### 4. Редактировать данные пользователя
**Endpoint**: `PATCH /users/:userId`

**Описание**: Обновить данные пользователя (displayName и/или login). Все поля опциональны.

**Authentication**: Required ✓

**Path Parameters**:
- `userId` (string, required): ULID идентификатор пользователя

**Request Body**:
```json
{
  "displayName": "string (optional)",
  "login": "string (optional)"
}
```

**Parameters**:
- `displayName` (string, optional): Новое отображаемое имя пользователя
- `login` (string, optional): Новый логин пользователя (должен быть уникальным)

**Response (200 OK)**:
```json
{
  "ok": true
}
```

**Response (400 Bad Request)**:
```json
{
  "ok": false,
  "error": "login_already_exists"
}
```

или

```json
{
  "ok": false,
  "error": "no_fields_to_update"
}
```

**Response (404 Not Found)**:
```json
{
  "ok": false,
  "error": "user_not_found"
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
- `200` - Пользователь успешно обновлен
- `400` - Некорректные данные (логин занят или нет полей для обновления)
- `401` - Не авторизован
- `404` - Пользователь не найден
- `500` - Внутренняя ошибка сервера

---

### 5. Удалить пользователя (мягкое удаление)
**Endpoint**: `DELETE /users/:userId`

**Описание**: Выполнить soft delete пользователя (пользователь остается в БД, но отмечается как удаленный). Пользователь не сможет логиниться.

**Authentication**: Required ✓

**Path Parameters**:
- `userId` (string, required): ULID идентификатор пользователя

**Request Body**: Не требуется (пустое тело или не отправлять)

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
- `200` - Пользователь успешно удален
- `401` - Не авторизован
- `500` - Внутренняя ошибка сервера

---

### 6. Заблокировать пользователя
**Endpoint**: `POST /users/:userId/block`

**Описание**: Заблокировать пользователя. Заблокированный пользователь не сможет логиниться.

**Authentication**: Required ✓

**Path Parameters**:
- `userId` (string, required): ULID идентификатор пользователя

**Request Body**: Не требуется (пустое тело или не отправлять)

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
- `200` - Пользователь успешно заблокирован
- `401` - Не авторизован
- `500` - Внутренняя ошибка сервера

---

### 7. Разблокировать пользователя
**Endpoint**: `POST /users/:userId/unblock`

**Описание**: Разблокировать пользователя. После разблокировки пользователь сможет логиниться.

**Authentication**: Required ✓

**Path Parameters**:
- `userId` (string, required): ULID идентификатор пользователя

**Request Body**: Не требуется (пустое тело или не отправлять)

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
- `200` - Пользователь успешно разблокирован
- `401` - Не авторизован
- `500` - Внутренняя ошибка сервера

---

## Управление подписками

### 8. Активировать подписку пользователя
**Endpoint**: `POST /users/:userId/subscription/activate`

**Описание**: Активировать премиум-подписку для пользователя на указанное количество дней.

**Authentication**: Required ✓

**Path Parameters**:
- `userId` (string, required): ULID идентификатор пользователя

**Request Body**:
```json
{
  "days": 30
}
```

**Parameters**:
- `days` (integer, required): Количество дней подписки (от 1 до 3650)

**Response (200 OK)**:
```json
{
  "ok": true,
  "until": "ISO 8601 timestamp"
}
```

**Response Example**:
```json
{
  "ok": true,
  "until": "2026-02-09T12:00:00Z"
}
```

**Response (400 Bad Request)**:
```json
{
  "ok": false,
  "error": "invalid_days_value"
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
- `200` - Подписка успешно активирована
- `400` - Некорректное значение дней
- `401` - Не авторизован
- `500` - Внутренняя ошибка сервера

---

### 9. Деактивировать подписку пользователя
**Endpoint**: `POST /users/:userId/subscription/deactivate`

**Описание**: Деактивировать премиум-подписку пользователя.

**Authentication**: Required ✓

**Path Parameters**:
- `userId` (string, required): ULID идентификатор пользователя

**Request Body**: Не требуется (пустое тело или не отправлять)

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
- `200` - Подписка успешно деактивирована
- `401` - Не авторизован
- `500` - Внутренняя ошибка сервера

---

## Типы ошибок

### Общие ошибки
- `invalid_credentials` - Неверный логин или пароль
- `unauthorized` - Отсутствует или неверный токен доступа
- `user_not_found` - Пользователь не найден
- `login_already_exists` - Логин уже используется другим пользователем
- `no_fields_to_update` - Не указаны поля для обновления
- `invalid_days_value` - Некорректное значение количества дней подписки
- `server_error` - Внутренняя ошибка сервера

### HTTP Status Codes
- `200 OK` - Успешный запрос
- `400 Bad Request` - Некорректные данные запроса
- `401 Unauthorized` - Требуется аутентификация или токен недействителен
- `404 Not Found` - Ресурс не найден
- `500 Internal Server Error` - Внутренняя ошибка сервера

---

## Примеры использования (JavaScript/Fetch API)

### Вход в админ панель
```javascript
const login = async (loginCredential, password) => {
  const response = await fetch('http://localhost:8080/api/admin/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      login: loginCredential,
      password: password
    })
  });
  
  const data = await response.json();
  if (data.ok) {
    localStorage.setItem('adminToken', data.accessToken);
    return data.accessToken;
  } else {
    throw new Error(data.error);
  }
};
```

### Получить список пользователей
```javascript
const getUsers = async (type = null) => {
  const token = localStorage.getItem('adminToken');
  let url = 'http://localhost:8080/api/admin/users';
  
  if (type) {
    url += `?type=${type}`;
  }
  
  const response = await fetch(url, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    }
  });
  
  const data = await response.json();
  if (data.ok) {
    return data.items;
  } else {
    throw new Error(data.error);
  }
};

// С фильтрами
const searchUsers = async (filters) => {
  const token = localStorage.getItem('adminToken');
  const params = new URLSearchParams(filters);
  const url = `http://localhost:8080/api/admin/users?${params}`;
  
  const response = await fetch(url, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    }
  });
  
  const data = await response.json();
  if (data.ok) {
    return data.items;
  } else {
    throw new Error(data.error);
  }
};

// Примеры использования фильтров:
// searchUsers({ type: 'user', premium: 'true' })
// searchUsers({ search: 'john' })
// searchUsers({ blocked: 'true' })
```

### Получить информацию о пользователе
```javascript
const getUserDetails = async (userId) => {
  const token = localStorage.getItem('adminToken');
  
  const response = await fetch(
    `http://localhost:8080/api/admin/users/${userId}`,
    {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      }
    }
  );
  
  const data = await response.json();
  if (data.ok) {
    return data;
  } else {
    throw new Error(data.error);
  }
};
```

### Редактировать пользователя
```javascript
const updateUser = async (userId, updates) => {
  const token = localStorage.getItem('adminToken');
  
  const response = await fetch(
    `http://localhost:8080/api/admin/users/${userId}`,
    {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(updates)
    }
  );
  
  const data = await response.json();
  if (!data.ok) {
    throw new Error(data.error);
  }
};

// Примеры:
// updateUser('01ARZ3NDEKTSV4RRFFQ69G5FAV', { displayName: 'New Name' })
// updateUser('01ARZ3NDEKTSV4RRFFQ69G5FAV', { login: 'newlogin' })
// updateUser('01ARZ3NDEKTSV4RRFFQ69G5FAV', { displayName: 'New Name', login: 'newlogin' })
```

### Заблокировать пользователя
```javascript
const blockUser = async (userId) => {
  const token = localStorage.getItem('adminToken');
  
  const response = await fetch(
    `http://localhost:8080/api/admin/users/${userId}/block`,
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
```

### Активировать подписку
```javascript
const activateSubscription = async (userId, days) => {
  const token = localStorage.getItem('adminToken');
  
  const response = await fetch(
    `http://localhost:8080/api/admin/users/${userId}/subscription/activate`,
    {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ days })
    }
  );
  
  const data = await response.json();
  if (data.ok) {
    return data.until;
  } else {
    throw new Error(data.error);
  }
};

// Пример: активировать подписку на 30 дней
// activateSubscription('01ARZ3NDEKTSV4RRFFQ69G5FAV', 30)
```

### Деактивировать подписку
```javascript
const deactivateSubscription = async (userId) => {
  const token = localStorage.getItem('adminToken');
  
  const response = await fetch(
    `http://localhost:8080/api/admin/users/${userId}/subscription/deactivate`,
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
```

---

## Интеграция с Nuxt

### Composable для работы с админ API
```javascript
// composables/useAdminApi.ts
export const useAdminApi = () => {
  const config = useRuntimeConfig();
  const baseUrl = config.public.apiUrl;
  const token = useState('adminToken');

  const login = async (login: string, password: string) => {
    const { data } = await $fetch(`${baseUrl}/api/admin/login`, {
      method: 'POST',
      body: { login, password }
    });
    
    token.value = data.accessToken;
    return data;
  };

  const getUsers = async (filters?: {
    type?: string,
    search?: string,
    blocked?: string,
    deleted?: string,
    premium?: string
  }) => {
    const { data } = await $fetch(`${baseUrl}/api/admin/users`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token.value}`
      },
      query: filters || {}
    });
    
    return data.items;
  };

  const getUserDetails = async (userId: string) => {
    return await $fetch(
      `${baseUrl}/api/admin/users/${userId}`,
      {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      }
    );
  };

  const updateUser = async (userId: string, updates: {
    displayName?: string,
    login?: string
  }) => {
    return await $fetch(
      `${baseUrl}/api/admin/users/${userId}`,
      {
        method: 'PATCH',
        headers: {
          'Authorization': `Bearer ${token.value}`
        },
        body: updates
      }
    );
  };

  const blockUser = async (userId: string) => {
    return await $fetch(
      `${baseUrl}/api/admin/users/${userId}/block`,
      {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      }
    );
  };

  const unblockUser = async (userId: string) => {
    return await $fetch(
      `${baseUrl}/api/admin/users/${userId}/unblock`,
      {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      }
    );
  };

  const deleteUser = async (userId: string) => {
    return await $fetch(
      `${baseUrl}/api/admin/users/${userId}`,
      {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      }
    );
  };

  const activateSubscription = async (userId: string, days: number) => {
    return await $fetch(
      `${baseUrl}/api/admin/users/${userId}/subscription/activate`,
      {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token.value}`
        },
        body: { days }
      }
    );
  };

  const deactivateSubscription = async (userId: string) => {
    return await $fetch(
      `${baseUrl}/api/admin/users/${userId}/subscription/deactivate`,
      {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      }
    );
  };

  return {
    token,
    login,
    getUsers,
    getUserDetails,
    updateUser,
    blockUser,
    unblockUser,
    deleteUser,
    activateSubscription,
    deactivateSubscription
  };
};
```

---

## Заметки для разработки

1. **Token Expiration**: Access token действует 30 минут. Необходимо реализовать логику обновления или переавторизации.

2. **CORS**: API настроен с поддержкой CORS для указанных origin'ов.

3. **Rate Limiting**: На всем API применяется rate limiting (5 запросов в первое окно, 10 во второе).

4. **Валидация**: Все запросы валидируются на бэке. Убедитесь, что отправляемые данные имеют корректный формат.

5. **Безопасность**: Всегда передавайте token через заголовок Authorization, никогда не передавайте его в URL или body.

6. **Status Codes**: Четко обрабатывайте различные HTTP статус коды в приложении (401, 500 и т.д.).

---

## Потенциальные расширения API

На основе существующей архитектуры возможны следующие расширения:
- ~~Получить детальную информацию об отдельном пользователе (GET /users/:userId)~~ ✓ Реализовано
- ~~Редактирование данных пользователя (PATCH /users/:userId)~~ ✓ Реализовано
- ~~Поиск по пользователям (добавить фильтры в GET /users)~~ ✓ Реализовано
- ~~Управление подписками пользователей~~ ✓ Реализовано
- Просмотр логов действий в админ панели
- Статистика по типам пользователей
- Массовые операции над пользователями
- История изменений пользователей
- Расширенные фильтры и сортировки
