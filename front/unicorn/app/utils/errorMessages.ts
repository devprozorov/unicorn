export type TranslateFn = (key: string) => string

const errorMap: Record<string, string> = {
  // Auth
  auth_error: 'auth.errors.generic',
  invalid_session: 'auth.errors.invalidSession',
  invalid_mfa_state: 'auth.errors.invalidMfaState',
  invalid_code: 'auth.errors.invalidCode',
  INVALID_CODE: 'auth.errors.invalidCode',
  bad_request: 'auth.errors.badRequest',
  login_failed: 'auth.errors.loginFailed',
  register_failed: 'auth.errors.registerFailed',
  user_exists: 'auth.errors.userExists',
  user_not_found: 'auth.errors.userNotFound',
  weak_password: 'auth.errors.weakPassword',
  password_too_short: 'auth.errors.weakPassword',

  // MFA
  mfa_required: 'errors.mfaRequired',
  MFA_ENROLL_FAILED: 'errors.mfaEnrollFailed',

  // Generic
  load_failed: 'errors.loadFailed',
  save_failed: 'errors.saveFailed',
  upload_failed: 'errors.uploadFailed',
  delete_failed: 'errors.deleteFailed',
  update_failed: 'errors.updateFailed',
  create_failed: 'errors.createFailed',
  block_failed: 'errors.blockFailed',
  unblock_failed: 'errors.unblockFailed',
  activation_failed: 'errors.activationFailed',
  deactivation_failed: 'errors.deactivationFailed',
  payment_failed: 'errors.paymentFailed',

  // HTTP / access
  unauthorized: 'errors.unauthorized',
  forbidden: 'errors.forbidden',
  not_found: 'errors.notFound',
  rate_limited: 'errors.rateLimited',
  inbox_endpoint_404: 'errors.notFound'
}

export function getErrorMessage(code: unknown, t: TranslateFn, fallbackKey = 'errors.unknown') {
  if (!code) {
    const fallback = t(fallbackKey)
    return fallback === fallbackKey ? 'Something went wrong' : fallback
  }

  const raw = String(code)
  const normalized = raw.toLowerCase()

  if (raw.startsWith('errors.') || raw.startsWith('auth.')) {
    const translated = t(raw)
    return translated === raw ? raw : translated
  }

  if (normalized.startsWith('http_')) {
    const httpKey = normalized
    if (httpKey === 'http_401') return t('errors.unauthorized')
    if (httpKey === 'http_403') return t('errors.forbidden')
    if (httpKey === 'http_404') return t('errors.notFound')
    if (httpKey === 'http_429') return t('errors.rateLimited')
  }

  const key = errorMap[raw] || errorMap[normalized] || fallbackKey
  const translated = t(key)
  return translated === key ? raw : translated
}
