import React, { ChangeEvent, ReactElement, useContext, useState } from 'react';

import { APIContext } from '../api/api';

import Subject from '../components/Subject';

import { Typography, Layout, Row, Col, PageHeader, Input, Space } from 'antd';
import { SubjectDTO } from '../api/definitions';

import { useAsync } from 'react-async-hook';
const { Title, Paragraph } = Typography;

const { Text } = Typography;
const { Content } = Layout;

export function Subjects(): ReactElement {
  const api = useContext(APIContext);
  const [error, setError] = useState<string | undefined>(undefined);

  const [subjects, setSubjects] = useState<SubjectDTO[] | undefined>(undefined);
  const [search, setSearch] = useState<string>('');

  const onSearch = (el: ChangeEvent<HTMLInputElement>) => {
    setSearch(el.target.value);
  };

  useAsync(async () => {
    try {
      setSubjects(await api.services.readSubjects());
    } catch (e) {
      setError('Failed to load subjects.');
    }
  }, []);

  if (error !== undefined) {
    return (
      <Layout>
        <Text>{error}</Text>
      </Layout>
    );
  }

  const displaySubjects = subjects?.map((subject, index) => {
    if (subject.name.includes(search))
      return (
        <Col key={index} xxl={8} md={12} xs={24}>
          <Subject subject={subject} />
        </Col>
      );
    return undefined;
  });

  return (
    <Content style={{ padding: '2em 0' }}>
      <Row>
        <Col xl={{ offset: 4, span: 16 }} lg={{ offset: 2, span: 20 }} span={24}>
          <Row justify="space-between">
            <Title>Tutors</Title>
            <Space>
              <Input.Search key="1" placeholder="Search for a subject" allowClear onChange={onSearch} />
            </Space>
          </Row>
          <Row>{displaySubjects}</Row>
        </Col>
      </Row>
    </Content>
  );
}
