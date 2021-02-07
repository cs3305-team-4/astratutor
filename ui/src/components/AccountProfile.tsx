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

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;
const { TextArea } = Input

export interface ProfileProps {
  uuid: string
  type: AccountType
}

export default function AccountProfile(props: ProfileProps) {
  const [profile, setProfile] = React.useState<ReadProfileDTO | undefined>(undefined)
  const [error, setError] = React.useState<string>("")

  const [editDesc, setEditDesc] = React.useState<boolean>()
  const [newDesc, setNewDesc] = React.useState<string>()

  const auth = React.useContext(AuthContext)

  React.useState(() => {
    const commitDesc = async (newDesc: string) => {
      try {
        const res = await fetchRest(`${config.apiUrl}/${props.type}s/${props.uuid}/profile`)
        const profile = await res.json() as ReadProfileDTO
        setProfile(profile)
      } catch(e) {
        setError(`could not load profile ${e}`)
      }
    }
  })

  useAsync(async () => {
    try {
      const res = await fetchRest(
        `${config.apiUrl}/${props.type}s/${props.uuid}/profile`, {
          headers: {
            'Authoriziation': `Bearer ${auth.bearerToken}`
          }
        }
      )
      const profile = await res.json() as ReadProfileDTO
      setProfile(profile)
    } catch(e) {
      setError(`could not load profile ${e}`)
    }
  }, [])

  if (profile === undefined) {
    return <h1>{error}</h1>
  }

  return (
    <PageHeader
      title={
        <>
          {profile.first_name + " " + profile.last_name}
        </>}
      className="site-page-header"
      subTitle={<>
        <Text type="secondary">
          Chemistry Teacher with 24yrs experience
        </Text>
        <Button size="small" style={{margin: "4px"}}>
          <EditOutlined/>
        </Button>
      </>}
      avatar={{ size: 64, src: profile.avatar }}
    >
      <Content>        
        <Row gutter={4}>
          <Col md={12} sm={24} xs={24}>
            <Title level={5}>
              Description
              <Button size="small" style={{margin: "4px"}}>
                <EditOutlined onClick={()=>setEditDesc(!editDesc)}/>
              </Button>
            </Title>
            {
              !editDesc 
                ? 
              (<Paragraph>{ profile.description }</Paragraph>) 
                : 
              (<>
                <TextArea placeholder="textarea with clear icon" allowClear onChange={onChange} />
                <EditOutlined onClick={()=>setEditDesc(!editDesc)}/>
              </>)
            }
          </Col>
          <Col md={12} sm={24} xs={24}>
            <Title level={5}>
              Availability
              <Button size="small" style={{margin: "4px"}}>
                <EditOutlined/>
              </Button>
            </Title>
            <Paragraph>
              MIT
            </Paragraph>
          </Col>
          <Col md={12} sm={24} xs={24}>
            <Title level={5}>
              Qualifications
              <Button size="small" style={{margin: "4px"}}>
                <EditOutlined/>
              </Button>
            </Title>
            <Paragraph>
              BSc Computer Science
            </Paragraph>
          </Col>
          <Col md={12} sm={24} xs={24}>
            <Title level={5}>
              Work Experience
              <Button size="small" style={{margin: "4px"}}>
                <EditOutlined/>
              </Button>
            </Title>
            <Paragraph>
              MIT
            </Paragraph>
          </Col>
        </Row>
      </Content>
  </PageHeader>)
}