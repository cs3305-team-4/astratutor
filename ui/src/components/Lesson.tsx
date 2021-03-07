import React, { useState } from 'react';

import { useHistory } from 'react-router';
import { UserAvatar } from './UserAvatar';

import { PaymentIntent, PaymentMethod, CreatePaymentMethodCardData } from '@stripe/stripe-js';

import { useAsync } from 'react-async-hook';

import {
  Typography,
  Layout,
  Row,
  Col,
  Avatar,
  PageHeader,
  Button,
  Card,
  Statistic,
  Divider,
  Modal,
  Form,
  Input,
  DatePicker,
  Radio,
  RadioChangeEvent,
} from 'antd';

import { CreditCardFilled, CreditCardOutlined } from '@ant-design/icons';

import { Elements, useStripe, CardElement, useElements } from '@stripe/react-stripe-js';

import {
  LessonResponseDTO,
  ProfileRequestDTO,
  ProfileResponseDTO,
  LessonRequestStage,
  SubjectTaughtDTO,
  AccountType,
} from '../api/definitions';

import { APIContext } from '../api/api';
import { useForm } from 'antd/lib/form/Form';
import { Link } from 'react-router-dom';

const { Title } = Typography;
const { Content } = Layout;

export interface LessonProps {
  // Lesson object
  lesson: LessonResponseDTO;

  // Profile of the person who isn't the currently signed in user
  // i.e the student if the person signed in is the tutor
  otherProfile: ProfileResponseDTO;

  onUpdate(lesson: LessonResponseDTO, otherProfile: ProfileRequestDTO): void;
}

