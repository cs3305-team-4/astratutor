import React from 'react';
import styled from 'styled-components';

import countries from 'i18n-iso-countries';
import locales from 'i18n-iso-countries/langs/en.json';

import { useAsync } from 'react-async-hook';

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
  Upload,
  Modal,
  Select,
} from 'antd';

import { EditOutlined, UserAddOutlined, UserOutlined } from '@ant-design/icons';

import { AccountType, ProfileRequestDTO } from '../api/definitions';
import { APIContext } from '../api/api';
import DefaultAvatar from '../assets/default_avatar.png';

import { Redirect, useHistory } from 'react-router-dom';
countries.registerLocale(locales);

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;
const { TextArea } = Input;

export function CreateProfileForm(): React.ReactElement {
  const [error, setError] = React.useState<string>('');
  const api = React.useContext(APIContext);
  const history = useHistory();

  const onFinish = (values: ProfileRequestDTO): Promise<void> => {
    return new Promise<void>((resolve: (value: void | PromiseLike<void>) => void, reject: (reason?: any) => void) => {
      try {
        const img = new Image();
        img.onload = async (e) => {
          const canvas = document.createElement('canvas');
          canvas.width = 256;
          canvas.height = 256;
          const ctx = canvas.getContext('2d');
          ctx.drawImage(img, 0, 0, canvas.width, canvas.height);

          await api.services.createProfileByAccount(api.account, {
            ...values,
            avatar: canvas.toDataURL('image/jpeg', 0.8),
          });

          history.push('/account/profile');
        };

        img.src = DefaultAvatar;

        img.onerror = (e: string | Event) => {
          Modal.error({
            title: 'Error',
            content: `Could not create profile: ${e}`,
          });
        };
      } catch (e) {
        Modal.error({
          title: 'Error',
          content: `Could not create profile: ${e}`,
        });
      }
    });
  };

  return (
    <PageHeader
      title={
        <>
          <Title level={3}>
            <UserAddOutlined /> Create Profile
          </Title>
        </>
      }
      className="site-page-header"
    >
      <Content>
        <Row gutter={4}>
          <Col md={12} sm={24} xs={24}>
            <Form
              layout="vertical"
              name="create-profile"
              onFinish={onFinish}
              initialValues={{
                country: 'IE',
              }}
            >
              <UserAddOutlined
                style={{
                  display: 'block',
                  margin: '0 auto',
                  fontSize: '4rem',
                  padding: '2rem',
                  color: 'rgb(200,200,200)',
                }}
              />
              <Form.Item name="first_name" rules={[{ required: true, message: 'Please input your first name!' }]}>
                <Input placeholder="First Name" />
              </Form.Item>
              <Form.Item name="last_name" rules={[{ required: true, message: 'Please input your last name!' }]}>
                <Input placeholder="Last Name" />
              </Form.Item>
              <Form.Item name="city" rules={[{ required: true, message: 'Please input your city!' }]}>
                <Input placeholder="City" />
              </Form.Item>
              <Form.Item name="country" rules={[{ required: true, message: 'Please input your country!' }]}>
                <Select style={{ width: '100%' }}>
                  {Object.keys(countries.getNames('en', { select: 'official' })).map((value: string, index: number) => (
                    <Select.Option key={index} value={value}>
                      {countries.getNames('en', { select: 'official' })[value]}
                    </Select.Option>
                  ))}
                </Select>
              </Form.Item>
              <Form.Item>
                <Button style={{ width: '100%' }} type="primary" htmlType="submit">
                  Create Profile
                </Button>
              </Form.Item>
            </Form>
            <span style={{ color: 'red' }}>{error}</span>
          </Col>
        </Row>
      </Content>
    </PageHeader>
  );
}
