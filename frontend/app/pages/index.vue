<template>
  <div class="flex min-h-dvh flex-col items-center justify-center gap-4 p-4">
    <UAuthForm
      :schema="schema"
      title="Login to Pocketvue"
      :fields="fields"
      :loading-auto="true"
      :submit="{
        size: 'lg',
        variant: 'subtle',
        label: 'Login'
      }"
      class="w-full max-w-sm"
      icon="i-custom-logo"
      :loading="loading"
      :providers="providers"
      @submit="onSubmit"
    >
      <template #footer>
        Don't have an account?
        <ULink to="/register" class="text-primary font-medium">Sign up</ULink>.
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

const { login, oauthLogin } = useAuth()

const fields = ref([
  {
    name: 'email',
    type: 'email' as const,
    autoComplete: 'email',
    label: 'Email',
    size: 'lg' as const,
    autocomplete: 'email'
  },
  {
    name: 'password',
    type: 'password' as const,
    autoComplete: 'current-password',
    label: 'Password',
    size: 'lg' as const,
    autocomplete: 'current-password'
  }
])
const loading = ref(false)
const schema = z.object({
  email: z.email('Invalid email'),
  password: z
    .string('Password is required')
    .min(8, 'Must be at least 8 characters')
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
  loading.value = true
  const success = await login(payload.data.email, payload.data.password)
  if (!success) {
    loading.value = false
  }
}

const handleOAuthLogin = async (provider: 'github' | 'google') => {
  loading.value = true
  const success = await oauthLogin(provider)
  if (!success) {
    loading.value = false
  }
}
</script>
