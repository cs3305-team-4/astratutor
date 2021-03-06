import React, { useEffect, useState } from 'react';

import { useHistory } from 'react-router';
import { UserAvatar } from './UserAvatar';

import {
  Typography,
  Layout,
  Row,
  Col,
  Avatar,
  PageHeader,
  Button,
  Statistic,
  Divider,
  Modal,
  Form,
  Input,
  DatePicker,
  message,
} from 'antd';

import {
  LessonResponseDTO,
  ProfileRequestDTO,
  ProfileResponseDTO,
  LessonRequestStage,
  AccountType,
} from '../api/definitions';

import { APIContext } from '../api/api';
import { useForm } from 'antd/lib/form/Form';
import { Link, useLocation } from 'react-router-dom';

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
  const history = useHistory();
  const profile = props.otherProfile;
  const lesson = props.lesson;

  const [showDenyModal, setShowDenyModal] = useState<boolean>(false);
  const [denyForm] = useForm();

  const [showCancelModal, setShowCancelModal] = useState<boolean>(false);
  const [cancelForm] = useForm();

  const [showRescheduleModal, setShowRescheduleModal] = useState<boolean>(false);
  const [rescheduleForm] = useForm();

  const query = new URLSearchParams(useLocation().search);

  const reload = async () => {
    props.onUpdate(await api.services.readLessonByAccountId(props.lesson.id), props.otherProfile);
  };

  const buttons = [];

  const requestPendingButton = (
    <>
      <Button type="text" disabled style={{ margin: '0.2rem' }}>
        Request Pending
      </Button>
    </>
  );

  const acceptButton = (
    <>
      <Button
        type="primary"
        style={{ margin: '0.2rem' }}
        onClick={async () => {
          try {
            await api.services.updateLessonStageAccept(props.lesson.id);
            await reload();
          } catch (e) {
            message.error('Failed to accept lesson! Please try again later.');
          }
        }}
      >
        Accept
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
            message.error('Failed to deny lesson! Please try again later.');
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
              message.error('Failed to cancel lesson! Please try again later.');
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
              message.error('Failed to reschedule lesson! Please try again later.');
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

  switch (props.lesson.request_stage) {
    case LessonRequestStage.Rescheduled:
    case LessonRequestStage.Requested:
      if (api.account.id !== lesson.request_stage_changer_id) {
        // If you arent the account who requested the lesson than you can accept/ deny
        buttons.push(acceptButton, denyButton);
      } else {
        // If you are the one who requested the lesson than you can cancel or see its pending
        buttons.push(cancelButton, requestPendingButton);
      }
      break;
    case LessonRequestStage.Accepted:
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
      buttons.push(cancelButton);
      break;
    default:
      buttons.push(rescheduleButton);
      break;
  }

  useEffect(() => {
    console.log(profile);
    if (query.has('reschedule') && query.get('reschedule') === lesson.id) {
      console.log(lesson.id);
      setShowRescheduleModal(true);
    }
  }, []);

  return (
    <PageHeader
      title={
        <>
          <Title level={5}>
            <Link to={api.account?.type === AccountType.Student && `/tutors/${profile.account_id}/profile`}>
              <UserAvatar props={{ size: 96 }} profile={profile}></UserAvatar>
            </Link>
          </Title>
          {`${profile.first_name} ${profile.last_name}`}
        </>
      }
      className="site-page-header-responsive"
      extra={[
        <Row key="stats" gutter={32} style={{ marginTop: '1rem' }} align="top" justify="start">
          <Col>
            <Statistic title="Subject" value={`L.C. - English`} />
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
