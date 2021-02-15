import React, { useContext, useState } from 'react';
import moment, { Moment } from 'moment';

import { useHistory } from 'react-router';

import { Modal, ModalProps, Button, Form, Select, Input, Typography, DatePicker, TimePicker, Row, Col } from 'antd';

import { Availability } from './Availability';
import { AccountType, ProfileResponseDTO, SubjectTaughtDTO } from '../api/definitions';
import { APIContext } from '../api/api';
import { useAsync } from 'react-async-hook';

const { Title } = Typography;
const { TextArea } = Input;

export interface RequestLessonModalProps extends ModalProps {
  profile: ProfileResponseDTO;
  type: AccountType;
}

export function RequestLessonModal(props: RequestLessonModalProps): React.ReactElement {
  //const subjects = ['Leaving Certificate - English', 'Leaving Certificate - Irish'];

  const [tutorSubjects, setTutorSubjects] = useState<SubjectTaughtDTO[] | undefined>(undefined);

  const api = React.useContext(APIContext);
  const history = useHistory();

  interface FormModel {
    start_time: Moment;
    subject: string;
    lesson_detail: string;
  }

  useAsync(async () => {
    setTutorSubjects(await api.services.readTutorSubjectsByAccountId(props.profile.account_id));
  }, []);

  const onFinish = async (values: FormModel) => {
    console.log(values);

    const tutor_id = api.account.type === AccountType.Tutor ? api.account.id : props.profile.account_id;
    const student_id = api.account.type === AccountType.Student ? api.account.id : props.profile.account_id;

    console.log(tutor_id, student_id);
    try {
      await api.services.createLesson({
        start_time: values.start_time
          .set({
            minute: 0,
            second: 0,
            millisecond: 0,
          })
          .toISOString(),
        tutor_id,
        student_id,
        lesson_detail: values.lesson_detail,
      });

      Modal.info({
        title: 'Info',
        content: 'Your lesson request has been made, you will be notified when the tutor accepts or denies the request',
        onOk: () => history.push('/lessons'),
      });

      props.onOk(null);
    } catch (e) {
      Modal.error({
        title: 'Error',
        content: `Could not create lesson: ${e}`,
      });
    }
  };

  return (
    <Modal
      {...props}
      width={'960px'}
      title="Request Lesson"
      footer={[
        <Button form="request-lesson" key="submit" style={{ width: '100%' }} type="primary" htmlType="submit">
          Request
        </Button>,
      ]}
    >
      <Form
        onFinish={onFinish}
        initialValues={{ subject: tutorSubjects ? tutorSubjects[0].id : undefined }}
        layout="vertical"
        name="request-lesson"
        preserve={false}
      >
        <Row gutter={16}>
          <Col xs={24} sm={24} md={12}>
            <Row gutter={16}>
              <Col xs={24}>
                <Title level={5}>
                  {props.profile.first_name} {props.profile.last_name}&lsquo;s Availability
                </Title>
                {props.type == AccountType.Tutor && (
                  <Availability
                    hideUnavailable={true}
                    availability={props.profile.availability}
                    editable={false}
                  ></Availability>
                )}
              </Col>
              <Col xs={24}>
                <Form.Item name="start_time" rules={[{ required: true, message: 'Please select a date!' }]}>
                  <DatePicker style={{ width: '100%' }} size="large" showTime format={'YYYY-MM-DD HH:00:00'} />
                </Form.Item>
              </Col>
            </Row>
          </Col>
          <Col xs={24} sm={24} md={12}>
            <Title level={5}>Subject</Title>
            <Form.Item name="subject" rules={[{ required: true, message: 'Please select a subject!' }]}>
              <Select size="large" style={{ width: '100%' }}>
                {tutorSubjects?.map((subject, index) => (
                  <Select.Option key={index} value={subject.id}>
                    {subject.name}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
            <Title level={5}>What would you like the lesson to focus on?</Title>
            <Form.Item
              name="lesson_detail"
              rules={[{ required: true, message: 'Please enter a lesson request description!' }]}
            >
              <TextArea
                maxLength={1200}
                placeholder="Describe what you would like to achieve with the lesson"
                style={{ minHeight: '240px', margin: '0.5rem 0' }}
                size="large"
              />
            </Form.Item>
          </Col>
        </Row>
      </Form>
    </Modal>
  );
}
