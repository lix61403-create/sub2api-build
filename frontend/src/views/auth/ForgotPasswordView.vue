<template>
  <AuthLayout>
    <div class="space-y-6">
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ t('auth.forgotPasswordTitle') }}
        </h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
          {{ t('auth.forgotPasswordHint') }}
        </p>
      </div>

      <div v-if="isSuccess" class="space-y-6">
        <div class="rounded-xl border border-green-200 bg-green-50 p-6 dark:border-green-800/50 dark:bg-green-900/20">
          <div class="flex flex-col items-center gap-4 text-center">
            <div class="flex h-12 w-12 items-center justify-center rounded-full bg-green-100 dark:bg-green-800/50">
              <Icon name="checkCircle" size="lg" class="text-green-600 dark:text-green-400" />
            </div>
            <div>
              <h3 class="text-lg font-semibold text-green-800 dark:text-green-200">
                {{ t('auth.passwordResetSuccess') }}
              </h3>
              <p class="mt-2 text-sm text-green-700 dark:text-green-300">
                {{ t('auth.passwordResetSuccessHint') }}
              </p>
            </div>
          </div>
        </div>

        <div class="text-center">
          <router-link to="/login" class="btn btn-primary inline-flex items-center gap-2">
            <Icon name="login" size="md" />
            {{ t('auth.signIn') }}
          </router-link>
        </div>
      </div>

      <form v-else @submit.prevent="handleResetPassword" class="space-y-5">
        <div>
          <label for="email" class="input-label">
            {{ t('auth.emailLabel') }}
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
              <Icon name="mail" size="md" class="text-gray-400 dark:text-dark-500" />
            </div>
            <input
              id="email"
              v-model="formData.email"
              type="email"
              required
              autofocus
              autocomplete="email"
              :disabled="isSendingCode || isResetting || codeSent"
              class="input pl-11"
              :class="{ 'input-error': errors.email }"
              :placeholder="t('auth.emailPlaceholder')"
            />
          </div>
        </div>

        <div v-if="turnstileEnabled && turnstileSiteKey && shouldShowTurnstile">
          <TurnstileWidget
            ref="turnstileRef"
            :site-key="turnstileSiteKey"
            @verify="onTurnstileVerify"
            @expire="onTurnstileExpire"
            @error="onTurnstileError"
          />
        </div>

        <button
          v-if="!codeSent"
          type="button"
          :disabled="isSendingCode || (turnstileEnabled && !turnstileToken)"
          class="btn btn-primary w-full"
          @click="handleSendCode"
        >
          <svg
            v-if="isSendingCode"
            class="-ml-1 mr-2 h-4 w-4 animate-spin text-white"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          <Icon v-else name="mail" size="md" class="mr-2" />
          {{ isSendingCode ? t('auth.sendingResetCode') : t('auth.sendResetCode') }}
        </button>

        <template v-if="codeSent">
          <div class="rounded-xl border border-blue-200 bg-blue-50 p-4 text-sm text-blue-700 dark:border-blue-800/50 dark:bg-blue-900/20 dark:text-blue-300">
            {{ t('auth.resetCodeSentHint') }}
          </div>

          <div>
            <label for="verifyCode" class="input-label">
              {{ t('auth.verifyCode') }}
            </label>
            <div class="flex gap-2">
              <div class="relative min-w-0 flex-1">
                <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
                  <Icon name="shield" size="md" class="text-gray-400 dark:text-dark-500" />
                </div>
                <input
                  id="verifyCode"
                  v-model="formData.verifyCode"
                  type="text"
                  inputmode="numeric"
                  maxlength="6"
                  required
                  autocomplete="one-time-code"
                  :disabled="isResetting"
                  class="input pl-11"
                  :class="{ 'input-error': errors.verifyCode }"
                  :placeholder="t('auth.verifyCodePlaceholder')"
                />
              </div>
              <button
                type="button"
                class="btn btn-secondary shrink-0"
                :disabled="isSendingCode || countdown > 0 || (turnstileEnabled && showResendTurnstile && !turnstileToken)"
                @click="handleSendCode"
              >
                {{ resendButtonText }}
              </button>
            </div>
          </div>

          <div>
            <label for="password" class="input-label">
              {{ t('auth.newPassword') }}
            </label>
            <div class="relative">
              <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
                <Icon name="lock" size="md" class="text-gray-400 dark:text-dark-500" />
              </div>
              <input
                id="password"
                v-model="formData.password"
                :type="showPassword ? 'text' : 'password'"
                required
                autocomplete="new-password"
                :disabled="isResetting"
                class="input pl-11 pr-11"
                :class="{ 'input-error': errors.password }"
                :placeholder="t('auth.newPasswordPlaceholder')"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 flex items-center pr-3.5 text-gray-400 transition-colors hover:text-gray-600 dark:hover:text-dark-300"
              >
                <Icon v-if="showPassword" name="eyeOff" size="md" />
                <Icon v-else name="eye" size="md" />
              </button>
            </div>
          </div>

          <div>
            <label for="confirmPassword" class="input-label">
              {{ t('auth.confirmPassword') }}
            </label>
            <div class="relative">
              <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
                <Icon name="lock" size="md" class="text-gray-400 dark:text-dark-500" />
              </div>
              <input
                id="confirmPassword"
                v-model="formData.confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                required
                autocomplete="new-password"
                :disabled="isResetting"
                class="input pl-11 pr-11"
                :class="{ 'input-error': errors.confirmPassword }"
                :placeholder="t('auth.confirmPasswordPlaceholder')"
              />
              <button
                type="button"
                @click="showConfirmPassword = !showConfirmPassword"
                class="absolute inset-y-0 right-0 flex items-center pr-3.5 text-gray-400 transition-colors hover:text-gray-600 dark:hover:text-dark-300"
              >
                <Icon v-if="showConfirmPassword" name="eyeOff" size="md" />
                <Icon v-else name="eye" size="md" />
              </button>
            </div>
          </div>

          <button type="submit" :disabled="isResetting" class="btn btn-primary w-full">
            <svg
              v-if="isResetting"
              class="-ml-1 mr-2 h-4 w-4 animate-spin text-white"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            <Icon v-else name="checkCircle" size="md" class="mr-2" />
            {{ isResetting ? t('auth.resettingPassword') : t('auth.resetPassword') }}
          </button>
        </template>
      </form>
    </div>

    <template #footer>
      <p class="text-gray-500 dark:text-dark-400">
        {{ t('auth.rememberedPassword') }}
        <router-link
          to="/login"
          class="font-medium text-primary-600 transition-colors hover:text-primary-500 dark:text-primary-400 dark:hover:text-primary-300"
        >
          {{ t('auth.signIn') }}
        </router-link>
      </p>
    </template>
  </AuthLayout>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { AuthLayout } from '@/components/layout'