export default function Lesson(props: LessonProps): React.ReactElement {
  const api = React.useContext(APIContext);
  const stripe = useStripe();
  const elements = useElements();
  const history = useHistory();
  const profile = props.otherProfile;
  const lesson = props.lesson;

  const [showDenyModal, setShowDenyModal] = useState<boolean>(false);
  const [denyForm] = useForm();

  const [showCancelModal, setShowCancelModal] = useState<boolean>(false);
  const [cancelForm] = useForm();

  const [showRescheduleModal, setShowRescheduleModal] = useState<boolean>(false);
  const [rescheduleForm] = useForm();

  const [showPayModal, setShowPayModal] = useState<boolean>(false);
  const [selectedCard, setSelectedCard] = useState<number>(0);
  const [cards, setCards] = useState<PaymentMethod[]>([]);
  const [paymentIntent, setPaymentIntent] = useState<PaymentIntent | undefined>(undefined);
  const [showNewCard, setShowNewCard] = useState<boolean>(false);
  //const [subject, setSubject] = useState<SubjectResponseDTO | undefined>(undefined);

  const reload = async () => {
    props.onUpdate(await api.services.readLessonByAccountId(props.lesson.id), props.otherProfile);
  };

  useAsync(async () => {
    switch (props.lesson.request_stage) {
      case LessonRequestStage.PaymentRequired: {
        // Ask the server to refresh the paid status (used to check if the user paid their bill on return from checkout)
        // We do this by requesting the server to advance the lesson to scheduled
        await api.services.updateLessonStageScheduled(props.lesson.id);
        await reload();
        break;
      }
    }
  }, []);

  const buttons = [];

  const requestPendingButton = (
    <>
      <Button type="text" disabled style={{ margin: '0.2rem' }}>
        Request Pending
      </Button>
    </>
  );

  const paymentPendingButton = (
    <>
      <Button type="text" disabled style={{ margin: '0.2rem' }}>
        Awaiting Payment from Student
      </Button>
    </>
  );

  const acceptButton = (
    <>
      <Button
        type="primary"
        style={{ margin: '0.2rem' }}
        onClick={async () => {
          await api.services.updateLessonStagePaymentRequired(props.lesson.id);
          await reload();
        }}
      >
        Accept &amp; Request Payment
      </Button>
    </>
  );

  const denyButton = (
    <>
      <Button style={{ margin: '0.2rem' }} danger onClick={() => setShowDenyModal(true)}>
        Deny
      </Button>
      <Modal
        title="Deny Request"
        visible={showDenyModal}
        okText="Deny"
        okType="danger"
        onOk={async () => {
          try {
            await rescheduleForm.validateFields().then(async (values) => {
              await api.services.updateLessonStageDeny(props.lesson.id, values);
              await reload();
              setShowDenyModal(false);
            });
          } catch (e) {
            Modal.error({
              title: 'Error',
              content: `Failed to deny lesson! Please try again later.`,
            });
          }
        }}
        cancelText="Back"
        onCancel={() => setShowDenyModal(false)}
      >
        <Form form={denyForm} layout="vertical">
          <Form.Item
            label="Reason"
            name="reason"
            rules={[{ required: true, message: 'Please give a reason for denying the request' }]}
          >
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );

  const cancelButton = (
    <>
      <Button style={{ margin: '0.2rem' }} danger onClick={() => setShowCancelModal(true)}>
        Cancel
      </Button>
      <Modal
        title="Cancel Lesson"
        visible={showCancelModal}
        okText="Cancel"
        okType="danger"
        onOk={async () => {
          await rescheduleForm.validateFields().then(async (values) => {
            try {
              await api.services.updateLessonStageCancel(props.lesson.id, values);
              await reload();
              setShowCancelModal(false);
            } catch (e) {
              Modal.error({
                title: 'Error',
                content: `Failed to cancel lesson! Please try again later.`,
              });
            }
          });
        }}
        cancelText="Back"
        onCancel={() => setShowCancelModal(false)}
      >
        <Form form={cancelForm} layout="vertical">
          <Form.Item
            label="Reason"
            name="reason"
            rules={[{ required: true, message: 'Please give a reason for cancelling the request' }]}
          >
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );

  const rescheduleButton = (
    <>
      <Button style={{ margin: '0.2rem' }} type="dashed" onClick={() => setShowRescheduleModal(true)}>
        Reschedule
      </Button>
      <Modal
        title="Reschedule Lesson"
        visible={showRescheduleModal}
        okText="Reschedule"
        onOk={async () => {
          await rescheduleForm.validateFields().then(async (values) => {
            try {
              values.new_time = values.new_time
                .set({
                  minute: 0,
                  second: 0,
                  millisecond: 0,
                })
                .toISOString();
              await api.services.updateLessonStageReschedule(props.lesson.id, values);
              await reload();
              setShowRescheduleModal(false);
            } catch (e) {
              Modal.error({
                title: 'Error',
                content: `Failed to reschedule lesson! Please try again later.`,
              });
            }
          });
        }}
        cancelText="Back"
        onCancel={() => setShowRescheduleModal(false)}
      >
        <Form form={rescheduleForm} layout="vertical">
          <Form.Item label="New Time" name="new_time" rules={[{ required: true, message: 'Please select a date!' }]}>
            <DatePicker style={{ width: '100%' }} size="large" showTime format={'YYYY-MM-DD HH:00:00'} />
          </Form.Item>
          <Form.Item
            label="Reason"
            name="reason"
            rules={[{ required: true, message: 'Please give a reason for rescheduling the request' }]}
          >
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );

  const pay = async ({ cardName }: { cardName?: string }) => {
    const secret_id = await api.services.readLessonBillingPaymentIntentSecret(props.lesson.id);

    // if they selected the add a card option use that info or else the saved id of a previous payment method
    const payment_method:
      | string
      | Pick<CreatePaymentMethodCardData, 'card' | 'billing_details' | 'metadata' | 'payment_method'> =
      selectedCard === cards.length
        ? {
            card: elements.getElement(CardElement),
            billing_details: {
              name: cardName,
            },
          }
        : cards[selectedCard].id;

    console.log(selectedCard, payment_method, cards);
    try {
      const { error, paymentIntent } = await stripe.confirmCardPayment(secret_id, {
        payment_method,
        setup_future_usage: 'off_session',
        receipt_email: api.account.email,
      });

      if (error !== undefined) {
        Modal.error({
          title: 'Error',
          content: `Could not complete transaction : ${error.message}`,
        });
      } else {
        // Try to mark the lesson as scheduled and reload
        await api.services.updateLessonStageScheduled(props.lesson.id);
        await reload();

        Modal.info({
          title: 'Lesson Scheduled',
          content: `The lesson has now been paid for and scheduled, you can visit the Scheduled Lessons page to access information prior to the lesson`,
        });
      }

      setShowPayModal(false);
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not complete transaction : ${e}`,
      });
    }
  };

  const payButton = (
    <>
      <Button
        style={{ margin: '0.2rem' }}
        type="dashed"
        onClick={async () => {
          const cards = await api.services.readCardsByAccount(api.account.id);
          const secret_id = await api.services.readLessonBillingPaymentIntentSecret(props.lesson.id);
          setPaymentIntent(await (await stripe.retrievePaymentIntent(secret_id)).paymentIntent);
          setCards(cards);
          setSelectedCard(0);
          setShowPayModal(true);
        }}
      >
        Pay
      </Button>
      <Modal
        title="Pay for Lesson"
        visible={showPayModal}
        onCancel={() => setShowPayModal(false)}
        footer={[
          <Button form="pay" key="submit" style={{ width: '100%' }} type="primary" htmlType="submit">
            Pay
          </Button>,
        ]}
      >
        <Form onFinish={pay} layout="vertical" name="pay" preserve={false}>
          <Title level={5}>Payment</Title>
          <Typography.Paragraph>Select payment method to continue:</Typography.Paragraph>
          <Form.Item name="card" rules={[{ required: true, message: 'Please supply a card!' }]}>
            <Radio.Group
              value={selectedCard}
              onChange={(e: RadioChangeEvent) => {
                setSelectedCard(e.target.value);
              }}
            >
              {cards.map((pm: PaymentMethod, index: number) => (
                <Radio key={index} style={{ display: 'block' }} value={index}>
                  <CreditCardFilled style={{ marginRight: '0.5rem' }} />
                  {pm.card.brand.toUpperCase()} **** **** **** {pm.card.last4} &bull; {pm.card.exp_month}/
                  {pm.card.exp_year} &bull; {pm.billing_details.name}
                </Radio>
              ))}
              <Radio key={cards.length} style={{ display: 'block' }} value={cards.length}>
                <CreditCardOutlined style={{ marginRight: '0.5rem' }} />
                Use new card
              </Radio>
            </Radio.Group>
          </Form.Item>
          {selectedCard === cards.length && (
            <>
              <Form.Item
                label="Name on card"
                name="cardName"
                rules={[{ required: selectedCard === cards.length, message: 'Please supply the name on the card!' }]}
              >
                <Input />
              </Form.Item>
              <Form.Item>
                <CardElement
                  options={{
                    hidePostalCode: true,
                    style: {
                      base: {
                        color: 'black',
                        fontFamily:
                          "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji'",
                        fontSize: '16px',
                        '::placeholder': {
                          color: '#aab7c4',
                        },
                      },
                      invalid: {
                        color: '#9e2146',
                      },
                    },
                  }}
                />
              </Form.Item>
            </>
          )}
        </Form>

        {paymentIntent !== undefined && (
          <Title level={2} style={{ textAlign: 'center', margin: '1rem 0 0 0' }}>
            {paymentIntent.amount / 100} {paymentIntent.currency.toUpperCase()}
          </Title>
        )}
      </Modal>
    </>
  );

  switch (props.lesson.request_stage) {
    case LessonRequestStage.Rescheduled:
    case LessonRequestStage.Requested:
      if (api.account.id !== lesson.request_stage_changer_id) {
        // If you arent the account who requested the lesson than you can accept/ deny
        buttons.push(acceptButton, denyButton, rescheduleButton);
      } else {
        // If you are the one who requested the lesson than you can cancel or see its pending
        buttons.push(cancelButton, requestPendingButton);
      }
      break;

    case LessonRequestStage.Scheduled:
      buttons.push(
        <Button
          ghost
          type="primary"
          key="enter classroom"
          style={{ margin: '0.2rem' }}
          onClick={() => {
            history.push(`/lessons/${props.lesson.id}/lobby`);
          }}
        >
          Enter Classroom
        </Button>,
      );
      // Once the lesson is accepted either party can cancel or reschedule a lesson
      buttons.push(cancelButton, rescheduleButton);
      break;

    case LessonRequestStage.PaymentRequired:
      if (api.account.type == AccountType.Student) {
        buttons.push(payButton, cancelButton);
      } else {
        buttons.push(paymentPendingButton, cancelButton);
      }
      break;
  }

  return (
    <PageHeader
      title={
        <>
          <Title level={5}>
            <UserAvatar props={{ size: 96 }} profile={profile}></UserAvatar>
          </Title>
          {`${profile.first_name} ${profile.last_name}`}
        </>
      }
      className="site-page-header-responsive"
      extra={[
        <Row key="stats" gutter={32} style={{ marginTop: '1rem' }} align="top" justify="start">
          <Col>
            <Statistic title="Subject" value={props.otherProfile.s} />
          </Col>
          <Col>
            <Statistic
              title="Time"
              value={`${new Intl.DateTimeFormat('en-IE', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                weekday: 'long',
                hour: 'numeric',
                minute: 'numeric',
              }).format(new Date(props.lesson.start_time))}`}
            />
          </Col>
        </Row>,
        <Row key="buttons" gutter={16} align="bottom" justify="end" style={{ margin: '0.5rem 0.2rem' }}>
          {buttons}
        </Row>,
      ]}
    >
      <Typography>
        <Content>{lesson.lesson_detail}</Content>
      </Typography>
      <Divider />
    </PageHeader>
  );
}
