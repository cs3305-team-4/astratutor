import React from 'react';
import styled from 'styled-components';

import { useAsync } from 'react-async-hook';

import {
  Typography,
  Layout,
  Row,
  Col,
  Card,
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

import { BankOutlined, CreditCardFilled, DeleteFilled } from '@ant-design/icons';

import { PaymentMethod } from '@stripe/stripe-js';
import { Elements, useStripe } from '@stripe/react-stripe-js';
import {
  AccountType,
  ProfileRequestDTO,
  BillingPayeePayment,
  BillingPayeesPaymentsResponseDTO,
} from '../api/definitions';
import { APIContext } from '../api/api';
import DefaultAvatar from '../assets/default_avatar.png';

import { Redirect, useHistory } from 'react-router-dom';

const { Title, Paragraph, Text, Link } = Typography;
const { Header, Footer, Sider, Content } = Layout;
const { TextArea } = Input;

export function BillingStudent(): React.ReactElement {
  const stripe = useStripe();
  const [error, setError] = React.useState<string>('');

  const [payeesPayments, setPayeesPayments] = React.useState<BillingPayeePayment[] | undefined>(undefined);
  const [ready, setReady] = React.useState<boolean>(false);

  const [cards, setCards] = React.useState<PaymentMethod[]>([]);
  const api = React.useContext(APIContext);
  const history = useHistory();

  const redirectCardSetupSession = async () => {
    const id = await api.services.createCardSetupSession(api.account.id, {
      success_path: '/account/billing',
      cancel_path: '/account/billing',
    });
    stripe.redirectToCheckout({
      sessionId: id,
    });
  };

  const reload = async () => {
    setCards(await api.services.readCardsByAccount(api.account.id));
    setPayeesPayments(await api.services.readPayeesPayments(api.account.id));
  };

  useAsync(async () => {
    await reload();
    setReady(true);
  }, []);

  const removeCard = async (pm: PaymentMethod) => {
    try {
      await api.services.deleteCardByAccount(api.account.id, pm.id);
      await reload();
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not retrieve your cards: ${e}`,
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
        </Row>
        <Row gutter={16} style={{ margin: '1rem' }}>
          <Col md={24} sm={24} xs={24}>
            <Title level={5}>Saved Cards</Title>
          </Col>
          {cards.map((pm: PaymentMethod, index: number) => (
            <Col key={pm.id}>
              <Card
                title={
                  <>
                    <CreditCardFilled style={{ marginRight: '0.5rem' }} />
                    {pm.card.brand.toUpperCase()}
                    <Button style={{ marginLeft: '1rem' }} size="small" onClick={() => removeCard(pm)}>
                      <DeleteFilled />
                      Remove
                    </Button>
                  </>
                }
                bordered
                style={{ boxShadow: '0 0 4px rgba(0,0,0,0.35)' }}
              >
                **** **** **** {pm.card.last4} &bull; {pm.card.exp_month}/{pm.card.exp_year} <br />
                <br />
                {pm.billing_details.name}
              </Card>
            </Col>
          ))}
          {cards.length === 0 && (
            <Alert
              message="No Debit/Credit Cards Saved"
              description="You have not yet setup billing with a debit or credit card, you will need to do this once you request a lesson."
              type="warning"
              showIcon
              style={{ margin: '1rem 0' }}
            />
          )}
        </Row>
        <Row gutter={16} style={{ margin: '0 1rem' }}>
          <Col md={24}>
            <Button style={{ margin: '0.5em 0' }} onClick={redirectCardSetupSession}>
              Add Card
            </Button>
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
                { title: 'Remarks', key: 'remarks', dataIndex: 'remarks' },
                // { title: 'Actions', key: 'actions', dataIndex: 'actions' },
              ]}
              size="small"
              style={{ width: '100%' }}
              pagination={false}
              dataSource={payeesPayments.map((payment: BillingPayeePayment) => {
                return {
                  ...payment,
                  amount: `€${payment.amount / 100}`,
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
