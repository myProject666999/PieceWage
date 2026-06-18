import { useState, useEffect } from 'react'
import { Table, Button, Modal, Form, InputNumber, Select, Input, Tag, Space, message, Card, Typography } from 'antd'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons'
import { listProcessPrices, createProcessPrice, listProcessSteps, getEffectivePrice } from '../api'
import dayjs from 'dayjs'

const { Title } = Typography

const gradeMap = { STD: '标准', PRE: '特级', ADV: '高级' }
const gradeColorMap = { STD: 'blue', PRE: 'orange', ADV: 'red' }

export default function PriceManagement() {
  const [prices, setPrices] = useState([])
  const [steps, setSteps] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [previewPrice, setPreviewPrice] = useState(null)
  const [form] = Form.useForm()

  const fetchPrices = async (p = page) => {
    setLoading(true)
    try {
      const res = await listProcessPrices({ page: p, pageSize: 20 })
      setPrices(res.data?.list || [])
      setTotal(res.data?.total || 0)
    } catch (err) {
      message.error('获取单价列表失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchSteps = async () => {
    try {
      const res = await listProcessSteps({ page: 1, pageSize: 200 })
      setSteps(res.data?.list || [])
    } catch {}
  }

  useEffect(() => { fetchPrices(); fetchSteps() }, [])

  const handleCreate = async (values) => {
    try {
      await createProcessPrice(values)
      message.success('单价创建成功')
      setModalOpen(false)
      form.resetFields()
      fetchPrices()
    } catch (err) {
      message.error(err.message || '创建失败')
    }
  }

  const handlePreviewPrice = async () => {
    const processId = form.getFieldValue('processId')
    const gradeLevel = form.getFieldValue('gradeLevel') || 'STD'
    const effectiveDate = form.getFieldValue('effectiveDate')
    if (!processId || !effectiveDate) {
      message.warning('请先选择工序和生效日期')
      return
    }
    try {
      const res = await getEffectivePrice({ processId, gradeLevel, date: effectiveDate })
      setPreviewPrice(res.data)
      if (res.data) {
        message.info(`当前生效单价: ¥${res.data.unitPrice} (版本${res.data.versionNo})`)
      }
    } catch {
      message.info('该日期暂无生效单价')
      setPreviewPrice(null)
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    {
      title: '工序', dataIndex: 'process', width: 180,
      render: (p) => p ? `${p.processName} (${p.processCode})` : '-',
    },
    {
      title: '产品', dataIndex: 'process', width: 140,
      render: (p) => p?.product?.productName || '-',
    },
    { title: '版本号', dataIndex: 'versionNo', width: 80, align: 'center' },
    {
      title: '等级', dataIndex: 'gradeLevel', width: 80, align: 'center',
      render: (v) => <Tag color={gradeColorMap[v] || 'blue'}>{gradeMap[v] || v}</Tag>,
    },
    {
      title: '单价(元)', dataIndex: 'unitPrice', width: 100, align: 'right',
      render: (v) => <span style={{ fontWeight: 'bold', color: '#cf1322' }}>¥{v}</span>,
    },
    { title: '生效日期', dataIndex: 'effectiveDate', width: 110 },
    { title: '失效日期', dataIndex: 'expiryDate', width: 110, render: (v) => v || <Tag>长期</Tag> },
    { title: '说明', dataIndex: 'remark', ellipsis: true },
  ]

  return (
    <div>
      <Title level={4}>工序单价管理</Title>
      <Card>
        <Space style={{ marginBottom: 16 }}>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
            新增单价
          </Button>
        </Space>

        <Table rowKey="id" dataSource={prices} columns={columns}
          loading={loading} pagination={{
            current: page, total, pageSize: 20,
            onChange: (p) => { setPage(p); fetchPrices(p) },
          }} size="middle" scroll={{ x: 1200 }} />
      </Card>

      <Modal title="新增工序单价" open={modalOpen} onCancel={() => setModalOpen(false)}
        onOk={() => form.submit()} width={560}>
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="processId" label="工序" rules={[{ required: true, message: '请选择工序' }]}>
            <Select placeholder="选择工序" showSearch optionFilterProp="label"
              options={steps.map(s => ({ value: s.id, label: `${s.processName} (${s.processCode})` }))} />
          </Form.Item>
          <Form.Item name="gradeLevel" label="等级" initialValue="STD"
            rules={[{ required: true }]}>
            <Select options={[
              { value: 'STD', label: '标准(STD)' },
              { value: 'PRE', label: '特级(PRE)' },
              { value: 'ADV', label: '高级(ADV)' },
            ]} />
          </Form.Item>
          <Form.Item name="unitPrice" label="单件工价(元)"
            rules={[{ required: true, message: '请输入单价' }]}>
            <InputNumber min={0.0001} step={0.01} precision={4} style={{ width: '100%' }}
              placeholder="例如: 3.5000" />
          </Form.Item>
          <Form.Item name="effectiveDate" label="生效日期"
            rules={[{ required: true, message: '请输入生效日期' }]}>
            <Input placeholder="YYYY-MM-DD，例如: 2024-04-01" />
          </Form.Item>
          <Form.Item>
            <Button icon={<SearchOutlined />} onClick={handlePreviewPrice}>
              查看当前生效单价
            </Button>
            {previewPrice && (
              <span style={{ marginLeft: 12 }}>
                当前生效: <Tag color="green">¥{previewPrice.unitPrice} (V{previewPrice.versionNo})</Tag>
              </span>
            )}
          </Form.Item>
          <Form.Item name="remark" label="调整说明">
            <Input.TextArea rows={2} placeholder="例如: 原材料上涨调整+5分" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
