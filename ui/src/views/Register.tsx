import React, { useReducer, useState } from 'react';
import styled from 'styled-components';

import { Layout, Typography, Row, Col, Form, Input, Button, Checkbox, Radio, Space, Modal } from 'antd';
import { LockOutlined, MailOutlined, UserAddOutlined } from '@ant-design/icons';

import Config from '../config';
import { fetchRest } from '../api/rest';
import { APIContext } from '../api/api';
import DeskImg from '../assets/stock/desk-medium.jpg';
import { useHistory } from 'react-router-dom';
import { AccountType } from '../api/definitions';

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;

const StyledLayout = styled(Layout)`
  height: calc(100vh - 72px);
  background-image: url(${DeskImg});
  background-size: cover;
`;

export const Register: React.FunctionComponent = () => {
  const [accountType, setAccountType] = useState<AccountType>(AccountType.Student);
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [confirmPassword, setConfirmPassword] = useState<string>('');
  const [under16, setUnder16] = useState<boolean>(false);
  const [parentsEmail, setParentsEmail] = useState<string>('');

  const history = useHistory();
  const [error, setError] = useState<string>('');

  const onSubmit = async (values: any) => {
    console.log('Success:', values);

    let mixin = {};
    if (under16 === true) {
      mixin = {
        parents_email: values.parentsEmail,
      };
    }

    try {
      await api.services.createAccount({
        email: values.email,
        password: values.password,
        type: accountType,
        ...mixin,
      });

      history.push('/');
    } catch (e) {
      setError(`Registration failed: ${e.message}`);
    }
  };

  const FormItems = () => {
    // All default items in a register form
    const items = [
      <Form.Item name="email" key="email" rules={[{ required: true, message: 'Please input your email!' }]}>
        <Input size="large" prefix={<MailOutlined />} placeholder="Email" />
      </Form.Item>,
      <Form.Item name="password" key="password" rules={[{ required: true, message: 'Please enter your password!' }]}>
        <Input.Password size="large" prefix={<LockOutlined />} placeholder="Password" />
      </Form.Item>,
      <Form.Item key="confirm-password" rules={[{ required: true, message: 'Please confirm your password!' }]}>
        <Input.Password size="large" prefix={<LockOutlined />} placeholder="Confirm Password" />
      </Form.Item>,
    ];

    // Items specific to student registration
    if (accountType == 'student') {
      items.push(
        <Form.Item name="under16" valuePropName="checked" initialValue={under16}>
          <Checkbox onChange={() => setUnder16(!under16)} value={under16}>
            I am under the age of 16
          </Checkbox>
        </Form.Item>,
      );

      if (under16) {
        items.push(
          <Form.Item name="parentsEmail">
            <Input size="large" prefix={<MailOutlined />} placeholder="Parents Email" />
          </Form.Item>,
        );
      }
    }

    // Accept Terms
    items.push(
      <Form.Item name="terms" valuePropName="checked" initialValue={false}>
        <Checkbox>
          I accept the <a href="/terms">Terms and Conditions</a>
        </Checkbox>
      </Form.Item>,
    );

    return items;
  };

  const api = React.useContext(APIContext);
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
      <Row style={{ height: 'calc(100vh - 72px)' }} align="middle" justify="center">
        <Col md={8} sm={6} xs={0} />
        <Col md={6} sm={10} xs={24} style={{ padding: '2rem 4rem', backgroundColor: 'rgba(255,255,255,0.8)' }}>
          <UserAddOutlined
            style={{
              display: 'block',
              margin: '0 auto',
              fontSize: '4rem',
              padding: '2rem',
              color: 'rgb(200,200,200)',
            }}
          />
          <Form onFinish={onSubmit}>
            <Form.Item>
              <Row justify="center">
                <Radio.Group value={accountType}>
                  <Radio.Button onClick={() => setAccountType(AccountType.Student)} value="student">
                    Student
                  </Radio.Button>
                  <Radio.Button onClick={() => setAccountType(AccountType.Tutor)} value="tutor">
                    Tutor
                  </Radio.Button>
                </Radio.Group>
              </Row>
            </Form.Item>

            {FormItems()}

            <Form.Item>
              <Button style={{ width: '100%' }} type="primary" htmlType="submit">
                Register
              </Button>
            </Form.Item>
          </Form>
          <div style={{ color: 'red' }}>{error}</div>
        </Col>
        <Col md={8} sm={6} xs={0} />
      </Row>
    </StyledLayout>
  );
};
