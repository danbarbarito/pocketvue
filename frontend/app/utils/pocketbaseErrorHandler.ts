export interface PocketBaseError {
  data: Record<string, any>
  message: string
  status: number
}

export interface ParsedError {
  title: string
  message: string
  field?: string
}

type ErrorContext = 'login' | 'register' | 'general'

// Error message mappings
const ERROR_MESSAGES: Record<string, string> = {
  'Failed to authenticate.':
    'Invalid email or password. Please check your credentials and try again.',
  "The request doesn't satisfy the collection requirements to authenticate.":
    'Please verify your email address before logging in. Check your inbox for a verification link.',
  'Failed to create record.':
    'Registration failed. Please check your details and try again.',
  validation_not_unique:
    'This email address is already registered. Please try logging in instead.',
  validation_required: 'This field is required.',
  validation_invalid_email: 'Please enter a valid email address.',
  validation_length_out_of_range:
    'Password must be at least 8 characters long.',
  validation_values_mismatch: 'Passwords do not match.'
}

// Field-specific error messages
const FIELD_ERROR_MESSAGES: Record<string, Record<string, string>> = {
  email: {
    validation_not_unique:
      'An account with this email already exists. Please try logging in instead.',
    validation_invalid_email: 'Please enter a valid email address.'
  },
  password: {
    validation_length_out_of_range:
      'Password must be at least 8 characters long.'
  },
  passwordconfirm: {
    validation_values_mismatch: 'Passwords do not match. Please try again.'
  },
  name: {
    validation_required: 'Name is required.'
  }
}

function formatErrorMessage(message: string): string {
  return message
    .replace(/^validation_/, '')
    .replace(/_/g, ' ')
    .replace(/\b\w/g, l => l.toUpperCase())
}

function getFieldErrorMessage(
  fieldName: string,
  errorCode: string,
  errorMessage: string
): string {
  // Check global error messages
  if (ERROR_MESSAGES[errorCode]) {
    return ERROR_MESSAGES[errorCode]
  }

  // Check field-specific messages
  const fieldMessages = FIELD_ERROR_MESSAGES[fieldName.toLowerCase()]
  if (fieldMessages?.[errorCode]) {
    return fieldMessages[errorCode]
  }

  // Fallback to formatted message
  return formatErrorMessage(errorMessage)
}

function getContextualTitle(context: ErrorContext, errorType: string): string {
  const titles: Record<ErrorContext, Record<string, string>> = {
    login: { field_error: 'Login Error', default: 'Login Failed' },
    register: {
      field_error: 'Registration Error',
      default: 'Registration Failed'
    },
    general: { default: 'Error' }
  }

  return titles[context][errorType] || titles[context].default || 'Error'
}

function normalizeError(error: any): PocketBaseError | null {
  // PocketBase ClientResponseError
  if (error?.response && error?.message && error?.status !== undefined) {
    return {
      data: error.response.data || error.response,
      message: error.response.message || error.message,
      status: error.response.status || error.status
    }
  }

  // Direct PocketBase error
  if (error?.data && error?.message && error?.status !== undefined) {
    return error as PocketBaseError
  }

  // Axios-style error
  if (error?.response?.data) {
    return error.response.data
  }

  return null
}

/**
 * Parse a PocketBase error and return a user-friendly error object
 */
export function parsePocketBaseError(
  error: any,
  context: ErrorContext = 'general'
): ParsedError {
  const pbError = normalizeError(error)

  // Handle non-PocketBase errors
  if (!pbError) {
    return {
      title: getContextualTitle(context, 'error'),
      message:
        error?.message || 'An unexpected error occurred. Please try again.'
    }
  }

  // Extract field data (handle nested structure)
  let fieldData = pbError.data
  if (pbError.data?.data && typeof pbError.data.data === 'object') {
    fieldData = pbError.data.data
  }

  // Handle field-specific errors
  if (fieldData && Object.keys(fieldData).length > 0) {
    const firstField = Object.keys(fieldData)[0]
    const fieldError = firstField ? fieldData[firstField] : null

    if (fieldError?.code && fieldError?.message) {
      return {
        title: getContextualTitle(context, 'field_error'),
        message: getFieldErrorMessage(
          firstField,
          fieldError.code,
          fieldError.message
        ),
        field: firstField
      }
    }
  }

  // Handle general errors
  if (pbError.message) {
    return {
      title: getContextualTitle(context, 'general'),
      message: ERROR_MESSAGES[pbError.message] || pbError.message
    }
  }

  // Fallback
  return {
    title: getContextualTitle(context, 'error'),
    message: 'An unexpected error occurred. Please try again.'
  }
}

/**
 * Get just the error message
 */
export function getPocketBaseErrorMessage(
  error: any,
  context: ErrorContext = 'general'
): string {
  return parsePocketBaseError(error, context).message
}

/**
 * Check if error is related to email verification
 */
export function isEmailVerificationError(error: any): boolean {
  const message = error?.message || error?.response?.data?.message || ''
  return message.includes('collection requirements to authenticate')
}

/**
 * Check if error is related to duplicate email
 */
export function isDuplicateEmailError(error: any): boolean {
  const data = error?.data || error?.response?.data?.data || {}
  return data.email?.code === 'validation_not_unique'
}
