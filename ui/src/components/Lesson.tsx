import React, { useState } from 'react';

import { useHistory } from 'react-router';

import { Typography, Layout, Row, Col, Avatar, PageHeader, Button, Statistic, Divider, Modal, Form, Input } from 'antd';

import { LessonResponseDTO, ProfileRequestDTO, ProfileResponseDTO, LessonRequestStage } from '../api/definitions';

import { APIContext } from '../api/api';
import { useForm } from 'antd/lib/form/Form';

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

  const reload = async () => {
    props.onUpdate(await api.services.readLessonByAccountId(props.lesson.id), props.otherProfile);
  };

  const buttons = [];
  buttons.push(
    <Button
      key="enter classroom"
      style={{ margin: '0.2rem' }}
      onClick={() => {
        history.push(`/lessons/${props.lesson.id}/classroom`);
      }}
    >
      Enter Classroom
    </Button>,
  );

  switch (props.lesson.request_stage) {
    case LessonRequestStage.Requested:
      if (api.account.id !== lesson.requester_id) {
        buttons.push(
          <Button
            type="primary"
            key="enter classroom"
            style={{ margin: '0.2rem' }}
            onClick={async () => {
              await api.services.updateLessonStageAccept(props.lesson.id);
              await reload();
            }}
          >
            Accept
          </Button>,
        );

        buttons.push(
          <>
            <Button style={{ margin: '0.2rem' }} onClick={() => setShowDenyModal(true)}>
              Deny
            </Button>
            <Modal
              title="Deny Request"
              visible={showDenyModal}
              okText="Deny"
              okType="danger"
              onOk={async () => {
                await api.services.updateLessonStageDeny(props.lesson.id, { reason: denyForm.getFieldValue('reason') });
                await reload();
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
          </>,
        );
      } else {
        buttons.push(
          <Button type="primary" disabled style={{ margin: '0.2rem' }}>
            Request Pending
          </Button>,
        );
      }
      break;

    case LessonRequestStage.Accepted:
      buttons.push(
        <>
          <Button style={{ margin: '0.2rem' }} onClick={() => setShowCancelModal(true)}>
            Cancel
          </Button>
          <Modal
            title="Cancel Lesson"
            visible={showCancelModal}
            okText="Cancel"
            okType="danger"
            onOk={async () => {
              await cancelForm.validateFields().then(async (values) => {
                await api.services.updateLessonStageCancel(props.lesson.id, values);
                await reload();
              });
            }}
            cancelText="Back"
            onCancel={() => setShowCancelModal(false)}
          >
            <Form form={cancelForm} layout="vertical">
              <Form.Item
                label="Reason"
                name="reason"
                rules={[{ required: true, message: 'Please give a reason for cancelling the lesson' }]}
              >
                <Input />
              </Form.Item>
            </Form>
          </Modal>
        </>,
      );
      break;
  }

  return (
    <PageHeader
      title={
        <>
          <Title level={5}>
            <Avatar src={profile.avatar} size={96}></Avatar>
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
        <Row key="buttons" gutter={16} align="top" justify="end" style={{ margin: '0.5rem 0.2rem' }}>
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
