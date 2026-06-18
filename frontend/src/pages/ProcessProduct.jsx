import { useState, useEffect } from 'react'
import { Card, Table, Button, Modal, Form, Input, Select, InputNumber, Tag, Space, message, Typography, Tabs } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { listProducts, createProduct, listProcessSteps, createProcessStep } from '../api'

const { Title } = Typography

export default function ProcessProduct() {
  const [activeTab, setActiveTab] = useState('products')
  const [products, setProducts] = useState([])
  const [steps, setSteps] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [form] = Form.useForm()

  const fetchProducts = async (p = 1) => {
    setLoading(true)
    try {
      const res = await listProducts({ page: p, pageSize: 20 })
      setProducts(res.data?.list || [])
      setTotal(res.data?.total || 0)
    } catch { message.error('获取产品列表失败') }
    finally { setLoading(false) }
  }

  const fetchSteps = async (p = 1) => {
    setLoading(true)
    try {
      const res = await listProcessSteps({ page: p, pageSize: 20 })
      setSteps(res.data?.list || [])
      setTotal(res.data?.total || 0)
    } catch { message.error('获取工序列表失败') }
    finally { setLoading(false) }
  }

  useEffect(() => {
    if (activeTab === 'products') fetchProducts()
    else fetchSteps()
  }, [activeTab])

  const handleCreateProduct = async (values) => {
    try {
      await createProduct(values)
      message.success('产品创建成功')
      setModalOpen(false); form.resetFields(); fetchProducts()
    } catch (err) { message.error(err.message || '创建失败') }
  }

  const handleCreateStep = async (values) => {
    try {
      await createProcessStep(values)
      message.success('工序创建成功')
      setModalOpen(false); form.resetFields(); fetchSteps()
    } catch (err) { message.error(err.message || '创建失败') }
  }

  const productColumns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '产品编号', dataIndex: 'productCode', width: 120 },
    { title: '产品名称', dataIndex: 'productName', width: 160 },
    { title: '规格型号', dataIndex: 'spec', ellipsis: true },
    { title: '状态', dataIndex: 'status', width: 80,
      render: (v) => v === 1 ? <Tag color="green">启用</Tag> : <Tag color="red">停用</Tag> },
  ]

  const stepColumns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '工序编号', dataIndex: 'processCode', width: 130 },
    { title: '工序名称', dataIndex: 'processName', width: 130 },
    { title: '产品', dataIndex: 'product', width: 140,
      render: (p) => p?.productName || '-' },
    { title: '难度', dataIndex: 'difficulty', width: 80, align: 'center',
      render: (v) => {
        const map = { 1: '简单', 2: '普通', 3: '复杂', 4: '高难' }
        const colorMap = { 1: 'green', 2: 'blue', 3: 'orange', 4: 'red' }
        return <Tag color={colorMap[v]}>{map[v] || v}</Tag>
      }},
    { title: '班组共享', dataIndex: 'isShared', width: 90, align: 'center',
      render: (v) => v === 1 ? <Tag color="purple">共享</Tag> : <Tag>独立</Tag> },
    { title: '描述', dataIndex: 'description', ellipsis: true },
  ]

  return (
    <div>
      <Title level={4}>工序产品管理</Title>
      <Card>
        <Tabs activeKey={activeTab} onChange={(k) => { setActiveTab(k); setPage(1) }}
          items={[
            { key: 'products', label: '产品管理' },
            { key: 'steps', label: '工序管理' },
          ]}
          tabBarExtraContent={
            <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
              新增{activeTab === 'products' ? '产品' : '工序'}
            </Button>
          }
        />

        {activeTab === 'products' ? (
          <Table rowKey="id" dataSource={products} columns={productColumns}
            loading={loading} size="middle"
            pagination={{ current: page, total, pageSize: 20, onChange: (p) => { setPage(p); fetchProducts(p) } }} />
        ) : (
          <Table rowKey="id" dataSource={steps} columns={stepColumns}
            loading={loading} size="middle" scroll={{ x: 900 }}
            pagination={{ current: page, total, pageSize: 20, onChange: (p) => { setPage(p); fetchSteps(p) } }} />
        )}
      </Card>

      <Modal title={`新增${activeTab === 'products' ? '产品' : '工序'}`}
        open={modalOpen} onCancel={() => setModalOpen(false)}
        onOk={() => form.submit()} width={500}>
        {activeTab === 'products' ? (
          <Form form={form} layout="vertical" onFinish={handleCreateProduct}>
            <Form.Item name="productCode" label="产品编号" rules={[{ required: true }]}>
              <Input placeholder="例如: P-1003" />
            </Form.Item>
            <Form.Item name="productName" label="产品名称" rules={[{ required: true }]}>
              <Input placeholder="例如: 精密齿轮C" />
            </Form.Item>
            <Form.Item name="spec" label="规格型号">
              <Input placeholder="例如: 直径100mm" />
            </Form.Item>
          </Form>
        ) : (
          <Form form={form} layout="vertical" onFinish={handleCreateStep}
            initialValues={{ difficulty: 1, isShared: 0 }}>
            <Form.Item name="processCode" label="工序编号" rules={[{ required: true }]}>
              <Input placeholder="例如: PRC-1003-01" />
            </Form.Item>
            <Form.Item name="processName" label="工序名称" rules={[{ required: true }]}>
              <Input placeholder="例如: 精磨" />
            </Form.Item>
            <Form.Item name="productId" label="所属产品" rules={[{ required: true }]}>
              <Select placeholder="选择产品" showSearch optionFilterProp="label"
                options={products.map(p => ({ value: p.id, label: p.productName }))} />
            </Form.Item>
            <Form.Item name="difficulty" label="难度等级" rules={[{ required: true }]}>
              <Select options={[
                { value: 1, label: '简单' }, { value: 2, label: '普通' },
                { value: 3, label: '复杂' }, { value: 4, label: '高难度' },
              ]} />
            </Form.Item>
            <Form.Item name="isShared" label="是否班组共享工序">
              <Select options={[
                { value: 0, label: '否(独立计件)' }, { value: 1, label: '是(班组共享)' },
              ]} />
            </Form.Item>
            <Form.Item name="description" label="工序描述">
              <Input.TextArea rows={2} />
            </Form.Item>
          </Form>
        )}
      </Modal>
    </div>
  )
}
