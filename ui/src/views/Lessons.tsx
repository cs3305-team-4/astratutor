import React, { useEffect } from 'react';
import styled from 'styled-components';

import { Badge, Layout, Menu } from 'antd';

import { MailOutlined, ScheduleOutlined, FileDoneOutlined, BankOutlined } from '@ant-design/icons';

import { APIContext } from '../api/api';
import { AccountType, LessonRequestStage, LessonResponseDTO, ProfileResponseDTO } from '../api/definitions';
import { useAsync } from 'react-async-hook';
import Lesson, { LessonProps } from '../components/Lesson';
import Link from 'antd/lib/typography/Link';
import { useLocation } from 'react-router-dom';

const { Sider, Content } = Layout;

enum Menus {
  Requests = 'Pending Requests',
  PaymentRequired = 'Payment Required',
  Scheduled = 'Scheduled',
  Completed = 'Completed',
  Cancelled = 'Cancelled',
  Denied = 'Denied',
}

const StyledLayout = styled(Layout)`
  background-color: white;
`;

const StyledSider = styled(Sider)`
  background-color: rgba(233, 233, 233);
`;

export function Lessons(): React.ReactElement {
  const api = React.useContext(APIContext);
  const [menu, setMenu] = React.useState<Menus>(Menus.Scheduled);

  const [lessonProps, setLessonProps] = React.useState<{ [uuid: string]: LessonProps }>({});

  const query = useLocation();

  useAsync(async () => {
    const lessons: LessonResponseDTO[] = await api.services.readLessonsByAccountId(api.account.id);

    const profiles: { [uuid: string]: ProfileResponseDTO } = {};

    const lprops: { [uuid: string]: LessonProps } = {};

    // go through the lessons and request any of the profiles we need to display each lesson
    for (const lesson of lessons) {
      if (!(lesson.student_id in profiles) && lesson.student_id !== api.account.id) {
        profiles[lesson.student_id] = await api.services.readProfileByAccountID(lesson.student_id, AccountType.Student);
      }

      if (!(lesson.tutor_id in profiles) && lesson.tutor_id !== api.account.id) {
        profiles[lesson.tutor_id] = await api.services.readProfileByAccountID(lesson.tutor_id, AccountType.Tutor);
      }

      lprops[lesson.id] = {
        lesson,
        otherProfile: lesson.student_id !== api.account.id ? profiles[lesson.student_id] : profiles[lesson.tutor_id],
        onUpdate: (lnew: LessonResponseDTO, otherProfile: ProfileResponseDTO) => {
          setLessonProps((lprops: { [uuid: string]: LessonProps }) => {
            const lpropsnew = { ...lprops };
            lpropsnew[lnew.id].lesson = lnew;
            lpropsnew[lnew.id].otherProfile = otherProfile;
            return lpropsnew;
          });
        },
      };
    }

    setLessonProps(lprops);
  }, []);

  useEffect(() => {
    switch (query.pathname) {
      case '/lessons/requests':
        setMenu(Menus.Requests);
        break;
      case '/lessons/scheduled':
        setMenu(Menus.Scheduled);
        break;
      case '/lessons/required':
        setMenu(Menus.PaymentRequired);
        break;
      case '/lessons/completed':
        setMenu(Menus.Completed);
        break;
      case '/lessons/cancelled':
        setMenu(Menus.Cancelled);
        break;
      case '/lessons/denied':
        setMenu(Menus.Denied);
        break;
    }
  }, [query]);

  return (
    <StyledLayout>
      <StyledSider>
        <Menu selectedKeys={[menu]} mode="inline">
          <Menu.Item
            disabled={
              Object.values(lessonProps).filter(
                (v) =>
                  v.lesson.request_stage === LessonRequestStage.Requested ||
                  v.lesson.request_stage === LessonRequestStage.Rescheduled,
              ).length === 0
            }
            title={
              Object.values(lessonProps).filter(
                (v) =>
                  v.lesson.request_stage === LessonRequestStage.Requested ||
                  v.lesson.request_stage === LessonRequestStage.Rescheduled,
              ).length === 0
                ? 'No requests right now!'
                : ''
            }
            onClick={() => setMenu(Menus.Requests)}
            key={Menus.Requests}
            icon={<MailOutlined />}
          >
            <Badge
              style={{ background: '#1890ff' }}
              offset={[72, 7]}
              count={
                Object.values(lessonProps).filter(
                  (v) =>
                    v.lesson.request_stage === LessonRequestStage.Requested ||
                    v.lesson.request_stage === LessonRequestStage.Rescheduled,
                ).length
              }
            >
              Requests
            </Badge>
          </Menu.Item>
          <Menu.Item
            disabled={
              Object.values(lessonProps).filter((v) => v.lesson.request_stage === LessonRequestStage.PaymentRequired)
                .length === 0
            }
            title={
              Object.values(lessonProps).filter((v) => v.lesson.request_stage === LessonRequestStage.PaymentRequired)
                .length === 0
                ? 'No payments required right now!'
                : ''
            }
            onClick={() => setMenu(Menus.PaymentRequired)}
            key={Menus.PaymentRequired}
            icon={<BankOutlined />}
          >
            <Badge
              style={{ background: '#1890ff' }}
              offset={[15, 7]}
              count={
                Object.values(lessonProps).filter((v) => v.lesson.request_stage === LessonRequestStage.PaymentRequired)
                  .length
              }
            >
              Payment Required
            </Badge>
          </Menu.Item>
          <Menu.Item onClick={() => setMenu(Menus.Scheduled)} key={Menus.Scheduled} icon={<ScheduleOutlined />}>
            Scheduled
          </Menu.Item>
          <Menu.ItemGroup title="Past Lessons">
            <Menu.Item
              style={{ color: 'green' }}
              onClick={() => setMenu(Menus.Completed)}
              key={Menus.Completed}
              icon={<FileDoneOutlined />}
            >
              Completed
            </Menu.Item>
            <Menu.Item
              style={{ color: 'red' }}
              onClick={() => setMenu(Menus.Cancelled)}
              key={Menus.Cancelled}
              icon={<FileDoneOutlined />}
            >
              Cancelled
            </Menu.Item>
            <Menu.Item
              style={{ color: '#424242' }}
              onClick={() => setMenu(Menus.Denied)}
              key={Menus.Denied}
              icon={<FileDoneOutlined />}
            >
              Denied
            </Menu.Item>
          </Menu.ItemGroup>
        </Menu>
      </StyledSider>
      <Content style={{ minHeight: 'calc(100vh - 72px)' }}>
        {menu === Menus.Requests &&
          Object.keys(lessonProps).map((key: string, index: number) => {
            if (
              lessonProps[key].lesson.request_stage === LessonRequestStage.Requested ||
              lessonProps[key].lesson.request_stage === LessonRequestStage.Rescheduled
            ) {
              return <Lesson key={index} {...lessonProps[key]} />;
            }

            return <></>;
          })}
        {menu === Menus.PaymentRequired &&
          Object.keys(lessonProps).map((key: string, index: number) => {
            if (lessonProps[key].lesson.request_stage === LessonRequestStage.PaymentRequired) {
              return <Lesson key={index} {...lessonProps[key]} />;
            }

            return <></>;
          })}
        {menu === Menus.Scheduled &&
          Object.keys(lessonProps).map((key: string, index: number) => {
            if (lessonProps[key].lesson.request_stage === LessonRequestStage.Scheduled) {
              return <Lesson key={index} {...lessonProps[key]} />;
            }

            return <></>;
          })}
        {menu === Menus.Completed &&
          Object.keys(lessonProps).map((key: string, index: number) => {
            if (lessonProps[key].lesson.request_stage === LessonRequestStage.Completed) {
              return <Lesson key={index} {...lessonProps[key]} />;
            }

            return <></>;
          })}
        {menu === Menus.Cancelled &&
          Object.keys(lessonProps).map((key: string, index: number) => {
            if (lessonProps[key].lesson.request_stage === LessonRequestStage.Cancelled) {
              return <Lesson key={index} {...lessonProps[key]} />;
            }

            return <></>;
          })}
        {menu === Menus.Denied &&
          Object.keys(lessonProps).map((key: string, index: number) => {
            if (lessonProps[key].lesson.request_stage === LessonRequestStage.Denied) {
              return <Lesson key={index} {...lessonProps[key]} />;
            }

            return <></>;
          })}
        {/* {menu === Menus.Scheduled && (

        )}
        ap((lessonProp: LessonProps, index: number) => <Lesson key={index} {...lessonProp} />)}
        {menu === Menus.Completed && (

        )} */}
      </Content>
    </StyledLayout>
  );
}
