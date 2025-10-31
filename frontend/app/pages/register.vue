<template>
  <div class="flex min-h-dvh flex-col items-center justify-center gap-4 p-4">
    <UAuthForm
      :schema="schema"
      title="Create an account"
      :fields="fields"
      :loading-auto="true"
      :submit="{
        size: 'lg',
        variant: 'subtle',
        label: 'Create account'
      }"
      class="w-full max-w-sm"
      icon="i-custom-logo"
      :providers="providers"
      @submit="onSubmit"
    >
      <template #footer>
        Already have an account?
        <ULink to="/" class="text-primary font-medium">Login</ULink>.
      </template>
    </UAuthForm>
  </div>
</template>

<script setup lang="ts">
import * as z from 'zod'
import type { FormSubmitEvent, ButtonProps } from '@nuxt/ui'

definePageMeta({
  middleware: 'workspace'
})

const { register, oauthLogin } = useAuth()

const fields = [
  {
    name: 'name',
    type: 'text' as const,
    label: 'Full name',
    placeholder: 'Enter your name',
    required: true,
    size: 'lg' as const,
    autocomplete: 'name'
  },
  {
    name: 'email',
    type: 'email' as const,
    label: 'Email',
    placeholder: 'Enter your email',
    required: true,
    size: 'lg' as const,
    autocomplete: 'email'
  },
  {
    name: 'password',
    label: 'Password',
    type: 'password' as const,
    placeholder: 'Create a password',
    required: true,
    size: 'lg' as const,
    autocomplete: 'new-password'
  },
  {
    name: 'confirmPassword',
    label: 'Confirm password',
    type: 'password' as const,
    placeholder: 'Confirm your password',
    required: true,
    size: 'lg' as const,
    autocomplete: 'new-password'
  }
]

const schema = z
  .object({
    name: z.string().min(1, 'Name is required'),
    email: z.email('Invalid email'),
    password: z
      .string('Password is required')
      .min(8, 'Must be at least 8 characters'),
    confirmPassword: z
      .string('Confirm your password')
      .min(8, 'Must be at least 8 characters')
  })
  .superRefine((data, ctx) => {
    if (data.password !== data.confirmPassword) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Passwords don't match",
        path: ['confirmPassword']
      })
    }
  })

type Schema = z.output<typeof schema>

const providers = [
  {
    label: 'Google',
    size: 'lg' as const,
    icon: 'i-custom-google',
    ui: { leadingIcon: 'size-4' },
    variant: 'soft',
    onClick: () => {
      handleOAuthLogin('google')
    }
  },
  {
    label: 'GitHub',
    size: 'lg' as const,
    icon: 'i-custom-github',
    ui: { leadingIcon: 'size-4' },
    variant: 'soft',
    onClick: () => {
      handleOAuthLogin('github')
    }
  }
] satisfies ButtonProps[]

async function onSubmit(payload: FormSubmitEvent<Schema>) {
  await register(payload.data.email, payload.data.password, payload.data.name)
}

const handleOAuthLogin = async (provider: 'github' | 'google') => {
  await oauthLogin(provider)
}
</script>
