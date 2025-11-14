/**
 * Utility functions for date formatting
 * All dates are formatted in dd/mm/yyyy format for international users
 */

/**
 * Formats a date string to dd/mm/yyyy format
 * @param dateString - ISO date string or any valid date string
 * @returns Formatted date in dd/mm/yyyy format
 * @example
 * formatDate("2025-11-14T10:30:00Z") // "14/11/2025"
 */
export const formatDate = (dateString: string): string => {
  const date = new Date(dateString);
  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();
  return `${day}/${month}/${year}`;
};

/**
 * Formats a date string to dd/mm/yyyy HH:mm format
 * @param dateString - ISO date string or any valid date string
 * @returns Formatted date with time in dd/mm/yyyy HH:mm format
 * @example
 * formatDateTime("2025-11-14T10:30:00Z") // "14/11/2025 10:30"
 */
export const formatDateTime = (dateString: string): string => {
  const date = new Date(dateString);
  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${day}/${month}/${year} ${hours}:${minutes}`;
};

/**
 * Formats a date string to a short readable format (dd/mm/yyyy)
 * Alias for formatDate for clarity in some contexts
 * @param dateString - ISO date string or any valid date string
 * @returns Formatted date in dd/mm/yyyy format
 */
export const formatDateShort = (dateString: string): string => {
  return formatDate(dateString);
};

/**
 * Formats a date string to a long readable format with month name
 * @param dateString - ISO date string or any valid date string
 * @returns Formatted date like "14 November 2025"
 */
export const formatDateLong = (dateString: string): string => {
  const date = new Date(dateString);
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  };
  return date.toLocaleDateString('en-GB', options);
};
