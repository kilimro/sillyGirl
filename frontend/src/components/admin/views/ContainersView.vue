<script setup lang="ts">
import Button from "ant-design-vue/es/button";
import Form from "ant-design-vue/es/form";
import Input from "ant-design-vue/es/input";
import Modal from "ant-design-vue/es/modal";
import { Plus, RefreshCw, Trash2 } from "lucide-vue-next";
import Popconfirm from "ant-design-vue/es/popconfirm";
import Segmented from "ant-design-vue/es/segmented";
import Table from "ant-design-vue/es/table";
import Tag from "ant-design-vue/es/tag";
import Typography from "ant-design-vue/es/typography";
import { timestamp } from "../../../utils";
import { useAdminViewContext } from "../adminViewContext";

const {
  containerAddLabel,
  containerHelpText,
  containerKind,
  containerOptions,
  daidai,
  loadActiveContainerPanels,
  refreshActiveContainerPanels,
  openActiveContainerPanel,
  openDaidaiPanel,
  openQinglongPanel,
  page,
  qinglong,
  removeDaidaiPanel,
  removeQinglongPanel,
  saveDaidaiPanel,
  saveQinglongPanel,
  testDaidaiPanel,
  testQinglongPanel,
} = useAdminViewContext();
</script>

<template>
  <section v-if="page === 'containers'" class="panel">
    <div class="toolbar">
      <div class="toolbar-left">
        <Segmented v-model:value="containerKind" :options="containerOptions" />
        <Button type="primary" @click="openActiveContainerPanel()">
          <template #icon><Plus :size="16" /></template>{{ containerAddLabel }}
        </Button>
        <Button @click="refreshActiveContainerPanels">
          <template #icon><RefreshCw :size="16" /></template>刷新
        </Button>
      </div>
      <Typography.Text class="muted">{{ containerHelpText }}</Typography.Text>
    </div>

    <Table
      v-if="containerKind === 'qinglong'"
      row-key="id"
      :loading="qinglong.loading"
      :data-source="qinglong.rows"
      :pagination="{ total: qinglong.total, pageSize: 20 }"
    >
      <Table.Column title="#" :width="72">
        <template #default="{ index }">{{ index + 1 }}</template>
      </Table.Column>
      <Table.Column title="名称" data-index="name" :width="180">
        <template #default="{ record }">
          <Typography.Text strong>{{ record.name || record.address }}</Typography.Text>
        </template>
      </Table.Column>
      <Table.Column title="地址" data-index="address" ellipsis />
      <Table.Column title="Client ID" data-index="client_id" :width="220" ellipsis />
      <Table.Column title="状态" data-index="status" :width="120">
        <template #default="{ record }">
          <Tag :color="record.status === 'online' ? 'green' : 'default'">
            {{ record.status === "online" ? "在线" : "未检测" }}
          </Tag>
        </template>
      </Table.Column>
      <Table.Column title="最后检测" data-index="last_checked_at" :width="180">
        <template #default="{ text }">{{ timestamp(text) }}</template>
      </Table.Column>
      <Table.Column title="操作" :width="210">
        <template #default="{ record }">
          <Button type="text" @click="testQinglongPanel(record)">检测</Button>
          <Button type="text" @click="openQinglongPanel(record)">编辑</Button>
          <Popconfirm title="确认删除这个青龙面板？" @confirm="removeQinglongPanel(record)">
            <Button type="text" danger><Trash2 :size="16" /></Button>
          </Popconfirm>
        </template>
      </Table.Column>
    </Table>

    <Table
      v-else
      row-key="id"
      :loading="daidai.loading"
      :data-source="daidai.rows"
      :pagination="{ total: daidai.total, pageSize: 20 }"
    >
      <Table.Column title="#" :width="72">
        <template #default="{ index }">{{ index + 1 }}</template>
      </Table.Column>
      <Table.Column title="名称" data-index="name" :width="180">
        <template #default="{ record }">
          <Typography.Text strong>{{ record.name || record.address }}</Typography.Text>
        </template>
      </Table.Column>
      <Table.Column title="地址" data-index="address" ellipsis />
      <Table.Column title="App Key" data-index="app_key" :width="220" ellipsis />
      <Table.Column title="状态" data-index="status" :width="120">
        <template #default="{ record }">
          <Tag :color="record.status === 'online' ? 'green' : 'default'">
            {{ record.status === "online" ? "在线" : "未检测" }}
          </Tag>
        </template>
      </Table.Column>
      <Table.Column title="最后检测" data-index="last_checked_at" :width="180">
        <template #default="{ text }">{{ timestamp(text) }}</template>
      </Table.Column>
      <Table.Column title="操作" :width="210">
        <template #default="{ record }">
          <Button type="text" @click="testDaidaiPanel(record)">检测</Button>
          <Button type="text" @click="openDaidaiPanel(record)">编辑</Button>
          <Popconfirm title="确认删除这个呆呆面板？" @confirm="removeDaidaiPanel(record)">
            <Button type="text" danger><Trash2 :size="16" /></Button>
          </Popconfirm>
        </template>
      </Table.Column>
    </Table>
  </section>

  <Modal
    :open="!!qinglong.editing"
    title="青龙面板"
    width="720px"
    :confirm-loading="qinglong.saving"
    @cancel="qinglong.editing = null"
    @ok="saveQinglongPanel"
  >
    <Form layout="vertical">
      <Form.Item label="名称"><Input v-model:value="qinglong.form.name" placeholder="例如：主青龙" /></Form.Item>
      <Form.Item label="青龙地址" required>
        <Input v-model:value="qinglong.form.address" placeholder="http://127.0.0.1:5700" />
      </Form.Item>
      <Form.Item label="Client ID" required><Input v-model:value="qinglong.form.client_id" /></Form.Item>
      <Form.Item label="Client Secret" required>
        <Input.Password v-model:value="qinglong.form.client_secret" />
      </Form.Item>
      <Button @click="testQinglongPanel()" :loading="qinglong.testing">
        <template #icon><RefreshCw :size="16" /></template>检测连接
      </Button>
    </Form>
  </Modal>

  <Modal
    :open="!!daidai.editing"
    title="呆呆面板"
    width="720px"
    :confirm-loading="daidai.saving"
    @cancel="daidai.editing = null"
    @ok="saveDaidaiPanel"
  >
    <Form layout="vertical">
      <Form.Item label="名称"><Input v-model:value="daidai.form.name" placeholder="例如：主呆呆" /></Form.Item>
      <Form.Item label="呆呆面板地址" required>
        <Input v-model:value="daidai.form.address" placeholder="http://127.0.0.1:5701" />
      </Form.Item>
      <Form.Item label="App Key" required><Input v-model:value="daidai.form.app_key" /></Form.Item>
      <Form.Item label="App Secret" required>
        <Input.Password v-model:value="daidai.form.app_secret" />
      </Form.Item>
      <Button @click="testDaidaiPanel()" :loading="daidai.testing">
        <template #icon><RefreshCw :size="16" /></template>检测连接
      </Button>
    </Form>
  </Modal>
</template>
