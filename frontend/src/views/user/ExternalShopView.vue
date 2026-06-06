<template>
  <AppLayout>
    <div class="mx-auto flex h-[calc(100vh-7rem)] max-w-7xl flex-col gap-4">
      <div class="card flex flex-col gap-3 p-5 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p class="text-xs font-medium text-gray-400 dark:text-gray-500">充值商品</p>
          <h1 class="mt-1 text-xl font-bold text-gray-900 dark:text-white">Mpl API 充值中心</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            在下方商品页面完成支付。若页面无法正常显示，请使用新窗口打开。
          </p>
        </div>
        <a
          class="btn btn-primary shrink-0"
          :href="shopUrl"
          target="_blank"
          rel="noopener noreferrer"
        >
          新窗口打开
        </a>
      </div>

      <div class="card min-h-0 flex-1 overflow-hidden p-0">
        <iframe
          class="h-full w-full border-0 bg-white"
          :src="shopUrl"
          title="Mpl API 充值商品"
          loading="eager"
          referrerpolicy="strict-origin-when-cross-origin"
          sandbox="allow-forms allow-scripts allow-same-origin allow-popups allow-popups-to-escape-sandbox allow-top-navigation-by-user-activation"
        />
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import AppLayout from '@/components/layout/AppLayout.vue'

const shopUrl = 'https://catfk.com/shop/api.mplvip.xyz'
</script>
