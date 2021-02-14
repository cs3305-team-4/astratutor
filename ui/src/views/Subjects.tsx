import React, { ChangeEvent, ReactElement, useEffect, useState } from 'react';
import styled from 'styled-components';

import { Link } from 'react-router-dom';

import Subject from '../components/Subject';

import { Typography, Layout, Card, Row, Col, PageHeader, Input } from 'antd';
import { ReadSubjectsDTO, SubjectDTO } from '../api/definitions';

import { useAsync } from 'react-async-hook';

import { fetchRest } from '../api/rest';
import config from '../config';

const { Title, Paragraph, Text } = Typography;
const { Header, Footer, Sider, Content } = Layout;

export function Subjects(): ReactElement {
  const [error, setError] = useState<string | undefined>(undefined);

  const [subjects, setSubjects] = useState<SubjectDTO[] | undefined>(undefined);
  const [search, setSearch] = useState<string>('');

  const onSearch = (el: ChangeEvent<HTMLInputElement>) => {
    setSearch(el.target.value);
  };

  useEffect(() => {
    setSubjects([
      {
        name: 'Maths',
        slug: 'maths',
        image: '',
      },
      {
        name: 'English Higher Level',
        slug: 'english-higher-level',
        image: '',
      },
      {
        name: 'History',
        slug: 'history',
        image: '',
      },
      {
        name: 'Chemistry',
        slug: 'chemistry',
        image: '',
      },
      {
        name: 'Engineering',
        slug: 'engineering',
        image: '',
      },
      {
        name: 'Computer Science',
        slug: 'computer-science',
        image: '',
      },
      {
        name: 'Chemical Engineering',
        slug: 'chemical-engineering',
        image: '',
      },
      {
        name: 'Biology',
        slug: 'biology',
        image: '',
      },
      {
        name: 'Arts',
        slug: 'arts',
        image: '',
      },
    ]);
  }, []);

  // useAsync(async () => {
  //   try {
  //     const res = await fetchRest(`${config.apiUrl}/subjects`);
  //     const subjects = (await res.json()) as ReadSubjectsDTO;
  //     setSubjects(subjects);
  //   } catch (e) {
  //     setError('Failed to load subjects.');
  //   }
  // }, []);

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
  });

  return (
    <Content style={{ padding: '2em 0' }}>
      <Row>
        <Col xl={{ offset: 4, span: 16 }} lg={{ offset: 2, span: 20 }} span={24}>
          <PageHeader
            title="Subjects"
            extra={[<Input.Search key="1" placeholder="Search for a subject" allowClear onChange={onSearch} />]}
          />
          <Row gutter={[16, 16]}>{displaySubjects}</Row>
        </Col>
      </Row>
    </Content>
  );
}
