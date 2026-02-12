# unicorn-auth (Gin + MongoDB)

Безопасный auth-сервис:
- register / login
- TOTP (2FA) enroll + enable + verify
- refresh (rotation) / logout
- change-password (при включенном MFA требует TOTP)
- защита от NoSQL injection (строгие DTO + DisallowUnknownFields + ручные фильтры)
- rate limit (in-memory; для прод лучше Redis)
- refresh токен в HttpOnly Secure cookie

## Запуск

```bash
cp .env.example .env
# заполни JWT_HS256_SECRET и TOTP_ENC_KEY_B64

go mod tidy
go run ./cmd/server
```

## Быстрый тест (curl)

Register:
```bash
curl -s -X POST http://localhost:8080/api/auth/register   -H 'Content-Type: application/json'   -d '{"login":"timur","password":"StrongPass123!@#","displayName":"Timur","type":"user"}'
```

Login:
```bash
curl -i -s -X POST http://localhost:8080/api/auth/login   -H 'Content-Type: application/json'   -d '{"login":"timur","password":"StrongPass123!@#"}'
```

Refresh:
```bash
curl -i -s -X POST http://localhost:8080/api/auth/refresh   --cookie "refresh_token=..."
```

Protected endpoint example (change password):
```bash
curl -s -X POST http://localhost:8080/api/auth/change-password   -H "Authorization: Bearer <ACCESS>"   -H 'Content-Type: application/json'   -d '{"currentPassword":"StrongPass123!@#","newPassword":"EvenStrongerPass123!@#"}'
```
next information will be later

test

