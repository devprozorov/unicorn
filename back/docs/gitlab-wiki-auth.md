# AUTH

> Базовый URL: `http://localhost:8080`

## Авторизация (JWT access) и refresh-cookie

- **Access token** — возвращается в JSON (`accessToken`). Его нужно отправлять в заголовке:
  - `Authorization: Bearer <accessToken>`
- **Refresh token** — хранится в `HttpOnly` cookie `refresh`. Браузер отправляет его автоматически (если CORS `credentials=true`).

### Типы аккаунтов
- `type = "user"` — соискатель
- `type = "company"` — компания

### TOTP
- **Да, TOTP секрет у каждого пользователя свой** (уникальный).
- Включение TOTP: `/api/auth/totp/enroll` → `/api/auth/totp/enable`
- Важно: создавать вакансии/резюме/отклики/чат можно **только** если у аккаунта включён TOTP.

---

## Endpoints

### POST /api/auth/register
Создаёт пользователя и сразу выдаёт access + refresh-cookie.

**Body**
```json
{
  "login": "timur",
  "password": "StrongPass123",
  "displayName": "Timur",
  "type": "user"
}
```

**Response**
```json
{ "ok": true, "accessToken": "..." }
```

---

### POST /api/auth/login
Если у пользователя **нет TOTP** → выдаёт access + refresh-cookie.  
Если **есть TOTP** → возвращает `mfaRequired=true` и `mfaToken` (нужно подтвердить через `/api/auth/totp/verify`).

**Body**
```json
{ "login": "timur", "password": "StrongPass123" }
```

**Response (без TOTP)**
```json
{ "ok": true, "accessToken": "..." }
```

**Response (с TOTP)**
```json
{ "ok": true, "mfaRequired": true, "mfaToken": "..." }
```

---

### POST /api/auth/totp/verify
Второй шаг логина, когда TOTP включён.

**Body**
```json
{ "mfaToken": "...", "code": "123456" }
```

**Response**
```json
{ "ok": true, "accessToken": "..." }
```

---

### POST /api/auth/refresh
Обновляет access token по refresh-cookie и **ротации** refresh.  
Cookie должен быть отправлен автоматически браузером (или через curl с `-b/-c`).

**Response**
```json
{ "ok": true, "accessToken": "..." }
```

---

### POST /api/auth/logout
Ревокает текущую refresh-сессию и очищает cookie.

---

### GET /api/auth/me
Проверка токена, возвращает ник и тип аккаунта.

**Headers**
- `Authorization: Bearer <accessToken>`

**Response**
```json
{ "ok": true, "displayName": "Timur", "type": "user" }
```

---

### POST /api/auth/change-password
**Headers**
- `Authorization: Bearer <accessToken>`

**Body**
```json
{ "oldPassword": "OldPass123", "newPassword": "NewPass12345" }
```

---

### POST /api/auth/totp/enroll
Выдаёт новый секрет TOTP (pending), сохраняет в БД на 10 минут.

**Headers**
- `Authorization: Bearer <accessToken>`

**Response**
```json
{
  "ok": true,
  "secret": "BASE32SECRET",
  "otpauth": "otpauth://totp/...",
  "expiresIn": 600
}
```

---

### POST /api/auth/totp/enable
Подтверждает TOTP кодом и включает MFA.

**Headers**
- `Authorization: Bearer <accessToken>`

**Body**
```json
{ "code": "123456" }
```

---

## Дополнительные модули (кратко)

### Профиль
- `GET /api/profile/:userId`
- `GET /api/profile/me`
- `POST /api/profile/me`
- `PATCH /api/profile/me`

### Поиск компаний
- `GET /api/companies/search?q=...&location=...&industry=...`

### Вакансии (компания + TOTP)
- `GET /api/vacancies`
- `GET /api/vacancies/:id`
- `POST /api/vacancies`
- `PATCH /api/vacancies/:id`
- `DELETE /api/vacancies/:id`

### Резюме (user + TOTP)
- `GET /api/resumes/my`
- `POST /api/resumes`
- `PATCH /api/resumes/:id`
- `DELETE /api/resumes/:id`

### Отклики (TOTP)
- `POST /api/applications` (user)
- `GET /api/applications/inbox` (company)
- `POST /api/applications/:id/accept` (company)
- `POST /api/applications/:id/reject` (company)

### Чат (TOTP, только в рамках applicationId)
- `GET /api/chat/:applicationId/messages`
- `POST /api/chat/:applicationId/messages`

### Админ (отдельная коллекция admins)
- `POST /api/admin/login`
- `GET /api/admin/users?type=user|company`
- `DELETE /api/admin/users/:userId`
- `POST /api/admin/users/:userId/block`
- `POST /api/admin/users/:userId/unblock`