import Icon from '@/components/icons/Icon.vue'
import TurnstileWidget from '@/components/TurnstileWidget.vue'
import { useAppStore } from '@/stores'
import { getPublicSettings, resetPasswordWithCode, sendPasswordResetCode } from '@/api/auth'
import { buildAuthErrorMessage } from '@/utils/authError'

const { t } = useI18n()
const appStore = useAppStore()

const isSendingCode = ref<boolean>(false)
const isResetting = ref<boolean>(false)
const isSuccess = ref<boolean>(false)
const codeSent = ref<boolean>(false)
const showResendTurnstile = ref<boolean>(false)
const countdown = ref<number>(0)
const showPassword = ref<boolean>(false)
const showConfirmPassword = ref<boolean>(false)

const turnstileEnabled = ref<boolean>(false)
const turnstileSiteKey = ref<string>('')
const turnstileRef = ref<InstanceType<typeof TurnstileWidget> | null>(null)
const turnstileToken = ref<string>('')

let countdownTimer: ReturnType<typeof setInterval> | null = null

const formData = reactive({
  email: '',
  verifyCode: '',
  password: '',
  confirmPassword: ''
})

const errors = reactive({
  email: '',
  verifyCode: '',
  password: '',
  confirmPassword: '',
  turnstile: ''
})

const validationToastMessage = computed(
  () => errors.email || errors.verifyCode || errors.password || errors.confirmPassword || errors.turnstile || ''
)

const resendButtonText = computed(() => {
  if (isSendingCode.value) {
    return t('auth.sendingResetCode')
  }
  if (countdown.value > 0) {
    return t('auth.resendCodeCountdown', { seconds: countdown.value })
  }
  return t('auth.resendCode')
})

const shouldShowTurnstile = computed(() => !codeSent.value || showResendTurnstile.value)

watch(validationToastMessage, (value, previousValue) => {
  if (value && value !== previousValue) {
    appStore.showError(value)
  }
})

