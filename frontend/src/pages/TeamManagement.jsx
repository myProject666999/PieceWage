import { useState, useEffect } from 'react'
import { Card, Table, Button, Modal, Form, Input, Select, Tag, Space, message, Typography } from 'antd'
import { PlusOutlined, TeamOutlined } from '@ant-design/icons'
import { listTeams, createTeam, getTeamMembers } from '../api'

const { Title } = Typography

export default function TeamManagement() {
  const [teams, setTeams] = useState([])
  const [members, setMembers] = useState([])
  const [modalOpen, setModalOpen] = useState(false)
  const [memberModalOpen, setMemberModalOpen] = useState(false)
  const [selectedTeam, setSelectedTeam] = useState(null)
  const [form] = Form.useForm()

  const fetchTeams = async () => {
    try {
      const res = await listTeams()
      setTeams(res.data || [])
    } catch { message.error('获取班组列表失败') }
  }

  useEffect(() => { fetchTeams() }, [])

  const handleCreate = async (values) => {
    try {
      await createTeam(values)
      message.success('班组创建成功')
      setModalOpen(false); form.resetFields(); fetchTeams()
    } catch (err) { message.error(err.message || '创建失败') }
  }

  const handleViewMembers = async (team) => {
    setSelectedTeam(team)
    try {
      const res = await getTeamMembers(team.id)
      setMembers(res.data || [])
      setMemberModalOpen(true)
    } catch { message.error('获取成员失败') }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '班组编号', dataIndex: 'teamCode', width: 120 },
    { title: '班组名称', dataIndex: 'teamName', width: 160 },
    { title: '状态', dataIndex: 'status', width: 80,
      render: (v) => v === 1 ? <Tag color="green">启用</Tag> : <Tag color="red">停用</Tag> },
    {
      title: '操作', width: 120,
      render: (_, r) => (
        <Button size="small" icon={<TeamOutlined />} onClick={() => handleViewMembers(r)}>
          查看成员
        </Button>
      ),
    },
  ]

  const memberColumns = [
    { title: 'ID', dataIndex: 'userId', width: 60 },
    { title: '工人ID', dataIndex: 'userId', width: 100 },
    { title: '入组日期', dataIndex: 'joinDate', width: 120 },
    { title: '离组日期', dataIndex: 'leaveDate', width: 120,
      render: (v) => v || <Tag color="green">在组</Tag> },
  ]

  return (
    <div>
      <Title level={4}>班组管理</Title>
      <Card>
        <Space style={{ marginBottom: 16 }}>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
            新建班组
          </Button>
        </Space>
        <Table rowKey="id" dataSource={teams} columns={columns} size="middle" />
      </Card>

      <Modal title="新建班组" open={modalOpen} onCancel={() => setModalOpen(false)}
        onOk={() => form.submit()}>
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="teamName" label="班组名称" rules={[{ required: true }]}>
            <Input placeholder="例如: 冲压二组" />
          </Form.Item>
          <Form.Item name="teamCode" label="班组编号" rules={[{ required: true }]}>
            <Input placeholder="例如: TM003" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal title={selectedTeam ? `${selectedTeam.teamName} - 班组成员` : '班组成员'}
        open={memberModalOpen} onCancel={() => setMemberModalOpen(false)} footer={null}>
        <Table rowKey="id" dataSource={members} columns={memberColumns} size="small" pagination={false} />
      </Modal>
    </div>
  )
}
