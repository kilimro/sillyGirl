<script setup lang="ts">
import Button from "ant-design-vue/es/button";
import { Edit3, Plus, RefreshCw, Trash2 } from "lucide-vue-next";
import Form from "ant-design-vue/es/form";
import Input from "ant-design-vue/es/input";
import InputNumber from "ant-design-vue/es/input-number";
import Modal from "ant-design-vue/es/modal";
import Popconfirm from "ant-design-vue/es/popconfirm";
import Segmented from "ant-design-vue/es/segmented";
import Select from "ant-design-vue/es/select";
import Switch from "ant-design-vue/es/switch";
import Tabs, { TabPane } from "ant-design-vue/es/tabs";
import Table from "ant-design-vue/es/table";
import Typography from "ant-design-vue/es/typography";
import message from "ant-design-vue/es/message";
import { timestamp } from "../../../utils";
import { useAdminViewContext } from "../adminViewContext";

const {
  carry,
  changeCarryPlatform,
  loadActiveMessageTool,
  loadCarry,
  loadReplies,
  messageBuckets,
  messageToolAddLabel,
  messageToolHelpText,
  messageToolKind,
  messageToolOptions,
  msgState,
  openActiveMessageTool,
  openCarry,
  openMessage,
  openReply,
  optionMap,
  page,
  recordOptions,
  removeCarry,
  removeMessageRow,
  removeReply,
  replies,
  saveCarry,
  saveMessageRow,
  saveReply,
  scripts,
} = useAdminViewContext();
</script>

