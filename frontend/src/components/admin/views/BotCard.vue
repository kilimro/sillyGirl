<script setup lang="ts">
import Button from "ant-design-vue/es/button";
import Tag from "ant-design-vue/es/tag";
import Typography from "ant-design-vue/es/typography";
import {
  Antenna,
  CircleCheck,
  CircleX,
  Pause,
  Play,
  Settings,
} from "lucide-vue-next";
import { platformIcon } from "./platformIcons";

defineProps<{
  record: {
    platform: string;
    label: string;
    online?: boolean;
    manageable?: boolean;
    count?: number;
    bots_id?: string[];
  };
  enabled: boolean;
}>();

const emit = defineEmits<{ toggle: [value: boolean]; open: [] }>();
</script>

<template>
  <article
    class="bot-card"
    :class="{
      'is-online': record.online,
      'is-disabled': !enabled,
    }"
  >
    <header class="bot-card-header">
      <span class="bot-card-avatar" v-html="platformIcon(record.platform)" style="display:flex;align-items:center;justify-content:center;width:32px;height:32px;"></span>
      <div class="bot-card-heading">
        <strong>{{ record.label }}</strong>
        <span>{{ record.platform === 'yyw' ? 'MianPro 专用' : record.platform }}</span>
      </div>
      <div class="bot-card-actions">
        <template v-if="record.manageable !== false">
          <Button
            v-if="enabled"
            class="bot-card-toggle bot-card-pause"
            type="text"
            shape="circle"
            :title="`暂停 ${record.label}`"
            :aria-label="`暂停 ${record.label}`"
            @click="emit('toggle', false)"
          >
            <template #icon><Pause :size="18" /></template>
          </Button>
          <Button
            v-else
            class="bot-card-toggle bot-card-play"
            type="text"
            shape="circle"
            :title="`开启 ${record.label}`"
            :aria-label="`开启 ${record.label}`"
            @click="emit('toggle', true)"
          >
            <template #icon><Play :size="18" /></template>
          </Button>
        </template>
        <Button
          class="bot-card-settings"
          type="text"
          shape="circle"
          :title="`${record.label}设置`"
          :aria-label="`${record.label}设置`"
          @click="emit('open')"
        >
          <template #icon><Settings :size="18" /></template>
        </Button>
      </div>
    </header>

    <div class="bot-card-status-list">
      <div class="bot-card-status">
        <CircleCheck
          v-if="enabled"
          class="bot-status-enabled"
          :size="20"
          aria-hidden="true"
        />
        <CircleX
          v-else
          class="bot-status-disabled"
          :size="20"
          aria-hidden="true"
        />
        <span>启用状态</span>
        <strong
          :class="enabled ? 'bot-status-enabled' : 'bot-status-disabled'"
        >
          {{ enabled ? "已启用" : "未启用" }}
        </strong>
      </div>
      <div class="bot-card-status">
        <Antenna
          :class="record.online ? 'bot-status-online' : 'bot-status-offline'"
          :size="20"
          aria-hidden="true"
        />
        <span>连接状态</span>
        <strong
          :class="record.online ? 'bot-status-online' : 'bot-status-offline'"
        >
          {{ record.online ? "已连接" : "未连接" }}
        </strong>
      </div>
    </div>

    <div class="bot-card-meta">
      <div>
        <span>实例数</span><strong>{{ record.count || 0 }}</strong>
      </div>
      <div>
        <span>Bot ID</span>
        <Typography.Text class="mono bot-card-bot-ids">{{
          record.bots_id?.length ? record.bots_id.join(", ") : "-"
        }}</Typography.Text>
      </div>
      <div>
        <span>类型</span><Tag color="blue">内置 BOT</Tag>
      </div>
    </div>
  </article>
</template>
