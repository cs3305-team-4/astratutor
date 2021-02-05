
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
  MailOutlined
} from "@ant-design/icons";

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const StyledLayout = styled(Layout)`
  padding: 2em 0;
`;

const Login: React.FunctionComponent = () => {

  const onSubmit = (values: any) => {
    console.log('Success:', values);
  }

  return (
    <StyledLayout>
      <Row justify="center">
        <Col span={8}>
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
              <Button type="primary" htmlType="submit"> 
                Submit
              </Button>
              <Button type="link" href="/register">
                Register
              </Button>
            </Form.Item>
          </Form>
        </Col>
      </Row>
    </StyledLayout>
  )
}

export default Login;