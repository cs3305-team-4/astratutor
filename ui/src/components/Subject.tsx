import React from 'react';
import styled from 'styled-components';

import { Link } from 'react-router-dom';

import { Typography, Layout, Card, Row, Col, Image } from 'antd';
import { SubjectDTO } from '../api/definitions';

const { Title, Paragraph, Text } = Typography;
const { Header, Footer, Sider, Content } = Layout;

export interface SubjectProps {
  subject: SubjectDTO;
}

export default function Subject(props: SubjectProps): ReactElement {
  return (
    <Link to={`/subjects/tutors?filter=${props.subject.slug}`}>
      <Row
        style={{
          height: '150px',
          overflow: 'hidden',
          backgroundColor: 'white',
          border: '1px solid rgba(0,0,0,0.1)',
          borderRadius: 'px',
        }}
        align="middle"
      >
        <Col span={24} style={{ padding: '20px' }}>
          <Title>{props.subject.name}</Title>
        </Col>
        {/* <Col span={8}>
          <img height="150px" src={props.subject.image} alt={props.subject.name} />
        </Col> */}
      </Row>
    </Link>
  );
}
