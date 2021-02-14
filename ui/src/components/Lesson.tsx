import React from 'react';
import moment from 'moment';

import { useHistory } from 'react-router';

import { Typography, Layout, Row, Col, Avatar, PageHeader, Button, Statistic, Divider } from 'antd';

import { LessonResponseDTO, ProfileRequestDTO, ProfileResponseDTO, LessonRequestStage } from '../api/definitions';

import { APIContext } from '../api/api';

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
          <Button
            style={{ margin: '0.2rem' }}
            onClick={async () => {
              await api.services.updateLessonStageDeny(props.lesson.id, 'Lesson denied');
              await reload();
            }}
          >
            Deny
          </Button>,
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
        <Button
          style={{ margin: '0.2rem' }}
          onClick={async () => {
            await api.services.updateLessonStageCancel(props.lesson.id, 'Lesson cancelled');
            await reload();
          }}
        >
          Cancel
        </Button>,
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
            <Statistic title="Time" value={`${moment(props.lesson.start_time).format('LLLL')}`} />
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
