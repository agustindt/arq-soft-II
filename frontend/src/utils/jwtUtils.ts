/**
 * JWT Utilities
 * Funciones para decodificar y verificar tokens JWT
 */

interface JWTPayload {
  user_id: number;
  email: string;
  username: string;
  role: string;
  exp: number; // Expiration timestamp in seconds
  iat: number; // Issued at timestamp in seconds
  nbf: number; // Not before timestamp in seconds
}

/**
 * Decodifica un token JWT sin verificar la firma
 * Solo para uso en el cliente para verificar expiración
 */
export function decodeJWT(token: string): JWTPayload | null {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) {
      return null;
    }

    const payload = parts[1];
    const decoded = JSON.parse(atob(payload));
    return decoded as JWTPayload;
  } catch (error) {
    console.error('Error decoding JWT:', error);
    return null;
  }
}

/**
 * Verifica si un token JWT ha expirado
 * @param token - El token JWT a verificar
 * @returns true si el token ha expirado, false si aún es válido
 */
export function isTokenExpired(token: string): boolean {
  const decoded = decodeJWT(token);
  if (!decoded || !decoded.exp) {
    return true;
  }

  // exp está en segundos, Date.now() está en milisegundos
  const currentTime = Date.now() / 1000;
  return decoded.exp < currentTime;
}

/**
 * Obtiene el tiempo restante en segundos antes de que expire el token
 * @param token - El token JWT a verificar
 * @returns segundos restantes o 0 si ya expiró
 */
export function getTokenTimeRemaining(token: string): number {
  const decoded = decodeJWT(token);
  if (!decoded || !decoded.exp) {
    return 0;
  }

  const currentTime = Date.now() / 1000;
  const remaining = decoded.exp - currentTime;
  return remaining > 0 ? remaining : 0;
}

/**
 * Verifica si el token está próximo a expirar (menos de 5 minutos)
 * @param token - El token JWT a verificar
 * @returns true si expira en menos de 5 minutos
 */
export function isTokenExpiringSoon(token: string): boolean {
  const timeRemaining = getTokenTimeRemaining(token);
  const FIVE_MINUTES = 5 * 60; // 5 minutos en segundos
  return timeRemaining > 0 && timeRemaining < FIVE_MINUTES;
}

/**
 * Obtiene el rol del usuario desde el token
 * @param token - El token JWT
 * @returns el rol del usuario o null
 */
export function getUserRoleFromToken(token: string): string | null {
  const decoded = decodeJWT(token);
  return decoded?.role || null;
}
