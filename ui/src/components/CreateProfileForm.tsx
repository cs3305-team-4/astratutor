
import React from 'react';
import styled from 'styled-components';


import {
  useAsync
} from "react-async-hook"

import {
  Typography,
  Layout,
  Row,
  Col,
  Avatar,
  PageHeader,
  Input,
  Button,
  Statistic,
  Form,
  Upload
} from "antd";

import ImgCrop from 'antd-img-crop'

import {
  EditOutlined, UserAddOutlined, UserOutlined
} from "@ant-design/icons"

import config from '../config'
import { fetchRest } from "../api/rest"
import { AccountType, CreateProfileDTO, ReadProfileDTO } from  "../api/definitions"
import { UploadFile } from 'antd/lib/upload/interface';
import DefaultAvatar from "../assets/default_avatar.png"

import { AuthContext } from "../api/auth"
import { Redirect, useHistory } from 'react-router-dom';

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;
const { TextArea } = Input

export interface ProfileProps {
  uuid: string
  type: AccountType
}

export default function CreateProfileForm(props: ProfileProps) {
  const [error, setError] = React.useState<string>("")
  const auth = React.useContext(AuthContext)
  const history = useHistory()

  const [createProfile, setCreateProfile] = React.useState<CreateProfileDTO>({
    first_name: "",
    last_name: "",
    avatar: "",
    city: "",
    country: "",

  })

  const onSubmitCreateProfile = async (values: CreateProfileDTO) => {
    try {
      const res = await fetchRest(
        `${config.apiUrl}/${props.type}s/${props.uuid}/profile`, {
          method: "POST",
          body: JSON.stringify(values),
          headers: {
            "Authorization": `Bearer ${auth.bearerToken}`
          }
        }
      )

      history.push("/account")
    } catch(e) {
      setError(`could not create profile: ${e}`)
    }
  }


  return (
    <PageHeader
      title={<>
        <Title level={3}><UserAddOutlined/> Create Profile</Title>
      </>}
      className="site-page-header"
    >
      <Content>
        <Row gutter={4}>
          <Col md={12} sm={24} xs={24}>
            <Form
              layout="vertical"
              name="create-profile"
              onFinish={onSubmitCreateProfile}
            >
              <UserAddOutlined style={{ display: "block", margin: "0 auto", fontSize: "4rem", padding: "2rem", color: "rgb(200,200,200)"}} />
              <Form.Item
                name="first_name"
                rules={[{ required: true, message: 'Please input your first name!' }]}
              >
                <Input 
                  placeholder="First Name"
                />
              </Form.Item>
              <Form.Item
                name="last_name"
                rules={[{ required: true, message: 'Please input your last name!' }]}
              >
                <Input
                  placeholder="Last Name"
                />
              </Form.Item>
              <Form.Item
                name="city"
                rules={[{ required: true, message: 'Please input your city!' }]}
              >
                <Input
                  placeholder="City"
                />
              </Form.Item>
              <Form.Item
                name="country"
                rules={[{ required: true, message: 'Please input your country!' }]}
              >
                <Input
                  placeholder="Country"
                />
              </Form.Item>
              <Form.Item>
                <Button style={{width: "100%"}} type="primary" htmlType="submit"> 
                  Create Profile
                </Button>
              </Form.Item>
            </Form>
            <span style={{color: "red"}}>
              { error }
            </span>
          </Col>
        </Row>
      </Content>
    </PageHeader>
  )
}