onMounted(async () => {
  try {
    const settings = await getPublicSettings()
    turnstileEnabled.value = settings.turnstile_enabled
    turnstileSiteKey.value = settings.turnstile_site_key || ''
  } catch (error) {
    console.error('Failed to load public settings:', error)
  }
})

onBeforeUnmount(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
})

function normalizeEmail(): string {
  return formData.email.trim().toLowerCase()
}

function startCountdown(seconds: number): void {
  countdown.value = seconds

  if (countdownTimer) {
    clearInterval(countdownTimer)
  }

  countdownTimer = setInterval(() => {
    if (countdown.value > 0) {
      countdown.value--
      return
    }
    if (countdownTimer) {
      clearInterval(countdownTimer)
      countdownTimer = null
    }
  }, 1000)
}

function resetTurnstile(): void {
  if (turnstileRef.value) {
    turnstileRef.value.reset()
  }
  turnstileToken.value = ''
}

function onTurnstileVerify(token: string): void {
  turnstileToken.value = token
  errors.turnstile = ''
}

function onTurnstileExpire(): void {
  turnstileToken.value = ''
  errors.turnstile = t('auth.turnstileExpired')
}

function onTurnstileError(): void {
  turnstileToken.value = ''
  errors.turnstile = t('auth.turnstileFailed')
}

function validateEmailForCode(): boolean {
  errors.email = ''
  errors.turnstile = ''

  if (!formData.email.trim()) {
    errors.email = t('auth.emailRequired')
    return false
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email.trim())) {
    errors.email = t('auth.invalidEmail')
    return false
  }
  if (turnstileEnabled.value && shouldShowTurnstile.value && !turnstileToken.value) {
    errors.turnstile = t('auth.completeVerification')
    return false
  }
  return true
}

function validateResetForm(): boolean {
  errors.verifyCode = ''
  errors.password = ''
  errors.confirmPassword = ''

  let isValid = true
  const verifyCode = formData.verifyCode.trim()

  if (!verifyCode) {
    errors.verifyCode = t('auth.codeRequired')
    isValid = false
  } else if (!/^\d{6}$/.test(verifyCode)) {
    errors.verifyCode = t('auth.invalidCode')
    isValid = false
  }

  if (!formData.password) {
    errors.password = t('auth.passwordRequired')
    isValid = false
  } else if (formData.password.length < 6) {
    errors.password = t('auth.passwordMinLength')
    isValid = false
  }

  if (!formData.confirmPassword) {
    errors.confirmPassword = t('auth.confirmPasswordRequired')
    isValid = false
  } else if (formData.password !== formData.confirmPassword) {
    errors.confirmPassword = t('auth.passwordsDoNotMatch')
    isValid = false
  }

  return isValid
}

async function handleSendCode(): Promise<void> {
  if (turnstileEnabled.value && codeSent.value && !showResendTurnstile.value) {
    showResendTurnstile.value = true
    return
  }

  if (countdown.value > 0 || !validateEmailForCode()) {
    return
  }

  isSendingCode.value = true

  try {
    const response = await sendPasswordResetCode({
      email: normalizeEmail(),
      turnstile_token: turnstileEnabled.value ? turnstileToken.value : undefined
    })

    formData.email = normalizeEmail()
    codeSent.value = true
    showResendTurnstile.value = false
    startCountdown(response.countdown || 60)
    appStore.showSuccess(t('auth.resetCodeSent'))
    resetTurnstile()
  } catch (error: unknown) {
    const message = buildAuthErrorMessage(error, {
      fallback: t('auth.sendResetCodeFailed')
    })
    appStore.showError(message)
    resetTurnstile()
  } finally {
    isSendingCode.value = false
  }
}

async function handleResetPassword(): Promise<void> {
  if (!codeSent.value) {
    await handleSendCode()
    return
  }
  if (!validateResetForm()) {
    return
  }

  isResetting.value = true

  try {
    await resetPasswordWithCode({
      email: normalizeEmail(),
      verify_code: formData.verifyCode.trim(),
      new_password: formData.password
    })

    isSuccess.value = true
    appStore.showSuccess(t('auth.passwordResetSuccess'))
  } catch (error: unknown) {
    const message = buildAuthErrorMessage(error, {
      fallback: t('auth.resetPasswordFailed')
    })
    appStore.showError(message)
  } finally {
    isResetting.value = false
  }
}
</script>
