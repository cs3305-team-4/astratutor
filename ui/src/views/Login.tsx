
import React from 'react';
import styled from 'styled-components';

import {
  Layout,
  Typography,
  Row,
  Col,
  Form, 
  Input, 
  Button, 
  Checkbox
} from "antd";
import {
  LockOutlined,
  MailOutlined,
  UserOutlined
} from "@ant-design/icons";

import DeskImg from "../assets/stock/desk-medium.jpg"

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;



const StyledLayout = styled(Layout)`
  height: calc(100vh - 72px);
  background-image: url(${DeskImg});
  background-size: cover;
`;

const Login: React.FunctionComponent = () => {

  const onSubmit = (values: any) => {
    console.log('Success:', values);
  }

  return (
    <StyledLayout>
      <Content>
        <Row
          style={{height: "calc(100vh - 72px)"}}
          align="middle"
          justify="center"
        >
          <Col md={10} sm={6} xs={0}/>
          <Col md={4} sm={10} xs={24} style={{padding: "1rem", backgroundColor: "rgba(255,255,255,0.8)"}}>
            <UserOutlined style={{ display: "block", margin: "0 auto", fontSize: "4rem", padding: "2rem", color: "rgb(200,200,200)"}} />
            <Form
              layout="vertical"
              name="login"
              onFinish={onSubmit}
            >
              <Form.Item
                name="email"
                rules={[{ required: true, message: 'Please input your email!' }]}
              >
                <Input 
                  prefix={<MailOutlined />}
                  placeholder="Email"
                />
              </Form.Item>
              <Form.Item
                name="password"
                rules={[{ required: true, message: 'Please input your password!' }]}
              >
                <Input.Password 
                  prefix={<LockOutlined />}
                  placeholder="Password"
                />
              </Form.Item>
              <Form.Item
                name="remember"
                valuePropName="checked"
              >
                <Checkbox>Remember me</Checkbox>
              </Form.Item>
              <Form.Item>
                <Button style={{width: "100%"}} type="primary" htmlType="submit"> 
                  Log in
                </Button>
              </Form.Item>
            </Form>
          </Col>
          <Col md={10} sm={6} xs={0}/>
        </Row>
      </Content>
    </StyledLayout>
  )
}

export default Login;