<template>
  <section v-if="page === 'message-tools'" class="panel">
    <div class="toolbar">
      <div class="toolbar-left">
        <Segmented
          v-model:value="messageToolKind"
          :options="messageToolOptions"
        />
        <Button type="primary" @click="openActiveMessageTool()"
          ><template #icon><Plus :size="16" /></template
          >{{ messageToolAddLabel }}</Button
        >
        <Button @click="loadActiveMessageTool"
          ><template #icon><RefreshCw :size="16" /></template>刷新</Button
        >
      </div>
      <Typography.Text class="muted">{{ messageToolHelpText }}</Typography.Text>
    </div>

    <Table
      v-if="messageToolKind === 'carry'"
      row-key="chat_id"
      :data-source="carry.rows"
      :pagination="{ total: carry.total, pageSize: 20, onChange: loadCarry }"
    >
      <Table.Column title="#" data-index="id" :width="64" />
      <Table.Column title="平台" data-index="platform" :width="100" />
      <Table.Column title="群号" data-index="chat_id" :width="160" />
      <Table.Column title="备注" data-index="remark" />
      <Table.Column title="操作" :width="150"
        ><template #default="{ record }"
          ><Button type="text" @click="openCarry(record)">编辑</Button
          ><Popconfirm title="确认删除？" @confirm="removeCarry(record)"
            ><Button
              type="text"
              danger
              :title="`删除转发群组 ${record.chat_id}`"
              :aria-label="`删除转发群组 ${record.chat_id}`"
              ><Trash2 :size="16" /></Button></Popconfirm></template
      ></Table.Column>
    </Table>

    <Table
      v-else-if="messageToolKind === 'reply'"
      :row-key="(row: any) => String(row.id)"
      :data-source="replies.rows"
      :pagination="{
        total: replies.total,
        pageSize: 20,
        onChange: loadReplies,
      }"
    >
      <Table.Column title="#" data-index="index" :width="64" />
      <Table.Column title="关键词" data-index="keyword" :width="220" />
      <Table.Column title="回复内容" data-index="value" ellipsis />
      <Table.Column title="对象" data-index="number" :width="140" />
      <Table.Column title="优先级" data-index="priority" :width="90" />
      <Table.Column title="创建时间" data-index="created_at" :width="180">
        <template #default="{ text }">{{ timestamp(text) }}</template>
      </Table.Column>
      <Table.Column title="操作" :width="150">
        <template #default="{ record }">
          <Button
            type="text"
            :title="`编辑回复 ${record.keyword}`"
            :aria-label="`编辑回复 ${record.keyword}`"
            @click="openReply(record)"
            ><Edit3 :size="16"
          /></Button>
          <Popconfirm title="确认删除？" @confirm="removeReply(record)"
            ><Button
              type="text"
              danger
              :title="`删除回复 ${record.keyword}`"
              :aria-label="`删除回复 ${record.keyword}`"
              ><Trash2 :size="16" /></Button
          ></Popconfirm>
        </template>
      </Table.Column>
    </Table>

    <template v-else>
      <Tabs v-model:active-key="msgState.active">
        <TabPane
          v-for="[key, item] in Object.entries(messageBuckets)"
          :key="key"
          :tab="item.label"
        />
      </Tabs>
      <Table row-key="key" :data-source="msgState.rows">
        <Table.Column title="号码" data-index="key" :width="220" />
        <Table.Column title="平台" data-index="platform" :width="140" />
        <Table.Column title="说明" data-index="desc" />
        <Table.Column title="启用" data-index="enable" :width="90"
          ><template #default="{ text }">{{
            text ? "是" : "否"
          }}</template></Table.Column
        >
        <Table.Column title="操作" :width="150"
          ><template #default="{ record }"
            ><Button type="text" @click="openMessage(record)">编辑</Button
            ><Popconfirm title="确认删除？" @confirm="removeMessageRow(record)"
              ><Button
                type="text"
                danger
                :title="`删除 ${record.key}`"
                :aria-label="`删除 ${record.key}`"
                ><Trash2 :size="16" /></Button></Popconfirm></template
        ></Table.Column>
      </Table>
    </template>
  </section>

  <Modal
    :open="!!replies.editing"
    title="回复规则"
    @cancel="replies.editing = null"
    @ok="saveReply"
  >
    <Form layout="vertical"
      ><Form.Item label="关键词/正则" html-for="reply-keyword"
        help="消息内容匹配到这个关键词时自动回复，支持正则表达式（如 ^你好$）"
        ><Input
          id="reply-keyword"
          name="reply-keyword"
          v-model:value="replies.form.keyword" /></Form.Item
      ><Form.Item label="回复内容" html-for="reply-content"
        help="匹配成功后自动回复的文本内容"
        ><Input.TextArea
          id="reply-content"
          name="reply-content"
          v-model:value="replies.form.value"
          :rows="6" /></Form.Item
      ><Form.Item label="限定用户/群号" html-for="reply-number"
        help="只对这个用户或群号生效，留空则全局生效"
        ><Input
          id="reply-number"
          name="reply-number"
          v-model:value="replies.form.number" /></Form.Item
      ><Form.Item label="平台" html-for="reply-platforms"
        help="限定哪个平台生效（微信/QQ/飞书等），不选则所有平台生效"
        ><Select
          id="reply-platforms"
          v-model:value="replies.form.platforms"
          mode="tags" /></Form.Item
      ><Form.Item label="优先级" html-for="reply-priority"
        help="数字越大越先匹配，同时匹配多条规则时按优先级从高到低执行"
        ><InputNumber
          id="reply-priority"
          v-model:value="replies.form.priority"
          style="width: 100%" /></Form.Item
    ></Form>
  </Modal>

  <Modal
    :open="!!carry.editing"
    title="搬运群组"
    width="820px"
    @cancel="carry.editing = null"
    @ok="saveCarry"
  >
    <Form layout="vertical">
      <Form.Item label="平台" html-for="carry-platform" required
        help="限定哪个平台的群参与搬运（微信/QQ等）"
      >
        <Select
          id="carry-platform"
          v-model:value="carry.form.platform"
          :options="optionMap(carry.selects.platforms)"
          @change="changeCarryPlatform"
        />
      </Form.Item>
      <Form.Item label="群号" html-for="carry-chat-id" required
        help="要监听的群聊 ID，该群的消息会被采集并按规则转发到其他地方"
      >
        <Input
          id="carry-chat-id"
          name="carry-chat-id"
          v-model:value="carry.form.chat_id"
        />
      </Form.Item>
      <Form.Item label="备注" html-for="carry-remark"
        help="这个搬运任务的备注说明，方便识别"
      >
        <Input.TextArea
          id="carry-remark"
          name="carry-remark"
          v-model:value="carry.form.remark"
          :rows="2"
        />
      </Form.Item>
      <Form.Item label="工作机器人" html-for="carry-bots"
        help="限定哪些机器人参与搬运（多选），不选则所有机器人都参与"
      >
        <Select
          id="carry-bots"
          v-model:value="carry.form.bots_id"
          mode="multiple"
          :options="optionMap(carry.selects.bots_id)"
        />
      </Form.Item>
      <Form.Item label="处理脚本" html-for="carry-scripts"
        help="消息到达后调用的脚本（节点/Python），自定义处理后转发到目标平台"
      >
        <Select
          id="carry-scripts"
          v-model:value="carry.form.scripts"
          mode="multiple"
          :options="recordOptions(carry.selects.scripts)"
        />
      </Form.Item>
    </Form>
  </Modal>

  <Modal
    :open="!!msgState.editing"
    :title="messageBuckets[msgState.active].label"
    @cancel="msgState.editing = null"
    @ok="saveMessageRow"
  >
    <Form layout="vertical"
      ><Form.Item
        :label="msgState.active === 'private' ? '用户 ID' : '群号'"
        html-for="message-rule-key"
        :help="msgState.active === 'private' ? '要监听的用户 ID（私聊）' : '要监听的群号（群聊）'"
        ><Input
          id="message-rule-key"
          name="message-rule-key"
          v-model:value="msgState.form.key"
          :disabled="!!msgState.editing?.value" /></Form.Item
      ><Form.Item label="平台" html-for="message-rule-platform"
        help="哪个平台的用户/群（微信/QQ/飞书等）"
        ><Select
          id="message-rule-platform"
          v-model:value="msgState.form.platform"
          :options="msgState.platforms" /></Form.Item
      ><Form.Item label="说明" html-for="message-rule-description"
        help="这个监听规则的备注说明，方便识别用途"
        ><Input
          id="message-rule-description"
          name="message-rule-description"
          v-model:value="msgState.form.desc" /></Form.Item
      ><Form.Item label="启用" html-for="message-rule-enabled"
        help="开关：启用后才会监听该用户/群的消息"
        ><Switch
          id="message-rule-enabled"
          v-model:checked="msgState.form.enable"
          aria-label="启用消息规则" /></Form.Item
    ></Form>
  </Modal>
</template>
