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

import {
  AccountType,
  BillingPayerPayment,
  BillingPayersPaymentsResponseDTO,
  BillingPayoutInfoResponseDTO,
  ProfileRequestDTO,
} from '../api/definitions';
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
  const [payoutInfo, setPayoutInfo] = React.useState<BillingPayoutInfoResponseDTO | undefined>(undefined);
  const [payersPayments, setPayersPayments] = React.useState<BillingPayerPayment[] | undefined>(undefined);
  const [ready, setReady] = React.useState<boolean>(false);
  const api = React.useContext(APIContext);
  const history = useHistory();

  const reload = async () => {
    const onboard = await api.services.readTutorBillingOnboardStatus(api.account.id);
    setBillingOnboard(onboard);

    if (!onboard) {
      setBillingOnboardUrl(await api.services.readTutorBillingOnboardUrl(api.account.id));
      setReady(true);
    } else {
      const reqsMet = await api.services.readTutorBillingRequirementsMetStatus(api.account.id);
      setBillingRequirementsMet(reqsMet);
      setBillingPanelUrl(await api.services.readTutorBillingPanelUrl(api.account.id));
      setPayoutInfo(await api.services.readPayoutInfo(api.account.id));
      setPayersPayments(await api.services.readPayersPayments(api.account.id));
      setReady(true);
    }
  };

  useAsync(async () => {
    await reload();
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

  const payout = async () => {
    try {
      await api.services.createPayout(api.account.id);
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not payout: ${e}`,
      });
    }
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
              value={`€${payoutInfo !== undefined ? payoutInfo.payout_balance / 100 : '-.--'}`}
            />
          </Col>
          {/* <Col md={8} sm={24} xs={24}>
            <Statistic valueStyle={{ fontSize: '4rem' }} title="Total This Month" value={'€0.00'} />
          </Col>
          <Col md={8} sm={24} xs={24}>
            <Statistic valueStyle={{ fontSize: '4rem' }} title="Total This Year" value={'€0.00'} />
          </Col> */}
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
                <Button
                  type="primary"
                  onClick={payout}
                  style={{ margin: '0.5em' }}
                  disabled={payoutInfo !== undefined && payoutInfo.payout_balance <= 0}
                >
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
                { title: 'Description', key: 'description', dataIndex: 'description' },
                { title: 'Date', key: 'date', dataIndex: 'date' },
                { title: 'Amount', key: 'amount', dataIndex: 'amount' },
                { title: 'Available for Payout', key: 'available_for_payout', dataIndex: 'available_for_payout' },
                { title: 'Paid Out', key: 'paid_out', dataIndex: 'paid_out' },
                { title: 'Actions', key: 'actions', dataIndex: 'actions' },
              ]}
              size="small"
              style={{ width: '100%' }}
              pagination={false}
              dataSource={payersPayments.map((payment: BillingPayerPayment) => {
                return {
                  ...payment,
                  amount: `€${payment.amount / 100}`,
                  available_for_payout: payment.available_for_payout ? '\u2713' : '\u2717',
                  paid_out: payment.paid_out ? '\u2713' : '\u2717',
                  date: new Intl.DateTimeFormat('en-IE', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric',
                    weekday: 'long',
                    hour: 'numeric',
                    minute: 'numeric',
                  }).format(new Date(payment.date)),
                  actions: <></>,
                };
              })}
            ></Table>
          </Col>
        </Row>
      </Content>
    </Typography>
  );
}
