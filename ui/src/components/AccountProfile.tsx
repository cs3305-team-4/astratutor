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
  Button,
  Image
} from "antd";

import {
  EditOutlined
} from "@ant-design/icons"

import config from '../config'
import { fetchRest } from "../api/rest"
import { ProfileDTO } from  "../api/definitions"

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;

export interface ProfileProps {
  uuid: string
  type: string
}

export default function AccountProfile(props: ProfileProps) {
  const [profile, setProfile] = React.useState<ProfileDTO | undefined>(undefined)
  const [error, setError] = React.useState<string>("")

  useAsync(async () => {
    try {
      const res = await fetchRest(`${config.apiUrl}/${props.type}/${props.uuid}/profile`)
      const profile = await res.json() as ProfileDTO
      setProfile(profile)
    } catch(e) {
      setError(`could not load profile ${e}`)
    }
  }, [])

  // if (profile == undefined) {
  //   return (
  //     <Layout>
  //     </Layout>
  //   )
  // }

  return (
    <PageHeader
      title="Oisin Canty"
      className="site-page-header"
      subTitle={<>
        <Text type="secondary">
          Chemistry Teacher with 24yrs experience
        </Text>
        <Button size="small" style={{margin: "4px"}}>
          <EditOutlined/>
        </Button>
      </>}
      avatar={{ size: 64, src: 'https://zos.alipayobjects.com/rmsportal/ODTLcjxAfvqbxHnVXCYX.png' }}
    >
      <Content>
        <Row>
          <Col md={12} sm={24} xs={24}>
            <Title level={5}>
              Description
              <Button size="small" style={{margin: "4px"}}>
                <EditOutlined/>
              </Button>
            </Title>
            <Paragraph>
              Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis id tortor lectus. Nullam eu nisi et est pretium hendrerit. Etiam sed bibendum metus. Etiam ut lacinia lorem. Quisque fermentum tristique eros, ac lacinia diam consequat ac. Aenean eros ipsum, interdum non massa nec, laoreet iaculis odio. Nam volutpat justo vitae orci lacinia, sed mattis justo gravida. Aliquam consequat placerat libero, ac imperdiet mi pellentesque commodo. Aliquam finibus diam iaculis ipsum pharetra vehicula. Aliquam non nulla eu mi pharetra facilisis fermentum ut lorem. Proin tincidunt turpis et placerat gravida. Nulla eget posuere nulla, non euismod enim. Donec ex nisl, ultricies cursus odio id, sollicitudin dictum dui. Nunc pharetra iaculis tellus. Donec finibus urna semper, convallis velit at, condimentum mauris.
            </Paragraph>
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
        </Row>
        <Row>
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