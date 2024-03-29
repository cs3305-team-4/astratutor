import React from 'react';
import styled from 'styled-components';

import { Link, useHistory } from 'react-router-dom';

import { Layout, Typography, Row, Col, Form, Input, Button, Checkbox } from 'antd';

import { LockOutlined, MailOutlined, UserOutlined } from '@ant-design/icons';

import Config from '../config';
import { fetchRest } from '../api/rest';
import { APIContext } from '../api/api';
import { LoginRequestDTO, LoginResponseDTO } from '../api/definitions';
import DeskImg from '../assets/stock/desk-medium.jpg';

const { Title, Paragraph, Text } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const StyledLayout = styled(Layout)`
  height: calc(100vh - 72px);
  background-image: url(${DeskImg});
  background-size: cover;
`;

export const Login: React.FunctionComponent = () => {
  const api = React.useContext(APIContext);
  const [error, setError] = React.useState('');
  const history = useHistory();

  const onSubmit = async (values: LoginRequestDTO) => {
    try {
      const res = await api.services.login(values);
      api.loginFromJwt(res.jwt);

      history.push('/');
    } catch (e) {
      setError(`Login failed: ${e.message}`);
    }
  };

  if (api.isLoggedIn()) {
    return (
      <StyledLayout>
        <Content>
          <Row style={{ height: 'calc(100vh - 72px)' }} align="middle" justify="center">
            <Col md={9} sm={6} xs={0} />
            <Col md={6} sm={10} xs={24} style={{ padding: '1rem', backgroundColor: 'rgba(255,255,255,0.8)' }}>
              You are already logged in. Click <Link to="/">here</Link> to continue
            </Col>
            <Col md={9} sm={6} xs={0} />
          </Row>
        </Content>
      </StyledLayout>
    );
  }

  return (
    <StyledLayout>
      <Content>
        <Row style={{ height: 'calc(100vh - 72px)' }} align="middle" justify="center">
          <Col md={9} sm={6} xs={0} />
          <Col md={6} sm={10} xs={24} style={{ padding: '2rem 4rem', backgroundColor: 'rgba(255,255,255,0.8)' }}>
            <Form layout="vertical" name="login" onFinish={onSubmit}>
              <UserOutlined
                style={{
                  display: 'block',
                  margin: '0 auto',
                  fontSize: '4rem',
                  padding: '2rem',
                  color: 'rgb(200,200,200)',
                }}
              />
              <Form.Item name="email" rules={[{ required: true, message: 'Please input your email!' }]}>
                <Input size="large" prefix={<MailOutlined />} placeholder="Email" />
              </Form.Item>
              <Form.Item name="password" rules={[{ required: true, message: 'Please input your password!' }]}>
                <Input.Password size="large" prefix={<LockOutlined />} placeholder="Password" />
              </Form.Item>
              <Form.Item>
                <Button style={{ width: '100%' }} type="primary" htmlType="submit">
                  Log in
                </Button>
              </Form.Item>
            </Form>
            <div style={{ color: 'red' }}>{error}</div>
          </Col>
          <Col md={9} sm={6} xs={0} />
        </Row>
      </Content>
    </StyledLayout>
  );
};
