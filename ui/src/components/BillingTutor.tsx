import React from 'react';
import styled from 'styled-components';

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
  Alert,
  Table,
  Skeleton,
} from 'antd';

import { BankOutlined } from '@ant-design/icons';

import { AccountType, ProfileRequestDTO } from '../api/definitions';
import { APIContext } from '../api/api';
import DefaultAvatar from '../assets/default_avatar.png';

import { Redirect, useHistory } from 'react-router-dom';

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;
const { TextArea } = Input;

export function BillingTutor(): React.ReactElement {
  const [error, setError] = React.useState<string>('');

  const [billingOnboard, setBillingOnboard] = React.useState<boolean>(true);
  const [billingOnboardUrl, setBillingOnboardUrl] = React.useState<string>('');
  const [billingRequirementsMet, setBillingRequirementsMet] = React.useState<boolean>(true);
  const [billingPanelUrl, setBillingPanelUrl] = React.useState<string>('');
  const [ready, setReady] = React.useState<boolean>(false);
  const api = React.useContext(APIContext);
  const history = useHistory();

  useAsync(async () => {
    const onboard = await api.services.readTutorBillingOnboardStatus(api.account.id);
    setBillingOnboard(onboard);

    if (!onboard) {
      setBillingOnboardUrl(await api.services.readTutorBillingOnboardUrl(api.account.id));
      setReady(true);
    } else {
      const reqsMet = await api.services.readTutorBillingRequirementsMetStatus(api.account.id);
      setBillingRequirementsMet(reqsMet);
      setBillingPanelUrl(await api.services.readTutorBillingPanelUrl(api.account.id));
      setReady(true);
    }
  }, []);

  const redirectBillingOnboard = () => {
    window.location.href = billingOnboardUrl;
  };

  const redirectBillingPanelAccount = () => {
    window.open(billingPanelUrl + '#/account');
  };

  const redirectBillingPanel = () => {
    window.open(billingPanelUrl + '#/account');
  };

  if (!ready) {
    return (
      <>
        <Skeleton />
      </>
    );
  }

  return (
    <Typography>
      <Content>
        <Row gutter={16} style={{ margin: '1rem' }}>
          <Title level={3}>
            <BankOutlined /> Billing
          </Title>
          <Col md={24} sm={24} xs={24}>
            {!billingOnboard && (
              <>
                <Alert
                  message="Information Required"
                  description="You have not yet setup billing and linked a bank account, you will not be able to request payouts until this is completed."
                  type="warning"
                  showIcon
                  style={{ margin: '1rem 0' }}
                />
              </>
            )}
            {!billingRequirementsMet && (
              <>
                <Alert
                  message="Information Required"
                  description="Additional information is required to continue payouts, you will not be able to request payouts until this is completed"
                  type="warning"
                  showIcon
                  style={{ margin: '1rem 0' }}
                />
              </>
            )}
          </Col>
          <Col md={8} sm={24} xs={24}>
            <Statistic
              style={{ margin: '0 auto' }}
              valueStyle={{ fontSize: '4rem' }}
              title="Balance For Payout"
              value={'€30.87'}
            />
          </Col>
          <Col md={8} sm={24} xs={24}>
            <Statistic valueStyle={{ fontSize: '4rem' }} title="Total This Month" value={'€0.00'} />
          </Col>
          <Col md={8} sm={24} xs={24}>
            <Statistic valueStyle={{ fontSize: '4rem' }} title="Total This Year" value={'€0.00'} />
          </Col>
          <Col md={24} sm={24} xs={24}>
            {!billingOnboard && (
              <>
                <Button type="primary" onClick={redirectBillingOnboard}>
                  Setup Billing
                </Button>
              </>
            )}
            {billingOnboard && billingRequirementsMet && (
              <>
                <Button type="primary" onClick={redirectBillingOnboard} style={{ margin: '0.5em' }}>
                  Payout
                </Button>
                <Button style={{ margin: '0.5em' }} onClick={redirectBillingPanelAccount}>
                  Modify Linked Bank Account
                </Button>
              </>
            )}
            {!billingRequirementsMet && (
              <>
                <Button type="primary" disabled onClick={redirectBillingOnboard} style={{ margin: '0.5em' }}>
                  Payout
                </Button>
                <Button style={{ margin: '0.5em' }} onClick={redirectBillingPanel}>
                  Provide Required Information
                </Button>
              </>
            )}
          </Col>
        </Row>
        <Row gutter={16} style={{ margin: '1rem' }}>
          <Title level={3}>Invoices</Title>
          <Col md={24} sm={24} xs={24}>
            <Table
              locale={{
                emptyText: 'No invoices available',
              }}
              columns={[
                { title: 'Lesson', key: 'degree', dataIndex: 'degree' },
                { title: 'Date', key: 'field', dataIndex: 'field' },
                { title: 'Amount', key: 'school', dataIndex: 'school' },
                { title: 'Available for Payout', key: 'verified', dataIndex: 'verified' },
                { title: 'Paid Out', key: 'verified', dataIndex: 'verified' },
                { title: '', key: 'delete', dataIndex: 'delete' },
              ]}
              size="small"
              style={{ width: '100%' }}
              pagination={false}
              // dataSource={profile.qualifications.map((quali: QualificationResponseDTO) => {
              //   return {
              //     degree: quali.degree,
              //     field: quali.field,
              //     school: quali.school,
              //     verified: quali.verified ? '\u2713' : '\u2717',
              //     delete: editQualis ? (
              //       <Button onClick={() => deleteQuali(quali.id)} style={{ margin: '0 0.5rem' }} size="small">
              //         <DeleteOutlined />
              //         Remove
              //       </Button>
              //     ) : (
              //       <></>
              //     ),
              // //   };
              // })}
            ></Table>
          </Col>
        </Row>
      </Content>
    </Typography>
  );
}
